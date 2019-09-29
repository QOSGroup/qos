# Parameters

```go
type Params struct {
	MaxValidatorCnt             int64         `json:"max_validator_cnt"`          // max validator counts
	ValidatorVotingStatusLen    int64         `json:"voting_status_len"`          // voting window length
	ValidatorVotingStatusLeast  int64         `json:"voting_status_least"`        // min vote times in voting window
	ValidatorSurvivalSecs       int64         `json:"survival_secs"`              // inactive survive seconds
	DelegatorUnbondFrozenHeight int64         `json:"unbond_frozen_height"`       // unbond and redelegation completing height
	MaxEvidenceAge              time.Duration `json:"max_evidence_age"`           // max evidence age
	SlashFractionDoubleSign     types.Dec     `json:"slash_fraction_double_sign"` // double sign slash fraction
	SlashFractionDowntime       types.Dec     `json:"slash_fraction_downtime"`    // downtime slash fraction 
}
```

All the parameters can be changed by [governance proposal](../gov).