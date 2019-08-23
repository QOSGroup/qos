package mapper

import (
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	bmapper "github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
)

func defaultContext() context.Context {
	mapperMap := make(map[string]bmapper.IMapper)

	cdc := baseabci.MakeQBaseCodec()
	qtypes.RegisterCodec(cdc)

	accountMapper := bacc.NewAccountMapper(nil, qtypes.ProtoQOSAccount)
	accountMapper.SetCodec(cdc)
	acountKey := accountMapper.GetStoreKey()
	mapperMap[bacc.AccountMapperName] = accountMapper

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(acountKey, btypes.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := context.NewContext(cms, abci.Header{}, false, log.NewNopLogger(), mapperMap)
	return ctx
}

func genTestAccount(addr btypes.Address, qos, qsc int64) qtypes.QOSAccount {
	return qtypes.QOSAccount{
		BaseAccount: bacc.BaseAccount{
			AccountAddress: addr,
			Publickey:      nil,
			Nonce:          0,
		},
		QOS: btypes.NewInt(qos),
		QSCs: qtypes.QSCs{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(qsc),
			},
		},
	}
}

func TestAccountInvariant(t *testing.T) {
	ctx := defaultContext()
	approveMapper := GetMapper(ctx)

	addr1 := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
	addr2 := btypes.Address(ed25519.GenPrivKey().PubKey().Address())

	account := genTestAccount(addr1, 100, 100)
	approveMapper.SetAccount(&account)
	_, coins, broken := AccountInvariant("bank")(ctx)
	require.False(t, broken)
	require.Equal(t, coins.AmountOf(qtypes.QOSCoinName), account.QOS)
	require.Equal(t, coins.AmountOf("qstar"), account.QSCs.AmountOf("qstar"))

	account = genTestAccount(addr1, -100, 100)
	approveMapper.SetAccount(&account)
	_, coins, broken = AccountInvariant("bank")(ctx)
	require.True(t, broken)
	require.Equal(t, coins.AmountOf(qtypes.QOSCoinName), account.QOS)
	require.Equal(t, coins.AmountOf("qstar"), account.QSCs.AmountOf("qstar"))

	account = genTestAccount(addr1, 100, 100)
	approveMapper.SetAccount(&account)
	account = genTestAccount(addr2, 100, 100)
	approveMapper.SetAccount(&account)
	_, coins, broken = AccountInvariant("bank")(ctx)
	require.False(t, broken)
	require.Equal(t, coins.AmountOf(qtypes.QOSCoinName), btypes.NewInt(200))
	require.Equal(t, coins.AmountOf("qstar"), btypes.NewInt(200))

	account = genTestAccount(addr1, 100, 100)
	approveMapper.SetAccount(&account)
	account = genTestAccount(addr2, -100, 100)
	approveMapper.SetAccount(&account)
	_, coins, broken = AccountInvariant("bank")(ctx)
	require.True(t, broken)
	require.True(t, coins.AmountOf(qtypes.QOSCoinName).IsZero())
	require.Equal(t, coins.AmountOf("qstar"), btypes.NewInt(200))
}
