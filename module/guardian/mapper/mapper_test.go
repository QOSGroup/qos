package mapper

import (
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	bmapper "github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/guardian/types"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
)

func defaultContext() context.Context {
	mapperMap := make(map[string]bmapper.IMapper)

	cdc := baseabci.MakeQBaseCodec()
	qtypes.RegisterCodec(cdc)

	guardianMapper := NewMapper()
	guardianMapper.SetCodec(cdc)
	guardianKey := guardianMapper.GetStoreKey()
	mapperMap[MapperName] = guardianMapper

	accountMapper := bacc.NewAccountMapper(nil, qtypes.ProtoQOSAccount)
	accountMapper.SetCodec(cdc)
	acountKey := accountMapper.GetStoreKey()
	mapperMap[bacc.AccountMapperName] = accountMapper

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(guardianKey, btypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(acountKey, btypes.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := context.NewContext(cms, abci.Header{}, false, log.NewNopLogger(), mapperMap)
	return ctx
}

func TestGuardianMapper_GetGuardian(t *testing.T) {
	ctx := defaultContext()
	mapper := GetMapper(ctx)

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
	mapper := GetMapper(ctx)

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
	mapper := GetMapper(ctx)

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
