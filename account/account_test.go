package account

import (
	"testing"

	"github.com/QOSGroup/qbase/account"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/types"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/store"
	"github.com/QOSGroup/qbase/mapper"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	abci "github.com/tendermint/tendermint/abci/types"
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
	coinList := types.QSCs{
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

func defaultContext(key store.StoreKey, mapperMap map[string]mapper.IMapper) context.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, store.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := context.NewContext(cms, abci.Header{}, false, log.NewNopLogger(), mapperMap)
	return ctx
}
func TestAccountMapperGetSet(t *testing.T) {
	seedMapper := account.NewAccountMapper(ProtoQOSAccount)
	seedMapper.SetCodec(cdc)

	mapperMap := make(map[string]mapper.IMapper)
	mapperMap[seedMapper.Name()] = seedMapper

	ctx := defaultContext(seedMapper.GetStoreKey(), mapperMap)

	mapper, _ := ctx.Mapper(account.AccountMapperName).(*account.AccountMapper)
	for i:=0; i < 100; i++ {
		_, pubkey, addr := keyPubAddr()

		// 没有存过该addr，取出来应为nil
		acc := mapper.GetAccount(addr)
		require.Nil(t, acc)

		qosacc := mapper.NewAccountWithAddress(addr).(*QOSAccount)
		require.NotNil(t, qosacc)
		require.Equal(t, addr, qosacc.GetAddress())
		require.EqualValues(t, nil, qosacc.GetPubicKey())
		require.EqualValues(t, 0, qosacc.GetNonce())

		// 新的account尚未存储，依然取出nil
		require.Nil(t, mapper.GetAccount(addr))

		nonce := int64(20)
		qosacc.SetNonce(nonce)
		qosacc.SetPublicKey(pubkey)
		qosacc.SetQOS(btypes.NewInt(100))
		qosacc.SetQSC(types.NewQSC("QSC1", btypes.NewInt(1234)))
		qosacc.SetQSC(types.NewQSC("QSC2", btypes.NewInt(5678)))
		// 存储account
		mapper.SetAccount(qosacc)

		// 将account以地址取出并验证
		qosacc = mapper.GetAccount(addr).(*QOSAccount)
		require.NotNil(t, qosacc)
		require.Equal(t, nonce, qosacc.GetNonce())

	}
	//批量处理特定前缀存储的账户
	mapper.IterateAccounts(func(acc account.Account) bool{
		bz := mapper.EncodeAccount(acc)
		acc1 := mapper.DecodeAccount(bz)
		require.Equal(t, acc, acc1)
		return false
	})
}
