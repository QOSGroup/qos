# Parameters

```go
type Params struct {
	// params of normal level
	// DepositParams
	NormalMinDeposit             btypes.BigInt `json:"normal_min_deposit"`               //  Minimum deposit for a proposal to enter voting period.
	NormalMinProposerDepositRate qtypes.Dec    `json:"normal_min_proposer_deposit_rate"` //  Minimum deposit rate for proposer to submit a proposal.
	NormalMaxDepositPeriod       time.Duration `json:"normal_max_deposit_period"`        //  Maximum period for Atom holders to deposit on a proposal.
	// VotingParams
	NormalVotingPeriod time.Duration `json:"normal_voting_period"` //  Length of the voting period.
	// TallyParams
	NormalQuorum    qtypes.Dec `json:"normal_quorum"`    //  Minimum percentage of total stake needed to vote for a result to be considered valid
	NormalThreshold qtypes.Dec `json:"normal_threshold"` //  Minimum propotion of Yes votes for proposal to pass.
	NormalVeto      qtypes.Dec `json:"normal_veto"`      //  Minimum value of Veto votes to Total votes ratio for proposal to be vetoed.
	NormalPenalty   qtypes.Dec `json:"normal_penalty"`   //  Penalty if validator does not vote
	// BurnRate
	NormalBurnRate qtypes.Dec `json:"normal_burn_rate"` // Deposit burning rate when proposals pass or reject.

	// params of important level
	// DepositParams
	ImportantMinDeposit             btypes.BigInt `json:"important_min_deposit"`               //  Minimum deposit for a proposal to enter voting period.
	ImportantMinProposerDepositRate qtypes.Dec    `json:"important_min_proposer_deposit_rate"` //  Minimum deposit rate for proposer to submit a proposal.
	ImportantMaxDepositPeriod       time.Duration `json:"important_max_deposit_period"`        //  Maximum period for Atom holders to deposit on a proposal.
	// VotingParams
	ImportantVotingPeriod time.Duration `json:"important_voting_period"` //  Length of the voting period.
	// TallyParams
	ImportantQuorum    qtypes.Dec `json:"important_quorum"`    //  Minimum percentage of total stake needed to vote for a result to be considered valid
	ImportantThreshold qtypes.Dec `json:"important_threshold"` //  Minimum propotion of Yes votes for proposal to pass.
	ImportantVeto      qtypes.Dec `json:"important_veto"`      //  Minimum value of Veto votes to Total votes ratio for proposal to be vetoed.
	ImportantPenalty   qtypes.Dec `json:"important_penalty"`   //  Penalty if validator does not vote
	// BurnRate
	ImportantBurnRate qtypes.Dec `json:"important_burn_rate"` // Deposit burning rate when proposals pass or reject.

	// params of critical level
	// DepositParams
	CriticalMinDeposit             btypes.BigInt `json:"critical_min_deposit"`               //  Minimum deposit for a proposal to enter voting period.
	CriticalMinProposerDepositRate qtypes.Dec    `json:"critical_min_proposer_deposit_rate"` //  Minimum deposit rate for proposer to submit a proposal.
	CriticalMaxDepositPeriod       time.Duration `json:"critical_max_deposit_period"`        //  Maximum period for Atom holders to deposit on a proposal.
	// VotingParams
	CriticalVotingPeriod time.Duration `json:"critical_voting_period"` //  Length of the voting period.
	// TallyParams
	CriticalQuorum    qtypes.Dec `json:"critical_quorum"`    //  Minimum percentage of total stake needed to vote for a result to be considered valid
	CriticalThreshold qtypes.Dec `json:"critical_threshold"` //  Minimum propotion of Yes votes for proposal to pass.
	CriticalVeto      qtypes.Dec `json:"critical_veto"`      //  Minimum value of Veto votes to Total votes ratio for proposal to be vetoed.
	CriticalPenalty   qtypes.Dec `json:"critical_penalty"`   //  Penalty if validator does not vote
	// BurnRate
	CriticalBurnRate qtypes.Dec `json:"critical_burn_rate"` // Deposit burning rate when proposals pass or reject.
}
```

All the parameters can be changed by [governance proposal](../gov).