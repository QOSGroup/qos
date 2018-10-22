package txs

import (
	baccount "github.com/QOSGroup/qbase/account"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

// todo: 依赖ctx中store操作，稍后更新(暂模拟)
func GetAccount(addr btypes.Address) (acc *account.QOSAccount) {
	//accmapper := baccount.NewAccountMapper(account.ProtoQOSAccount)
	//addrKey := baccount.AddressStoreKey(addr)
	//if !accmapper.Get(addrKey, &acc) {
	//	return nil
	//}
	//
	//return

	baseacc := baccount.BaseAccount{
		addr,
		ed25519.GenPrivKey().PubKey(),
		uint64(2),
	}

	qscList := []*types.QSC{
		{"qsc1", btypes.NewInt(100)},
		{"qsc2", btypes.NewInt(200)},
		{"qsc3", btypes.NewInt(100)},
		{"qsc4", btypes.NewInt(200)},
		{"qsc5", btypes.NewInt(100)},
	}

	acc = &account.QOSAccount{
		baseacc,
		btypes.NewInt(80000),
		qscList,
	}

	return
}

// todo: 暂模拟，稍后实现
func CreateAccount(addr btypes.Address) (acc *account.QOSAccount) {
	acc = GetAccount(addr)
	//if acc != nil {
	//	err = errors.New("Error; address()  account existed. ")
	//	return acc, err
	//}

	//accmapper := baccount.NewAccountMapper(account.ProtoQOSAccount)
	//addrKey := baccount.AddressStoreKey(addr)

	return acc
}

// todo: 暂模拟
func GetBanker(qscname string) (ret *account.QOSAccount) {
	baseacc := baccount.BaseAccount{[]byte("baseaccount1"),
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
