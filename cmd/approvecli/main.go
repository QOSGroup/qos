package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	bacc "github.com/QOSGroup/qbase/account"
	btxs "github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/app"
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
	addr := flag.String("addr", "", "input account addr(bech32)")

	flag.Parse()

	http := client.NewHTTP("tcp://127.0.0.1:26657", "/websocket")

	switch *mode {
	case "tx": // 预授权
		approveHandle(http, cdc, *command, *from, *to, *prikey, *nonce, *coins)
	case "account": // 查询账户
		queryAccount(http, cdc, *addr)
	case "approve": // 查询授权
		queryApprove(http, cdc, *from, *to)
	default:
		fmt.Println("what you want?")
	}
}

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
		qscs := types.QSCS{}
		for _, val := range coinAndAmounts {
			coinAndAmount := strings.Split(val, ",")
			amount, _ := strconv.ParseInt(coinAndAmount[1], 10, 64)
			qsc := types.QSC{
				Name:   coinAndAmount[0],
				Amount: btypes.NewInt(amount),
			}
			qscs = append(qscs, qsc)
		}
		approve := &types.Approve{
			From:  fromAddr,
			To:    toAddr,
			Coins: qscs,
		}
		var stdTx btxs.TxStd
		switch command {
		case "create":
			tx := txs.TxApproveCreate{approve,}
			stdTx.ITx = &tx
		case "increase":
			tx := txs.TxApproveIncrease{approve,}
			stdTx.ITx = &tx
		case "decrease":
			tx := txs.TxApproveDecrease{approve,}
			stdTx.ITx = &tx
		case "use":
			tx := txs.TxApproveUse{approve,}
			stdTx.ITx = &tx
		}
		stdTx.ChainID = "qos-chain"
		stdTx.MaxGas = btypes.NewInt(int64(0))
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
		approve := &types.ApproveCancel{
			From: fromAddr,
			To:   toAddr,
		}
		tx := txs.TxApproveCancel{approve,}
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

func queryAccount(http *client.HTTP, cdc *amino.Codec, addr string) {
	if addr == "" {
		panic("usage: -m=acc -addr=xxx")
	}
	address, _ := btypes.GetAddrFromBech32(addr)
	key := bacc.AddressStoreKey(address)
	result, err := http.ABCIQuery("/store/acc/key", key)
	if err != nil {
		panic(err)
	}

	queryValueBz := result.Response.GetValue()
	var acc *account.QOSAccount
	cdc.UnmarshalBinaryBare(queryValueBz, &acc)

	fmt.Println(fmt.Sprintf("query addr is %s = %v", addr, acc))
}

func queryApprove(http *client.HTTP, cdc *amino.Codec, from string, to string) {
	if from == "" || to == "" {
		panic("usage: -m=approve -from=xxx -to=xxx")
	}
	key := fmt.Sprintf("from:[%s]/to:[%s]", from, to)
	result, err := http.ABCIQuery("/store/approve/key", []byte(key))
	if err != nil {
		panic(err)
	}

	queryValueBz := result.Response.GetValue()
	var approve types.Approve
	cdc.UnmarshalBinaryBare(queryValueBz, &approve)

	fmt.Println(fmt.Sprintf("query addr is from:[%s]/to:[%s] = %v", from, to, approve))
}
