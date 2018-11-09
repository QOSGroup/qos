package main

import (
	"flag"
	"fmt"
	baccount "github.com/QOSGroup/qbase/account"
	btxs "github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/app"
	"github.com/QOSGroup/qos/mapper"
	"github.com/QOSGroup/qos/txs"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/bech32"
	"github.com/tendermint/tendermint/rpc/client"
	"math/rand"
)

//cli端： 	qos\cmd\qosappcli.go
//qosd: 	qos\cmd\qosd\main.go
//nonce: 	参数中含nonce值，每次执行需修改nonce+=1

//1, qos初始化(qosd)		init --chain-id=qos
//2, qos启动(qosd)			start --with-tendermint=true
//3, 发送TxCreateQSC(cli端)	-m=txcreateqsc -pathqsc=d:\qsc.crt -pathbank=d:\banker.crt -qscchainid=qcptest -qoschainid=qos -maxgas=100 -nonce=1 -privkeycreator=rDwWppdGKFCv0wUxFqVID87GI/CFwLbL9p6EM6ug5brPbkXQoZMIH9+Rgi1/vFcNJUHp88fKZDNFdEif8dg73A==
//	 	3.1, pathqsc/pathbank: qsc/banker的CA文件路径
// 			 github.com/QOSGroup/kepler/examples/v1  (qsc.crt, banker.crt)
//		3.2, qscchainid:  联盟链chainid
//		3.3, qoschainid:  公链chainid
//		3.4, maxgas: 期望最大gas花费
//		3.5, privkeycreator: creator的 private key.
//4, 发送TxIssue(cli端)		-m=txissue -qscname=QSC -nonce=1 -qoschainid=qos -maxgas=100 -privkeybank=maD8NeYMqx6fHWHCiJdkV4/B+tDXFIpY4LX4vhrdmAYIKC67z/lpRje4NAN6FpaMBWuIjhWcYeI5HxMh2nTOQg==
//		4.1, qscname需和banker中的qscname相同，区分大小写
//		4.2, qoschainid:  公链chainid
//		4.3, privkeybank: banker的privatekey
//--------------------------------
//查询创建的联盟链信息
//	cli端:	-m=qscquery -qscname=QSC
//		qscname: 要查询的联盟链名称
//
//查询账户信息(步骤2,3,4之后都可以执行查询账户信息，验证tx结果)
//
//cli端查询creator		-m=accquery -addr=address1auug9tjmkm00w36savxjywmj0sjccaam3pvjfu
//cli端查询banker		-m=accquery -addr=address1l7d3dc26adk9gwzp777s3a9p5tprn7m43p99cg
//cli端查询acc1			-m=accquery -addr=address1zsqzn6wdecyar6c6nzem3e8qss2ws95csr8d0r
//cli端查询acc2			-m=accquery -addr=address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355
//cli端查询acc3			-m=accquery -addr=address1y9r4pjjnvkmpvw46de8tmwunw4nx4qnz2ax5ux

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
	privkeycreator := flag.String("privkeycreator", "", "private key of creator")

	// issue
	amount := flag.Int64("amount", 100000, "coins send to banker")
	qscname := flag.String("qscname", "qsc1", "qsc name")
	nonce := flag.Int64("nonce", 0, "value of nonce")
	privkeybank := flag.String("privkeybank", "", "private key of creator")

	// txstd
	maxgas := flag.Int64("maxgas", 0, "maxgas for txstd")

	flag.Parse()

	http := client.NewHTTP("tcp://127.0.0.1:26657", "/websocket")

	switch *mode {
	case "accquery": // 账户查询
		queryAccount(http, cdc, addr)
	case "qscquery":
		queryQscInfo(http, cdc, *qscname)
	case "txissue":
		addr, privkey := parseJsonPrivkey(cdc, *privkeybank)
		stdTxIssue(http, cdc, *qscname, btypes.NewInt(*amount), addr, privkey, *nonce, *qoschainid, *maxgas)
	case "txcreateqsc":
		if !cmdcheck(*mode, *privkeycreator, *pathbank, *pathqsc) {
			return
		}

		caQsc := txs.FetchCA(*pathqsc)
		caBanker := txs.FetchCA(*pathbank)
		acc := []txs.AddrCoin{}
		accprivkeys := []string {
			"vAeIlHuWjvz/JmyGcB46ZHfCZdXCYuRogqxDgjYUM5wNwKIyIYQBs9VZxGyD9FS5J4XvZntnUaTtoGsEl7+3hg==",
			"31PlT2p6UICjV63dG7Nh3Mh9W0b+7FAEU+KOAxyNbZ29rwqNzxQJlQPh59tZpbS1EdIT6TE5N6L72se9BUe9iw==",
			"9QkouVPl29N2v1lBO1+azUDqm38fAgs6d3Xo8DcnCus7xjMqsavhc190xCGzZuXcjapUahi7Y7v2DD4hzVCAsQ==",
		}

		for _, v := range accprivkeys {
			addracc, _ := parseJsonPrivkey(cdc, v)
			acc = append(acc, txs.AddrCoin{addracc, btypes.NewInt(int64(rand.Int()%1000 + 100))})
		}

		addr, privkey := parseJsonPrivkey(cdc, *privkeycreator)
		stdTxCreateQSC(http, cdc, caQsc, caBanker, addr, privkey, &acc, *extrate, "", *qscchainid, *qoschainid, *maxgas, *nonce)
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

func queryQscInfo(http *client.HTTP, cdc *amino.Codec, qscname string) (err error) {
	key := fmt.Sprintf("qsc/[%s]", qscname)
	result, err := http.ABCIQuery("store/base/key", []byte(key))
	if err != nil {
		panic(err)
	}

	var qcpinfo mapper.QscInfo
	queryValueBz := result.Response.GetValue()
	if len(queryValueBz) == 0 {
		fmt.Printf("Chain (%s) not exist!", qscname)
		return nil
	}

	err = cdc.UnmarshalBinaryBare(queryValueBz, &qcpinfo)
	if err != nil {
		panic(err)
	}

	addr,err := bech32.ConvertAndEncode("address", qcpinfo.BankAddr)
	fmt.Printf("qscname:%s \nbankeraddr:%s \nchainid: %s\n", qcpinfo.Qscname, addr, qcpinfo.ChainID)

	return nil
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

func parseJsonPrivkey(cdc *amino.Codec, jsprivkey string) (addr []byte, privkey ed25519.PrivKeyEd25519) {
	privstr := fmt.Sprintf(` {
 			 	"type": "tendermint/PrivKeyEd25519",
 			 	"value": "%s"
 			}`, jsprivkey)
	err := cdc.UnmarshalJSON([]byte(privstr), &privkey)
	if err != nil {
		panic("parse json privkey error!")
	}

	addr = privkey.PubKey().Address()

	return
}

func cmdcheck(mode string, option...string) bool {
	switch mode {
	case "txcreateqsc":
		//*privkeybank, *pathbank, *pathqsc
		for idx,v := range option {
			switch idx {
			case 0:
				if len(v) < 64 {
					fmt.Print("invalide private key!")
					return false
				}
			case 1,2:
				if len(v) < 6 {
					fmt.Print("invalide path!")
					return false
				}
			}
		}
	default:
		fmt.Printf("cmd (%s) not support!", mode)
		return false
	}

	return true
}