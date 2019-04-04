package gov

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qos/types"
	"time"
)

const (
	// Default period for deposits & voting
	DefaultPeriod = /*86400 **/ 2 * time.Minute // 2 days
)

type GenesisState struct {
	StartingProposalID uint64 `json:"starting_proposal_id"`
	Params             Params `json:"params"`
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
		Params: Params{
			MinDeposit:       10,
			MaxDepositPeriod: DefaultPeriod,
			VotingPeriod:     DefaultPeriod,
			Quorum:           types.NewDecWithPrec(334, 3),
			Threshold:        types.NewDecWithPrec(5, 1),
			Veto:             types.NewDecWithPrec(334, 3),
			Penalty:          types.ZeroDec(),
		},
	}
}

// ValidateGenesis
func ValidateGenesis(data GenesisState) error {
	threshold := data.Params.Threshold
	if threshold.IsNegative() || threshold.GT(types.OneDec()) {
		return fmt.Errorf("Governance vote threshold should be positive and less or equal to one, is %s",
			threshold.String())
	}

	veto := data.Params.Veto
	if veto.IsNegative() || veto.GT(types.OneDec()) {
		return fmt.Errorf("Governance vote veto threshold should be positive and less or equal to one, is %s",
			veto.String())
	}

	if data.Params.MaxDepositPeriod > data.Params.VotingPeriod {
		return fmt.Errorf("Governance deposit period should be less than or equal to the voting period (%ds), is %ds",
			data.Params.VotingPeriod, data.Params.MaxDepositPeriod)
	}

	if data.Params.MinDeposit <= 0 {
		return fmt.Errorf("Governance deposit amount must be a valid sdk.Coins amount, is %v",
			data.Params.MinDeposit)
	}

	return nil
}

// InitGenesis - store genesis parameters
func InitGenesis(ctx context.Context, data GenesisState) {
	mapper := GetGovMapper(ctx)
	err := mapper.setInitialProposalID(ctx, data.StartingProposalID)
	if err != nil {
		panic(err)
	}
	mapper.SetParams(ctx, data.Params)
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx context.Context) GenesisState {
	mapper := GetGovMapper(ctx)
	startingProposalID, _ := mapper.peekCurrentProposalID(ctx)
	params := mapper.GetParams(ctx)

	return GenesisState{
		StartingProposalID: startingProposalID,
		Params:             params,
	}
}
