package account

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/store"
	"testing"

	"github.com/QOSGroup/qbase/account"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)
func keyPubAddr() (crypto.PrivKey, crypto.PubKey, btypes.Address) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := btypes.Address(pub.Address())
	return key, pub, addr
}

func genNewAccount() (qosAccount QOSAccount){
	_, pub, addr := keyPubAddr()
	coinList := []*types.QSC{
		types.NewQSC("QSC1", btypes.NewInt(1234)),
		types.NewQSC("QSC2", btypes.NewInt(5678)),
	}
	qosAccount = QOSAccount{
		account.BaseAccount{addr, pub, 0},
		btypes.NewInt(5380394853),
		coinList,
	}
	return
}

func TestQOSAccountEditing(t *testing.T) {
	qsc := types.NewQSC("QSC1", btypes.NewInt(1234))
	qosAccount := genNewAccount()
	//test getter
	qsc1 := qosAccount.GetQSC("QSC1")
	require.Equal(t, qsc, qsc1)

	//modify coin and test setter
	qsc1.SetAmount(btypes.NewInt(4321))
	err := qosAccount.SetQSC(qsc1)
	require.Nil(t, err)
	qsc1 = qosAccount.GetQSC("QSC1")
	require.NotEqual(t, qsc, qsc1)

	//test remove
	err = qosAccount.RemoveQSCByName("QSC2")
	require.Nil(t, err)
	qsc1 = qosAccount.GetQSC("QSC2")
	require.Nil(t, qsc1)
}

func TestAccountMarshal(t *testing.T) {
	qosAccount := genNewAccount()

	qosAccount_json, err := cdc.MarshalJSON(qosAccount)
	require.Nil(t, err)

	another_qosAdd := QOSAccount{}
	err = cdc.UnmarshalJSON(qosAccount_json, &another_qosAdd)
	require.Nil(t, err)
	require.Equal(t, qosAccount, another_qosAdd)

	qosAccount_binary, err := cdc.MarshalBinary(qosAccount)
	require.Nil(t, err)

	another_qosAdd = QOSAccount{}
	err = cdc.UnmarshalBinary(qosAccount_binary, &another_qosAdd)
	require.Nil(t, err)
	require.Equal(t, qosAccount, another_qosAdd)

	another_qosAdd = QOSAccount{}
	another_json := []byte{}
	err = cdc.UnmarshalBinary(qosAccount_binary[:len(qosAccount_binary)/2], &another_json)
	require.NotNil(t, err)

}

func defaultContext(key store.StoreKey) context.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, store.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := context.NewContext(cms, abci.Header{}, false, log.NewNopLogger())
	return ctx
}

func TestAccountMapperGetSet(t *testing.T) {
	mapper := account.NewAccountMapper(cdc, ProtoQOSAccount)
	ctx := defaultContext(mapper.GetStoreKey())

	for i:=0; i < 1; i++ {
		_, pubkey, addr := keyPubAddr()

		// 没有存过该addr，取出来应为nil
		acc := mapper.GetAccount(ctx, addr)
		require.Nil(t, acc)

		acc = mapper.NewAccountWithAddress(ctx, addr).(*QOSAccount)
		require.NotNil(t, acc)
		require.Equal(t, addr, acc.GetAddress())
		require.EqualValues(t, nil, acc.GetPubicKey())
		require.EqualValues(t, 0, acc.GetNonce())

		// 新的account尚未存储，依然取出nil
		require.Nil(t, mapper.GetAccount(ctx, addr))

		nonce := uint64(20)
		acc.SetNonce(nonce)
		acc.SetPublicKey(pubkey)
		acc.(*QOSAccount).SetQOS(btypes.NewInt(100))
		acc.(*QOSAccount).SetQSC(types.NewQSC("QSC1", btypes.NewInt(1234)))
		acc.(*QOSAccount).SetQSC(types.NewQSC("QSC2", btypes.NewInt(5678)))
		// 存储account
		mapper.SetAccount(ctx, acc)

		// 将account以地址取出并验证
		acc = mapper.GetAccount(ctx, addr).(*QOSAccount)
		require.NotNil(t, acc)
		require.Equal(t, nonce, acc.GetNonce())

	}
	//批量处理特定前缀存储的账户
	mapper.IterateAccounts(ctx, func(acc account.Account) bool{
		bz := mapper.EncodeAccount(acc)
		acc1 := mapper.DecodeAccount(bz)
		require.Equal(t, acc, acc1)
		return false
	})
}
