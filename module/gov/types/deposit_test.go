package types

import (
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

func TestDeposit_Equals(t *testing.T) {
	addr1 := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
	addr2 := btypes.Address(ed25519.GenPrivKey().PubKey().Address())

	cases := []struct {
		input1   Deposit
		input2   Deposit
		expected bool
	}{
		{Deposit{addr1, 1, 1},
			Deposit{addr1, 1, 1},
			true},
		{Deposit{addr1, 1, 1},
			Deposit{addr2, 1, 1},
			false},
		{Deposit{addr1, 1, 1},
			Deposit{addr1, 2, 1},
			false},
		{Deposit{addr1, 1, 1},
			Deposit{addr1, 1, 2},
			false},
	}

	for tcIndex, tc := range cases {
		require.Equal(t, tc.input1.Equals(tc.input2), tc.expected, "tc #%d", tcIndex)
	}
}