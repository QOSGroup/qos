package types

import (
	"time"
)

const (
	// Default period for deposits & voting
	DefaultPeriod = /*86400*/ 60 * 2 * time.Second // 2 days
)

type GenesisProposal struct {
	Proposal Proposal  `json:"proposal"`
	Deposits []Deposit `json:"deposits"`
	Votes    []Vote    `json:"votes"`
}

type GenesisState struct {
	StartingProposalID int64             `json:"starting_proposal_id"`
	Params             Params            `json:"params"`
	Proposals          []GenesisProposal `json:"proposals"`
}

func NewGenesisState(startingProposalID int64, params Params) GenesisState {
	return GenesisState{
		StartingProposalID: startingProposalID,
		Params:             params,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		StartingProposalID: 1,
		Params:             DefaultParams(),
	}
}

// ValidateGenesis
func ValidateGenesis(data GenesisState) error {
	// validate params
	err := data.Params.Validate()
	if err != nil {
		return err
	}

	return nil
}
