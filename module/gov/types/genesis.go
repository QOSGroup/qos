package types

import (
	"time"
)

const (
	// Default period for deposits & voting
	DefaultDepositPeriod = 7 * 24 * time.Hour  // 7 days
	DefaultVotingPeriod  = 14 * 24 * time.Hour // 14 days
)

type GenesisProposal struct {
	Proposal Proposal  `json:"proposal"` // 提议
	Deposits []Deposit `json:"deposits"` // 质押
	Votes    []Vote    `json:"votes"`    // 投票
}

// 创世状态
type GenesisState struct {
	StartingProposalID int64             `json:"starting_proposal_id"` // 下一个提议ID
	Params             Params            `json:"params"`               // 提议相关参数
	Proposals          []GenesisProposal `json:"proposals"`            // 提议
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
