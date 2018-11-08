package main

import (
	"flag"
	"fmt"
	baccount "github.com/QOSGroup/qbase/account"
	btxs "github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/app"
	"github.com/QOSGroup/qos/test"
	"github.com/QOSGroup/qos/txs"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/rpc/client"
)

//cli端： 	qos\cmd\qosappcli.go
//qosd: 	qos\cmd\qosd\main.go
//nonce: 	参数中含nonce值，每次执行需修改nonce+=1

//1, qos初始化(qosd)		init --chain-id=qos
//2, qos启动(qosd)			start --with-tendermint=true
//3, 发送TxCreateQSC(cli端)	-m=txcreateqsc -pathqsc=d:\qsc.crt -pathbank=d:\banker.crt -qscchainid=qsctest -qoschainid=qos -maxgas=100 -nonce=1
//	 	3.1, pathqsc & pathbank 分别为qsc和banker的CA文件路径
//	 	3.2, example: D:\banker.crt
// 		3.3, 参考: github.com/QOSGroup/kepler/examples/v1  (qsc.crt, banker.crt)
//4, 发送TxIssue(cli端)		-m=txissue -qscname=QSC -nonce=1 -chainid=qsctest -maxgas=100
//		4.1, qscname需和banker中的qscname相同，区分大小写
//--------------------------------
//查询账户信息(步骤2,3,4之后都可以执行查询账户信息，验证tx结果)
//
//cli端查询banker		-m=accquery -addr=address1l7d3dc26adk9gwzp777s3a9p5tprn7m43p99cg
//cli端查询acc1			-m=accquery -addr=address1zsqzn6wdecyar6c6nzem3e8qss2ws95csr8d0r

func main() {
	cdc := app.MakeCodec()

	mode := flag.String("m", "", "client mode: get/send")
	qoschainid := flag.String("qoschainid", "qos", "chainid of qos, used in tx")

	// query
	addr := flag.String("addr", "", "input account addr(bech32)")

	// createQSC
	extrate := flag.String("extrate", "1:280.0000", "extrate: qos/qscxxx")
	pathqsc := flag.String("pathqsc", "", "path of CA(qsc)")
	pathbank := flag.String("pathbank", "", "path of CA(banker)")
	qscchainid := flag.String("qscchainid", "", "chainid of qsc, used in tx")

	// issue
	amount := flag.Int64("amount", 100000, "coins send to banker")
	qscname := flag.String("qscname", "qsc1", "qsc name")
	nonce := flag.Int64("nonce", 0, "value of nonce")

	// txstd
	maxgas := flag.Int64("maxgas", 0, "maxgas for txstd")

	flag.Parse()

	http := client.NewHTTP("tcp://127.0.0.1:26657", "/websocket")

	switch *mode {
	case "accquery": // 账户查询
		queryAccount(http, cdc, addr)
	case "txissue":
		accary := test.InitKeys(cdc)
		privkey := accary[1].PrivKey
		addr := accary[1].Acc.GetAddress()
		stdTxIssue(http, cdc, *qscname, btypes.NewInt(*amount), addr, privkey, *nonce, *qoschainid, *maxgas)
	case "txcreateqsc":
		caQsc := txs.FetchCA(*pathqsc)
		caBanker := txs.FetchCA(*pathbank)
		acc := []txs.AddrCoin{}
		accary := test.InitKeys(cdc)

		for i := 2; i < 5; i++ {
			acc = append(acc, txs.AddrCoin{accary[i].Acc.GetAddress(), btypes.NewInt(int64(i + 100))})
		}

		stdTxCreateQSC(http, cdc, caQsc, caBanker, accary[0].Acc.GetAddress(), accary[0].PrivKey, &acc, *extrate, "", *qscchainid, *qoschainid, *maxgas, *nonce)
	default:
		fmt.Printf("%s doen't support now!", *mode)
	}
}

func stdTxCreateQSC(http *client.HTTP, cdc *amino.Codec, caqsc *[]byte, cabank *[]byte,
	createaddr btypes.Address, privkey crypto.PrivKey, accs *[]txs.AddrCoin,
	extrate string, dsp string, qscchainid string, qoschainid string, maxgas int64, nonce int64) {

	var accsigner []*accsign
	accsigner = append(accsigner, &accsign{
		privkey,
		nonce,
	})

	tx := txs.NewCreateQsc(cdc, caqsc, cabank, qscchainid, createaddr,accs, extrate, dsp)
	if tx == nil {
		panic("createqsc error!")
	}

	txstd := genTxStd(cdc, tx, qoschainid, maxgas, accsigner)
	if txstd == nil {
		panic("gentxstd error!")
	}

	broadcastTxStd(http, cdc, txstd)
}

func stdTxIssue(http *client.HTTP, cdc *amino.Codec, qscname string, amount btypes.BigInt, bankaddr btypes.Address,
	privkey crypto.PrivKey, nonce int64, chainid string, maxgas int64) {
	tx := txs.NewTxIssueQsc(qscname, amount, bankaddr)
	if tx == nil {
		panic("create txissue error")
	}

	accsigner := []*accsign{}
	accsigner = append(accsigner, &accsign{privkey, nonce})
	txstd := genTxStd(cdc, tx, chainid, maxgas, accsigner)

	broadcastTxStd(http, cdc, txstd)
}

// 查询账户状态
func queryAccount(http *client.HTTP, cdc *amino.Codec, addr *string) (acc *account.QOSAccount) {
	if *addr == "" {
		panic("usage: -m=accquery -addr=xxx")
	}
	address, _ := btypes.GetAddrFromBech32(*addr)
	key := baccount.AddressStoreKey(address)
	result, err := http.ABCIQuery("/store/acc/key", key)
	if err != nil {
		panic(err)
	}

	queryValueBz := result.Response.GetValue()
	cdc.UnmarshalBinaryBare(queryValueBz, &acc)

	jsacc, _ := cdc.MarshalJSONIndent(acc, "", " ")
	fmt.Println(fmt.Sprintf("query addr is %s = %s", *addr, jsacc))

	return
}

type accsign struct {
	privkey crypto.PrivKey `json:"privkey"`
	nonce   int64          `json:"nonce"`
}

func genTxStd(cdc *amino.Codec, itx btxs.ITx, chainid string, maxgas int64, accsigner []*accsign) (txstd *btxs.TxStd) {
	if len(accsigner) == 0 {
		return nil
	}

	tx := btxs.NewTxStd(itx, chainid, btypes.NewInt(maxgas))

	for _, acc := range accsigner {
		signdata, _ := tx.SignTx(acc.privkey, acc.nonce)
		tx.Signature = append(tx.Signature, btxs.Signature{
			Pubkey:    acc.privkey.PubKey(),
			Signature: signdata,
			Nonce:     acc.nonce,
		})
	}

	return tx
}

func broadcastTxStd(http *client.HTTP, cdc *amino.Codec, txstd *btxs.TxStd) {
	tx, err := cdc.MarshalBinaryBare(txstd)
	if err != nil {
		panic("use cdc encode object fail")
	}

	result, err := http.BroadcastTxSync(tx)
	if err != nil {
		fmt.Println(err)
		panic("BroadcastTxSync err")
	}

	fmt.Println(fmt.Sprintf("tx result:  %v", result))
}
