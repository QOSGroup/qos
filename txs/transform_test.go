package txs

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/store"
	"github.com/QOSGroup/qbase/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	abci "github.com/tendermint/tendermint/abci/types"

	"strconv"
	"testing"
)

func newAddrTrans(n int) (ret []AddrTrans){
	for i := 0 ; i < n ; i++ {
		addrtrans := AddrTrans{[]byte("address" + strconv.Itoa(i)),
		types.NewInt(int64(i)),
		"qsc" + strconv.Itoa(i)}
		ret = append(ret, addrtrans)
	}

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
	sender[2].Amount.Add(types.NewInt(9))
	txTran := NewTransform(&sender, &receiver)
	if txTran == nil  {
		t.Error("NewTranform error!")
		return
	}

	if !txTran.ValidateData(){
		t.Error("Invalidate Transform")
		return
	}

	key := store.NewKVStoreKey(t.Name())
	ctx := defaultContext(key)
	result := txTran.Exec(ctx)
	if result.Code != types.ABCICodeOK {
		fmt.Printf("Execute Error: %d", result.Code)
		return
	}

	fmt.Print("Excute successful!")
	return
}