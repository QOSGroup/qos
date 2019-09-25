# ABCI

## EndBlocker

通胀奖励计算：
```go
// 首次奖励区块高度，由于第一块的时间取的是genesis.json的创建时间，在计算的奖励值上存在较大误差，这里取的是第二块的时间
firstBlockTime := mintMapper.GetFirstBlockTime()

// 平均出块时间
blockTimeAvg := (currentBlockTime.Unix() - firstBlockTime) / (height - BeginRewardHeight)

// 当前通胀阶段预计生成区块数
blocks := int64(currentPhrase.EndTime.Sub(currentBlockTime).Seconds()) / blockTimeAvg

// 当前通胀平均块奖励
rewardPerBlock := currentPhrase.TotalAmount.Sub(currentPhrase.AppliedAmount).DivRaw(blocks)
```

`rewardPerBlock`会累加到`distribution`模块中的`PreDistributionQOS`，与每块中收集到的`tx gas`一起作为挖矿奖励分发给验证节点及其质押账户。