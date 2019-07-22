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

func TestTxDeposit_ValidateData(t *testing.T) {
	ctx := defaultContext()
	initGenesis(ctx, types.DefaultGenesisState())
	accountMapper := baseabci.GetAccountMapper(ctx)
	addr := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
	accountMapper.SetAccount(qtypes.NewQOSAccount(addr, btypes.NewInt(20), nil))
	params.GetMapper(ctx).RegisterParamSet(&types.Params{})
	proposal := NewTxProposal("p1", "p1", addr, 10)
	proposal.Exec(ctx)

	cases := []struct {
		input *TxDeposit
		valid bool
	}{
		{NewTxDeposit(1, addr, 10), true},
		{NewTxDeposit(2, addr, 10), false},
		{NewTxDeposit(1, nil, 10), false},
		{NewTxDeposit(1, addr, 0), false},
	}

	for tcIndex, tc := range cases {
		err := tc.input.ValidateData(ctx)
		require.Equal(t, tc.valid, err == nil, "tc #%d", tcIndex)
	}
}

func TestTxDeposit_Exec(t *testing.T) {
	ctx := defaultContext()
	initGenesis(ctx, types.DefaultGenesisState())
	accountMapper := baseabci.GetAccountMapper(ctx)
	addr := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
	accountMapper.SetAccount(qtypes.NewQOSAccount(addr, btypes.NewInt(20), nil))
	params.GetMapper(ctx).RegisterParamSet(&types.Params{})
	proposal := NewTxProposal("p1", "p1", addr, 10)
	proposal.Exec(ctx)

	govMapper := mapper.GetMapper(ctx)
	deposit, exists := govMapper.GetDeposit(1, addr)
	require.True(t, true, exists)
	require.NotNil(t, deposit)
	require.Equal(t, uint64(10), deposit.Amount)

	tx := NewTxDeposit(1, addr, 10)
	tx.Exec(ctx)

	deposit, exists = govMapper.GetDeposit(1, addr)
	require.True(t, true, exists)
	require.NotNil(t, deposit)
	require.Equal(t, uint64(20), deposit.Amount)
}
