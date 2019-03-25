package gov

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qos/types"
	"time"
)

const (
	// Default period for deposits & voting
	DefaultPeriod = 86400 * 2 * time.Second // 2 days
)

type GenesisState struct {
	StartingProposalID uint64        `json:"starting_proposal_id"`
	DepositParams      DepositParams `json:"deposit_params"`
	VotingParams       VotingParams  `json:"voting_params"`
	TallyParams        TallyParams   `json:"tally_params"`
}

func NewGenesisState(startingProposalID uint64, dp DepositParams, vp VotingParams, tp TallyParams) GenesisState {
	return GenesisState{
		StartingProposalID: startingProposalID,
		DepositParams:      dp,
		VotingParams:       vp,
		TallyParams:        tp,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		StartingProposalID: 1,
		DepositParams: DepositParams{
			MinDeposit:       10000000,
			MaxDepositPeriod: DefaultPeriod,
		},
		VotingParams: VotingParams{
			VotingPeriod: DefaultPeriod,
		},
		TallyParams: TallyParams{
			Quorum:    types.NewDecWithPrec(334, 3),
			Threshold: types.NewDecWithPrec(5, 1),
			Veto:      types.NewDecWithPrec(334, 3),
		},
	}
}

// ValidateGenesis
func ValidateGenesis(data GenesisState) error {
	threshold := data.TallyParams.Threshold
	if threshold.IsNegative() || threshold.GT(types.OneDec()) {
		return fmt.Errorf("Governance vote threshold should be positive and less or equal to one, is %s",
			threshold.String())
	}

	veto := data.TallyParams.Veto
	if veto.IsNegative() || veto.GT(types.OneDec()) {
		return fmt.Errorf("Governance vote veto threshold should be positive and less or equal to one, is %s",
			veto.String())
	}

	if data.DepositParams.MaxDepositPeriod > data.VotingParams.VotingPeriod {
		return fmt.Errorf("Governance deposit period should be less than or equal to the voting period (%ds), is %ds",
			data.VotingParams.VotingPeriod, data.DepositParams.MaxDepositPeriod)
	}

	if data.DepositParams.MinDeposit <= 0 {
		return fmt.Errorf("Governance deposit amount must be a valid sdk.Coins amount, is %v",
			data.DepositParams.MinDeposit)
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
	mapper.setDepositParams(data.DepositParams)
	mapper.setVotingParams(data.VotingParams)
	mapper.setTallyParams(data.TallyParams)
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx context.Context) GenesisState {
	mapper := GetGovMapper(ctx)
	startingProposalID, _ := mapper.peekCurrentProposalID(ctx)
	depositParams := mapper.GetDepositParams()
	votingParams := mapper.GetVotingParams()
	tallyParams := mapper.GetTallyParams()

	return GenesisState{
		StartingProposalID: startingProposalID,
		DepositParams:      depositParams,
		VotingParams:       votingParams,
		TallyParams:        tallyParams,
	}
}
