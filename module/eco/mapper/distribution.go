package mapper

import (
	"fmt"
	"reflect"

	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/mapper"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/eco/types"
	pmapper "github.com/QOSGroup/qos/module/params"
	qtypes "github.com/QOSGroup/qos/types"
)

/*
	Distribution Mapper 收益分配相关Mapper:

keys:
	Validator:
	 1. 历史计费点汇总收益信息: validatorHistoryPeriodSummaryPrefixKey
	 2. 当前计费点收益信息: validatorCurrentPeriodSummaryPrefixKey
	Delegator:
	 1. Delegator收益计算信息: delegatorEarningsStartInfoPrefixKey
	 2. Delegator某高度下是否发放收益信息: delegatorPeriodIncomePrefixKey: 仅在对delegator发放收益时删除
*/

func BuildDistributionStoreQueryPath() []byte {
	return []byte(fmt.Sprintf("/store/%s/key", types.DistributionMapperName))
}

//初始化validator历史计费点汇总收益,当前计费点收益信息.
func (mapper *DistributionMapper) InitValidatorPeriodSummaryInfo(valAddr btypes.Address) types.ValidatorCurrentPeriodSummary {
	mapper.Set(types.BuildValidatorHistoryPeriodSummaryKey(valAddr, uint64(0)), qtypes.ZeroFraction())
	current := types.ValidatorCurrentPeriodSummary{
		Period: 1,
		Fees:   btypes.ZeroInt(),
	}
	mapper.Set(types.BuildValidatorCurrentPeriodSummaryKey(valAddr), current)
	return current
}

//清空validator收益分配相关信息
func (mapper *DistributionMapper) DeleteValidatorPeriodSummaryInfo(valAddr btypes.Address) {
	periodPrifixKey := append(types.GetValidatorHistoryPeriodSummaryPrefixKey(), valAddr...)
	iter := btypes.KVStorePrefixIterator(mapper.GetStore(), periodPrifixKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		mapper.Del(iter.Key())
	}

	k := types.BuildValidatorCurrentPeriodSummaryKey(valAddr)
	mapper.Del(k)
}

//首次Delegate:
//1. first delegate
//2. unbond all , then delegate
func (mapper *DistributionMapper) InitDelegatorIncomeInfo(ctx context.Context, valAddr, deleAddr btypes.Address, bondTokens, currHeight uint64) {
	//初始化delegaotr 收益计算信息
	vcps, _ := mapper.GetValidatorCurrentPeriodSummary(valAddr)
	params := mapper.GetParams(ctx)

	key := types.BuildDelegatorEarningStartInfoKey(valAddr, deleAddr)

	startInfo := types.DelegatorEarningsStartInfo{
		PreviousPeriod:        vcps.Period - 1,
		BondToken:             bondTokens,
		CurrentStartingHeight: currHeight,
		FirstDelegateHeight:   currHeight,
		HistoricalRewardFees:  btypes.ZeroInt(),
	}

	//delegator unbond全部后,又重新delegate
	var info types.DelegatorEarningsStartInfo
	exsits := mapper.Get(key, &info)
	if exsits { //保留delegator历史收益,不计算阶段内收益
		startInfo.HistoricalRewardFees = startInfo.HistoricalRewardFees.Add(info.HistoricalRewardFees)
		startInfo.FirstDelegateHeight = info.FirstDelegateHeight
	}

	mapper.Set(key, startInfo)

	//发放收益高度
	if !exsits {
		incomeHeight := currHeight + params.DelegatorsIncomePeriodHeight
		mapper.Set(types.BuildDelegatorPeriodIncomeKey(valAddr, deleAddr, incomeHeight), true)
	}

}

//删除delegator收益计算信息
//todo: 某高度下发放收益信息没有删除
func (mapper *DistributionMapper) DeleteDelegatorIncomeInfo(valAddr, deleAddr btypes.Address) {
	k := types.BuildDelegatorEarningStartInfoKey(valAddr, deleAddr)
	mapper.Del(k)
}

