package gov

import (
	"github.com/QOSGroup/qbase/baseabci"
	btypes "github.com/QOSGroup/qbase/types"
	gtypes "github.com/QOSGroup/qos/module/gov/types"
	"github.com/QOSGroup/qos/module/params"
	"github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

func TestTxVote_ValidateData(t *testing.T) {
	ctx := defaultContext()
	InitGenesis(ctx, DefaultGenesisState())
	accountMapper := baseabci.GetAccountMapper(ctx)
	addr := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
	accountMapper.SetAccount(types.NewQOSAccount(addr, btypes.NewInt(20), nil))
	params.GetMapper(ctx).RegisterParamSet(&Params{})
	proposal := NewTxProposal("p1", "p1", addr, 1)
	proposal.Exec(ctx)

	cases := []struct {
		input *TxVote
		valid bool
	}{
		{NewTxVote(1, addr, gtypes.OptionYes), false},
		{NewTxVote(1, addr, gtypes.OptionYes), true},
		{NewTxVote(1, addr, gtypes.VoteOption(0xff)), false},
		{NewTxVote(1, nil, gtypes.OptionYes), false},
	}

	for tcIndex, tc := range cases {
		err := tc.input.ValidateData(ctx)
		require.Equal(t, tc.valid, err == nil, "tc #%d", tcIndex)
		if tcIndex == 0 {
			tx := NewTxDeposit(1, addr, 10)
			tx.Exec(ctx)
		}
	}
}

func TestTxVote_Exec(t *testing.T) {
	ctx := defaultContext()
	InitGenesis(ctx, DefaultGenesisState())
	accountMapper := baseabci.GetAccountMapper(ctx)
	addr := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
	accountMapper.SetAccount(types.NewQOSAccount(addr, btypes.NewInt(20), nil))
	params.GetMapper(ctx).RegisterParamSet(&Params{})
	proposal := NewTxProposal("p1", "p1", addr, 10)
	proposal.Exec(ctx)

	tx := NewTxVote(1, addr, gtypes.OptionYes)
	tx.Exec(ctx)

	vote, exists := GetGovMapper(ctx).GetVote(1, addr)
	require.True(t, exists)
	require.NotNil(t, vote)
	require.Equal(t, gtypes.OptionYes, vote.Option)

}
