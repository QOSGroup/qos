package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	btxs "github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/app"
	"github.com/QOSGroup/qos/txs"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/rpc/client"
	"strconv"
	"strings"
)

func main() {
	cdc := app.MakeCodec()

	mode := flag.String("m", "", "client mode: tx/qcptx/qcpseq/qcpquery")
	tx := flag.String("tx", "", "client tx")
	prikeys := flag.String("prikey", "", "client prikeys")
	nonces := flag.String("nonce", "", "client nonce")
	fromChain := flag.String("fromchain", "", "input qcp fromchain")
	toChain := flag.String("tochain", "", "input qcp tochain")
	qcpPriKey := flag.String("qcpprikey", "", "input qcp prikey")
	qcpseq := flag.Int64("qcpseq", 0, "input qcp sequence")

	flag.Parse()

	http := client.NewHTTP("tcp://127.0.0.1:26657", "/websocket")
	switch *mode {
	case "tx":
		transfer(http, cdc, tx, prikeys, nonces)
	case "qcptx":
		qcpTransfer(http, cdc, tx, prikeys, nonces, fromChain, toChain, qcpPriKey, qcpseq)
	case "qcpseq":
		qcpSeq(http, cdc, fromChain)
	case "qcpquery":
		qcpQuery(http, cdc, fromChain, qcpseq)
	default:
		fmt.Println("invalid command")
	}
}

//-m=tx -prikey=0xa328891040ae9b773bcd30005235f99a8d62df03a89e4f690f9fa03abb1bf22715fc9ca05613f2d8061492e9f8149510b5b67d340d199ff24f34c85dbbbd7e0df780e9a6cc -nonce=0 -tx='{"type":"qos/txs/TransferTx","value":{"senders":[{"addr":"address1k0m8ucnqug974maa6g36zw7g2wvfd4sug6uxay","qos":"10","qscs":[{"coin_name":"qstar","amount":"100"}]}],"receivers":[{"addr":"address1uwxz03zryaqxmaw88pd0reck5qdj07k5k8l2au","qos":"10","qscs":[{"coin_name":"qstar","amount":"50"}]},{"addr":"address1wzle536jwgn8e6evrwjxl4daacg00wd4tahc32","qos":"0","qscs":[{"coin_name":"qstar","amount":"50"}]}]}}'
func transfer(http *client.HTTP, cdc *amino.Codec, tx *string, prikeys *string, nonces *string) {
	if *tx == "" || *prikeys == "" || *nonces == "" {
		panic("usage: -m=tx -tx=xxx -prikeys=xxx,xxx -nonces=xxx,xxx")
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

//-m=qcptx -prikey=0xa328891040ae9b773bcd30005235f99a8d62df03a89e4f690f9fa03abb1bf22715fc9ca05613f2d8061492e9f8149510b5b67d340d199ff24f34c85dbbbd7e0df780e9a6cc,0xa3288910401696447c6e716e0d88492d9d4e6351114b5d36d541138bfd4b08b340c138468a3f7a523ce67dd15c599053bad424760064b8f1537c7bfce28fb30bddd5fb1f77 -nonce=1,0 -tx='{"type":"qos/txs/TransferTx","value":{"senders":[{"addr":"address1k0m8ucnqug974maa6g36zw7g2wvfd4sug6uxay","qos":"10","qscs":[{"coin_name":"qstar","amount":"100"}]},{"addr":"address1uwxz03zryaqxmaw88pd0reck5qdj07k5k8l2au","qos":"10","qscs":[{"coin_name":"qstar","amount":"50"}]}],"receivers":[{"addr":"address1wzle536jwgn8e6evrwjxl4daacg00wd4tahc32","qos":"20","qscs":[{"coin_name":"qstar","amount":"150"}]}]}}' -fromchain=qstar -tochain=qos-chain -qcpprikey=0xa3288910405746e29aeec7d5ed56fac138b215e651e3244e6d995f25cc8a74c40dd1ef8d2e8ac876faaa4fb281f17fb9bebb08bc14e016c3a88c6836602ca97595ae32300b -qcpseq=1
func qcpTransfer(http *client.HTTP, cdc *amino.Codec, tx *string, prikeys *string, nonces *string, fromChain *string, toChain *string, caPriHex *string, qcpseq *int64) {
	if *tx == "" || *nonces == "" || *prikeys == "" || *fromChain == "" || *toChain == "" || *caPriHex == "" || *qcpseq <= 0 {
		panic("usage: -m=qcptx -tx=xxx -prikeys=xxx,xxx -nonces=xxx,xxx -fromchain=xxx -tochain=xx -qcpprikey=xxx -qcpseq=xxx")
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
		stdTx.Signature = append(stdTx.Signature, btxs.Signature{
			Pubkey:    priKey.PubKey(),
			Signature: signature,
			Nonce:     n,
		})
	}

	qcpTx := btxs.NewTxQCP(stdTx, *fromChain, *toChain, *qcpseq, 0, 0, false, "")
	caHex, _ := hex.DecodeString((*caPriHex)[2:])
	var caPriKey ed25519.PrivKeyEd25519
	cdc.MustUnmarshalBinaryBare(caHex, &caPriKey)
	sig, _ := qcpTx.SignTx(caPriKey)
	qcpTx.Sig.Nonce = *qcpseq
	qcpTx.Sig.Signature = sig
	qcpTx.Sig.Pubkey = caPriKey.PubKey()

	bz, err := cdc.MarshalBinaryBare(qcpTx)
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

//-m=qcpseq -fromchain=qstar
func qcpSeq(http *client.HTTP, cdc *amino.Codec, chainid *string) {
	// in sequence
	keyIn := fmt.Sprintf("sequence/in/%s", *chainid)
	inResult, err := http.ABCIQuery("/store/qcp/key", []byte(keyIn))
	if err != nil {
		panic(err)
	}
	var in int64
	if inResult.Response.GetValue() != nil {
		cdc.UnmarshalBinaryBare(inResult.Response.GetValue(), &in)
	}

	// out sequence
	keyOut := fmt.Sprintf("sequence/out/%s", *chainid)
	outResult, err := http.ABCIQuery("/store/qcp/key", []byte(keyOut))
	if err != nil {
		panic(err)
	}
	var out int64
	if outResult.Response.GetValue() != nil {
		cdc.UnmarshalBinaryBare(outResult.Response.GetValue(), &out)
	}

	fmt.Println(fmt.Sprintf("query chain is %s, sequence in/out: %d/%d", *chainid, in, out))
}

//-m=qcpseq -fromchain=qstar -qcpseq=1
func qcpQuery(http *client.HTTP, cdc *amino.Codec, chainid *string, qcpseq *int64) {
	key := fmt.Sprintf("tx/out/%s/%d", *chainid, *qcpseq)
	result, err := http.ABCIQuery("/store/qcp/key", []byte(key))
	if err != nil {
		panic(err)
	}

	var tx btxs.TxQcp
	if result.Response.GetValue() != nil {
		cdc.UnmarshalBinaryBare(result.Response.GetValue(), &tx)
	}

	json, _ := cdc.MarshalJSON(tx)
	fmt.Println(fmt.Sprintf("query chain is %s, tx out[%d] is %s", *chainid, *qcpseq, json))
}
