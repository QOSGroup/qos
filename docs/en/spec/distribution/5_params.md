# Parameters

```go
type Params struct {
	ProposerRewardRate           qtypes.Dec `json:"proposer_reward_rate"`           // proposer reward rate
	CommunityRewardRate          qtypes.Dec `json:"community_reward_rate"`          // community reward rate
	DelegatorsIncomePeriodHeight int64      `json:"delegator_income_period_height"` // reward period
	GasPerUnitCost               int64      `json:"gas_per_unit_cost"`              // 1 QOS converts the amount of Gas
}
```

All the parameters can be changed by [governance proposal](../gov).