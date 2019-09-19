package mapper

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/distribution/types"
	"github.com/QOSGroup/qos/module/stake"
)

var _ stake.Hooks = (*StakingHooks)(nil)

type StakingHooks struct{}

func NewStakingHooks() *StakingHooks {
	return &StakingHooks{}
}

func (hooks *StakingHooks) HookMapper() string {
	return stake.ModuleName
}

// 创建validator时初始化分配信息
func (hooks *StakingHooks) AfterValidatorCreated(ctx context.Context, val btypes.ValAddress) {
	GetMapper(ctx).InitValidatorPeriodSummaryInfo(val)
}

// 删除validator时分配处理逻辑
func (hooks *StakingHooks) BeforeValidatorRemoved(ctx context.Context, val btypes.ValAddress) {
	dm := GetMapper(ctx)
	sm := stake.GetMapper(ctx)

	validator, exists := sm.GetValidator(val)
	if !exists {
		panic(fmt.Sprintf("validator %s not exists", val))
	}

	//1. validator的汇总收益增加
	endPeriod := dm.IncrementValidatorPeriod(validator)

	//2. 计算所有delegator的收益信息,并返回delegator绑定的token
	prefixKey := append(types.GetDelegatorEarningsStartInfoPrefixKey(), val...)
	iter := btypes.KVStorePrefixIterator(dm.GetStore(), prefixKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var info types.DelegatorEarningsStartInfo
		dm.BaseMapper.DecodeObject(iter.Value(), &info)

		unbondToken := info.BondToken
		_, delAddr := types.GetDelegatorEarningStartInfoAddr(iter.Key())
		rewards := dm.CalculateRewardsBetweenPeriod(val, info.PreviousPeriod, endPeriod, unbondToken)

		info.BondToken = btypes.ZeroInt()
		info.CurrentStartingHeight = ctx.BlockHeight()
		info.PreviousPeriod = endPeriod
		info.HistoricalRewardFees = info.HistoricalRewardFees.Add(rewards)

		dm.Set(types.BuildDelegatorEarningStartInfoKey(val, delAddr), info)

		// 删除delegate数据,增加unbond数据
		sm.DelDelegationInfo(delAddr, val)

		//unbond height
		unbondHeight := sm.GetParams(ctx).DelegatorUnbondFrozenHeight + ctx.BlockHeight()
		sm.AddUnbondingDelegation(stake.NewUnbondingInfo(delAddr, val, ctx.BlockHeight(), unbondHeight, unbondToken))
	}

	//删除validator汇总收益数据
	dm.DeleteValidatorPeriodSummaryInfo(val)
}

// 创建delegation时初始化分配信息
func (hooks *StakingHooks) AfterDelegationCreated(ctx context.Context, val btypes.ValAddress, del btypes.AccAddress) {
	delegation, exists := stake.GetMapper(ctx).GetDelegationInfo(del, val)
	if !exists {
		panic(fmt.Sprintf("delegation from %s to %s not exists", del, val))
	}
	GetMapper(ctx).InitDelegatorIncomeInfo(ctx, val, del, delegation.Amount, ctx.BlockHeight())
}

// 更新绑定tokens时分配处理逻辑
func (hooks *StakingHooks) BeforeDelegationModified(ctx context.Context, val btypes.ValAddress, del btypes.AccAddress, updateAmount btypes.BigInt) {
	dm := GetMapper(ctx)
	sm := stake.GetMapper(ctx)
	validator, exists := sm.GetValidator(val)
	if !exists {
		panic(fmt.Sprintf("validator %s not exists", val))
	}
	delegation, exists := stake.GetMapper(ctx).GetDelegationInfo(del, val)
	if !exists {
		panic(fmt.Sprintf("delegation from %s to %s not exists", del, val))
	}
	err := dm.ModifyDelegatorTokens(validator, delegation.DelegatorAddr, updateAmount, ctx.BlockHeight())
	if err != nil {
		panic(fmt.Sprintf("modify delegation from %s to %s error: %v", del, val, err))
	}
}

// validator惩罚后操作
func (hooks *StakingHooks) AfterValidatorSlashed(ctx context.Context, slashedTokens btypes.BigInt) {
	GetMapper(ctx).AddToCommunityFeePool(slashedTokens)
}
