# ABCI

## EndBlocker

Inflation reward calculation:
```go
// Since the time of the first block is the creation time of genesis.json, there is a large error in the calculated bonus value, here we save the time of the second block.
firstBlockTime := mintMapper.GetFirstBlockTime()

// average block time
blockTimeAvg := (currentBlockTime.Unix() - firstBlockTime) / (height - BeginRewardHeight)

// the number of blocks expected to be generated in the current inflation phase
blocks := int64(currentPhrase.EndTime.Sub(currentBlockTime).Seconds()) / blockTimeAvg

// current inflation average block reward
rewardPerBlock := currentPhrase.TotalAmount.Sub(currentPhrase.AppliedAmount).DivRaw(blocks)
```

`rewardPerBlock` will be added to `PreDistributionQOS` in `distribution` module, together with the `tx gas` collected in each block as the mining reward.