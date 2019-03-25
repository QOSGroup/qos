package gov

import (
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	ecomapper "github.com/QOSGroup/qos/module/eco/mapper"
	ecotypes "github.com/QOSGroup/qos/module/eco/types"
	gtypes "github.com/QOSGroup/qos/module/gov/types"
	"github.com/QOSGroup/qos/types"
)

func tally(ctx context.Context, mapper *GovMapper, proposal gtypes.Proposal) (passes bool, tallyResults gtypes.TallyResult) {
	results := make(map[gtypes.VoteOption]uint64)
	results[gtypes.OptionYes] = 0
	results[gtypes.OptionAbstain] = 0
	results[gtypes.OptionNo] = 0
	results[gtypes.OptionNoWithVeto] = 0

	totalVotingPower := uint64(0)
	totalSystemPower := uint64(0)
	currValidators := make(map[string]ecotypes.Validator)

	// fetch all the bonded validators, insert them into currValidators
	validatorMapper := ctx.Mapper(ecotypes.ValidatorMapperName).(*ecomapper.ValidatorMapper)
	delegatorMapper := ctx.Mapper(ecotypes.DelegationMapperName).(*ecomapper.DelegationMapper)
	iterator := validatorMapper.IteratorValidatrorByVoterPower(false)
	defer iterator.Close()
	var key []byte
	for ; iterator.Valid(); iterator.Next() {
		key = iterator.Key()
		valAddr := btypes.Address(key[9:])
		if validator, exists := validatorMapper.GetValidator(valAddr); exists {
			currValidators[validator.GetValidatorAddress().String()] = validator
			totalSystemPower = totalSystemPower + validator.BondTokens
		}
	}

	// iterate over all the votes
	votesIterator := mapper.GetVotes(ctx, proposal.ProposalID)
	defer votesIterator.Close()
	for ; votesIterator.Valid(); votesIterator.Next() {
		vote := &gtypes.Vote{}
		mapper.GetCodec().MustUnmarshalBinaryLengthPrefixed(votesIterator.Value(), vote)

		valAddrStr := btypes.Address(vote.Voter).String()
		// iterate over all delegations from voter, deduct from any delegated-to validators
		delegatorMapper.IterateDelegationsInfo(vote.Voter, func(delegation ecotypes.DelegationInfo) {
			if _, ok := currValidators[valAddrStr]; ok {
				totalVotingPower += delegation.Amount
				results[vote.Option] += delegation.Amount
			}
		})

		mapper.deleteVote(ctx, vote.ProposalID, vote.Voter)
	}

	tallyParams := mapper.GetTallyParams()
	tallyResults = gtypes.TallyResult{
		Yes:        results[gtypes.OptionYes],
		Abstain:    results[gtypes.OptionAbstain],
		No:         results[gtypes.OptionNo],
		NoWithVeto: results[gtypes.OptionNoWithVeto],
	}

	// If no one votes, proposal fails
	if totalVotingPower == results[gtypes.OptionAbstain] {
		return false, tallyResults
	}

	//if more than 1/3 of voters abstain, proposal fails
	if 3*totalVotingPower < totalSystemPower {
		return false, tallyResults
	}

	// If more than 1/3 of voters veto, proposal fails
	if types.NewDec(int64(results[gtypes.OptionNoWithVeto])).Quo(types.NewDec(int64(totalVotingPower))).GT(tallyParams.Veto) {
		return false, tallyResults
	}

	// If more than 1/2 of non-abstaining voters vote Yes, proposal passes
	if types.NewDec(int64(results[gtypes.OptionYes])).Quo(types.NewDec(int64(totalVotingPower - results[gtypes.OptionAbstain]))).GT(tallyParams.Threshold) {
		return true, tallyResults
	}

	// If more than 1/2 of non-abstaining voters vote No, proposal fails

	return false, tallyResults
}
