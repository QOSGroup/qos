package mint

import (
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/distribution"
	"github.com/QOSGroup/qos/module/mint/mapper"
	"github.com/QOSGroup/qos/module/mint/types"

	abci "github.com/tendermint/tendermint/abci/types"
)

//BeginBlocker: 挖矿奖励
func BeginBlocker(ctx context.Context, req abci.RequestBeginBlock) {
	height := uint64(ctx.BlockHeight())
	currentBlockTime := ctx.BlockHeader().Time.UTC().Unix()

	mintMapper := mapper.GetMapper(ctx)
	currentInflationPhrase, exist := mintMapper.GetCurrentInflationPhrase(uint64(currentBlockTime))
	if exist == false || currentInflationPhrase.TotalAmount == 0 {
		return
	}

	if height == uint64(1) {
		mintMapper.SetFirstBlockTime(currentBlockTime)
	}

	// for the first block, assuming average block time is 5s
	blockTimeAvg := uint64(5)
	if height > 1 {
		firstBlockTime := mintMapper.GetFirstBlockTime()
		blockTimeAvg = uint64(currentBlockTime-firstBlockTime) / (height - 1)
	}

	totalQOSAmount := currentInflationPhrase.TotalAmount
	blocks := (uint64(currentInflationPhrase.EndTime.UTC().Unix()) - uint64(currentBlockTime)) / blockTimeAvg
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
			mintMapper.MintQOS(uint64(currentBlockTime), rewardPerBlock)
			distributionMapper := distribution.GetMapper(ctx)
			distributionMapper.AddPreDistributionQOS(btypes.NewInt(int64(rewardPerBlock)))

			ctx.EventManager().EmitEvent(
				btypes.NewEvent(
					types.EventTypeMint,
					btypes.NewAttribute(types.AttributeKeyHeight, string(height)),
					btypes.NewAttribute(types.AttributeKeyTokens, string(rewardPerBlock)),
				),
			)
		}
	}
}
