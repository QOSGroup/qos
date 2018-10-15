package txs

import (
	"fmt"
	"testing"
	"github.com/QOSGroup/qbase/store"
	"github.com/QOSGroup/qbase/types"
	go_amino "github.com/tendermint/go-amino"
)


func TestNewCreateQsc(t *testing.T) {
	key := store.NewKVStoreKey(t.Name())
	ctx := defaultContext(key)
	cdc := go_amino.NewCodec()
	caQsc := FetchQscCA()
	caBanker := FetchBankerCA()
	creater := []byte("creator_qsc10")
	accInit := []AddrCoin{
		AddrCoin{[]byte("qscAddress1"), types.NewInt(3000)},
		AddrCoin{[]byte("qscAddress2"), types.NewInt(2000)},
	}

	txCreateQsc := NewCreateQsc(cdc,caQsc, caBanker, creater, &accInit,"1:280.0000","createQsc")
	if txCreateQsc == nil {
		t.Error("new TxCreateQsc error!")
		return
	}

	if !txCreateQsc.ValidateData() {
		t.Error("TxCreateQsc Invalidate error")
		return
	}

	result := txCreateQsc.Exec(ctx)
	if result.Code != types.ABCICodeOK {
		fmt.Printf("TxCreateQsc execute Error: %d", result.Code)
		return
	}

	fmt.Print("TxCreateQsc excute successful!")
	return
}
