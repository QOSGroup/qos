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

		mapper.DeleteProposal(proposalID)
		mapper.DeleteDeposits(ctx, proposalID) // delete any associated deposits (burned)

		resTags = resTags.AppendTags(btypes.NewTags(TagProposalID, types.Uint64ToBigEndian(proposalID), TagProposalResult, ProposalResultDropped))

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
			mapper.RefundDeposits(ctx, activeProposal.ProposalID, true)
			activeProposal.Status = gtypes.StatusPassed
			tagValue = ProposalResultPassed
			Execute(ctx, activeProposal, logger)
			break
		case gtypes.REJECT:
			mapper.RefundDeposits(ctx, activeProposal.ProposalID, true)
			activeProposal.Status = gtypes.StatusRejected
			tagValue = ProposalResultRejected
			break
		case gtypes.REJECTVETO:
			mapper.DeleteDeposits(ctx, activeProposal.ProposalID)
			activeProposal.Status = gtypes.StatusRejected
			tagValue = ProposalResultRejected
			break
		}

		activeProposal.FinalTallyResult = tallyResults
		mapper.SetProposal(activeProposal)
		mapper.RemoveFromActiveProposalQueue(activeProposal.VotingEndTime, activeProposal.ProposalID)

		logger.Info(
			fmt.Sprintf(
				"proposal %d (%s) tallied; result: %v",
				activeProposal.ProposalID, activeProposal.GetTitle(), proposalResult,
			),
		)

		resTags = resTags.AppendTags(btypes.NewTags(TagProposalID, types.Uint64ToBigEndian(proposalID), TagProposalResult, tagValue))

		penalty := mapper.GetParams(ctx).Penalty
		if penalty.GT(types.ZeroDec()) {
			validatorMapper := ecomapper.GetValidatorMapper(ctx)
			validators := validatorMapper.GetActiveValidatorSet(false)
			for _, val := range validators {
				if _, ok := votingValidators[val.String()]; !ok {
					if validator, exists := validatorMapper.GetValidator(val); exists && validator.BondHeight < activeProposal.VotingStartHeight {
						e := slash(ctx, validator, penalty, activeProposal.ProposalID)
						if e != nil {
							logger.Error("slash validator error", "e", e, "validator", validator.GetValidatorAddress().String())
						}
					}
				}
			}
		}

		mapper.DeleteValidatorSet(proposalID)
	}

	return resTags
}

// TODO slash
func slash(ctx context.Context, validator ecotypes.Validator, penalty types.Dec, proposalID uint64) error {

	log := ctx.Logger().With("module", "module/gov")

	log.Debug("slash validator", "proposalId", proposalID, "validator", validator.GetValidatorAddress().String())

	validatorMapper := ecomapper.GetValidatorMapper(ctx)
	delegationMapper := ecomapper.GetDelegationMapper(ctx)
	distributionMapper := ecomapper.GetDistributionMapper(ctx)
	var delegations []ecotypes.DelegationInfo
	delegationMapper.IterateDelegationsValDeleAddr(validator.GetValidatorAddress(), func(valAddr btypes.Address, delAddr btypes.Address) {
		if delegation, exists := delegationMapper.GetDelegationInfo(delAddr, valAddr); exists {
			delegations = append(delegations, delegation)
		}
	})

	totalSlashTokens := int64(0)
	height := uint64(ctx.BlockHeight())

	log.Debug("slash delegations", "delegations", delegations)

	for _, delegation := range delegations {
		bondTokens := int64(delegation.Amount)

		// 计算惩罚
		tokenSlashed := types.NewDec(bondTokens).Mul(penalty).TruncateInt64()

		if tokenSlashed >= bondTokens {
			tokenSlashed = bondTokens
		}

		// 修改delegator绑定收益
		remainTokens := uint64(bondTokens - tokenSlashed)
		if err := distributionMapper.ModifyDelegatorTokens(validator, delegation.DelegatorAddr, remainTokens, height); err != nil {
			return err
		}
		delegation.Amount = remainTokens

		// 更新delegation
		delegationMapper.SetDelegationInfo(delegation)

		log.Debug("slash validator's delegators", "delegator", delegation.DelegatorAddr.String(), "preToken", bondTokens, "slashToken", tokenSlashed, "remainTokens", remainTokens)

		totalSlashTokens += tokenSlashed
	}

	if uint64(totalSlashTokens) > validator.BondTokens {
		panic("slash token is greater then validator bondTokens")
	}

	// 更新validator
	updatedValidatorTokens := validator.BondTokens - uint64(totalSlashTokens)
	validatorMapper.ChangeValidatorBondTokens(validator, updatedValidatorTokens)

	log.Debug("slash validator bond tokens", "validator", validator.GetValidatorAddress().String(), "preTokens", validator.BondTokens, "slashTokens", totalSlashTokens, "afterTokens", updatedValidatorTokens)

	// slash放入社区费池
	distributionMapper.AddToCommunityFeePool(btypes.NewInt(totalSlashTokens))

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

	logger.Info("execute parameterChange", "proposal", proposal.ProposalID)

	return nil
}

func executeTaxUsage(ctx context.Context, proposal gtypes.Proposal, logger log.Logger) error {
	proposalContent := proposal.ProposalContent.(*gtypes.TaxUsageProposal)
	distributionMapper := ecomapper.GetDistributionMapper(ctx)
	accountMapper := ctx.Mapper(account.AccountMapperName).(*account.AccountMapper)
	acc := accountMapper.GetAccount(proposalContent.DestAddress)
	var account *types.QOSAccount
	if acc == nil {
		account = types.NewQOSAccountWithAddress(proposalContent.DestAddress)
	} else {
		account = acc.(*types.QOSAccount)
	}
	feePool := distributionMapper.GetCommunityFeePool()
	qos := types.NewDec(feePool.Int64()).Mul(proposalContent.Percent).TruncateInt()
	account.MustPlusQOS(qos)
	accountMapper.SetAccount(account)

	distributionMapper.SetCommunityFeePool(feePool.Sub(qos))

	logger.Info("execute taxUsage", "proposal", proposal.ProposalID)

	return nil
}
