package mint

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/distribution"
	"github.com/QOSGroup/qos/module/mint/mapper"
	"github.com/QOSGroup/qos/module/mint/types"

	abci "github.com/tendermint/tendermint/abci/types"
)

// 首次奖励区块高度，由于第一块的时间取的是genesis.json的创建时间，在计算的奖励值上存在较大误差，这里取第二块开始计算挖矿奖励
const BeginRewardHeight = 2

//BeginBlocker: 挖矿奖励
func BeginBlocker(ctx context.Context, req abci.RequestBeginBlock) {
	logger := ctx.Logger().With("module", "module/mint")
	height := ctx.BlockHeight()
	currentBlockTime := ctx.BlockHeader().Time.UTC()

	mintMapper := mapper.GetMapper(ctx)
	distributionMapper := distribution.GetMapper(ctx)

	// 通胀规则
	phrases := mintMapper.MustGetInflationPhrases()

	// 当前通胀校验
	currentPhrase, exists := phrases.GetPhrase(currentBlockTime)
	if !exists || !currentPhrase.TotalAmount.GT(btypes.ZeroInt()) || currentPhrase.AppliedAmount.Equal(currentPhrase.TotalAmount) {
		return
	}

	// 处理前一通胀阶段未完整发行情况，剩余转到社区账户
	if currentPhrase.AppliedAmount.Equal(btypes.ZeroInt()) {
		if prePhrase, exists := phrases.GetPrePhrase(currentBlockTime); exists {
			if !prePhrase.AppliedAmount.Equal(prePhrase.TotalAmount) {
				prePhraseLeft := prePhrase.TotalAmount.Sub(prePhrase.AppliedAmount)
				phrases = phrases.ApplyQOS(prePhrase.EndTime, prePhraseLeft)
				distributionMapper.AddToCommunityFeePool(prePhraseLeft)
				logger.Info(fmt.Sprintf("%d pre-phrase left %s", height, prePhraseLeft.String()))
			}
		}

	}

	if height <= BeginRewardHeight {
		// 前两块块仅记录出块时间
		mintMapper.SetFirstBlockTime(currentBlockTime.Unix())
	} else {
		// 计算平均出快时间
		firstBlockTime := mintMapper.GetFirstBlockTime()
		blockTimeAvg := (currentBlockTime.Unix() - firstBlockTime) / (height - BeginRewardHeight)

		// 计算挖矿奖励
		blocks := int64(currentPhrase.EndTime.Sub(currentBlockTime).Seconds()) / blockTimeAvg
		rewardPerBlock := currentPhrase.TotalAmount.Sub(currentPhrase.AppliedAmount).DivRaw(blocks)

		logger.Info(fmt.Sprintf("%d reward per block, firstBlockTime:%d, blockTimeAvg:%d, blocks:%d, reward:%s",
			height, firstBlockTime, blockTimeAvg, blocks, rewardPerBlock.String()))

		if rewardPerBlock.GT(btypes.ZeroInt()) {
			// 保存通胀发行更新
			mintMapper.AddAllTotalMintQOSAmount(rewardPerBlock)
			phrases := phrases.ApplyQOS(currentPhrase.EndTime, rewardPerBlock)
			mintMapper.SetInflationPhrases(phrases)

			// 挖矿奖励保存至待分配
			distributionMapper.AddPreDistributionQOS(rewardPerBlock)

			// 发送事件
			ctx.EventManager().EmitEvent(
				btypes.NewEvent(
					types.EventTypeMint,
					btypes.NewAttribute(types.AttributeKeyHeight, string(height)),
					btypes.NewAttribute(types.AttributeKeyTokens, rewardPerBlock.String()),
				),
			)

			// metrics
			mintMapper.Metrics.TotalAppliedQOS.Set(float64(mintMapper.GetAllTotalMintQOSAmount().Int64()))
			mintMapper.Metrics.MintPerBlockQOS.Set(float64(rewardPerBlock.Int64()))
			mintMapper.Metrics.GasFeePerBlockQOS.Set(float64(distributionMapper.GetPreDistributionQOS().Sub(rewardPerBlock).Int64()))
		}
	}
}
