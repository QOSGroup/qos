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
	log := ctx.Logger()
	height := uint64(ctx.BlockHeight())

	mintMapper := ctx.Mapper(staketypes.MintMapperName).(*stakemapper.MintMapper)
	currentInflationPhrase, exist := mintMapper.GetCurrentInflationPhrase()
	if exist == false || currentInflationPhrase.TotalAmount == 0 {
		return
	}

	if height == uint64(1) {
		mintMapper.SetFirstBlockTime(ctx.BlockHeader().Time.UTC().Unix())
	}

	// for the first block, assuming average block time is 5s
	blockTimeAvg := uint64(5)
	if height > 1 {
		firstBlockTime := mintMapper.GetFirstBlockTime()
		//average block time = (last block time - first block time) / (last height - 1)
		blockTimeAvg = uint64(ctx.BlockHeader().Time.UTC().Unix()-firstBlockTime) / (height - 1)
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
			log.Debug("block mint", "height", height, "mint", rewardPerBlock)
			distributionMapper := stakemapper.GetDistributionMapper(ctx)
			distributionMapper.AddPreDistributionQOS(btypes.NewInt(int64(rewardPerBlock)))
		}
	}
}
