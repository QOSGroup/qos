package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	btxs "github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/app"
	"github.com/QOSGroup/qos/txs"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/rpc/client"
	"strconv"
	"strings"
)

//-prikey=0xa328891040ae9b773bcd30005235f99a8d62df03a89e4f690f9fa03abb1bf22715fc9ca05613f2d8061492e9f8149510b5b67d340d199ff24f34c85dbbbd7e0df780e9a6cc
//-nonce=0
//-tx={\"type\":\"qos/txs/TransferTx\",\"value\":{\"senders\":[{\"addr\":\"address1k0m8ucnqug974maa6g36zw7g2wvfd4sug6uxay\",\"qos\":\"10\",\"qscs\":[{\"coin_name\":\"qstar\",\"amount\":\"100\"}]}],\"receivers\":[{\"addr\":\"address103eak408d4yp944wv58epp3neyah8z5dlwyzg4\",\"qos\":\"10\",\"qscs\":[{\"coin_name\":\"qstar\",\"amount\":\"50\"}]},{\"addr\":\"address1wzle536jwgn8e6evrwjxl4daacg00wd4tahc32\",\"qos\":\"0\",\"qscs\":[{\"coin_name\":\"qstar\",\"amount\":\"50\"}]}]}}
func main() {
	cdc := app.MakeCodec()

	tx := flag.String("tx", "", "client tx")
	prikeys := flag.String("prikey", "", "client prikeys")
	nonces := flag.String("nonce", "", "client nonce")

	flag.Parse()

	http := client.NewHTTP("tcp://127.0.0.1:26657", "/websocket")

	if *tx == "" {
		panic("usage: -tx=xxx(json)")
	}

	ps := strings.Split(*prikeys, ",")
	ns := strings.Split(*nonces, ",")

	transferTx := txs.TransferTx{}
	err := cdc.UnmarshalJSON([]byte(*tx), &transferTx)
	if err != nil {
		panic(err)
	}

	stdTx := btxs.NewTxStd(&transferTx, "qos-chain", types.ZeroInt())
	for i, p := range ps {
		priHex, _ := hex.DecodeString(p[2:])
		var priKey ed25519.PrivKeyEd25519
		cdc.MustUnmarshalBinaryBare(priHex, &priKey)
		n, _ := strconv.ParseInt(ns[i], 10, 64)
		signature, _ := stdTx.SignTx(priKey, n)
		stdTx.Signature = []btxs.Signature{{
			Pubkey:    priKey.PubKey(),
			Signature: signature,
			Nonce:     n,
		}}
	}
	bz, err := cdc.MarshalBinaryBare(stdTx)
	if err != nil {
		panic("use cdc encode object fail")
	}
	_, err = http.BroadcastTxSync(bz)
	if err != nil {
		fmt.Println(err)
		panic("BroadcastTxSync err")
	}
	fmt.Println("send tx success")
}
