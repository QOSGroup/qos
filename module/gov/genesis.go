package gov

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qos/module/gov/mapper"
	"github.com/QOSGroup/qos/module/gov/types"
)

func InitGenesis(ctx context.Context, data types.GenesisState) {
	// 数据校验，主要是治理参数校验
	err := types.ValidateGenesis(data)
	if err != nil {
		panic(err)
	}
	mapper := mapper.GetMapper(ctx)
	// 初始化提议参数
	mapper.SetParams(ctx, data.Params)
	// 初始化提议数据
	err = mapper.SetInitialProposalID(data.StartingProposalID)
	if err != nil {
		panic(err)
	}
	for _, proposal := range data.Proposals {
		switch proposal.Proposal.Status {
		case types.StatusDepositPeriod:
			mapper.InsertInactiveProposalQueue(proposal.Proposal.DepositEndTime, proposal.Proposal.ProposalID)
		case types.StatusVotingPeriod:
			mapper.InsertActiveProposalQueue(proposal.Proposal.VotingEndTime, proposal.Proposal.ProposalID)
		}
		for _, deposit := range proposal.Deposits {
			mapper.SetDeposit(deposit.ProposalID, deposit.Depositor, deposit)
		}
		for _, vote := range proposal.Votes {
			mapper.SetVote(vote.ProposalID, vote.Voter, vote)
		}
		mapper.SetProposal(proposal.Proposal)
	}
}

// 状态数据导出
func ExportGenesis(ctx context.Context) types.GenesisState {
	mapper := mapper.GetMapper(ctx)
	startingProposalID, _ := mapper.PeekCurrentProposalID()
	params := mapper.GetParams(ctx)
	proposals := mapper.GetProposals()
	var genesisProposals []types.GenesisProposal
	for _, proposal := range proposals {
		var deposits []types.Deposit
		var votes []types.Vote
		proposalID := proposal.ProposalID
		depositsIterator := mapper.GetDeposits(proposalID)
		for ; depositsIterator.Valid(); depositsIterator.Next() {
			var deposit types.Deposit
			mapper.GetCodec().MustUnmarshalBinaryBare(depositsIterator.Value(), &deposit)
			deposits = append(deposits, deposit)
		}
		depositsIterator.Close()
		votesIterator := mapper.GetVotes(proposalID)
		for ; votesIterator.Valid(); votesIterator.Next() {
			var vote types.Vote
			mapper.GetCodec().MustUnmarshalBinaryBare(votesIterator.Value(), &vote)
			votes = append(votes, vote)
		}
		votesIterator.Close()
		genesisProposals = append(genesisProposals, types.GenesisProposal{proposal, deposits, votes})
	}

	return types.GenesisState{
		StartingProposalID: startingProposalID,
		Params:             params,
		Proposals:          genesisProposals,
	}
}
