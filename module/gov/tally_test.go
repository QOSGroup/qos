package gov

import (
	"github.com/QOSGroup/qbase/baseabci"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/eco/mapper"
	ecotypes "github.com/QOSGroup/qos/module/eco/types"
	gtypes "github.com/QOSGroup/qos/module/gov/types"
	"github.com/QOSGroup/qos/module/params"
	"github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

func TestTally(t *testing.T) {
	ctx := defaultContext()
	params.GetMapper(ctx).RegisterParamSet(&Params{})
	InitGenesis(ctx, DefaultGenesisState())
	govMapper := GetGovMapper(ctx)
	validatorMapper := mapper.GetValidatorMapper(ctx)
	delegationMapper := mapper.GetDelegationMapper(ctx)
	accountMapper := baseabci.GetAccountMapper(ctx)

	pubKeys := []crypto.PubKey{
		ed25519.GenPrivKey().PubKey(),
		ed25519.GenPrivKey().PubKey(),
		ed25519.GenPrivKey().PubKey(),
	}
	var addrs []btypes.Address
	for i, pub := range pubKeys {
		addr := btypes.Address(pub.Address())
		addrs = append(addrs, addr)
		accountMapper.SetAccount(types.NewQOSAccount(addr, btypes.NewInt(1000), nil))

		validator := ecotypes.Validator{
			Name:            string(i),
			Owner:           addr,
			ValidatorPubKey: pub,
			BondTokens:      1000,
			Description:     string(i),
			Status:          ecotypes.Active,
			MinPeriod:       uint64(0),
			BondHeight:      uint64(ctx.BlockHeight()),
		}

		delegationInfo := ecotypes.NewDelegationInfo(validator.Owner, validator.GetValidatorAddress(), validator.BondTokens, false)
		delegationMapper.SetDelegationInfo(delegationInfo)
		validatorMapper.CreateValidator(validator)
	}

	// no one votes, proposal fails
	proposalTx := NewTxProposal("p1", "p1", addrs[0], 10)
	result, _ := proposalTx.Exec(ctx)
	require.Equal(t, result.Code, btypes.CodeOK)
	id := uint64(1)
	proposal, _ := govMapper.GetProposal(id)
	passes, tallyResult, validatorSet := tally(ctx, govMapper, proposal)
	require.Equal(t, passes, gtypes.REJECT)
	require.Equal(t, int64(0), tallyResult.Abstain)
	require.Equal(t, 0, len(validatorSet))

	// more than 1/3 of voters abstain, proposal fails
	proposalTx = NewTxProposal("p2", "p2", addrs[0], 10)
	result, _ = proposalTx.Exec(ctx)
	require.Equal(t, result.Code, btypes.CodeOK)
	id = 2
	govMapper.AddVote(id, addrs[0], gtypes.OptionAbstain)
	govMapper.AddVote(id, addrs[1], gtypes.OptionAbstain)
	proposal, _ = govMapper.GetProposal(id)
	passes, tallyResult, validatorSet = tally(ctx, govMapper, proposal)
	require.Equal(t, passes, gtypes.REJECT)
	require.Equal(t, int64(2000), tallyResult.Abstain)
	require.Equal(t, 2, len(validatorSet))
	require.NotNil(t, validatorSet[string(addrs[0])])
	require.NotNil(t, validatorSet[string(addrs[1])])

	// more than 1/3 of voters veto, proposal fails
	proposalTx = NewTxProposal("p3", "p3", addrs[0], 10)
	result, _ = proposalTx.Exec(ctx)
	require.Equal(t, result.Code, btypes.CodeOK)
	id = 3
	govMapper.AddVote(id, addrs[0], gtypes.OptionNoWithVeto)
	govMapper.AddVote(id, addrs[1], gtypes.OptionNoWithVeto)
	govMapper.AddVote(id, addrs[2], gtypes.OptionYes)
	proposal, _ = govMapper.GetProposal(id)
	passes, tallyResult, validatorSet = tally(ctx, govMapper, proposal)
	require.Equal(t, passes, gtypes.REJECTVETO)
	require.Equal(t, int64(1000), tallyResult.Yes)
	require.Equal(t, int64(2000), tallyResult.NoWithVeto)
	require.Equal(t, 3, len(validatorSet))
	require.NotNil(t, validatorSet[string(addrs[0])])
	require.NotNil(t, validatorSet[string(addrs[1])])
	require.NotNil(t, validatorSet[string(addrs[2])])

	// more than 1/2 of non-abstaining voters vote Yes, proposal passes
	proposalTx = NewTxProposal("p4", "p4", addrs[0], 10)
	result, _ = proposalTx.Exec(ctx)
	require.Equal(t, result.Code, btypes.CodeOK)
	id = 4
	govMapper.AddVote(id, addrs[0], gtypes.OptionYes)
	govMapper.AddVote(id, addrs[1], gtypes.OptionNo)
	govMapper.AddVote(id, addrs[2], gtypes.OptionYes)
	proposal, _ = govMapper.GetProposal(id)
	passes, tallyResult, validatorSet = tally(ctx, govMapper, proposal)
	require.Equal(t, passes, gtypes.PASS)
	require.Equal(t, int64(2000), tallyResult.Yes)
	require.Equal(t, int64(1000), tallyResult.No)
	require.Equal(t, 3, len(validatorSet))
	require.NotNil(t, validatorSet[string(addrs[0])])
	require.NotNil(t, validatorSet[string(addrs[1])])
	require.NotNil(t, validatorSet[string(addrs[2])])

	// more than 1/2 of non-abstaining voters vote No, proposal fails
	proposalTx = NewTxProposal("p5", "p5", addrs[0], 10)
	result, _ = proposalTx.Exec(ctx)
	require.Equal(t, result.Code, btypes.CodeOK)
	id = 5
	govMapper.AddVote(id, addrs[0], gtypes.OptionNo)
	govMapper.AddVote(id, addrs[1], gtypes.OptionNo)
	govMapper.AddVote(id, addrs[2], gtypes.OptionYes)
	proposal, _ = govMapper.GetProposal(id)
	passes, tallyResult, validatorSet = tally(ctx, govMapper, proposal)
	require.Equal(t, passes, gtypes.REJECT)
	require.Equal(t, int64(2000), tallyResult.No)
	require.Equal(t, int64(1000), tallyResult.Yes)
	require.Equal(t, 3, len(validatorSet))
	require.NotNil(t, validatorSet[string(addrs[0])])
	require.NotNil(t, validatorSet[string(addrs[1])])
	require.NotNil(t, validatorSet[string(addrs[2])])
}
