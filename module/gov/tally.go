package gov

import (
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	ecomapper "github.com/QOSGroup/qos/module/eco/mapper"
	ecotypes "github.com/QOSGroup/qos/module/eco/types"
	gtypes "github.com/QOSGroup/qos/module/gov/types"
	"github.com/QOSGroup/qos/types"
)

func tally(ctx context.Context, mapper *GovMapper, proposal gtypes.Proposal) (passes gtypes.ProposalResult, tallyResults gtypes.TallyResult, validators map[string]bool) {
	results := make(map[gtypes.VoteOption]int64)
	results[gtypes.OptionYes] = 0
	results[gtypes.OptionAbstain] = 0
	results[gtypes.OptionNo] = 0
	results[gtypes.OptionNoWithVeto] = 0

	totalVotingPower := btypes.ZeroInt()
	totalSystemPower := btypes.ZeroInt()
	currValidators := make(map[string]ecotypes.Validator)

	// 统计当前验证节点信息
	validatorMapper := ctx.Mapper(ecotypes.ValidatorMapperName).(*ecomapper.ValidatorMapper)
	delegatorMapper := ctx.Mapper(ecotypes.DelegationMapperName).(*ecomapper.DelegationMapper)
	iterator := validatorMapper.IteratorValidatorByVoterPower(false)
	defer iterator.Close()
	var key []byte
	for ; iterator.Valid(); iterator.Next() {
		key = iterator.Key()
		valAddr := btypes.Address(key[9:])
		if validator, exists := validatorMapper.GetValidator(valAddr); exists {
			currValidators[validator.GetValidatorAddress().String()] = validator
			totalSystemPower = totalSystemPower.Add(btypes.NewInt(int64(validator.BondTokens)))
		}
	}

	// 统计投票人对应验证节点信息，一个验证节点任意delegator投过票，则算这个验证节点参与了投票
	votesIterator := mapper.GetVotes(ctx, proposal.ProposalID)
	defer votesIterator.Close()
	for ; votesIterator.Valid(); votesIterator.Next() {
		vote := &gtypes.Vote{}
		mapper.DecodeObject(votesIterator.Value(), vote)

		delegatorMapper.IterateDelegationsInfo(vote.Voter, func(delegation ecotypes.DelegationInfo) {
			if _, ok := currValidators[delegation.ValidatorAddr.String()]; ok {
				totalVotingPower = totalVotingPower.Add(btypes.NewInt(int64(delegation.Amount)))
				results[vote.Option] += int64(delegation.Amount)
				validators[delegation.ValidatorAddr.String()] = true
			}
		})
	}

	params := mapper.GetParams()
	tallyResults = gtypes.TallyResult{
		Yes:        results[gtypes.OptionYes],
		Abstain:    results[gtypes.OptionAbstain],
		No:         results[gtypes.OptionNo],
		NoWithVeto: results[gtypes.OptionNoWithVeto],
	}

	// If no one votes, proposal fails
	if totalVotingPower.Int64() == results[gtypes.OptionAbstain] {
		return gtypes.REJECT, tallyResults, validators
	}

	//if more than 1/3 of voters abstain, proposal fails
	if totalVotingPower.MulRaw(3).LT(totalSystemPower) {
		return gtypes.REJECT, tallyResults, validators
	}

	// If more than 1/3 of voters veto, proposal fails
	if types.NewDec(int64(results[gtypes.OptionNoWithVeto])).Quo(types.NewDecFromInt(totalVotingPower)).GT(params.Veto) {
		return gtypes.REJECTVETO, tallyResults, validators
	}

	// If more than 1/2 of non-abstaining voters vote Yes, proposal passes
	if types.NewDec(int64(results[gtypes.OptionYes])).Quo(types.NewDecFromInt(totalVotingPower.Sub(btypes.NewInt(results[gtypes.OptionAbstain])))).GT(params.Threshold) {
		return gtypes.PASS, tallyResults, validators
	}

	// If more than 1/2 of non-abstaining voters vote No, proposal fails

	return gtypes.REJECT, tallyResults, validators
}
