package types

import (
	"testing"

	btypes "github.com/QOSGroup/qbase/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func TestStakeKey(t *testing.T) {

	key := ed25519.GenPrivKey()
	valAddr := btypes.ValAddress(key.PubKey().Address())

	k := BuildValidatorVoteInfoKey(valAddr)
	addr := GetValidatorVoteInfoAddr(k)
	require.Equal(t, valAddr, addr)

	index := int64(10086)
	k = BuildValidatorVoteInfoInWindowKey(index, valAddr)

	i, addr := GetValidatorVoteInfoInWindowIndexAddr(k)
	require.Equal(t, index, i)
	require.Equal(t, valAddr, addr)

	power := int64(1228)
	bz := BuildValidatorByVotePower(power, valAddr)

	vp, va, e := ParseValidatorVotePowerKey(bz)

	require.Nil(t, e)
	require.Equal(t, int64(1228), vp)
	require.Equal(t, valAddr, va)
}
