package mapper

import (
	"fmt"

	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/eco/types"
	qtypes "github.com/QOSGroup/qos/types"
)

func BuildDistributionStoreQueryPath() []byte {
	return []byte(fmt.Sprintf("/store/%s/key", types.DistributionMapperName))
}

func (mapper *DistributionMapper) InitValidatorPeriodSummaryInfo(valAddr btypes.Address) types.ValidatorCurrentPeriodSummary {
	//初始化validator历史收益,当前周期汇总收益.
	mapper.Set(types.BuildValidatorHistoryPeriodSummaryKey(valAddr, uint64(0)), qtypes.ZeroFraction())
	current := types.ValidatorCurrentPeriodSummary{
		Period: 1,
		Fees:   btypes.ZeroInt(),
	}
	mapper.Set(types.BuildValidatorCurrentPeriodSummaryKey(valAddr), current)
	return current
}

//清空validator相关的收益汇总信息
func (mapper *DistributionMapper) DeleteValidatorPeriodSummaryInfo(valAddr btypes.Address) {
	periodPrifixKey := append(types.GetValidatorHistoryPeriodSummaryPrefixKey(), valAddr...)
	iter := store.KVStorePrefixIterator(mapper.GetStore(), periodPrifixKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		mapper.Del(iter.Key())
	}

	k := types.BuildValidatorCurrentPeriodSummaryKey(valAddr)
	mapper.Del(k)
}

func (mapper *DistributionMapper) InitDelegatorIncomeInfo(valAddr, deleAddr btypes.Address, bondTokens, currHeight uint64) {
	//初始化delegaotr starting , 发放收益信息
	vcps, _ := mapper.GetValidatorCurrentPeriodSummary(valAddr)
	params := mapper.GetParams()

	key := types.BuildDelegatorEarningStartInfoKey(valAddr, deleAddr)

	startInfo := types.DelegatorEarningsStartInfo{
		PreviousPeriod:       vcps.Period - 1,
		BondToken:            bondTokens,
		StartingHeight:       currHeight,
		HistoricalRewardFees: btypes.ZeroInt(),
	}

	//delegator unbond全部后,又重新delegate
	var info types.DelegatorEarningsStartInfo
	exsits := mapper.Get(key, &info)
	if exsits { //保留delegator历史收益,不计算阶段内收益
		startInfo.HistoricalRewardFees = startInfo.HistoricalRewardFees.Add(info.HistoricalRewardFees)
	}

	mapper.Set(key, startInfo)

	//发放收益高度
	incomeHeight := currHeight + params.DelegatorsIncomePeriodHeight
	mapper.Set(types.BuildDelegatorPeriodIncomeKey(valAddr, deleAddr, incomeHeight), true)
}

//删除delegator收益信息
//todo: 发放收益高度信息没有删除
func (mapper *DistributionMapper) DeleteDelegatorIncomeInfo(valAddr, deleAddr btypes.Address) {
	k := types.BuildDelegatorEarningStartInfoKey(valAddr, deleAddr)
	mapper.Del(k)
}

//增加validator的周期
//1. 保存历史周期
//2. 更新当前周期,返回上一周期数值
func (mapper *DistributionMapper) incrementValidatorPeriod(validator types.Validator) uint64 {
	valAddr := validator.GetValidatorAddress()
	//初始化delegaotr starting , 发放收益信息
	vcps, exsits := mapper.GetValidatorCurrentPeriodSummary(valAddr)
	if !exsits {
		vcps = mapper.InitValidatorPeriodSummaryInfo(valAddr)
	}

	var currentFraction qtypes.Fraction
	if validator.BondTokens == uint64(0) {
		communityFee := mapper.GetCommunityFeePool()
		communityFee = communityFee.Add(vcps.Fees)
		mapper.Set(types.BuildCommunityFeePoolKey(), communityFee)

		currentFraction = qtypes.ZeroFraction()
	} else {
		currentFraction = qtypes.Fraction{
			Numer:   vcps.Fees,
			Denomin: btypes.NewInt(int64(validator.BondTokens)),
		}
	}

	historySummaryFrac := mapper.GetValidatorHistoryPeriodSummary(valAddr, vcps.Period-1)
	//保存当前周期汇总数据
	mapper.Set(types.BuildValidatorHistoryPeriodSummaryKey(valAddr, vcps.Period), historySummaryFrac.Add(currentFraction).GCD())
	newPeriod := vcps.Period + 1
	mapper.Set(types.BuildValidatorCurrentPeriodSummaryKey(valAddr), types.ValidatorCurrentPeriodSummary{
		Period: newPeriod,
		Fees:   btypes.ZeroInt(),
	})

	return vcps.Period
}

//修改delegator绑定的token:
//1. 增加validator的周期
//2. 计算delegator在之前的周期的收益
//3. 追加该收益到delegator start信息中,并更新start信息
func (mapper *DistributionMapper) ModifyDelegatorTokens(validator types.Validator, deleAddr btypes.Address, updatedToken, blockHeight uint64) error {
	valAddr := validator.GetValidatorAddress()
	info, exsits := mapper.GetDelegatorEarningStartInfo(valAddr, deleAddr)
	if !exsits {
		return fmt.Errorf("DelegatorEarningStartInfo not exsist. deleAddr: %s, valAddr: %s ", deleAddr, valAddr)
	}

	endPeriod := mapper.incrementValidatorPeriod(validator)
	rewards := mapper.calculateRewardsBetweenPeriod(valAddr, info.PreviousPeriod, endPeriod, info.BondToken)

	//修改delegator start信息: 该阶段的delegator收益到周期后再发放
	info.BondToken = updatedToken
	info.StartingHeight = blockHeight
	info.PreviousPeriod = endPeriod
	info.HistoricalRewardFees = info.HistoricalRewardFees.Add(rewards)

	mapper.Set(types.BuildDelegatorEarningStartInfoKey(valAddr, deleAddr), info)
	return nil
}

