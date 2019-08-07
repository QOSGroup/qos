package types

import (
	"fmt"
	"time"

	"github.com/QOSGroup/qos/types"
)

const (
	// Default period for deposits & voting
	DefaultPeriod = /*86400*/60 * 2 * time.Second // 2 days
)

type GenesisProposal struct {
	Proposal Proposal  `json:"proposal"`
	Deposits []Deposit `json:"deposits"`
	Votes    []Vote    `json:"votes"`
}

type GenesisState struct {
	StartingProposalID uint64            `json:"starting_proposal_id"`
	Params             Params            `json:"params"`
	Proposals          []GenesisProposal `json:"proposals"`
}

func NewGenesisState(startingProposalID uint64, params Params) GenesisState {
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
	threshold := data.Params.Threshold
	if threshold.IsNegative() || threshold.GT(types.OneDec()) {
		return fmt.Errorf("governance vote threshold should be positive and less or equal to one, is %s",
			threshold.String())
	}

	veto := data.Params.Veto
	if veto.IsNegative() || veto.GT(types.OneDec()) {
		return fmt.Errorf("governance vote veto threshold should be positive and less or equal to one, is %s",
			veto.String())
	}

	if data.Params.MaxDepositPeriod > data.Params.VotingPeriod {
		return fmt.Errorf("governance deposit period should be less than or equal to the voting period (%ds), is %ds",
			data.Params.VotingPeriod, data.Params.MaxDepositPeriod)
	}

	if data.Params.MinDeposit <= 0 {
		return fmt.Errorf("governance deposit amount must be a valid sdk.Coins amount, is %v",
			data.Params.MinDeposit)
	}

	return nil
}
