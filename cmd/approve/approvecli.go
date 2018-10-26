package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	btxs "github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/app"
	"github.com/QOSGroup/qos/mapper"
	"github.com/QOSGroup/qos/txs"
	"github.com/QOSGroup/qos/types"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/rpc/client"
	"strconv"
	"strings"
)

func main() {
	cdc := app.MakeCodec()

	mode := flag.String("m", "", "client mode: approve,...")
	command := flag.String("c", "", "client command, for approve: create,increase,decrease,use,cancel")
	from := flag.String("from", "", "input from addr")
	prikey := flag.String("prikey", "", "input signer prikey")
	nonce := flag.Int64("nonce", 0, "input sender nonce")
	to := flag.String("to", "", "input to addr")
	coins := flag.String("coins", "", "input coinname,coinamount;coinname,coinamount")

	flag.Parse()

	http := client.NewHTTP("tcp://127.0.0.1:26657", "/websocket")

	switch *mode {
	case "tx": // 预授权
		approveHandle(http, cdc, *command, *from, *to, *prikey, *nonce, *coins)
	case "approve": // 查询授权
		queryApprove(http, cdc, *from, *to)
	default:
		fmt.Println("what you want?")
	}
}

// 创建
//-m=tx -c=create -from=address1k0m8ucnqug974maa6g36zw7g2wvfd4sug6uxay -to=address103eak408d4yp944wv58epp3neyah8z5dlwyzg4 -prikey=0xa328891040ae9b773bcd30005235f99a8d62df03a89e4f690f9fa03abb1bf22715fc9ca05613f2d8061492e9f8149510b5b67d340d199ff24f34c85dbbbd7e0df780e9a6cc -coins=qos,10;qstar,100 -nonce=0
// 增加
//-m=tx -c=increase -from=address1k0m8ucnqug974maa6g36zw7g2wvfd4sug6uxay -to=address103eak408d4yp944wv58epp3neyah8z5dlwyzg4 -prikey=0xa328891040ae9b773bcd30005235f99a8d62df03a89e4f690f9fa03abb1bf22715fc9ca05613f2d8061492e9f8149510b5b67d340d199ff24f34c85dbbbd7e0df780e9a6cc -coins=qos,10;qstar,100 -nonce=1
// 减少
//-m=tx -c=decrease -from=address1k0m8ucnqug974maa6g36zw7g2wvfd4sug6uxay -to=address103eak408d4yp944wv58epp3neyah8z5dlwyzg4 -prikey=0xa328891040ae9b773bcd30005235f99a8d62df03a89e4f690f9fa03abb1bf22715fc9ca05613f2d8061492e9f8149510b5b67d340d199ff24f34c85dbbbd7e0df780e9a6cc -coins=qstar,100 -nonce=2
// 使用
//-m=tx -c=use -from=address1k0m8ucnqug974maa6g36zw7g2wvfd4sug6uxay -to=address103eak408d4yp944wv58epp3neyah8z5dlwyzg4 -prikey=0xa3288910405746e29aeec7d5ed56fac138b215e651e3244e6d995f25cc8a74c40dd1ef8d2e8ac876faaa4fb281f17fb9bebb08bc14e016c3a88c6836602ca97595ae32300b -coins=qstar,100 -nonce=0
// 取消
//-m=tx -c=cancel -from=address1k0m8ucnqug974maa6g36zw7g2wvfd4sug6uxay -to=address103eak408d4yp944wv58epp3neyah8z5dlwyzg4 -prikey=0xa328891040ae9b773bcd30005235f99a8d62df03a89e4f690f9fa03abb1bf22715fc9ca05613f2d8061492e9f8149510b5b67d340d199ff24f34c85dbbbd7e0df780e9a6cc -nonce=3
func approveHandle(http *client.HTTP, cdc *amino.Codec, command string, from string, to string, prihex string, nonce int64, coinStr string) {
	if from == "" || to == "" || prihex == "" || nonce < 0 || (command != "cancel" && coinStr == "") {
		panic("usage: -m=approve -c=create/increase/decrease/use/cancel -from=xxx -to=xxx -coin=xxx,xxx;xxx,xxx -prikey=xxx -nonce=xxx(>=0)")
	}
	fromAddr, _ := btypes.GetAddrFromBech32(from)
	toAddr, _ := btypes.GetAddrFromBech32(to)
	priHex, _ := hex.DecodeString(prihex[2:])
	var priKey ed25519.PrivKeyEd25519
	cdc.MustUnmarshalBinaryBare(priHex, &priKey)
	var bz []byte
	var err error
	if command != "cancel" { // 创建、增加、减少、使用授权
		coinAndAmounts := strings.Split(coinStr, ";")
		qscs := []*types.QSC{}
		qos := btypes.BigInt{}
		for _, val := range coinAndAmounts {
			coinAndAmount := strings.Split(val, ",")
			amount, _ := strconv.ParseInt(coinAndAmount[1], 10, 64)
			if coinAndAmount[0] == "qos" {
				qos = btypes.NewInt(amount)
			} else {
				qsc := types.QSC{
					Name:   coinAndAmount[0],
					Amount: btypes.NewInt(amount),
				}
				qscs = append(qscs, &qsc)
			}
		}
		approve := types.Approve{
			From: fromAddr,
			To:   toAddr,
			QOS:  qos,
			QSCs: qscs,
		}
		var stdTx *btxs.TxStd
		switch command {
		case "create":
			stdTx = btxs.NewTxStd(txs.ApproveCreateTx{approve,}, "qos-chain", btypes.NewInt(0))
		case "increase":
			stdTx = btxs.NewTxStd(txs.ApproveIncreaseTx{approve,}, "qos-chain", btypes.NewInt(0))
		case "decrease":
			stdTx = btxs.NewTxStd(txs.ApproveDecreaseTx{approve,}, "qos-chain", btypes.NewInt(0))
		case "use":
			stdTx = btxs.NewTxStd(txs.ApproveUseTx{approve,}, "qos-chain", btypes.NewInt(0))
		}
		signature, _ := stdTx.SignTx(priKey, nonce)
		stdTx.Signature = []btxs.Signature{btxs.Signature{
			Pubkey:    priKey.PubKey(),
			Signature: signature,
			Nonce:     nonce,
		}}
		bz, err = cdc.MarshalBinaryBare(stdTx)
		if err != nil {
			panic("use cdc encode object fail")
		}
	} else { // 取消授权
		approve := types.ApproveCancel{
			From: fromAddr,
			To:   toAddr,
		}
		tx := txs.ApproveCancelTx{approve,}
		stdTx := btxs.NewTxStd(&tx, "qos-chain", btypes.NewInt(int64(0)))
		signature, _ := stdTx.SignTx(priKey, nonce)
		stdTx.Signature = []btxs.Signature{btxs.Signature{
			Pubkey:    priKey.PubKey(),
			Signature: signature,
			Nonce:     nonce,
		}}
		bz, err = cdc.MarshalBinaryBare(stdTx)
		if err != nil {
			panic("use cdc encode object fail")
		}
	}
	_, err = http.BroadcastTxSync(bz)
	if err != nil {
		fmt.Println(err)
		panic("BroadcastTxSync err")
	}
	fmt.Println("send tx success")
}

// 查询预授权
//-m=approve -from=address1k0m8ucnqug974maa6g36zw7g2wvfd4sug6uxay -to=address103eak408d4yp944wv58epp3neyah8z5dlwyzg4
func queryApprove(http *client.HTTP, cdc *amino.Codec, from string, to string) {
	if from == "" || to == "" {
		panic("usage: -m=approve -from=xxx -to=xxx")
	}
	key := mapper.BuildApproveKey(from, to)
	result, err := http.ABCIQuery("/store/approve/key", key)
	if err != nil {
		panic(err)
	}

	queryValueBz := result.Response.GetValue()
	var approve types.Approve
	cdc.UnmarshalBinaryBare(queryValueBz, &approve)

	json, _ := cdc.MarshalJSON(approve)
	fmt.Println(fmt.Sprintf("query addr is from:[%s]/to:[%s] = %s", from, to, string(json)))
}
