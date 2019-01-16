package mint

import (
	"fmt"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	stakemapper "github.com/QOSGroup/qos/module/eco/mapper"
	staketypes "github.com/QOSGroup/qos/module/eco/types"
	"github.com/QOSGroup/qos/types"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
)

//BeginBlocker: 挖矿奖励
func BeginBlocker(ctx context.Context, req abci.RequestBeginBlock) {
	height := uint64(ctx.BlockHeight())
	mintMapper := ctx.Mapper(staketypes.MintMapperName).(*stakemapper.MintMapper)
	currentInflationPhrase, exist := mintMapper.GetCurrentInflationPhrase()
	if exist == false || currentInflationPhrase.TotalAmount == 0 {
		return
	}

	// for the first block, assuming average block time is 5s
	blockTimeAvg := uint64(5)
	if height > 1 {
		//average block time = (last block time - first block time) / (last height - 1)
		blockTimeAvg = uint64(ctx.BlockHeader().Time.UTC().Sub(ctx.WithBlockHeight(1).BlockHeader().Time.UTC()).Seconds()) / (height - 1)
	}

	totalQOSAmount := currentInflationPhrase.TotalAmount
	blocks := (uint64(currentInflationPhrase.EndTime.UTC().Unix()) - uint64(time.Now().UTC().Unix()))/blockTimeAvg

	//totalBlock := mintMapper.GetParams().TotalBlock
	appliedQOSAmount := currentInflationPhrase.AppliedAmount

	if appliedQOSAmount >= totalQOSAmount {
		return
	}
	if blocks <= 0 {
		return
	}

	if ctx.BlockHeight() > 1 {
		rewardPerBlock := (totalQOSAmount - appliedQOSAmount) / blocks

		if rewardPerBlock > 0 {
			//distributionMapper := ctx.Mapper(DistributionMapperName).(*DistributionMapper)
			//distributionMapper.AddPreDistributionQOS(btypes.NewInt(int64(rewardPerBlock)))
			rewardVoteValidator(ctx, req, rewardPerBlock)
		}
	}
}

//基于投票的挖矿奖励: 10QOS*valVotePower/totalVotePower
func rewardVoteValidator(ctx context.Context, req abci.RequestBeginBlock, rewardPerBlock uint64) {

	logger := ctx.Logger()

	mintMapper := ctx.Mapper(staketypes.MintMapperName).(*stakemapper.MintMapper)
	accountMapper := baseabci.GetAccountMapper(ctx)
	validatorMapper := stakemapper.GetValidatorMapper(ctx)

	totalVotePower := int64(0)
	for _, val := range req.LastCommitInfo.Votes {
		if val.SignedLastBlock {
			totalVotePower += val.Validator.Power
		}
	}

	if totalVotePower <= int64(0) {
		logger.Error(fmt.Sprintf("totalVotePower: %d lte 0", totalVotePower))
		return
	}

	actualAppliedQOSAccount := btypes.NewInt(0)

	for _, val := range req.LastCommitInfo.Votes {
		if val.SignedLastBlock {
			//reward
			addr := btypes.Address(val.Validator.Address)
			validator, exsits := validatorMapper.GetValidator(addr)
			if !exsits {
				logger.Error(fmt.Sprintf("validator: %s not exsits", addr.String()))
				continue
			}

			acc := accountMapper.GetAccount(validator.Owner)
			if qosAcc, ok := acc.(*types.QOSAccount); ok {
				rewardQos := calRewardQos(val.Validator.Power, totalVotePower, rewardPerBlock)
				logger.Debug(fmt.Sprintf("address: %s add vote reward: %s", qosAcc.GetAddress().String(), rewardQos))
				qosAcc.SetQOS(qosAcc.GetQOS().NilToZero().Add(rewardQos))
				accountMapper.SetAccount(acc)

				actualAppliedQOSAccount = actualAppliedQOSAccount.Add(rewardQos)
			}
		}
	}

	logger.Info("mint reward", "predict", rewardPerBlock, "actual", actualAppliedQOSAccount.Int64())
	mintMapper.AddAppliedQOSAmount(uint64(actualAppliedQOSAccount.Int64()))

}

func calRewardQos(valVotePower int64, totalVotePower int64, rewardPerBlock uint64) btypes.BigInt {
	t := btypes.NewInt(int64(rewardPerBlock)).Mul(btypes.NewInt(valVotePower))
	return t.Div(btypes.NewInt(totalVotePower))
}
