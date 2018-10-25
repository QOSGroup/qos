package main

import (
	"flag"
	"fmt"
	"github.com/tendermint/tendermint/types"

	baccount "github.com/QOSGroup/qbase/account"
	btxs "github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/app"
	"github.com/QOSGroup/qos/txs"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/rpc/client"
)

func main() {
	cdc := app.MakeCodec()

	mode := flag.String("m", "", "client mode: get/send")
	addr := flag.String("addr", "", "input account addr(bech32)")
	receiver := flag.String("receiver", "", "input receive addr")
	amount := flag.Int64("amount", 0, "input amount")

	//createQSC
	extrate := flag.String("extrate", "1:280.0000", "extrate: qos/qscxxx")

	//txstd
	chainid := flag.String("chainid", "", "chainid, used in txstd")
	maxgas := flag.Int64("maxgas", 0, "maxgas for txstd")

	flag.Parse()

	http := client.NewHTTP("tcp://127.0.0.1:26657", "/websocket")

	switch *mode {
	case "accquery": // 账户查询
		queryAccount(http, cdc, addr)
	case "txtransfer":
		//acc := queryAccount(http, cdc, sender)
		//if acc == nil {
		//	fmt.Printf("account %s not exist!", sender)
		//	break
		//}
		stdTransfer(http, cdc, nil, receiver, *amount)
	//case "txissue":
	//	stdTxIssue(http, cdc)
	case "txcreateqsc":
		caQsc := txs.FetchQscCA()
		caBanker := txs.FetchBankerCA()
		acc := []txs.AddrCoin{}
		for i := 1; i < 6; i++ {
			acc = append(acc, txs.AddrCoin{[]byte(initkeys[i].PubKey().Address()), btypes.NewInt(int64(i+100))})
		}
		stdTxCreateQSC(http, cdc, caQsc, caBanker, []byte(initkeys[0].PubKey().Address()), initkeys[0], &acc, *extrate, "", *chainid, *maxgas)
	//case "qcptransfer": // QCP交易
	//	qcpTransfer(http, cdc, sender, prikey, receiver, coinStr, nonce, chainId, qcpPriKey, qcpseq, isresult)
	case "getpubkey":		//账户信息生成（test only: -m=getpubkey）
		var prvkey ed25519.PrivKeyEd25519
		for i:=0; i<len(initkeys); i++ {
			prvkey = initkeys[i]
			pbkey := prvkey.PubKey()
			addr := prvkey.PubKey().Address()
			jsonpbkey,_ := cdc.MarshalJSON(pbkey)
			skey := fmt.Sprintf("pbkey %d : %v", i, pbkey)
			jsskey := fmt.Sprintf("jsonpbkey %d : %s", i, jsonpbkey)
			saddr := fmt.Sprintf("addr %d : %v", i, addr)
			println(skey)
			println(jsskey)
			println(saddr)
		}
	default:
		fmt.Printf("%s doen't support now!", *mode)
	}
}

func stdTxCreateQSC(http *client.HTTP, cdc *amino.Codec, caqsc *[]byte, cabank *[]byte,
	createaddr btypes.Address, privkey ed25519.PrivKeyEd25519, accs *[]txs.AddrCoin,
	extrate string, dsp string, chainid string, maxgas int64) {

	var accsigner []*accsign
	accsigner = append(accsigner, &accsign{
		privkey,
		1,
	})

	tx := txs.NewCreateQsc(cdc, caqsc, cabank, createaddr, accs, extrate, dsp)
	if tx == nil {
		panic("createqsc error!")
	}

	txstd := genTxStd(cdc, tx, chainid, maxgas, accsigner)
	if txstd == nil {
		panic("gentxstd error!")
	}

	broadcastTxStd(http, cdc, txstd)
}

func stdTransfer(http *client.HTTP, cdc *amino.Codec, accsender *account.QOSAccount,  receiver *string, amount int64) {
	if *receiver == "" || amount == 0 {
		panic("usage: -m=txtransfer -from=xxx -to=xxx -amount=xxx -nonce=xxx(>=0)")
	}
	//senderAddr, _ := btypes.GetAddrFromBech32(*sender)
	receiverAddr, _ := btypes.GetAddrFromBech32(*receiver)
	txStd := genStdSendTx(cdc, accsender, receiverAddr, btypes.NewInt(amount), "qos")

	broadcastTxStd(http, cdc, txStd)
}

//func stdTxIssue(http *client.HTTP, cdc *amino.Codec, qscname string, amount btypes.BigInt, chainid string, maxgas int64) {
//	tx := txs.NewTxIssueQsc(qscname, amount)
//	if tx == nil {
//		panic("create txissue error")
//	}
//
//	var acc []*accsign
//	accsign{tx.GetSigner()}
//	banker := tx.GetSigner()
//	txstd := genTxStd(cdc, tx, chainid, maxgas, banker)
//}


