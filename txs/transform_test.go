package txs

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/store"
	"github.com/QOSGroup/qbase/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"

	"strconv"
	"testing"
)

func newAddrTrans(n int) (ret *[]AddrTrans) {
	var buf []AddrTrans

	for i := 1; i < 1+n; i++ {
		addrtrans := AddrTrans{[]byte("address" + strconv.Itoa(i)),
			types.NewInt(int64(i)),
			"qsc" + strconv.Itoa(i)}
		buf = append(buf, addrtrans)
	}
	ret = &buf

	return
}

func defaultContext(key store.StoreKey) context.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, store.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := context.NewContext(cms, abci.Header{}, false, nil, nil)
	return ctx
}

func TestNewTransform(t *testing.T) {
	sender := newAddrTrans(3)
	receiver := newAddrTrans(5)

	//Total Amount of sender & receiver is equal
	(*sender)[2].Amount = (*sender)[2].Amount.Add(types.NewInt(9))
	txTran := NewTransform(sender, receiver)
	require.NotNil(t, txTran)
	isvalid := txTran.ValidateData()
	require.Equal(t, isvalid, true)

	key := store.NewKVStoreKey(t.Name())
	ctx := defaultContext(key)
	result := txTran.Exec(ctx)
	require.Equal(t, result.Code, types.ABCICodeOK)

	return
}
