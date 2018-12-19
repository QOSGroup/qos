package miner

import (
	"github.com/QOSGroup/qbase/context"

	abci "github.com/tendermint/tendermint/abci/types"
)

//BeginBlocker: 挖矿奖励
func BeginBlocker(ctx context.Context, req abci.RequestBeginBlock) {

	if ctx.BlockHeight() > 1 {
		rewardVoteValidator(ctx, req)
	}

}

//基于投票的挖矿奖励: 10QOS*valVotePower/totalVotePower
func rewardVoteValidator(ctx context.Context, req abci.RequestBeginBlock) {

	// logger := ctx.Logger()
	// accountMapper := baseabci.GetAccountMapper(ctx)
	// validatorMapper := ctx.Mapper(validator.ValidatorMapperName).(*validator.ValidatorMapper)

	// totalVotePower := int64(0)
	// for _, val := range req.LastCommitInfo.Validators {
	// 	if val.SignedLastBlock {
	// 		totalVotePower += val.Validator.Power
	// 	}
	// }

	// if totalVotePower <= int64(0) {
	// 	logger.Error(fmt.Sprintf("totalVotePower: %d lte 0", totalVotePower))
	// 	return
	// }

	// for _, val := range req.LastCommitInfo.Validators {
	// 	if val.SignedLastBlock {
	// 		//reward
	// 		consAddress := types.Address(val.Validator.Address)

	// 		qVal, exsits := validatorMapper.GetByConsAddress(consAddress)
	// 		if !exsits {
	// 			logger.Error(fmt.Sprintf("consAddress: %s not exsits", consAddress))
	// 			continue
	// 		}

	// 		acc := accountMapper.GetAccount(qVal.Operator)

	// 		if qosAcc, ok := acc.(*qacc.QOSAccount); ok {
	// 			rewardQos := calRewardQos(val.Validator.Power, totalVotePower)
	// 			logger.Debug(fmt.Sprintf("address: %s add vote reward: %s", qosAcc.GetAddress().String(), rewardQos))
	// 			qosAcc.SetQOS(qosAcc.GetQOS().NilToZero().Add(rewardQos))
	// 			accountMapper.SetAccount(acc)
	// 		}
	// 	}
	// }

}

// func calRewardQos(valVotePower int64, totalVotePower int64) types.BigInt {
// 	t := types.NewInt(qtypes.BlockReward).Mul(types.NewInt(valVotePower))
// 	return t.Div(types.NewInt(totalVotePower))
// }
