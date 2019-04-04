package gov

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	gtypes "github.com/QOSGroup/qos/module/gov/types"
	"github.com/QOSGroup/qos/types"
	"time"
)

const (
	// Default period for deposits & voting
	DefaultPeriod = 86400 * 2 * time.Second // 2 days
)

type GenesisProposal struct {
	Proposal gtypes.Proposal  `json:"proposal"`
	Deposits []gtypes.Deposit `json:"deposits"`
	Votes    []gtypes.Vote    `json:"votes"`
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
	err := ValidateGenesis(data)
	if err != nil {
		panic(err)
	}
	mapper := GetGovMapper(ctx)
	err = mapper.setInitialProposalID(ctx, data.StartingProposalID)
	if err != nil {
		panic(err)
	}
	mapper.SetParams(ctx, data.Params)
	for _, proposal := range data.Proposals {
		switch proposal.Proposal.Status {
		case gtypes.StatusDepositPeriod:
			mapper.InsertInactiveProposalQueue(proposal.Proposal.DepositEndTime, proposal.Proposal.ProposalID)
		case gtypes.StatusVotingPeriod:
			mapper.InsertActiveProposalQueue(proposal.Proposal.VotingEndTime, proposal.Proposal.ProposalID)
		}
		for _, deposit := range proposal.Deposits {
			mapper.setDeposit(ctx, deposit.ProposalID, deposit.Depositor, deposit)
		}
		for _, vote := range proposal.Votes {
			mapper.setVote(vote.ProposalID, vote.Voter, vote)
		}
		mapper.SetProposal(proposal.Proposal)
	}
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx context.Context) GenesisState {
	mapper := GetGovMapper(ctx)
	startingProposalID, _ := mapper.peekCurrentProposalID()
	params := mapper.GetParams(ctx)
	proposals := mapper.GetProposals()
	var genesisProposals []GenesisProposal
	for _, proposal := range proposals {
		var deposits []gtypes.Deposit
		var votes []gtypes.Vote
		proposalID := proposal.ProposalID
		depositsIterator := mapper.GetDeposits(proposalID)
		defer depositsIterator.Close()
		for ; depositsIterator.Valid(); depositsIterator.Next() {
			var deposit gtypes.Deposit
			mapper.GetCodec().MustUnmarshalBinaryBare(depositsIterator.Value(), &deposit)
			deposits = append(deposits, deposit)
		}
		votesIterator := mapper.GetVotes(proposalID)
		defer votesIterator.Close()
		for ; votesIterator.Valid(); votesIterator.Next() {
			var vote gtypes.Vote
			mapper.GetCodec().MustUnmarshalBinaryBare(votesIterator.Value(), &vote)
			votes = append(votes, vote)
		}
		genesisProposals = append(genesisProposals, GenesisProposal{proposal, deposits, votes})
	}

	return GenesisState{
		StartingProposalID: startingProposalID,
		Params:             params,
		Proposals:          genesisProposals,
	}
}

func PrepForZeroHeightGenesis(ctx context.Context) {
	mapper := GetGovMapper(ctx)
	proposals := mapper.GetProposalsFiltered(ctx, nil, nil, gtypes.StatusDepositPeriod, 0)
	for _, proposal := range proposals {
		proposalID := proposal.ProposalID
		mapper.RefundDeposits(ctx, proposalID)
		mapper.DeleteProposal(ctx, proposalID)
	}

	proposals = mapper.GetProposalsFiltered(ctx, nil, nil, gtypes.StatusVotingPeriod, 0)
	for _, proposal := range proposals {
		proposalID := proposal.ProposalID
		mapper.RefundDeposits(ctx, proposalID)
		mapper.DeleteVotes(proposalID)
		mapper.DeleteProposal(ctx, proposalID)
	}
}
