package eco

import (
	"errors"
	"fmt"

	"github.com/QOSGroup/qbase/store"
	"github.com/QOSGroup/qos/module/eco/mapper"
	"github.com/QOSGroup/qos/module/eco/types"

	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
)

//删除和validator相关的eco数据
//CONTRACT:
//delegator当前收益和收益发放信息数据不删除, 只是将bondTokens重置为0
//发放收益时,若delegator非validator的委托人, 或validator 不存在 则可以将delegator的收益相关数据删除
//发放收益时,validator的汇总数据可能会不存在
func (e Eco) RemoveValidator(valAddr btypes.Address) error {

	height := uint64(e.Context.BlockHeight())

	distributionMapper := e.DistributionMapper
	delegationMapper := e.DelegationMapper
	validatorMapper := e.ValidatorMapper
	voteInfoMapper := e.VoteInfoMapper

	stakeParams := validatorMapper.GetParams()

	// 删除validator相关数据
	validator, ok := validatorMapper.KickValidator(valAddr)
	if !ok {
		return fmt.Errorf("validator:%s not exsits", valAddr)
	}

	//1. validator的汇总收益增加
	endPeriod := distributionMapper.IncrementValidatorPeriod(validator)

	//2. 计算所有delegator的收益信息,并将delegator绑定的token置为0
	prefixKey := append(types.GetDelegatorEarningsStartInfoPrefixKey(), valAddr...)
	iter := store.KVStorePrefixIterator(distributionMapper.GetStore(), prefixKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var info types.DelegatorEarningsStartInfo
		distributionMapper.BaseMapper.DecodeObject(iter.Value(), &info)

		unbondToken := info.BondToken
		_, deleAddr := types.GetDelegatorEarningStartInfoAddr(iter.Key())
		rewards := distributionMapper.CalculateRewardsBetweenPeriod(valAddr, info.PreviousPeriod, endPeriod, unbondToken)

		info.BondToken = uint64(0)
		info.CurrentStartingHeight = height
		info.PreviousPeriod = endPeriod
		info.HistoricalRewardFees = info.HistoricalRewardFees.Add(rewards)

		distributionMapper.Set(types.BuildDelegatorEarningStartInfoKey(valAddr, deleAddr), info)

		// 删除delegate数据,增加unbond数据
		delegationMapper.DelDelegationInfo(deleAddr, valAddr)

		//unbond height
		unbondHeight := uint64(stakeParams.DelegatorUnbondReturnHeight) + height
		delegationMapper.AddDelegatorUnbondingQOSatHeight(unbondHeight, deleAddr, unbondToken)
	}

	//删除validator汇总收益数据
	distributionMapper.DeleteValidatorPeriodSummaryInfo(valAddr)

	//删除validator 投票数据
	voteInfoMapper.DelValidatorVoteInfo(valAddr)
	voteInfoMapper.ClearValidatorVoteInfoInWindow(valAddr)

	return nil
}

func (e Eco) DelegateValidator(validator types.Validator, delegatorAddr btypes.Address, delegateAmount uint64, isCompound bool, needMinusAccountQOS bool) error {

	height := uint64(e.Context.BlockHeight())

	distributionMapper := e.DistributionMapper
	delegationMapper := e.DelegationMapper
	validatorMapper := e.ValidatorMapper

	valAddr := validator.GetValidatorAddress()

	if needMinusAccountQOS {
		//0. delegator账户扣减QOS, amount:qos = 1:1
		decrQOS := btypes.NewInt(int64(delegateAmount))
		if err := DecrAccountQOS(e.Context, delegatorAddr, decrQOS); err != nil {
			return err
		}
	}

	//获取delegation信息,若delegation信息不存在,则初始化degelator收益信息
	delegatedAmount := uint64(0)
	info, exsits := delegationMapper.GetDelegationInfo(delegatorAddr, valAddr)
	if !exsits {
		distributionMapper.InitDelegatorIncomeInfo(valAddr, delegatorAddr, uint64(0), height)
		info = types.NewDelegationInfo(delegatorAddr, valAddr, uint64(0), false)
	} else {
		delegatedAmount = info.Amount
	}

	updatedAmount := delegatedAmount + delegateAmount
	//1. validator增加周期 , 计算周期段内delegator收益,并更新收益信息
	if err := distributionMapper.ModifyDelegatorTokens(validator, delegatorAddr, updatedAmount, height); err != nil {
		return err
	}

	//2. 更新delegation信息
	info.Amount = updatedAmount
	info.IsCompound = isCompound
	delegationMapper.SetDelegationInfo(info)

	//3. 更新validator的bondTokens, amount:token = 1:1
	validatorAddToken := delegateAmount
	updatedValidatorTokens := validator.BondTokens + validatorAddToken
	validatorMapper.ChangeValidatorBondTokens(validator, updatedValidatorTokens)

	return nil
}

func (e Eco) UnbondValidator(validator types.Validator, delegatorAddr btypes.Address, isUnbondAll bool, unbondAmount uint64, isRedelegate bool) error {

	height := uint64(e.Context.BlockHeight())

	distributionMapper := e.DistributionMapper
	delegationMapper := e.DelegationMapper
	validatorMapper := e.ValidatorMapper

	valAddr := validator.GetValidatorAddress()
	info, _ := delegationMapper.GetDelegationInfo(delegatorAddr, valAddr)

	if isUnbondAll {
		unbondAmount = info.Amount
	}

	if info.Amount < unbondAmount {
		return errors.New("unbond amount overflow")
	}

	//1. 计算当前delegator收益
	updatedTokens := info.Amount - unbondAmount
	if err := distributionMapper.ModifyDelegatorTokens(validator, delegatorAddr, updatedTokens, height); err != nil {
		return err
	}

	//2. 更新delegation信息
	info.Amount = info.Amount - unbondAmount
	delegationMapper.SetDelegationInfo(info)

	if !isRedelegate {
		//3. 增加unbond信息
		stakeParams := validatorMapper.GetParams()
		unbondHeight := uint64(stakeParams.DelegatorUnbondReturnHeight) + height
		delegationMapper.AddDelegatorUnbondingQOSatHeight(unbondHeight, delegatorAddr, unbondAmount)
	}

	//4. 更新validator的bondTokens, amount:token = 1:1
	validatorMinusToken := unbondAmount
	updatedValidatorTokens := validator.BondTokens - validatorMinusToken
	validatorMapper.ChangeValidatorBondTokens(validator, updatedValidatorTokens)

	return nil
}

type Eco struct {
	Context            context.Context
	DistributionMapper *mapper.DistributionMapper
	DelegationMapper   *mapper.DelegationMapper
	ValidatorMapper    *mapper.ValidatorMapper
	VoteInfoMapper     *mapper.VoteInfoMapper
}

func GetEco(ctx context.Context) Eco {

	distributionMapper := mapper.GetDistributionMapper(ctx)
	delegationMapper := mapper.GetDelegationMapper(ctx)
	validatorMapper := mapper.GetValidatorMapper(ctx)
	voteInfoMapper := mapper.GetVoteInfoMapper(ctx)

	return Eco{
		Context:            ctx,
		DistributionMapper: distributionMapper,
		DelegationMapper:   delegationMapper,
		ValidatorMapper:    validatorMapper,
		VoteInfoMapper:     voteInfoMapper,
	}
}
