package distribution

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/distribution/mapper"
	"github.com/QOSGroup/qos/module/distribution/types"
	"github.com/QOSGroup/qos/module/stake"
	qtypes "github.com/QOSGroup/qos/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

//beginBlocker根据Vote信息进行QOS分配: mint+tx fee
func BeginBlocker(ctx context.Context, req abci.RequestBeginBlock) {
	sm := stake.GetMapper(ctx)

	totalPower, denomTotalPower := int64(0), int64(0)
	validators := make([]stake.Validator, 0, len(req.LastCommitInfo.GetVotes()))

	//获得所有符合分配条件的validator(投票 + active)
	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		totalPower += voteInfo.Validator.Power
		if !voteInfo.SignedLastBlock {
			continue
		}

		valAddr := voteInfo.Validator.Address

		v, exists := sm.GetValidator(valAddr)
		if !(exists && v.IsActive()) {
			continue
		}

		denomTotalPower += int64(v.GetBondTokens())
		validators = append(validators, v)
	}

	distributionMapper := mapper.GetMapper(ctx)

	if ctx.BlockHeight() > 1 {
		previousProposer := distributionMapper.GetLastBlockProposer()
		allocateQOS(ctx, previousProposer, denomTotalPower, validators)
	}

	consAddr := btypes.Address(req.Header.ProposerAddress)
	distributionMapper.SetLastBlockProposer(consAddr)
}

//对delegator的收益进行发放,并决定是否有下一次收益
func EndBlocker(ctx context.Context, req abci.RequestEndBlock) {
	height := uint64(req.Height)
	dm := mapper.GetMapper(ctx)

	prefixKey := types.BuildDelegatorPeriodIncomePrefixKey(height)
	iter := btypes.KVStorePrefixIterator(dm.GetStore(), prefixKey)
	validatorMap := make(map[string][]btypes.Address)
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		valAddr, delAddr, _ := types.GetDelegatorPeriodIncomeHeightAddr(key)
		mapKey := valAddr.String()
		validatorMap[mapKey] = append(validatorMap[mapKey], delAddr)
		dm.Del(key) //删除当前高度收益发放信息
	}
	iter.Close()

	params := dm.GetParams(ctx)
	for k, delegators := range validatorMap {
		valAddr, _ := btypes.GetAddrFromBech32(k)
		distributeEarningByValidator(ctx, valAddr, delegators, height, params.DelegatorsIncomePeriodHeight)

		//获取validator委托人的最小计费点周期, 删除validator历史计费点周期
		minPeriod := dm.GetValidatorMinPeriodFromDelegators(valAddr)
		dm.ClearValidatorHistoryPeroid(valAddr, minPeriod)
	}
}

//按周期分配收益:
//1. 计算delegator该周期的收益
//2. 判断delegator是否有下一周期的收益: validator存在且delegation中token大于0
//3. 若无下一周期收益,则将收益返还至delegator账户,删除delegation信息,delegator收益信息
//4. 若有下一周期收益
//      1. 若不复投,则收益直接返还至delegator账户,生成下一周期收益发放信息
//      2. 若复投, 则更新委托信息
//5.  更新validator totalpower信息
func distributeEarningByValidator(ctx context.Context, valAddr btypes.Address, delegators []btypes.Address, blockHeight, periodHeightParam uint64) {

	dm := mapper.GetMapper(ctx)
	sm := stake.GetMapper(ctx)
	am := baseabci.GetAccountMapper(ctx)

	m := make(map[string]struct{})
	addCompoundTokens := uint64(0)

	//0. 获取validator
	validator, exists := sm.GetValidator(valAddr)
	if !exists {
		//validator不存在时, 获取delegator当前收益信息, 将收益直接返还账户中,并删除当前delegator信息
		for _, delAddr := range delegators {
			if info, exists := dm.GetDelegatorEarningStartInfo(valAddr, delAddr); exists {
				delegator := am.GetAccount(delAddr).(*qtypes.QOSAccount)
				delegator.PlusQOS(info.HistoricalRewardFees.NilToZero())
				am.SetAccount(delegator)
				dm.MinusValidatorEcoFeePool(valAddr, info.HistoricalRewardFees.NilToZero())
				ctx.EventManager().EmitEvent(
					btypes.NewEvent(
						types.EventTypeDelegatorReward,
						btypes.NewAttribute(types.AttributeKeyTokens, info.HistoricalRewardFees.NilToZero().String()),
						btypes.NewAttribute(types.AttributeKeyValidator, valAddr.String()),
						btypes.NewAttribute(types.AttributeKeyDelegator, delAddr.String()),
					),
				)
				dm.DelDelegatorEarningStartInfo(valAddr, delAddr)
				sm.DelDelegationInfo(delAddr, valAddr)
			}
		}
		return
	}

	//1. validator汇总收益增加计费周期
	endPeriod := dm.IncrementValidatorPeriod(validator)

	//2. 处理delegator收益信息
	for _, deleAddr := range delegators {
		if _, ok := m[deleAddr.String()]; ok {
			continue
		}

		m[deleAddr.String()] = struct{}{}
		addTokens := distributeDelegatorEarning(ctx, validator, endPeriod, deleAddr, blockHeight, periodHeightParam)
		addCompoundTokens = addCompoundTokens + addTokens
	}

	if addCompoundTokens > 0 {
		//更新validator bondTokens
		updatedTokens := validator.GetBondTokens() + addCompoundTokens
		sm.ChangeValidatorBondTokens(validator, updatedTokens)
	}
}