//增加validator收益计费点
//1. 保存历史计费点收益汇总信息
//2. 更新当前计费点收益信息,返回上一计费点数值
func (mapper *DistributionMapper) IncrementValidatorPeriod(validator types.Validator) uint64 {
	valAddr := validator.GetValidatorAddress()
	vcps, exsits := mapper.GetValidatorCurrentPeriodSummary(valAddr)
	if !exsits {
		vcps = mapper.InitValidatorPeriodSummaryInfo(valAddr)
	}

	var currentFraction qtypes.Fraction
	if validator.BondTokens == uint64(0) {
		communityFee := mapper.GetCommunityFeePool()
		communityFee = communityFee.Add(vcps.Fees)
		mapper.SetCommunityFeePool(communityFee)

		currentFraction = qtypes.ZeroFraction()
	} else {
		currentFraction = qtypes.NewFractionFromBigInt(vcps.Fees, btypes.NewInt(int64(validator.BondTokens)))
	}

	historySummaryFrac := mapper.GetValidatorHistoryPeriodSummary(valAddr, vcps.Period-1)
	//保存当前计费点历史汇总数据
	mapper.Set(types.BuildValidatorHistoryPeriodSummaryKey(valAddr, vcps.Period), historySummaryFrac.Add(currentFraction))

	//增加当前计费点,更新数据
	newPeriod := vcps.Period + 1
	mapper.Set(types.BuildValidatorCurrentPeriodSummaryKey(valAddr), types.ValidatorCurrentPeriodSummary{
		Period: newPeriod,
		Fees:   btypes.ZeroInt(),
	})

	return vcps.Period
}

//修改delegator绑定的token:
//1. 增加validator的计费点
//2. 计算delegator在两次计费点间的收益
//3. 追加该收益到delegator 收益计算信息中
func (mapper *DistributionMapper) ModifyDelegatorTokens(validator types.Validator, deleAddr btypes.Address, updatedToken, blockHeight uint64) error {
	valAddr := validator.GetValidatorAddress()
	info, exsits := mapper.GetDelegatorEarningStartInfo(valAddr, deleAddr)
	if !exsits {
		return fmt.Errorf("DelegatorEarningStartInfo not exsist. deleAddr: %s, valAddr: %s ", deleAddr, valAddr)
	}

	endPeriod := mapper.IncrementValidatorPeriod(validator)
	rewards := mapper.CalculateRewardsBetweenPeriod(valAddr, info.PreviousPeriod, endPeriod, info.BondToken)

	//修改delegator 收益计算信息: 该区间收益在到达发放收益高度时发放
	//firstDelegateHeight不变
	info.BondToken = updatedToken
	info.CurrentStartingHeight = blockHeight
	info.PreviousPeriod = endPeriod
	info.HistoricalRewardFees = info.HistoricalRewardFees.Add(rewards)

	mapper.Set(types.BuildDelegatorEarningStartInfoKey(valAddr, deleAddr), info)
	return nil
}

func (mapper *DistributionMapper) GetValidatorMinPeriodFromDelegators(valAddr btypes.Address) uint64 {
	prefixKey := append(types.GetDelegatorEarningsStartInfoPrefixKey(), valAddr...)

	minPeriod := uint64(0)
	i := int64(0)

	iter := btypes.KVStorePrefixIterator(mapper.GetStore(), prefixKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var info types.DelegatorEarningsStartInfo
		mapper.BaseMapper.DecodeObject(iter.Value(), &info)

		if i == int64(0) {
			minPeriod = info.PreviousPeriod
			i = i + 1
		}

		if info.PreviousPeriod < minPeriod {
			minPeriod = info.PreviousPeriod
		}
	}

	return minPeriod
}

//删除validator历史计费点信息,额外保留2个历史数据
func (mapper *DistributionMapper) ClearValidatorHistoryPeroid(valAddr btypes.Address, minPeroid uint64) {
	if minPeroid <= uint64(2) {
		return
	}

	mapper.IteratorWithKV(types.BuildValidatorHistoryPeriodSummaryPrefixKey(valAddr), func(key []byte, value []byte) (stop bool) {
		_, p := types.GetValidatorHistoryPeriodSummaryAddrPeriod(key)

		if p >= (minPeroid - 2) {
			return true
		}

		mapper.Del(key)
		return false
	})
}

//计算delegator在计费点区间的收益
func (mapper *DistributionMapper) CalculateDelegatorPeriodRewards(valAddr, deleAddr btypes.Address, endPeriod, blockHeight uint64) (btypes.BigInt, error) {
	info, exsits := mapper.GetDelegatorEarningStartInfo(valAddr, deleAddr)
	if !exsits {
		return btypes.BigInt{}, fmt.Errorf("DelegatorEarningStartInfo not exsist. deleAddr: %s, valAddr: %s ", deleAddr, valAddr)
	}

	rewards := mapper.CalculateRewardsBetweenPeriod(valAddr, info.PreviousPeriod, endPeriod, info.BondToken)
	historicalRewards := info.HistoricalRewardFees
	//清空delegator start中的汇总信息
	info.CurrentStartingHeight = blockHeight
	info.PreviousPeriod = endPeriod
	info.HistoricalRewardFees = btypes.ZeroInt()

	totalRewards := historicalRewards.Add(rewards)
	info.LastIncomeCalHeight = blockHeight
	info.LastIncomeCalFees = totalRewards

	mapper.Set(types.BuildDelegatorEarningStartInfoKey(valAddr, deleAddr), info)

	return totalRewards, nil
}

