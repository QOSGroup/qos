package guardian

import (
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/guardian/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

func TestGuardianMapper_GetGuardian(t *testing.T) {
	ctx := defaultContext()
	mapper := GetGuardianMapper(ctx)

	addr1 := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
	g2, exists := mapper.GetGuardian(addr1)

	require.False(t, exists)

	g1 := types.NewGuardian("g1", types.Genesis, addr1, nil)
	mapper.AddGuardian(*g1)
	g2, exists = mapper.GetGuardian(addr1)

	require.True(t, exists)
	require.True(t, g1.Equals(g2))
}

func TestGuardianMapper_DeleteGuardian(t *testing.T) {
	ctx := defaultContext()
	mapper := GetGuardianMapper(ctx)

	addr1 := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
	g1 := types.NewGuardian("g1", types.Genesis, addr1, nil)
	mapper.AddGuardian(*g1)

	g2, exists := mapper.GetGuardian(addr1)
	require.True(t, exists)
	require.True(t, g1.Equals(g2))

	mapper.DeleteGuardian(addr1)
	_, exists = mapper.GetGuardian(addr1)
	require.False(t, exists)
}

func TestGuardianMapper_GuardiansIterator(t *testing.T) {
	ctx := defaultContext()
	mapper := GetGuardianMapper(ctx)

	addr1 := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
	g1 := types.NewGuardian("g1", types.Genesis, addr1, nil)
	mapper.AddGuardian(*g1)

	g2, exists := mapper.GetGuardian(addr1)
	require.True(t, exists)
	require.True(t, g1.Equals(g2))

	iterator := mapper.GuardiansIterator()
	defer iterator.Close()
	cnt := 0
	for ; iterator.Valid(); iterator.Next() {
		var g2 types.Guardian
		mapper.GetCodec().MustUnmarshalBinaryBare(iterator.Value(), &g2)
		require.True(t, g1.Equals(g2))
		cnt++
	}

	require.Equal(t, cnt, 1)
}
