package miner

import (
	"fmt"

	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/types"

	qacc "github.com/QOSGroup/qos/account"
	qtypes "github.com/QOSGroup/qos/types"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

//BeginBlocker: 挖矿奖励
func BeginBlocker(ctx context.Context, req abci.RequestBeginBlock) {

	if ctx.BlockHeight() > 1 {
		accountMapper := baseabci.GetAccountMapper(ctx)
		rewardVoteValidator(accountMapper, req, ctx.Logger())
	}

}

//基于投票的挖矿奖励: 10QOS*valVotePower/totalVotePower
func rewardVoteValidator(accMapper *account.AccountMapper, req abci.RequestBeginBlock, logger log.Logger) {

	totalVotePower := int64(0)
	for _, val := range req.LastCommitInfo.Validators {
		if val.SignedLastBlock {
			totalVotePower += val.Validator.Power
		}
	}

	if totalVotePower <= int64(0) {
		logger.Error(fmt.Sprintf("totalVotePower: %d lte 0", totalVotePower))
		return
	}

	for _, val := range req.LastCommitInfo.Validators {
		if val.SignedLastBlock {
			//reward
			addr := types.Address(val.Validator.Address)
			acc := accMapper.GetAccount(addr)

			if qosAcc, ok := acc.(*qacc.QOSAccount); ok {
				rewardQos := calRewardQos(val.Validator.Power, totalVotePower)
				logger.Debug(fmt.Sprintf("address: %s add vote reward: %s", addr.String(), rewardQos))
				qosAcc.SetQOS(qosAcc.GetQOS().NilToZero().Add(rewardQos))
				accMapper.SetAccount(acc)
			}
		}
	}

}

func calRewardQos(valVotePower int64, totalVotePower int64) types.BigInt {
	t := types.NewInt(qtypes.BlockReward).Mul(types.NewInt(valVotePower))
	return t.Div(types.NewInt(totalVotePower))
}
