package types

import (
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

func TestGuardian_Equals(t *testing.T) {
	addr1 := btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	addr2 := btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	cases := []struct {
		input1   Guardian
		input2   Guardian
		expected bool
	}{
		{Guardian{"g1", Genesis, addr1, nil},
			Guardian{"g1", Genesis, addr1, nil},
			true},
		{Guardian{"g1", Genesis, addr1, nil},
			Guardian{"g2", Genesis, addr1, nil},
			false},
		{Guardian{"g1", Genesis, addr1, nil},
			Guardian{"g1", Ordinary, addr1, nil},
			false},
		{Guardian{"g1", Genesis, addr1, nil},
			Guardian{"g1", Genesis, addr1, addr2},
			false},
		{Guardian{"g1", Genesis, addr1, nil},
			Guardian{"g1", Genesis, addr2, nil},
			false},
	}

	for tcIndex, tc := range cases {
		require.Equal(t, tc.input1.Equals(tc.input2), tc.expected, "tc #%d", tcIndex)
	}
}