// 生成链内交易SendTx
func genStdSendTx(cdc *amino.Codec, accsender *account.QOSAccount, receiver types.Address, amount btypes.BigInt,
	 coinName string) *btxs.TxStd {

	senders := []txs.AddrTrans{
		txs.AddrTrans{[]byte("address1k0m8ucnqug974maa6g36zw7g2wvfd4sug6uxay"), btypes.NewInt(5), coinName},
	}
	receivers := []txs.AddrTrans{
		txs.AddrTrans{[]byte("addrrecv1"),btypes.NewInt(2),coinName,},
		txs.AddrTrans{[]byte("addrrecv2"),btypes.NewInt(3),coinName,},
	}

	if accsender != nil {
		senders = []txs.AddrTrans{}
		senders = append(senders, txs.AddrTrans{accsender.GetAddress(),amount, coinName})
		//for _, sd := range accsender {
		//	senders = append(senders, txs.AddrTrans{sd.GetAddress(),
		//	amount, coinName})
		//}
	}

	txtran := txs.TxTransform{
		Senders: senders,
		Receivers: receivers,
	}

	tx := btxs.NewTxStd(txtran, "qos", btypes.NewInt(int64(100)))

	privKey := ed25519.GenPrivKey()
	signdata, _ := tx.SignTx(privKey, 1)

	tx.Signature = append(tx.Signature, btxs.Signature{
		Pubkey:    privKey.PubKey(),
		Signature: signdata,
		Nonce:     1,
	})

	return tx
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

	fmt.Println(fmt.Sprintf("query addr is %s = %v", *addr, *acc))

	return
}

//---------------------------------------------

var initkeys = [10]ed25519.PrivKeyEd25519 {
	[64]byte{223,215,229,180,200,255,139,128,87,60,254,209,46,229,221,210,209,210,226,195,17,40,69,13,207,72,28,73,81,196,152,35,172,16,174,28,34,178,225,143,206,250,7,55,242,70,251,188,198,236,95,140,22,219,255,118,135,191,218,23,249,99,101,6},
	[64]byte{200,158,4,178,74,149,16,125,49,160,166,236,43,157,85,138,209,61,135,184,11,94,66,22,17,29,166,27,184,125,185,43,13,127,165,49,42,127,252,95,122,105,53,58,34,102,154,26,72,2,150,106,98,98,5,187,163,75,171,181,167,102,62,227},
	[64]byte{48,131,120,96,178,66,88,8,2,248,200,39,147,168,117,66,165,56,94,245,33,176,207,206,102,33,170,85,87,203,143,40,67,167,238,146,174,211,15,149,60,204,123,53,150,0,184,42,197,43,167,236,104,211,205,122,223,9,94,134,45,98,157,162},
	[64]byte{221,65,21,10,37,155,27,249,28,126,116,63,105,228,81,116,211,121,182,207,208,77,93,212,41,213,19,45,70,54,140,171,36,129,236,222,44,16,30,222,82,62,118,168,73,66,110,24,29,177,28,57,117,80,178,202,208,124,36,213,18,242,174,210},
	[64]byte{96,107,233,217,67,226,33,171,232,230,209,47,210,43,191,189,90,126,100,222,206,85,224,114,200,245,62,239,164,43,216,142,66,143,43,51,7,51,20,26,171,123,236,161,52,9,188,27,64,117,22,97,240,201,72,1,110,6,193,10,73,5,53,73},
	[64]byte{188,75,146,48,162,78,211,30,236,12,119,203,200,171,212,58,97,52,194,71,179,189,239,145,188,234,219,11,163,240,45,49,131,171,161,18,47,199,255,67,66,236,163,64,163,24,25,47,190,198,64,215,161,101,67,40,118,149,54,28,238,29,53,174},
	[64]byte{92,206,248,242,6,150,45,189,15,223,132,130,152,207,39,109,93,126,37,94,224,173,79,70,69,13,92,144,141,152,22,45,62,200,76,68,186,105,45,148,183,200,92,170,42,189,99,85,204,58,17,211,221,5,23,105,196,94,9,24,40,187,144,182},
	[64]byte{218,18,255,188,152,182,15,60,32,144,210,60,23,141,97,210,205,29,238,179,14,47,62,140,214,99,6,104,120,13,142,106,128,94,27,219,115,242,196,112,78,138,20,115,62,204,99,228,246,252,35,246,146,161,60,127,103,9,44,125,111,165,180,158},
	[64]byte{180,58,147,27,85,218,159,78,196,255,122,62,247,123,248,11,160,43,63,89,177,50,155,139,211,46,157,209,15,234,165,137,57,223,227,189,27,221,0,80,249,162,147,167,152,62,158,157,225,224,229,122,18,56,160,138,75,159,22,234,138,210,254,120},
	[64]byte{237,59,33,88,47,210,114,111,5,249,200,49,243,254,156,123,118,70,92,30,204,162,25,152,110,51,85,77,75,106,39,199,79,13,165,131,70,76,32,229,202,108,255,131,203,30,185,176,181,156,117,102,112,209,70,3,133,182,67,197,214,35,40,18},
}

type accsign struct {
	privkey ed25519.PrivKeyEd25519	`json:"privkey"`
	nonce int64						`json:"nonce"`
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