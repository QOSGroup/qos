package gov

import (
	"time"

	"github.com/QOSGroup/qos/types"
)

// Params returns all of the governance params
type Params struct {
	// DepositParams
	MinDeposit       uint64        `json:"min_deposit"`        //  Minimum deposit for a proposal to enter voting period.
	MaxDepositPeriod time.Duration `json:"max_deposit_period"` //  Maximum period for Atom holders to deposit on a proposal. Initial value: 2 months

	// VotingParams
	VotingPeriod time.Duration `json:"voting_period"` //  Length of the voting period.

	// TallyParams
	Quorum    types.Dec `json:"quorum"`    //  Minimum percentage of total stake needed to vote for a result to be considered valid
	Threshold types.Dec `json:"threshold"` //  Minimum propotion of Yes votes for proposal to pass. Initial value: 0.5
	Veto      types.Dec `json:"veto"`      //  Minimum value of Veto votes to Total votes ratio for proposal to be vetoed. Initial value: 1/3
	Penalty   types.Dec `json:"penalty"`   //  Penalty if validator does not vote
}
