package txs

import (
	"github.com/QOSGroup/qbase/baseabci"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/gov/mapper"
	"github.com/QOSGroup/qos/module/gov/types"
	"github.com/QOSGroup/qos/module/params"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

func TestTxVote_ValidateData(t *testing.T) {
	ctx := defaultContext()
	initGenesis(ctx, types.DefaultGenesisState())
	accountMapper := baseabci.GetAccountMapper(ctx)
	addr := btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	accountMapper.SetAccount(qtypes.NewQOSAccount(addr, btypes.NewInt(20000000000), nil))
	params.GetMapper(ctx).RegisterParamSet(&types.Params{})
	proposal := NewTxProposal("p1", "p1", addr, btypes.NewInt(9000000000))
	proposal.Exec(ctx)

	cases := []struct {
		input *TxVote
		valid bool
	}{
		{NewTxVote(1, addr, types.OptionYes), true},
		{NewTxVote(1, addr, types.VoteOption(0xff)), false},
		{NewTxVote(1, nil, types.OptionYes), false},
	}

	for tcIndex, tc := range cases {
		err := tc.input.ValidateData(ctx)
		require.Equal(t, tc.valid, err == nil, "tc #%d", tcIndex)
		if tcIndex == 0 {
			tx := NewTxDeposit(1, addr, btypes.NewInt(1000000000))
			tx.Exec(ctx)
		}
	}
}

func TestTxVote_Exec(t *testing.T) {
	ctx := defaultContext()
	initGenesis(ctx, types.DefaultGenesisState())
	accountMapper := baseabci.GetAccountMapper(ctx)
	addr := btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	accountMapper.SetAccount(qtypes.NewQOSAccount(addr, btypes.NewInt(20), nil))
	params.GetMapper(ctx).RegisterParamSet(&types.Params{})
	proposal := NewTxProposal("p1", "p1", addr, btypes.NewInt(10))
	proposal.Exec(ctx)

	tx := NewTxVote(1, addr, types.OptionYes)
	tx.Exec(ctx)

	vote, exists := mapper.GetMapper(ctx).GetVote(1, addr)
	require.True(t, exists)
	require.NotNil(t, vote)
	require.Equal(t, types.OptionYes, vote.Option)

}