func distributeDelegatorEarning(ctx context.Context, validator stake.Validator, endPeriod uint64, delAddr btypes.Address, blockHeight, periodHeightParam uint64) uint64 {
	sm := stake.GetMapper(ctx)
	dm := mapper.GetMapper(ctx)
	am := baseabci.GetAccountMapper(ctx)

	valAddr := validator.GetValidatorAddress()
	rewards, err := dm.CalculateDelegatorPeriodRewards(valAddr, delAddr, endPeriod, blockHeight)
	if err != nil {
		return 0
	}
	rewards = rewards.NilToZero()

	delegationInfo, exists := sm.GetDelegationInfo(delAddr, valAddr)

	if !exists || delegationInfo.Amount == 0 {
		//已无委托关系,收益直接分配到delegator账户中
		delegator := am.GetAccount(delAddr).(*qtypes.QOSAccount)
		delegator.PlusQOS(rewards)
		am.SetAccount(delegator)
		dm.MinusValidatorEcoFeePool(valAddr, rewards)
		ctx.EventManager().EmitEvent(
			btypes.NewEvent(
				types.EventTypeDelegatorReward,
				btypes.NewAttribute(types.AttributeKeyTokens, rewards.String()),
				btypes.NewAttribute(types.AttributeKeyValidator, valAddr.String()),
				btypes.NewAttribute(types.AttributeKeyDelegator, delAddr.String()),
			),
		)
		dm.DelDelegatorEarningStartInfo(valAddr, delAddr)
		sm.DelDelegationInfo(delAddr, valAddr)
		return 0
	}

	//增加下一周期的收益发放信息
	nextIncomeHeight := blockHeight + periodHeightParam
	dm.Set(types.BuildDelegatorPeriodIncomeKey(valAddr, delAddr, nextIncomeHeight), true)

	//非复投,收益直接分配到delegator账户中
	if !delegationInfo.IsCompound {
		delegator := am.GetAccount(delAddr).(*qtypes.QOSAccount)
		delegator.PlusQOS(rewards)
		am.SetAccount(delegator)
		dm.MinusValidatorEcoFeePool(valAddr, rewards)
		ctx.EventManager().EmitEvent(
			btypes.NewEvent(
				types.EventTypeDelegatorReward,
				btypes.NewAttribute(types.AttributeKeyTokens, rewards.String()),
				btypes.NewAttribute(types.AttributeKeyValidator, valAddr.String()),
				btypes.NewAttribute(types.AttributeKeyDelegator, delAddr.String()),
			),
		)
		return 0
	}

	//复投
	addTokens := uint64(rewards.Int64())

	//更新delegation委托信息,更新delegate当前收益信息
	info, _ := dm.GetDelegatorEarningStartInfo(valAddr, delAddr)
	info.BondToken = info.BondToken + addTokens
	dm.Set(types.BuildDelegatorEarningStartInfoKey(valAddr, delAddr), info)

	delegationInfo.Amount = delegationInfo.Amount + addTokens
	sm.SetDelegationInfo(delegationInfo)

	//复投时validator收益池处理
	dm.MinusValidatorEcoFeePool(valAddr, rewards)
	ctx.EventManager().EmitEvent(
		btypes.NewEvent(
			types.EventTypeDelegatorReward,
			btypes.NewAttribute(types.AttributeKeyTokens, rewards.String()),
			btypes.NewAttribute(types.AttributeKeyValidator, valAddr.String()),
			btypes.NewAttribute(types.AttributeKeyDelegator, delAddr.String()),
		),
	)

	ctx.EventManager().EmitEvent(
		btypes.NewEvent(
			types.EventTypeDelegate,
			btypes.NewAttribute(types.AttributeKeyTokens, rewards.String()),
			btypes.NewAttribute(types.AttributeKeyValidator, valAddr.String()),
			btypes.NewAttribute(types.AttributeKeyDelegator, delAddr.String()),
		),
	)

	return addTokens
}