//计算delegator在截止周期前的收益
func (mapper *DistributionMapper) CalculateDelegatorPeriodRewards(valAddr, deleAddr btypes.Address, endPeriod, blockHeight uint64) (btypes.BigInt, error) {
	info, exsits := mapper.GetDelegatorEarningStartInfo(valAddr, deleAddr)
	if !exsits {
		return btypes.BigInt{}, fmt.Errorf("DelegatorEarningStartInfo not exsist. deleAddr: %s, valAddr: %s ", deleAddr, valAddr)
	}

	rewards := mapper.calculateRewardsBetweenPeriod(valAddr, info.PreviousPeriod, endPeriod, info.BondToken)
	historicalRewards := info.HistoricalRewardFees
	//清空delegator start中的汇总信息
	info.StartingHeight = blockHeight
	info.PreviousPeriod = endPeriod
	info.HistoricalRewardFees = btypes.ZeroInt()
	mapper.Set(types.BuildDelegatorEarningStartInfoKey(valAddr, deleAddr), info)

	return historicalRewards.Add(rewards), nil
}

//计算bondTokens在validator的两个周期内的收益
func (mapper *DistributionMapper) calculateRewardsBetweenPeriod(valAddr btypes.Address, startPeriod, endPeriod, bondTokens uint64) btypes.BigInt {

	if startPeriod > endPeriod {
		return btypes.ZeroInt()
	}

	if bondTokens == uint64(0) {
		return btypes.ZeroInt()
	}

	startFraction := mapper.GetValidatorHistoryPeriodSummary(valAddr, startPeriod)
	endFraction := mapper.GetValidatorHistoryPeriodSummary(valAddr, endPeriod)

	return (endFraction.Sub(startFraction)).MultiInt64(int64(bondTokens))
}

//-----------------------------------------------------------------

type DistributionMapper struct {
	*mapper.BaseMapper
}

var _ mapper.IMapper = (*DistributionMapper)(nil)

func NewDistributionMapper() *DistributionMapper {
	var distributionMapper = DistributionMapper{}
	distributionMapper.BaseMapper = mapper.NewBaseMapper(nil, types.DistributionMapperName)
	return &distributionMapper
}

func GetDistributionMapper(ctx context.Context) *DistributionMapper {
	return ctx.Mapper(types.DistributionMapperName).(*DistributionMapper)
}

func (mapper *DistributionMapper) Copy() mapper.IMapper {
	distributionMapper := &DistributionMapper{}
	distributionMapper.BaseMapper = mapper.BaseMapper.Copy()
	return distributionMapper
}

func (mapper *DistributionMapper) GetParams() (params types.DistributionParams) {
	exsits := mapper.Get(types.BuildDistributeParamsKey(), &params)
	if !exsits {
		params = types.DefaultDistributionParams()
	}
	return
}

func (mapper *DistributionMapper) GetValidatorCurrentPeriodSummary(valAddr btypes.Address) (vcps types.ValidatorCurrentPeriodSummary, exsits bool) {
	key := types.BuildValidatorCurrentPeriodSummaryKey(valAddr)
	exsits = mapper.Get(key, &vcps)
	return
}

func (mapper *DistributionMapper) GetCommunityFeePool() btypes.BigInt {
	var communityFeePool btypes.BigInt
	exsits := mapper.Get(types.BuildCommunityFeePoolKey(), &communityFeePool)
	if !exsits {
		return btypes.ZeroInt()
	}
	return communityFeePool
}

func (mapper *DistributionMapper) GetValidatorHistoryPeriodSummary(valAddr btypes.Address, period uint64) (frac qtypes.Fraction) {
	key := types.BuildValidatorHistoryPeriodSummaryKey(valAddr, period)
	exsits := mapper.Get(key, &frac)
	if !exsits {
		return qtypes.ZeroFraction()
	}
	return
}

func (mapper *DistributionMapper) GetDelegatorEarningStartInfo(valAddr, deleAddr btypes.Address) (info types.DelegatorEarningsStartInfo, exsits bool) {
	key := types.BuildDelegatorEarningStartInfoKey(valAddr, deleAddr)
	exsits = mapper.Get(key, &info)
	return
}

func (mapper *DistributionMapper) GetPreDistributionQOS() btypes.BigInt {
	var amount btypes.BigInt
	exsits := mapper.Get(types.BuildBlockDistributionKey(), &amount)
	if !exsits {
		return btypes.ZeroInt()
	}
	return amount
}

func (mapper *DistributionMapper) AddPreDistributionQOS(amount btypes.BigInt) {
	current := mapper.GetPreDistributionQOS()
	mapper.Set(types.BuildBlockDistributionKey(), current.Add(amount))
}

func (mapper *DistributionMapper) ClearPreDistributionQOS() {
	mapper.Set(types.BuildBlockDistributionKey(), btypes.ZeroInt())
}
