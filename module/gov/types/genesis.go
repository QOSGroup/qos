package types

import (
	"fmt"
	"github.com/QOSGroup/qos/types"
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
	for _, level := range ProposalLevels {
		levelParams := data.Params.GetLevelParams(level)
		threshold := levelParams.Threshold
		if threshold.IsNegative() || threshold.GT(types.OneDec()) {
			return fmt.Errorf("governance vote threshold should be positive and less or equal to one, is %s",
				threshold.String())
		}

		veto := levelParams.Veto
		if veto.IsNegative() || veto.GT(types.OneDec()) {
			return fmt.Errorf("governance vote veto threshold should be positive and less or equal to one, is %s",
				veto.String())
		}

		if levelParams.MaxDepositPeriod > levelParams.VotingPeriod {
			return fmt.Errorf("governance deposit period should be less than or equal to the voting period (%ds), is %ds",
				levelParams.VotingPeriod, levelParams.MaxDepositPeriod)
		}

		if levelParams.MinDeposit <= 0 {
			return fmt.Errorf("governance deposit amount must be a valid sdk.Coins amount, is %v",
				levelParams.MinDeposit)
		}
	}

	return nil
}
