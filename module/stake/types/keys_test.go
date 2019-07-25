package types

import (
	"testing"

	btypes "github.com/QOSGroup/qbase/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func TestStakeKey(t *testing.T) {

	key := ed25519.GenPrivKey()
	valAddr := btypes.Address(key.PubKey().Address())

	k := BuildValidatorVoteInfoKey(valAddr)
	addr := GetValidatorVoteInfoAddr(k)
	require.Equal(t, valAddr, addr)

	index := uint64(10086)
	k = BuildValidatorVoteInfoInWindowKey(index, valAddr)

	i, addr := GetValidatorVoteInfoInWindowIndexAddr(k)
	require.Equal(t, index, i)
	require.Equal(t, valAddr, addr)

}
