# 参数

```go
type Params struct {
	MaxValidatorCnt             int64         `json:"max_validator_cnt"`          // 最多验证节点数量
	ValidatorVotingStatusLen    int64         `json:"voting_status_len"`          // 投票窗口高度
	ValidatorVotingStatusLeast  int64         `json:"voting_status_least"`        // 最低投票高度
	ValidatorSurvivalSecs       int64         `json:"survival_secs"`              // inactive状态验证节点状态保持时间
	DelegatorUnbondFrozenHeight int64         `json:"unbond_frozen_height"`       // 解委托token锁定高度
	MaxEvidenceAge              time.Duration `json:"max_evidence_age"`           // 证据数据有效时长
	SlashFractionDoubleSign     types.Dec     `json:"slash_fraction_double_sign"` // 双签惩罚比例
	SlashFractionDowntime       types.Dec     `json:"slash_fraction_downtime"`    // 漏块惩罚比例
}
```

所有参数均可通过[治理](../gov)模块提交参数修改提议进行修改。