//计算bondTokens在validator的两个计费点区间的收益
func (mapper *DistributionMapper) CalculateRewardsBetweenPeriod(valAddr btypes.Address, startPeriod, endPeriod, bondTokens uint64) btypes.BigInt {

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

func (mapper *DistributionMapper) GetParams(ctx context.Context) (params types.DistributionParams) {
	pmapper.GetMapper(ctx).GetParamSet(&params)
	return
}

func (mapper *DistributionMapper) SetParams(ctx context.Context, params types.DistributionParams) {
	pmapper.GetMapper(ctx).SetParamSet(&params)
}

func (mapper *DistributionMapper) GetValidatorCurrentPeriodSummary(valAddr btypes.Address) (vcps types.ValidatorCurrentPeriodSummary, exsits bool) {
	key := types.BuildValidatorCurrentPeriodSummaryKey(valAddr)
	exsits = mapper.Get(key, &vcps)
	return
}

func (mapper *DistributionMapper) GetLastBlockProposer() btypes.Address {
	var previousProposer btypes.Address
	mapper.Get(types.BuildLastProposerKey(), &previousProposer)
	return previousProposer
}

func (mapper *DistributionMapper) SetLastBlockProposer(proposer btypes.Address) {
	mapper.Set(types.BuildLastProposerKey(), proposer)
}

func (mapper *DistributionMapper) GetCommunityFeePool() btypes.BigInt {
	var communityFeePool btypes.BigInt
	exsits := mapper.Get(types.BuildCommunityFeePoolKey(), &communityFeePool)
	if !exsits {
		return btypes.ZeroInt()
	}
	return communityFeePool
}

func (mapper *DistributionMapper) SetCommunityFeePool(communityFee btypes.BigInt) {
	mapper.Set(types.BuildCommunityFeePoolKey(), communityFee)
}

func (mapper *DistributionMapper) AddToCommunityFeePool(fee btypes.BigInt) {
	communityFee := mapper.GetCommunityFeePool()
	mapper.SetCommunityFeePool(communityFee.Add(fee))
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

func (mapper *DistributionMapper) DelDelegatorEarningStartInfo(valAddr, deleAddr btypes.Address) {
	key := types.BuildDelegatorEarningStartInfoKey(valAddr, deleAddr)
	mapper.Del(key)
}

func (mapper *DistributionMapper) GetPreDistributionQOS() btypes.BigInt {
	var amount btypes.BigInt
	exsits := mapper.Get(types.BuildBlockDistributionKey(), &amount)
	if !exsits {
		return btypes.ZeroInt()
	}
	return amount
}

func (mapper *DistributionMapper) SetPreDistributionQOS(amount btypes.BigInt) {
	mapper.Set(types.BuildBlockDistributionKey(), amount)
}

func (mapper *DistributionMapper) AddPreDistributionQOS(amount btypes.BigInt) {
	current := mapper.GetPreDistributionQOS()
	mapper.Set(types.BuildBlockDistributionKey(), current.Add(amount))
}

func (mapper *DistributionMapper) ClearPreDistributionQOS() {
	mapper.Set(types.BuildBlockDistributionKey(), btypes.ZeroInt())
}

//------------------------ genesis export

func (mapper *DistributionMapper) IteratorValidatorsHistoryPeriod(fn func(valAddr btypes.Address, period uint64, frac qtypes.Fraction)) {
	iter := btypes.KVStorePrefixIterator(mapper.GetStore(), types.GetValidatorHistoryPeriodSummaryPrefixKey())
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		value := iter.Value()
		valAddr, period := types.GetValidatorHistoryPeriodSummaryAddrPeriod(key)

		var frac qtypes.Fraction
		mapper.BaseMapper.DecodeObject(value, &frac)
		fn(valAddr, period, frac)
	}
}

func (mapper *DistributionMapper) IteratorValidatorsCurrentPeriod(fn func(btypes.Address, types.ValidatorCurrentPeriodSummary)) {
	iter := btypes.KVStorePrefixIterator(mapper.GetStore(), types.GetValidatorCurrentPeriodSummaryPrefixKey())
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		value := iter.Value()

		valAddr := types.GetValidatorCurrentPeriodSummaryAddr(key)

		var vcps types.ValidatorCurrentPeriodSummary
		mapper.BaseMapper.DecodeObject(value, &vcps)

		fn(valAddr, vcps)
	}
}

