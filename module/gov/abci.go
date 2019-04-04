package gov

import (
	"fmt"
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	ecomapper "github.com/QOSGroup/qos/module/eco/mapper"
	ecotypes "github.com/QOSGroup/qos/module/eco/types"
	gtypes "github.com/QOSGroup/qos/module/gov/types"
	"github.com/QOSGroup/qos/module/params"
	"github.com/QOSGroup/qos/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Called every block, process inflation, update validator set
func EndBlocker(ctx context.Context) btypes.Tags {
	logger := ctx.Logger().With("module", "module/gov")
	resTags := btypes.NewTags()

	mapper := GetGovMapper(ctx)
	inactiveIterator := mapper.InactiveProposalQueueIterator(ctx.BlockHeader().Time)
	defer inactiveIterator.Close()
	for ; inactiveIterator.Valid(); inactiveIterator.Next() {
		var proposalID uint64

		mapper.GetCodec().UnmarshalBinaryBare(inactiveIterator.Value(), &proposalID)
		inactiveProposal, ok := mapper.GetProposal(proposalID)
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
				mapper.GetParams(ctx).MinDeposit,
				inactiveProposal.TotalDeposit,
			),
		)
	}

	// fetch active proposals whose voting periods have ended (are passed the block time)
	activeIterator := mapper.ActiveProposalQueueIterator(ctx.BlockHeader().Time)
	defer activeIterator.Close()
	for ; activeIterator.Valid(); activeIterator.Next() {
		var proposalID uint64

		mapper.GetCodec().UnmarshalBinaryBare(activeIterator.Value(), &proposalID)
		activeProposal, ok := mapper.GetProposal(proposalID)
		if !ok {
			panic(fmt.Sprintf("proposal %d does not exist", proposalID))
		}

		proposalResult, tallyResults, votingValidators := tally(ctx, mapper, activeProposal)
		var tagValue string
		switch proposalResult {
		case gtypes.PASS:
			mapper.RefundDeposits(ctx, activeProposal.ProposalID)
			activeProposal.Status = gtypes.StatusPassed
			tagValue = ActionProposalPassed
			Execute(ctx, activeProposal, logger)
			break
		case gtypes.REJECT:
			mapper.RefundDeposits(ctx, activeProposal.ProposalID)
			activeProposal.Status = gtypes.StatusRejected
			tagValue = ActionProposalRejected
			break
		case gtypes.REJECTVETO:
			mapper.DeleteDeposits(ctx, activeProposal.ProposalID)
			activeProposal.Status = gtypes.StatusRejected
			tagValue = ActionProposalRejected
			break
		}

		activeProposal.FinalTallyResult = tallyResults
		mapper.SetProposal(activeProposal)
		mapper.RemoveFromActiveProposalQueue(activeProposal.VotingEndTime, activeProposal.ProposalID)

		logger.Info(
			fmt.Sprintf(
				"proposal %d (%s) tallied; passed: %v",
				activeProposal.ProposalID, activeProposal.GetTitle(), proposalResult,
			),
		)

		resTags = resTags.AppendTag(ProposalID, types.Uint64ToBigEndian(proposalID))
		resTags = resTags.AppendTag(ProposalResult, []byte(tagValue))

		penalty := mapper.GetParams(ctx).Penalty
		if penalty.GT(types.ZeroDec()) {
			validatorMapper := ecomapper.GetValidatorMapper(ctx)
			validators := validatorMapper.GetActiveValidatorSet(false)
			for _, val := range validators {
				if _, ok := votingValidators[val.String()]; !ok {
					if validator, exists := validatorMapper.GetValidator(val); exists && validator.BondHeight < activeProposal.VotingStartHeight {
						slash(ctx, validator, penalty)
					}
				}
			}
		}

		mapper.DeleteValidatorSet(proposalID)
	}

	return resTags
}

// TODO slash
func slash(ctx context.Context, validator ecotypes.Validator, penalty types.Dec) error {
	validatorMapper := ecomapper.GetValidatorMapper(ctx)
	delegationMapper := ecomapper.GetDelegationMapper(ctx)
	distributionMapper := ecomapper.GetDistributionMapper(ctx)
	var delegations []ecotypes.DelegationInfo
	delegationMapper.IterateDelegationsValDeleAddr(validator.GetValidatorAddress(), func(valAddr btypes.Address, delAddr btypes.Address) {
		if delegation, exists := delegationMapper.GetDelegationInfo(delAddr, valAddr); exists {
			delegations = append(delegations, delegation)
		}
	})

	totalSlash := types.ZeroDec()
	height := uint64(ctx.BlockHeight())
	for _, delegation := range delegations {
		// 计算惩罚
		amountSlashed := types.NewDec(int64(delegation.Amount)).Mul(penalty)
		totalSlash = totalSlash.Add(amountSlashed)

		// 计算当前delegator收益
		updatedTokens := types.NewDec(int64(delegation.Amount)).Sub(amountSlashed)
		delegation.Amount = uint64(updatedTokens.Int64())
		if err := distributionMapper.ModifyDelegatorTokens(validator, delegation.ValidatorAddr, delegation.Amount, height); err != nil {
			return err
		}

		// 更新delegation
		delegationMapper.SetDelegationInfo(delegation)
	}

	// 更新validator
	updatedValidatorTokens := uint64(types.NewDec(int64(validator.BondTokens)).Sub(totalSlash).Int64())
	validatorMapper.ChangeValidatorBondTokens(validator, updatedValidatorTokens)

	// slash放入社区费池
	distributionMapper.AddToCommunityFeePool(totalSlash.TruncateInt())

	return nil
}

func Execute(ctx context.Context, proposal gtypes.Proposal, logger log.Logger) error {
	switch proposal.GetProposalType() {
	case gtypes.ProposalTypeParameterChange:
		return executeParameterChange(ctx, proposal, logger)
	case gtypes.ProposalTypeTaxUsage:
		return executeTaxUsage(ctx, proposal, logger)
	}

	return nil
}

func executeParameterChange(ctx context.Context, proposal gtypes.Proposal, logger log.Logger) error {
	proposalContent := proposal.ProposalContent.(*gtypes.ParameterProposal)
	paramMapper := params.GetMapper(ctx)
	for _, param := range proposalContent.Params {
		paramSet, exists := paramMapper.GetModuleParamSet(param.Module)
		if !exists {
			panic(fmt.Sprintf("%s should exists", param.Module))
		}
		v, _ := paramSet.Validate(param.Key, param.Value)
		paramMapper.SetParam(param.Module, param.Key, v)
	}

	logger.Info("execute parameterChange, proposal: %d", proposal.ProposalID)

	return nil
}

func executeTaxUsage(ctx context.Context, proposal gtypes.Proposal, logger log.Logger) error {
	proposalContent := proposal.ProposalContent.(*gtypes.TaxUsageProposal)
	distributionMapper := ecomapper.GetDistributionMapper(ctx)
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
