package mapper

import (
	"github.com/QOSGroup/qbase/context"
	bmapper "github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
)

func defaultCodec() *amino.Codec {
	cdc := amino.NewCodec()
	cryptoAmino.RegisterAmino(cdc)
	return cdc
}

func defaultContext(key store.StoreKey, mapperMap map[string]bmapper.IMapper) context.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, store.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := context.NewContext(cms, abci.Header{}, false, log.NewNopLogger(), mapperMap)
	return ctx
}

func TestSaveCAPubKey(t *testing.T) {
	baseMapper := NewMainMapper()
	baseMapper.SetCodec(defaultCodec())
	storeKey := baseMapper.GetStoreKey()
	mapper := make(map[string]bmapper.IMapper)
	mapper[baseMapper.Name()] = baseMapper
	ctx := defaultContext(storeKey, mapper)
	baseMapper, _ = ctx.Mapper(baseMapper.Name()).(*MainMapper)

	origin := ed25519.GenPrivKey().PubKey()
	baseMapper.SetRootCA(origin)
	recover := baseMapper.GetRoot()
	require.Equal(t, origin, recover)
}
