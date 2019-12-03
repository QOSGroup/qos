package mapper

import (
	"github.com/QOSGroup/qbase/baseabci"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/gov/types"
	"github.com/QOSGroup/qos/module/params"
	"github.com/QOSGroup/qos/module/stake"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

func TestTally(t *testing.T) {
	ctx := defaultContext()
	params.GetMapper(ctx).RegisterParamSet(&types.Params{})
	initGenesis(ctx, types.DefaultGenesisState())
	govMapper := GetMapper(ctx)
	sm := stake.GetMapper(ctx)
	am := baseabci.GetAccountMapper(ctx)

	pubKeys := []crypto.PubKey{
		ed25519.GenPrivKey().PubKey(),
		ed25519.GenPrivKey().PubKey(),
		ed25519.GenPrivKey().PubKey(),
	}
	var addrs []btypes.AccAddress
	for i, pub := range pubKeys {
		addr := btypes.AccAddress(pub.Address())
		addrs = append(addrs, addr)
		am.SetAccount(qtypes.NewQOSAccount(addr, btypes.NewInt(1000), nil))

		validator := stake.Validator{
			Description:     stake.Description{Moniker: string(i)},
			Owner:           addr,
			OperatorAddress: btypes.ValAddress(addr),
			ConsPubKey:      pub,
			BondTokens:      btypes.NewInt(1000),
			Status:          stake.Active,
			MinPeriod:       int64(0),
			BondHeight:      ctx.BlockHeight(),
		}

		delegationInfo := stake.NewDelegationInfo(validator.Owner, validator.GetValidatorAddress(), validator.BondTokens, false)
		sm.SetDelegationInfo(delegationInfo)
		sm.CreateValidator(validator)
	}

	// no one votes, proposal fails
	textContent := types.NewTextProposal("p1", "p1", btypes.NewInt(10))
	proposal, err := govMapper.SubmitProposal(ctx, textContent)
	require.Nil(t, err)
	err, _ = govMapper.AddDeposit(ctx, proposal.ProposalID, addrs[0], textContent.Deposit)
	require.Nil(t, err)
	id := int64(1)
	proposal, exists := govMapper.GetProposal(id)
	require.True(t, exists)
	passes, tallyResult, validatorSet, _ := Tally(ctx, govMapper, proposal)
	require.Equal(t, passes, types.REJECT)
	require.Equal(t, btypes.ZeroInt(), tallyResult.Abstain)
	require.Equal(t, 0, len(validatorSet))

	// more than 1/3 of voters abstain, proposal fails
	textContent = types.NewTextProposal("p2", "p2", btypes.NewInt(10))
	proposal, err = govMapper.SubmitProposal(ctx, textContent)
	require.Nil(t, err)
	err, _ = govMapper.AddDeposit(ctx, proposal.ProposalID, addrs[0], textContent.Deposit)
	require.Nil(t, err)
	id = 2
	govMapper.AddVote(id, addrs[0], types.OptionAbstain)
	govMapper.AddVote(id, addrs[1], types.OptionAbstain)
	proposal, _ = govMapper.GetProposal(id)
	passes, tallyResult, validatorSet, _ = Tally(ctx, govMapper, proposal)
	require.Equal(t, passes, types.REJECT)
	require.Equal(t, btypes.NewInt(2000), tallyResult.Abstain)
	require.Equal(t, 2, len(validatorSet))
	require.NotNil(t, validatorSet[string(addrs[0])])
	require.NotNil(t, validatorSet[string(addrs[1])])

	// more than 1/3 of voters veto, proposal fails
	textContent = types.NewTextProposal("p3", "p3", btypes.NewInt(10))
	proposal, err = govMapper.SubmitProposal(ctx, textContent)
	require.Nil(t, err)
	err, _ = govMapper.AddDeposit(ctx, proposal.ProposalID, addrs[0], textContent.Deposit)
	require.Nil(t, err)
	id = 3
	govMapper.AddVote(id, addrs[0], types.OptionNoWithVeto)
	govMapper.AddVote(id, addrs[1], types.OptionNoWithVeto)
	govMapper.AddVote(id, addrs[2], types.OptionYes)
	proposal, _ = govMapper.GetProposal(id)
	passes, tallyResult, validatorSet, _ = Tally(ctx, govMapper, proposal)
	require.Equal(t, passes, types.REJECTVETO)
	require.Equal(t, btypes.NewInt(1000), tallyResult.Yes)
	require.Equal(t, btypes.NewInt(2000), tallyResult.NoWithVeto)
	require.Equal(t, 3, len(validatorSet))
	require.NotNil(t, validatorSet[string(addrs[0])])
	require.NotNil(t, validatorSet[string(addrs[1])])
	require.NotNil(t, validatorSet[string(addrs[2])])

	// more than 1/2 of non-abstaining voters vote Yes, proposal passes
	textContent = types.NewTextProposal("p4", "p4", btypes.NewInt(10))
	proposal, err = govMapper.SubmitProposal(ctx, textContent)
	require.Nil(t, err)
	err, _ = govMapper.AddDeposit(ctx, proposal.ProposalID, addrs[0], textContent.Deposit)
	require.Nil(t, err)
	id = 4
	govMapper.AddVote(id, addrs[0], types.OptionYes)
	govMapper.AddVote(id, addrs[1], types.OptionNo)
	govMapper.AddVote(id, addrs[2], types.OptionYes)
	proposal, _ = govMapper.GetProposal(id)
	passes, tallyResult, validatorSet, _ = Tally(ctx, govMapper, proposal)
	require.Equal(t, passes, types.PASS)
	require.Equal(t, btypes.NewInt(2000), tallyResult.Yes)
	require.Equal(t, btypes.NewInt(1000), tallyResult.No)
	require.Equal(t, 3, len(validatorSet))
	require.NotNil(t, validatorSet[string(addrs[0])])
	require.NotNil(t, validatorSet[string(addrs[1])])
	require.NotNil(t, validatorSet[string(addrs[2])])

	// more than 1/2 of non-abstaining voters vote No, proposal fails
	textContent = types.NewTextProposal("p5", "p5o", btypes.NewInt(10))
	proposal, err = govMapper.SubmitProposal(ctx, textContent)
	require.Nil(t, err)
	err, _ = govMapper.AddDeposit(ctx, proposal.ProposalID, addrs[0], textContent.Deposit)
	require.Nil(t, err)
	id = 5
	govMapper.AddVote(id, addrs[0], types.OptionNo)
	govMapper.AddVote(id, addrs[1], types.OptionNo)
	govMapper.AddVote(id, addrs[2], types.OptionYes)
	proposal, _ = govMapper.GetProposal(id)
	passes, tallyResult, validatorSet, _ = Tally(ctx, govMapper, proposal)
	require.Equal(t, passes, types.REJECT)
	require.Equal(t, btypes.NewInt(2000), tallyResult.No)
	require.Equal(t, btypes.NewInt(1000), tallyResult.Yes)
	require.Equal(t, 3, len(validatorSet))
	require.NotNil(t, validatorSet[string(addrs[0])])
	require.NotNil(t, validatorSet[string(addrs[1])])
	require.NotNil(t, validatorSet[string(addrs[2])])
}
