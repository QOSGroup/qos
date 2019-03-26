package gov

import (
	"fmt"
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/eco/mapper"
	gtypes "github.com/QOSGroup/qos/module/gov/types"
	"github.com/QOSGroup/qos/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Called every block, process inflation, update validator set
func EndBlocker(ctx context.Context) btypes.Tags {
	logger := ctx.Logger().With("module", "module/gov")
	resTags := btypes.NewTags()

	mapper := GetGovMapper(ctx)
	inactiveIterator := mapper.InactiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	defer inactiveIterator.Close()
	for ; inactiveIterator.Valid(); inactiveIterator.Next() {
		var proposalID uint64

		mapper.GetCodec().MustUnmarshalBinaryLengthPrefixed(inactiveIterator.Value(), &proposalID)
		inactiveProposal, ok := mapper.GetProposal(ctx, proposalID)
		if !ok {
			panic(fmt.Sprintf("proposal %d does not exist", proposalID))
		}

		mapper.DeleteProposal(ctx, proposalID)
		mapper.DeleteDeposits(ctx, proposalID) // delete any associated deposits (burned)

		resTags = resTags.AppendTag(ProposalID, types.Uint64ToBigEndian(proposalID))
		resTags = resTags.AppendTag(ProposalResult, []byte(ActionProposalDropped))

		logger.Info(
			fmt.Sprintf("proposal %d (%s) didn't meet minimum deposit of %d (had only %d); deleted",
				inactiveProposal.ProposalID,
				inactiveProposal.GetTitle(),
				mapper.GetDepositParams().MinDeposit,
				inactiveProposal.TotalDeposit,
			),
		)
	}

	// fetch active proposals whose voting periods have ended (are passed the block time)
	activeIterator := mapper.ActiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	defer activeIterator.Close()
	for ; activeIterator.Valid(); activeIterator.Next() {
		var proposalID uint64

		mapper.GetCodec().MustUnmarshalBinaryLengthPrefixed(activeIterator.Value(), &proposalID)
		activeProposal, ok := mapper.GetProposal(ctx, proposalID)
		if !ok {
			panic(fmt.Sprintf("proposal %d does not exist", proposalID))
		}

		passes, tallyResults := tally(ctx, mapper, activeProposal)

		var tagValue string
		if passes {
			mapper.RefundDeposits(ctx, activeProposal.ProposalID)
			activeProposal.Status = gtypes.StatusPassed
			tagValue = ActionProposalPassed
			Execute(ctx, activeProposal, logger)
		} else {
			mapper.DeleteDeposits(ctx, activeProposal.ProposalID)
			activeProposal.Status = gtypes.StatusRejected
			tagValue = ActionProposalRejected
		}

		activeProposal.FinalTallyResult = tallyResults
		mapper.SetProposal(ctx, activeProposal)
		mapper.RemoveFromActiveProposalQueue(ctx, activeProposal.VotingEndTime, activeProposal.ProposalID)

		logger.Info(
			fmt.Sprintf(
				"proposal %d (%s) tallied; passed: %v",
				activeProposal.ProposalID, activeProposal.GetTitle(), passes,
			),
		)

		resTags = resTags.AppendTag(ProposalID, types.Uint64ToBigEndian(proposalID))
		resTags = resTags.AppendTag(ProposalResult, []byte(tagValue))
	}

	return resTags
}

func Execute(ctx context.Context, proposal gtypes.Proposal, logger log.Logger) error {
	switch proposal.GetProposalType() {
	case gtypes.ProposalTypeParameterChange:
		return nil
	case gtypes.ProposalTypeTaxUsage:
		return executeTaxUsage(ctx, proposal, logger)
	}

	return nil
}

func executeTaxUsage(ctx context.Context, proposal gtypes.Proposal, logger log.Logger) error {
	proposalContent := proposal.ProposalContent.(*gtypes.TaxUsageProposal)
	distributionMapper := mapper.GetDistributionMapper(ctx)
	accountMapper := ctx.Mapper(account.AccountMapperName).(*account.AccountMapper)
	account := accountMapper.GetAccount(proposalContent.DestAddress).(*types.QOSAccount)
	if account == nil {
		account = types.NewQOSAccountWithAddress(proposalContent.DestAddress)
	}
	feePool := distributionMapper.GetCommunityFeePool()
	qos := types.NewDec(feePool.Int64()).Mul(proposalContent.Percent).TruncateInt()
	account.MustPlusQOS(qos)
	accountMapper.SetAccount(account)

	distributionMapper.SetCommunityFeePool(feePool.Sub(qos))

	logger.Info("execute taxUsage, proposal: %d", proposal.ProposalID)

	return nil
}
