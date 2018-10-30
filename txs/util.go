package txs

import (
	baccount "github.com/QOSGroup/qbase/account"
	bcontext "github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/account"
	qostest "github.com/QOSGroup/qos/test"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

//CA结构体
//todo: CA具体格式确定后会更改
type CA struct {
	Qcpname   string        `json:"qcpname"`
	Banker    bool          `json:"banker"`
	Pubkey    crypto.PubKey `json:"pubkey"`
	Signature []byte        `json:"signature"`
	Info      string        `json:"info"`
}

// 通过地址获取QOSAccount
func GetAccount(ctx bcontext.Context, addr btypes.Address) (acc *account.QOSAccount) {
	mapper := ctx.Mapper(baccount.AccountMapperName).(*baccount.AccountMapper)
	if mapper == nil {
		return nil
	}

	accbase := mapper.GetAccount(addr)
	if accbase == nil{
		return nil
	}
	acc = accbase.(*account.QOSAccount)

	return
}

// 通过地址创建QOSAccount
// 若账户存在，返回账户 & false
func CreateAndSaveAccount(ctx bcontext.Context, addr btypes.Address) (acc *account.QOSAccount, success bool) {
	mapper := ctx.Mapper(baccount.AccountMapperName).(*baccount.AccountMapper)
	if mapper == nil {
		return nil, false
	}

	accfind := mapper.GetAccount(addr).(*account.QOSAccount)
	if accfind != nil {
		return accfind, false
	}

	acc = mapper.NewAccountWithAddress(addr).(*account.QOSAccount)
	mapper.SetAccount(acc)

	return acc, true
}

func SaveAccount(ctx bcontext.Context, acc *account.QOSAccount) bool {
	mapper := ctx.Mapper(baccount.AccountMapperName).(*baccount.AccountMapper)
	if mapper == nil {
		return false
	}
	mapper.SetAccount(acc)

	return true
}

// todo: 暂模拟
var rootprivkey = ed25519.GenPrivKey()

func FetchQscCA() (caQsc *[]byte) {
	//pubkey := ed25519.GenPrivKey().PubKey()
	accandkey := qostest.InitKeys(cdc)

	ca := &CA{
		"qsc1",
		false,
		accandkey[0].PrivKey.PubKey(),
		[]byte{},
		"qsc ca data",
	}

	signdata := getsignCA(ca)
	ca.Signature, _ = rootprivkey.Sign(signdata)
	va, _ := cdc.MarshalBinaryBare(ca)
	caQsc = &va

	return
}

// todo: 暂模拟
func FetchBankerCA() (caBanker *[]byte) {
	accandkey := qostest.InitKeys(cdc)

	ca := &CA{
		"qsc1",
		true,
		accandkey[1].PrivKey.PubKey(),
		[]byte{},
		"qsc ca of banker",
	}

	signdata := getsignCA(ca)
	ca.Signature, _ = rootprivkey.Sign(signdata)
	va, _ := cdc.MarshalBinaryBare(ca)
	caBanker = &va

	return
}

// todo: 暂模拟
func getsignCA(ca *CA) (signdata []byte) {
	signdata = append(signdata, []byte(ca.Qcpname)...)
	signdata = append(signdata, btypes.Bool2Byte(ca.Banker)...)
	signdata = append(signdata, ca.Pubkey.Bytes()...)
	signdata = append(signdata, []byte(ca.Info)...)

	return
}

// todo: 暂模拟
func VerifyCA(pubkey crypto.PubKey, ca *CA) bool {
	if ca == nil || len(pubkey.Bytes()) == 0 {
		return false
	}

	return pubkey.VerifyBytes(getsignCA(ca), ca.Signature)
}

//func MakeTxStd(tx btxs.ITx, chainid string, maxgas int64) (txstd *btxs.TxStd) {
//	txstd = btxs.NewTxStd(tx, chainid, btypes.NewInt(maxgas))
//	signer := txstd.ITx.GetSigner()
//
//	// no signer, no signature
//	if signer == nil {
//		txstd.Signature = []btxs.Signature{}
//		return
//	}
//
//	// accmapper := baccount.NewAccountMapper(baccount.ProtoBaseAccount)
//	accmapper := baccount.NewAccountMapper(nil, account.ProtoQOSAccount)
//
//	// 填充 txstd.Signature[]
//	for _, sg := range signer {
//		prvKey := ed25519.GenPrivKey()
//		nonce, err := accmapper.GetNonce(baccount.AddressStoreKey(sg))
//		if err != nil {
//			return nil
//		}
//
//		signbyte, errsign := txstd.SignTx(prvKey, int64(nonce))
//		if signbyte == nil || errsign != nil {
//			return nil
//		}
//
//		signdata := btxs.Signature{
//			prvKey.PubKey(),
//			signbyte,
//			int64(nonce),
//		}
//
//		txstd.Signature = append(txstd.Signature, signdata)
//	}
//
//	return
//}