func (mapper *DistributionMapper) IteratorDelegatorsEarningStartInfo(fn func(btypes.Address, btypes.Address, types.DelegatorEarningsStartInfo)) {
	iter := btypes.KVStorePrefixIterator(mapper.GetStore(), types.GetDelegatorEarningsStartInfoPrefixKey())
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		value := iter.Value()

		valAddr, deleAddr := types.GetDelegatorEarningStartInfoAddr(key)

		var desi types.DelegatorEarningsStartInfo
		mapper.BaseMapper.DecodeObject(value, &desi)

		fn(valAddr, deleAddr, desi)
	}
}

func (mapper *DistributionMapper) DeleteDelegatorsEarningStartInfo() {
	iter := btypes.KVStorePrefixIterator(mapper.GetStore(), types.GetDelegatorEarningsStartInfoPrefixKey())
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		mapper.Del(iter.Key())
	}
}

func (mapper *DistributionMapper) IteratorDelegatorsIncomeHeight(fn func(btypes.Address, btypes.Address, uint64)) {
	iter := btypes.KVStorePrefixIterator(mapper.GetStore(), types.GetDelegatorPeriodIncomePrefixKey())
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		key := iter.Key()

		valAddr, deleAddr, height := types.GetDelegatorPeriodIncomeHeightAddr(key)
		fn(valAddr, deleAddr, height)
	}
}

func (mapper *DistributionMapper) DeleteDelegatorsIncomeHeight() {
	iter := btypes.KVStorePrefixIterator(mapper.GetStore(), types.GetDelegatorPeriodIncomePrefixKey())
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		mapper.Del(iter.Key())
	}
}

func (mapper *DistributionMapper) IteratorValidatorEcoFeePools(fn func(validatorAddr btypes.Address, pool types.ValidatorEcoFeePool)) {
	mapper.IteratorWithType(types.GetValidatorEcoFeePoolPrefixKey(), reflect.TypeOf(types.ValidatorEcoFeePool{}), func(key []byte, dataPtr interface{}) bool {
		sPtr := dataPtr.(*types.ValidatorEcoFeePool)
		ecoPool := *sPtr
		fn(types.GetValidatorEcoPoolAddress(key), ecoPool)
		return false
	})
}

func (mapper *DistributionMapper) AddValidatorEcoFeePool(validatorAddr btypes.Address, proposerReward, commissionReward, preDistributionReward btypes.BigInt) {
	pool := mapper.GetValidatorEcoFeePool(validatorAddr)

	pool.ProposerTotalRewardFee = pool.ProposerTotalRewardFee.Add(proposerReward)
	pool.CommissionTotalRewardFee = pool.CommissionTotalRewardFee.Add(commissionReward)
	pool.PreDistributeTotalRewardFee = pool.PreDistributeTotalRewardFee.Add(preDistributionReward)

	totalReward := proposerReward.Add(commissionReward).Add(preDistributionReward)
	pool.PreDistributeRemainTotalFee = pool.PreDistributeRemainTotalFee.Add(totalReward)

	mapper.SaveValidatorEcoFeePool(validatorAddr, pool)
}

func (mapper *DistributionMapper) MinusValidatorEcoFeePool(validatorAddr btypes.Address, bonusReward btypes.BigInt) {
	pool := mapper.GetValidatorEcoFeePool(validatorAddr)
	pool.PreDistributeRemainTotalFee = pool.PreDistributeRemainTotalFee.Sub(bonusReward)
	mapper.SaveValidatorEcoFeePool(validatorAddr, pool)
}

func (mapper *DistributionMapper) SaveValidatorEcoFeePool(validatorAddr btypes.Address, pool types.ValidatorEcoFeePool) {
	mapper.Set(types.BuildValidatorEcoFeePoolKey(validatorAddr), pool)
}

func (mapper *DistributionMapper) DeleteValidatorEcoFeePool(validatorAddr btypes.Address) {
	mapper.Del(types.BuildValidatorEcoFeePoolKey(validatorAddr))
}

func (mapper *DistributionMapper) GetValidatorEcoFeePool(validatorAddr btypes.Address) (pool types.ValidatorEcoFeePool) {
	key := types.BuildValidatorEcoFeePoolKey(validatorAddr)
	if exsits := mapper.Get(key, &pool); !exsits {
		pool = types.NewValidatorEcoFeePool()
	}
	return
}
