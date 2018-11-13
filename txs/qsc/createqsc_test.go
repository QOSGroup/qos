package qsc

import (
	"github.com/QOSGroup/qbase/store"
	"github.com/QOSGroup/qbase/types"
	"github.com/stretchr/testify/require"
	go_amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
	"testing"
)

func TestNewCreateQsc(t *testing.T) {
	key := store.NewKVStoreKey(t.Name())
	ctx := defaultContext(key)
	cdc := go_amino.NewCodec()
	cryptoAmino.RegisterAmino(cdc)
	RegisterCodec(cdc)

	rtprivkey := ed25519.GenPrivKey()
	qscpubkey := ed25519.GenPrivKey().PubKey()
	bankpubkey := ed25519.GenPrivKey().PubKey()
	caQsc := NewCertificate(cdc, "qsc", false, qscpubkey, rtprivkey.PubKey())
	caBanker := NewCertificate(cdc, "qsc", true, bankpubkey, rtprivkey.PubKey())

	accInit := []AddrCoin{
		AddrCoin{[]byte("qscAddress1"), types.NewInt(3000)},
		AddrCoin{[]byte("qscAddress2"), types.NewInt(2000)},
	}

	txCreateQsc := NewCreateQsc(cdc, &caQsc, &caBanker, "chainqsc1",
		[]byte("creator_qsc10"), &accInit, "1:280.0000", "createQsc")
	require.NotNil(t, txCreateQsc)
	txCreateQsc.ValidateData(ctx)

	//result,_ := txCreateQsc.Exec(ctx)
	//require.Equal(t, result.Code, types.ABCICodeOK)

	return
}
