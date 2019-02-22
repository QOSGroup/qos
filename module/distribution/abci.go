package distribution

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/eco"
	"github.com/QOSGroup/qos/module/eco/mapper"
	"github.com/QOSGroup/qos/module/eco/types"
	qtypes "github.com/QOSGroup/qos/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

//TODO: 分配时的精度问题

//beginblocker根据Vote信息进行QOS分配: mint+tx fee
func BeginBlocker(ctx context.Context, req abci.RequestBeginBlock) {

	totalPower, signedTotalPower := int64(0), int64(0)
	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		totalPower += voteInfo.Validator.Power
		if voteInfo.SignedLastBlock {
			signedTotalPower += voteInfo.Validator.Power
		}
	}

	distributionMapper := mapper.GetDistributionMapper(ctx)

	if ctx.BlockHeight() > 1 {
		previousProposer := distributionMapper.GetLastBlockProposer()
		allocateQOS(ctx, signedTotalPower, totalPower, previousProposer, req.LastCommitInfo.GetVotes())
	}

	consAddr := btypes.Address(req.Header.ProposerAddress)
	distributionMapper.SetLastBlockProposer(consAddr)
}

//endblocker对delegator的收益进行发放,并决定是否有下一次收益
func EndBlocker(ctx context.Context, req abci.RequestEndBlock) {

	height := uint64(req.Height)
	e := eco.GetEco(ctx)

	prefixKey := types.BuildDelegatorPeriodIncomePrefixKey(height)

	iter := store.KVStorePrefixIterator(e.DistributionMapper.GetStore(), prefixKey)
	validatorMap := make(map[string][]btypes.Address)
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		valAddr, deleAddr, _ := types.GetDelegatorPeriodIncomeHeightAddr(key)
		mapKey := valAddr.String()
		validatorMap[mapKey] = append(validatorMap[mapKey], deleAddr)
		e.DistributionMapper.Del(key) //删除当前高度收益发放信息
	}
	iter.Close()

	params := e.DistributionMapper.GetParams()
	for k, delegators := range validatorMap {
		valAddr, _ := btypes.GetAddrFromBech32(k)
		distributeEarningByValidator(e, valAddr, delegators, height, params.DelegatorsIncomePeriodHeight)
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
func distributeEarningByValidator(e eco.Eco, valAddr btypes.Address, delegators []btypes.Address, blockHeight, periodHeightParam uint64) {

	log := e.Context.Logger()

	m := make(map[string]struct{})
	addCompoundTokens := uint64(0)

	//0. 获取validator
	validator, exsits := e.ValidatorMapper.GetValidator(valAddr)
	if !exsits {
		log.Debug("distribute validator not exsits", "validator", valAddr.String())
		//validator不存在时, 获取delegator当前收益信息, 将收益直接返还账户中,并删除当前delegator信息
		for _, deleAddr := range delegators {
			if info, _exsits := e.DistributionMapper.GetDelegatorEarningStartInfo(valAddr, deleAddr); _exsits {
				eco.IncrAccountQOS(e.Context, deleAddr, info.HistoricalRewardFees.NilToZero())
				e.DistributionMapper.DelDelegatorEarningStartInfo(valAddr, deleAddr)
				e.DelegationMapper.DelDelegationInfo(deleAddr, valAddr)
			}
		}
		return
	}

	//1. validator汇总收益增加计费周期
	endPeriod := e.DistributionMapper.IncrementValidatorPeriod(validator)

	//2. 处理delegator收益信息
	for _, deleAddr := range delegators {
		if _, ok := m[deleAddr.String()]; ok {
			continue
		}

		m[deleAddr.String()] = struct{}{}
		addTokens := distributeDelegatorEarning(e, validator, endPeriod, deleAddr, blockHeight, periodHeightParam)
		addCompoundTokens = addCompoundTokens + addTokens
	}

	if addCompoundTokens > 0 {
		//更新validator bondTokens
		updatedTokens := validator.BondTokens + addCompoundTokens
		log.Debug("validator incr tokens", "validator", valAddr.String(), "addCompoundTokens", addCompoundTokens, "updatedTokens", updatedTokens)
		e.ValidatorMapper.ChangeValidatorBondTokens(validator, updatedTokens)
	}
}

func distributeDelegatorEarning(e eco.Eco, validator types.Validator, endPeriod uint64, deleAddr btypes.Address, blockHeight, periodHeightParam uint64) uint64 {

	valAddr := validator.GetValidatorAddress()

	log := e.Context.Logger()
	log.Debug("distribute delegator earning", "delegator", deleAddr.String(), "validator", valAddr.String(), "endPeriod", endPeriod, "height", blockHeight)

	rewards, err := e.DistributionMapper.CalculateDelegatorPeriodRewards(valAddr, deleAddr, endPeriod, blockHeight)
	if err != nil {
		log.Error("distribute delegator earning error", "delegator", deleAddr.String(), "error", err.Error())
		return 0
	}

	delegationInfo, exsits := e.DelegationMapper.GetDelegationInfo(deleAddr, valAddr)

	if !exsits || delegationInfo.Amount == 0 {
		//已无委托关系,收益直接分配到delegator账户中
		log.Debug("delegation not exsits. rewards to account", "rewards", rewards)
		eco.IncrAccountQOS(e.Context, deleAddr, rewards.NilToZero())
		e.DistributionMapper.DelDelegatorEarningStartInfo(valAddr, deleAddr)
		e.DelegationMapper.DelDelegationInfo(deleAddr, valAddr)
		return 0
	}

	//增加下一周期的收益发放信息
	nextIncomeHeight := blockHeight + periodHeightParam
	e.DistributionMapper.Set(types.BuildDelegatorPeriodIncomeKey(valAddr, deleAddr, nextIncomeHeight), true)

	//非复投,收益直接分配到delegator账户中
	if !delegationInfo.IsCompound {
		log.Debug("delegation is not compound. rewards to delegator account", "rewards", rewards)
		eco.IncrAccountQOS(e.Context, deleAddr, rewards.NilToZero())
		return 0
	}

	//复投
	addTokens := uint64(rewards.Int64())
	log.Debug("delegation is compound. rewards to delegation tokens", "addTokens", addTokens)

	//更新delegation委托信息,更新delegate当前收益信息
	info, _ := e.DistributionMapper.GetDelegatorEarningStartInfo(valAddr, deleAddr)
	info.BondToken = info.BondToken + addTokens
	e.DistributionMapper.Set(types.BuildDelegatorEarningStartInfoKey(valAddr, deleAddr), info)

	delegationInfo.Amount = delegationInfo.Amount + addTokens
	e.DelegationMapper.SetDelegationInfo(delegationInfo)

	return addTokens
}

// 2.  每块挖出的QOS数量:  `x%`proposer + `y%`validators + `z%`community
//        * `x%`proposer: 验证人获得的奖励,直接归属proposer
//        * `y%`validators: 根据每个validator的power占比平均分配
// 3.  validator奖励数 =  validator佣金 +  平分金额Fee
//        * validator佣金奖励: 佣金 = validator奖励数 * `commission rate`
//        * 平分金额Fee由validator,delegator根据各自绑定的stake平均分配
// 4.  validator的proposer奖励,佣金奖励 均按周期发放
//
func allocateQOS(ctx context.Context, signedTotalPower, totalPower int64, proposerAddr btypes.Address, votes []abci.VoteInfo) {

	e := eco.GetEco(ctx)
	log := ctx.Logger()

	params := e.DistributionMapper.GetParams()

	//获取待分配的QOS总量
	totalAmount := e.DistributionMapper.GetPreDistributionQOS()
	remainQOS := totalAmount
	e.DistributionMapper.ClearPreDistributionQOS()

	log.Debug("total rewards", "total rewards", totalAmount, "height", ctx.BlockHeight())
	//proposer奖励,直接归属proposer
	proposerRewards := params.ProposerRewardRate.MultiBigInt(totalAmount)
	proposerValidater, exsits := e.ValidatorMapper.GetValidator(proposerAddr)
	if !exsits {
		log.Error("proposer validator not exsits", "proposer", proposerAddr)
	} else {
		if info, exsits := e.DistributionMapper.GetDelegatorEarningStartInfo(proposerAddr, proposerValidater.Owner); exsits {
			log.Debug("reward proposer", "proposer", proposerAddr.String(), "owner", proposerValidater.Owner.String(), "rewards", proposerRewards)
			info.HistoricalRewardFees = info.HistoricalRewardFees.Add(proposerRewards)
			remainQOS = remainQOS.Sub(proposerRewards)
			e.DistributionMapper.Set(types.BuildDelegatorEarningStartInfoKey(proposerAddr, proposerValidater.Owner), info)
		}
	}

	//vote奖励
	votePercent := qtypes.OneFraction().Sub(params.ProposerRewardRate).Sub(params.CommunityRewardRate)
	for _, vote := range votes {
		votePowerFrac := qtypes.NewFraction(vote.Validator.Power, totalPower)
		rewards := votePowerFrac.Mul(votePercent).MultiBigInt(totalAmount)
		log.Debug("reward validator", "validator", btypes.Address(vote.Validator.Address).String(), "power", vote.Validator.Power, "total rewards", rewards)
		remainQOS = remainQOS.Sub(rewards)
		rewardToValidator(e, vote.Validator.Address, rewards, params.ValidatorCommissionRate)
	}

	//社区奖励
	communityFeePool := e.DistributionMapper.GetCommunityFeePool()
	communityFeePool = communityFeePool.Add(remainQOS)
	log.Debug("reward community", "rewards", remainQOS)
	e.DistributionMapper.SetCommunityFeePool(communityFeePool)
}

func rewardToValidator(e eco.Eco, valAddr btypes.Address, rewards btypes.BigInt, commissionRate qtypes.Fraction) {

	log := e.Context.Logger()

	commissionReward := commissionRate.MultiBigInt(rewards)
	sharedReward := rewards.Sub(commissionReward)

	validator, exsits := e.ValidatorMapper.GetValidator(valAddr)
	if !exsits {
		log.Error("reward validator, validator not exsits", "validator", valAddr.String())
		return
	}

	//validator 佣金收益
	if info, exsits := e.DistributionMapper.GetDelegatorEarningStartInfo(valAddr, validator.Owner); exsits {
		info.HistoricalRewardFees = info.HistoricalRewardFees.Add(commissionReward)
		e.DistributionMapper.Set(types.BuildDelegatorEarningStartInfoKey(valAddr, validator.Owner), info)
		log.Debug("reward validator commission", "validator", valAddr.String(), "commissionReward", commissionReward)
	}

	//delegator 共同收益
	if vcps, exsits := e.DistributionMapper.GetValidatorCurrentPeriodSummary(valAddr); exsits {
		vcps.Fees = vcps.Fees.Add(sharedReward)
		e.DistributionMapper.Set(types.BuildValidatorCurrentPeriodSummaryKey(valAddr), vcps)
		log.Debug("reward validator shared", "validator", valAddr.String(), "sharedReward", sharedReward)
	}
}