// 2.  每块挖出的QOS数量:  `x%`proposer + `y%`validators + `z%`community
//        * `x%`proposer: 验证人获得的奖励,直接归属proposer
//        * `y%`validators: 根据每个validator的power占比平均分配
// 3.  validator奖励数 =  validator佣金 +  平分金额Fee(漏签和inactive的validator不分配奖励)
//        * validator佣金奖励: 佣金 = validator奖励数 * `commission rate`
//        * 平分金额Fee由validator,delegator根据各自绑定的stake平均分配
// 4.  validator的proposer奖励,佣金奖励 均按周期发放
//
func allocateQOS(ctx context.Context, proposerAddr btypes.Address, denomTotalPower int64, validators []stake.Validator) {
	dm := mapper.GetMapper(ctx)
	sm := stake.GetMapper(ctx)

	params := dm.GetParams(ctx)

	//获取待分配的QOS总量
	totalAmount := dm.GetPreDistributionQOS()
	remainQOS := totalAmount
	dm.ClearPreDistributionQOS()

	//proposer奖励,直接归属proposer
	proposerRewards := params.ProposerRewardRate.MultiBigInt(totalAmount)
	proposerValidater, validatorExsits := sm.GetValidator(proposerAddr)
	proposerInfo, proposerInfoExsits := dm.GetDelegatorEarningStartInfo(proposerAddr, proposerValidater.Owner)

	if validatorExsits && proposerInfoExsits {
		proposerInfo.HistoricalRewardFees = proposerInfo.HistoricalRewardFees.Add(proposerRewards)
		remainQOS = remainQOS.Sub(proposerRewards)
		dm.Set(types.BuildDelegatorEarningStartInfoKey(proposerAddr, proposerValidater.Owner), proposerInfo)
		dm.AddValidatorEcoFeePool(proposerAddr, proposerRewards, btypes.ZeroInt(), btypes.ZeroInt())
		ctx.EventManager().EmitEvent(
			btypes.NewEvent(
				types.EventTypeProposerReward,
				btypes.NewAttribute(types.AttributeKeyTokens, proposerRewards.String()),
				btypes.NewAttribute(types.AttributeKeyValidator, proposerAddr.String()),
			),
		)
	}

	//vote奖励(漏签和inactive的validator不分配奖励)
	votePercent := qtypes.OneFraction().Sub(params.ProposerRewardRate).Sub(params.CommunityRewardRate)
	for _, validator := range validators {
		votePowerFrac := qtypes.NewFraction(int64(validator.BondTokens), denomTotalPower)
		rewards := votePowerFrac.Mul(votePercent).MultiBigInt(totalAmount)
		remainQOS = remainQOS.Sub(rewards)
		rewardToValidator(ctx, validator, rewards)
	}

	//社区奖励
	communityFeePool := dm.GetCommunityFeePool()
	communityFeePool = communityFeePool.Add(remainQOS)
	dm.SetCommunityFeePool(communityFeePool)
	ctx.EventManager().EmitEvent(
		btypes.NewEvent(
			types.EventTypeCommunity,
			btypes.NewAttribute(types.AttributeKeyTokens, communityFeePool.String()),
		),
	)
}

func rewardToValidator(ctx context.Context, validator stake.Validator, rewards btypes.BigInt) {
	dm := mapper.GetMapper(ctx)

	commissionReward := validator.Commission.Rate.MulInt(rewards).TruncateInt()
	sharedReward := rewards.Sub(commissionReward)

	valAddr := validator.GetValidatorAddress()
	dm.AddValidatorEcoFeePool(valAddr, btypes.ZeroInt(), commissionReward, sharedReward)

	//validator 佣金收益
	if info, exists := dm.GetDelegatorEarningStartInfo(valAddr, validator.Owner); exists {
		info.HistoricalRewardFees = info.HistoricalRewardFees.Add(commissionReward)
		dm.Set(types.BuildDelegatorEarningStartInfoKey(valAddr, validator.Owner), info)
		ctx.EventManager().EmitEvent(
			btypes.NewEvent(
				types.EventTypeCommission,
				btypes.NewAttribute(types.AttributeKeyTokens, commissionReward.String()),
				btypes.NewAttribute(types.AttributeKeyValidator, validator.GetValidatorAddress().String()),
			),
		)
	}

	//delegator 共同收益
	if vcps, exists := dm.GetValidatorCurrentPeriodSummary(valAddr); exists {
		vcps.Fees = vcps.Fees.Add(sharedReward)
		dm.Set(types.BuildValidatorCurrentPeriodSummaryKey(valAddr), vcps)
		ctx.EventManager().EmitEvent(
			btypes.NewEvent(
				types.EventTypeDelegatorRewards,
				btypes.NewAttribute(types.AttributeKeyTokens, sharedReward.String()),
				btypes.NewAttribute(types.AttributeKeyValidator, validator.GetValidatorAddress().String()),
			),
		)
	}

}
