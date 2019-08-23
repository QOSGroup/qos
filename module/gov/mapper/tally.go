package mapper

import (
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	gtypes "github.com/QOSGroup/qos/module/gov/types"
	"github.com/QOSGroup/qos/module/stake"
	"github.com/QOSGroup/qos/types"
)

func Tally(ctx context.Context, mapper *Mapper, proposal gtypes.Proposal) (passes gtypes.ProposalResult, tallyResults gtypes.TallyResult, validators map[string]bool, deductOption gtypes.DeductOption) {
	results := make(map[gtypes.VoteOption]int64)
	results[gtypes.OptionYes] = 0
	results[gtypes.OptionAbstain] = 0
	results[gtypes.OptionNo] = 0
	results[gtypes.OptionNoWithVeto] = 0

	totalVotingPower := btypes.ZeroInt()
	totalSystemPower := btypes.ZeroInt()
	validators = make(map[string]bool)
	currValidators := make(map[string]stake.Validator)

	// 统计当前验证节点信息
	sm := stake.GetMapper(ctx)
	iterator := sm.IteratorValidatorByVoterPower(false)
	defer iterator.Close()
	var key []byte
	for ; iterator.Valid(); iterator.Next() {
		key = iterator.Key()
		valAddr := btypes.Address(key[9:])
		if validator, exists := sm.GetValidator(valAddr); exists {
			currValidators[validator.GetValidatorAddress().String()] = validator
			totalSystemPower = totalSystemPower.Add(btypes.NewInt(int64(validator.BondTokens)))
		}
	}

	// 统计投票人对应验证节点信息，一个验证节点任意delegator投过票，则算这个验证节点参与了投票
	votesIterator := mapper.GetVotes(proposal.ProposalID)
	defer votesIterator.Close()
	for ; votesIterator.Valid(); votesIterator.Next() {
		vote := &gtypes.Vote{}
		mapper.DecodeObject(votesIterator.Value(), vote)

		sm.IterateDelegationsInfo(vote.Voter, func(delegation stake.Delegation) {
			if _, ok := currValidators[delegation.ValidatorAddr.String()]; ok {
				totalVotingPower = totalVotingPower.Add(btypes.NewInt(int64(delegation.Amount)))
				results[vote.Option] += int64(delegation.Amount)
				validators[delegation.ValidatorAddr.String()] = true
			}
		})
	}

	params := mapper.GetParams(ctx)
	tallyResults = gtypes.TallyResult{
		Yes:        results[gtypes.OptionYes],
		Abstain:    results[gtypes.OptionAbstain],
		No:         results[gtypes.OptionNo],
		NoWithVeto: results[gtypes.OptionNoWithVeto],
	}

	// If no one votes, proposal fails
	if totalVotingPower.Int64() == results[gtypes.OptionAbstain] {
		return gtypes.REJECT, tallyResults, validators, gtypes.DepositDeductNone
	}

	//if more than 1/3 of voters abstain, proposal fails
	if types.NewDecFromInt(totalVotingPower.Div(totalSystemPower)).LT(params.Quorum) {
		return gtypes.REJECT, tallyResults, validators, gtypes.DepositDeductPart
	}

	// If more than 1/3 of voters veto, proposal fails
	if types.NewDec(results[gtypes.OptionNoWithVeto]).Quo(types.NewDecFromInt(totalVotingPower)).GT(params.Veto) {
		return gtypes.REJECTVETO, tallyResults, validators, gtypes.DepositDeductAll
	}

	// If more than 1/2 of non-abstaining voters vote Yes, proposal passes
	if types.NewDec(results[gtypes.OptionYes]).Quo(types.NewDecFromInt(totalVotingPower.Sub(btypes.NewInt(results[gtypes.OptionAbstain])))).GT(params.Threshold) {
		return gtypes.PASS, tallyResults, validators, gtypes.DepositDeductNone
	}

	// If more than 1/2 of non-abstaining voters vote No, proposal fails

	return gtypes.REJECT, tallyResults, validators, gtypes.DepositDeductNone
}
