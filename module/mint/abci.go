package mint

import (
	"time"

	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	stakemapper "github.com/QOSGroup/qos/module/eco/mapper"
	staketypes "github.com/QOSGroup/qos/module/eco/types"

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
	blocks := (uint64(currentInflationPhrase.EndTime.UTC().Unix()) - uint64(time.Now().UTC().Unix())) / blockTimeAvg
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
			distributionMapper := stakemapper.GetDistributionMapper(ctx)
			distributionMapper.AddPreDistributionQOS(btypes.NewInt(int64(rewardPerBlock)))
		}
	}
}
