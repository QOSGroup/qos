package transfer

import (
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	bmapper "github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	transfertypes "github.com/QOSGroup/qos/module/transfer/types"
	"github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
)

func txTransferTestContext() context.Context {
	mapperMap := make(map[string]bmapper.IMapper)
	accountMapper := bacc.NewAccountMapper(nil, types.ProtoQOSAccount)
	accountMapper.SetCodec(cdc)
	acountKey := accountMapper.GetStoreKey()
	mapperMap[bacc.AccountMapperName] = accountMapper

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(acountKey, store.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := context.NewContext(cms, abci.Header{}, false, log.NewNopLogger(), mapperMap)
	return ctx
}

func TestTransferTx_ValidateData(t *testing.T) {
	ctx := txTransferTestContext()

	// 空
	tx := TxTransfer{
		Senders:   transfertypes.TransItems{},
		Receivers: transfertypes.TransItems{},
	}
	require.NotNil(t, tx.ValidateData(ctx))

	addr1 := ed25519.GenPrivKey().PubKey().Address().Bytes()
	addr2 := ed25519.GenPrivKey().PubKey().Address().Bytes()
	tx.Senders = append(tx.Senders, transfertypes.TransItem{addr1, btypes.NewInt(0), nil})
	tx.Receivers = append(tx.Receivers, transfertypes.TransItem{addr2, btypes.NewInt(0), nil})
	require.NotNil(t, tx.ValidateData(ctx))

	// 账户
	tx.Senders[0].QOS = btypes.NewInt(10)
	tx.Receivers[0].QOS = btypes.NewInt(10)
	require.NotNil(t, tx.ValidateData(ctx))
	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	accountMapper.SetAccount(accountMapper.NewAccountWithAddress(addr1))
	require.NotNil(t, tx.ValidateData(ctx))
	aac1 := accountMapper.GetAccount(addr1).(*types.QOSAccount)
	aac1.QOS = btypes.NewInt(100)
	accountMapper.SetAccount(aac1)
	require.Nil(t, tx.ValidateData(ctx))

	// 重复
	tx.Senders = append(tx.Senders, tx.Senders[0])
	require.NotNil(t, tx.ValidateData(ctx))

	// 金额
	tx.Senders = tx.Senders[0:1]
	tx.Receivers[0].QOS = btypes.NewInt(20)
	require.NotNil(t, tx.ValidateData(ctx))
	tx.Receivers[0].QOS = btypes.NewInt(10)
	require.Nil(t, tx.ValidateData(ctx))
}

func TestTransferTx_GetSigner(t *testing.T) {
	tx := TxTransfer{
		Senders: transfertypes.TransItems{
			{ed25519.GenPrivKey().PubKey().Address().Bytes(), btypes.NewInt(10), nil},
			{ed25519.GenPrivKey().PubKey().Address().Bytes(), btypes.NewInt(10), nil},
		},
		Receivers: transfertypes.TransItems{
			{ed25519.GenPrivKey().PubKey().Address().Bytes(), btypes.NewInt(20), nil},
		},
	}
	require.Equal(t, tx.GetSigner(), []btypes.Address{tx.Senders[0].Address, tx.Senders[1].Address})
}

func TestTransferTx_CalcGas(t *testing.T) {
	tx := TxTransfer{
		Senders: transfertypes.TransItems{
			{ed25519.GenPrivKey().PubKey().Address().Bytes(), btypes.NewInt(10), nil},
		},
		Receivers: transfertypes.TransItems{
			{ed25519.GenPrivKey().PubKey().Address().Bytes(), btypes.NewInt(10), nil},
		},
	}
	require.Equal(t, tx.CalcGas(), btypes.NewInt(0))
}

func TestTransferTx_GetGasPayer(t *testing.T) {
	tx := TxTransfer{
		Senders: transfertypes.TransItems{
			{ed25519.GenPrivKey().PubKey().Address().Bytes(), btypes.NewInt(10), nil},
			{ed25519.GenPrivKey().PubKey().Address().Bytes(), btypes.NewInt(10), nil},
		},
		Receivers: transfertypes.TransItems{
			{ed25519.GenPrivKey().PubKey().Address().Bytes(), btypes.NewInt(20), nil},
		},
	}
	require.Equal(t, tx.GetGasPayer(), tx.Senders[0].Address)
}

func TestTransferTx_GetSignData(t *testing.T) {
	tx := TxTransfer{
		Senders: transfertypes.TransItems{
			{ed25519.GenPrivKey().PubKey().Address().Bytes(), btypes.NewInt(10), nil},
			{ed25519.GenPrivKey().PubKey().Address().Bytes(), btypes.NewInt(10), nil},
		},
		Receivers: transfertypes.TransItems{
			{ed25519.GenPrivKey().PubKey().Address().Bytes(), btypes.NewInt(20), nil},
		},
	}

	ret := make([]byte, 0)
	for _, sender := range tx.Senders {
		ret = append(ret, sender.Address...)
		ret = append(ret, sender.QOS.String()...)
		ret = append(ret, sender.QSCs.String()...)
	}
	for _, receiver := range tx.Receivers {
		ret = append(ret, receiver.Address...)
		ret = append(ret, receiver.QOS.String()...)
		ret = append(ret, receiver.QSCs.String()...)
	}

	require.Equal(t, tx.GetSignData(), ret)
}
