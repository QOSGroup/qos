package txs

import (
	"github.com/QOSGroup/qbase/store"
	"github.com/QOSGroup/qbase/types"
	"github.com/stretchr/testify/require"
	go_amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
	"testing"
)

func TestNewCreateQsc(t *testing.T) {
	key := store.NewKVStoreKey(t.Name())
	ctx := defaultContext(key)
	cdc := go_amino.NewCodec()
	cryptoAmino.RegisterAmino(cdc)
	RegisterCodec(cdc)
	caQsc := FetchQscCA()
	caBanker := FetchBankerCA()

	accInit := []AddrCoin{
		AddrCoin{[]byte("qscAddress1"), types.NewInt(3000)},
		AddrCoin{[]byte("qscAddress2"), types.NewInt(2000)},
	}

	txCreateQsc := NewCreateQsc(cdc, caQsc, caBanker, []byte("creator_qsc10"), &accInit, "1:280.0000", "createQsc")
	require.NotNil(t, txCreateQsc)
	isvalid := txCreateQsc.ValidateData(ctx)
	require.Equal(t, isvalid, true)

	//result,_ := txCreateQsc.Exec(ctx)
	//require.Equal(t, result.Code, types.ABCICodeOK)

	return
}
