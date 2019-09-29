# 参数

分配模块参数：

```go
type Params struct {
	ProposerRewardRate           qtypes.Dec `json:"proposer_reward_rate"`           // 块提议者奖励比例
	CommunityRewardRate          qtypes.Dec `json:"community_reward_rate"`          // 社区奖励比例
	DelegatorsIncomePeriodHeight int64      `json:"delegator_income_period_height"` // 奖励发放周期
	GasPerUnitCost               int64      `json:"gas_per_unit_cost"`              // 1QOS折算Gas量
}
```

所有参数均可通过[治理](../gov)模块提交参数修改提议进行修改。