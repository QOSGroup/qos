package txs

import (
	baccount "github.com/QOSGroup/qbase/account"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	bcontext "github.com/QOSGroup/qbase/context"
	btxs "github.com/QOSGroup/qbase/txs"
)

// 通过地址获取QOSAccount
func GetAccount(ctx bcontext.Context, addr btypes.Address) (acc *account.QOSAccount) {
	mapper := ctx.Mapper(baccount.AccountMapperName).(*baccount.AccountMapper)
	if mapper == nil {
		return nil
	}

	acc = mapper.GetAccount(addr).(*account.QOSAccount)
	if acc == nil {
		return nil
	}

	return
}

// 通过地址创建QOSAccount
// 若账户存在，返回账户 & false
func CreateAccount(ctx bcontext.Context, addr btypes.Address) (acc *account.QOSAccount, success bool) {
	mapper := ctx.Mapper(baccount.AccountMapperName).(*baccount.AccountMapper)
	if mapper == nil {
		return nil, false
	}

	accfind := mapper.GetAccount(addr).(*account.QOSAccount)
	if accfind != nil {
		return accfind, false
	}

	acc = mapper.NewAccountWithAddress(addr).(*account.QOSAccount)

	return acc,true
}

// todo: 暂模拟
func GetBanker(qscname string) (ret *account.QOSAccount) {
	baseacc := baccount.BaseAccount{
		[]byte("baseaccount1"),
		ed25519.GenPrivKey().PubKey(),
		uint64(3),
	}

	qscList := []*types.QSC{
		&types.QSC{"qsc1", btypes.NewInt(100),},
		&types.QSC{"qsc2", btypes.NewInt(200),},
	}
	ret = &account.QOSAccount{baseacc,
		btypes.NewInt(10000),
		qscList,
	}

	return
}

func FetchQscCA() (caQsc *[]byte) {
	pubkey := ed25519.GenPrivKey().PubKey()

	ca := &CA{
		"qsc1",
		false,
		pubkey,
		"qsc ca data",
	}

	va, _ := cdc.MarshalBinaryBare(ca)
	caQsc = &va

	return
}
func FetchBankerCA() (caBanker *[]byte) {
	pubkey := ed25519.GenPrivKey().PubKey()

	ca := &CA{

		"qsc1",
		true,
		pubkey,
		"qsc ca of banker",
	}

	va, _ := cdc.MarshalBinaryBare(ca)
	caBanker = &va

	return
}

func MakeTxStd(tx btxs.ITx, chainid string, maxgas int64) (txstd *btxs.TxStd) {
	txstd = btxs.NewTxStd(tx, chainid, btypes.NewInt(maxgas))
	signer := txstd.ITx.GetSigner()

	// no signer, no signature
	if signer == nil {
		txstd.Signature = []btxs.Signature{}
		return
	}

	// accmapper := baccount.NewAccountMapper(baccount.ProtoBaseAccount)
	accmapper := baccount.NewAccountMapper(account.ProtoQOSAccount)

	// 填充 txstd.Signature[]
	for _, sg := range signer {
		prvKey := ed25519.GenPrivKey()
		nonce, err := accmapper.GetNonce(baccount.AddressStoreKey(sg))
		if err != nil {
			return nil
		}

		signbyte, errsign := txstd.SignTx(prvKey, int64(nonce))
		if signbyte == nil || errsign != nil {
			return nil
		}

		signdata := btxs.Signature{
			prvKey.PubKey(),
			signbyte,
			int64(nonce),
		}

		txstd.Signature = append(txstd.Signature, signdata)
	}

	return
}
