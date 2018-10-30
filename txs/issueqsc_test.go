package txs

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/store"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/mapper"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	bmapper "github.com/QOSGroup/qbase/mapper"
	bacc "github.com/QOSGroup/qbase/account"
	"testing"
)

func TestNewTxIssueQsc(t *testing.T) {
	var bankeraddr = []byte("address1zcjduepqdrq67yachuqs7w87syaq62x999r5dddtwehtamhymnpmxwew9j0srz0a36")
	txIssue := NewTxIssueQsc("qsc10", types.NewInt(10000), bankeraddr)
	require.NotNil(t, txIssue)

	key := store.NewKVStoreKey(t.Name())
	ctx := defaultContext(key)

	txIssue.ValidateData(ctx)
	//require.Equal(t, isvalid, true)

	//result, _ := txIssue.Exec(ctx)
	//require.Equal(t, result.Code, types.ABCICodeOK)

	return
}

func defaultContext(key store.StoreKey) context.Context {
	mapperMap := make(map[string]bmapper.IMapper)

	mainmapper := mapper.NewMainMapper()
	mainmapper.SetCodec(cdc)
	mainKey := mainmapper.GetStoreKey()
	mapperMap[mapper.GetMainStoreKey()] = mainmapper

	accountMapper := bacc.NewAccountMapper(nil, account.ProtoQOSAccount)
	accountMapper.SetCodec(cdc)
	accountKey := accountMapper.GetStoreKey()
	mapperMap[bacc.AccountMapperName] = accountMapper

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(mainKey, store.StoreTypeIAVL, db)
	cms.MountStoreWithDB(accountKey, store.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := context.NewContext(cms, abci.Header{}, false, nil, mapperMap)
	return ctx
}
