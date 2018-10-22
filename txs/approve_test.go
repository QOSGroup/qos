package txs

import (
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	bmapper "github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/mapper"
	"github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
)

var txApproveCreate ApproveCreateTx
var txApproveIncrease ApproveIncreaseTx
var txApproveDecrease ApproveDecreaseTx
var txApproveUse ApproveUseTx
var txApproveCancel ApproveCancelTx
var fromAccount = &account.QOSAccount{}
var toAccount = &account.QOSAccount{}

func approveTestInit() {
	fromPub := ed25519.GenPrivKey().PubKey()
	fromAddr := btypes.Address(fromPub.Address())
	toPub := ed25519.GenPrivKey().PubKey()
	toAddr := btypes.Address(toPub.Address())
	fromAccount.Qos = btypes.NewInt(100)
	fromAccount.QscList = []*types.QSC{
		{
			Name:   "qstar",
			Amount: btypes.NewInt(100),
		},
	}
	fromAccount.BaseAccount = bacc.BaseAccount{
		AccountAddress: fromAddr,
		Publickey:      fromPub,
		Nonce:          0,
	}
	toAccount.Qos = btypes.NewInt(0)
	toAccount.BaseAccount = bacc.BaseAccount{
		AccountAddress: toAddr,
		Publickey:      toPub,
		Nonce:          0,
	}
	approve1 := types.Approve{
		From: fromAddr,
		To:   toAddr,
		Qos:  btypes.NewInt(100),
		QscList: []*types.QSC{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	}
	approve2 := types.Approve{
		From: fromAddr,
		To:   toAddr,
		Qos:  btypes.NewInt(100),
		QscList: []*types.QSC{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	}
	approveCancel := types.ApproveCancel{
		From: fromAddr,
		To:   toAddr,
	}
	txApproveCreate = ApproveCreateTx{
		approve1,
	}
	txApproveIncrease = ApproveIncreaseTx{
		approve2,
	}
	txApproveDecrease = ApproveDecreaseTx{
		approve2,
	}
	txApproveUse = ApproveUseTx{
		approve2,
	}
	txApproveCancel = ApproveCancelTx{
		approveCancel,
	}
}

func txApproveTestContext() context.Context {
	approveTestInit()
	cdc := makeCodec()

	mapperMap := make(map[string]bmapper.IMapper)

	approveMapper := mapper.NewApproveMapper()
	approveMapper.SetCodec(cdc)
	approveKey := approveMapper.GetStoreKey()
	mapperMap[approveMapper.Name()] = approveMapper

	accountMapper := bacc.NewAccountMapper(account.ProtoQOSAccount)
	accountMapper.SetCodec(cdc)
	acountKey := accountMapper.GetStoreKey()
	mapperMap[accountMapper.Name()] = accountMapper

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(approveKey, store.StoreTypeIAVL, db)
	cms.MountStoreWithDB(acountKey, store.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := context.NewContext(cms, abci.Header{}, false, log.NewNopLogger(), mapperMap)
	return ctx
}

func TestTxApproveCreate_ValidateData(t *testing.T) {
	approveTestInit()

	ctx := txApproveTestContext()
	require.True(t, txApproveCreate.ValidateData(ctx))

	approveMapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
	err := approveMapper.SaveApprove(txApproveCreate.Approve)
	require.Nil(t, err)

	require.False(t, txApproveCreate.ValidateData(ctx))
}

func TestTxApproveCreate_Exec(t *testing.T) {
	approveTestInit()

	ctx := txApproveTestContext()
	result, cross := txApproveCreate.Exec(ctx)
	require.Nil(t, cross)
	require.Equal(t, result.Code, btypes.ABCICodeOK)

	approveMapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
	approve, exists := approveMapper.GetApprove(txApproveCreate.From, txApproveCreate.To)
	require.True(t, exists)
	require.True(t, txApproveCreate.Approve.Equals(approve))
}

func TestTxApproveIncrease_ValidateData(t *testing.T) {
	approveTestInit()

	ctx := txApproveTestContext()
	require.False(t, txApproveIncrease.ValidateData(ctx))

	approveMapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
	err := approveMapper.SaveApprove(txApproveCreate.Approve)
	require.Nil(t, err)

	require.True(t, txApproveIncrease.ValidateData(ctx))
}

func TestTxApproveIncrease_Exec(t *testing.T) {
	approveTestInit()

	ctx := txApproveTestContext()

	approveMapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
	err := approveMapper.SaveApprove(txApproveCreate.Approve)
	require.Nil(t, err)

	result, cross := txApproveIncrease.Exec(ctx)
	require.Nil(t, cross)
	require.Equal(t, result.Code, btypes.ABCICodeOK)

	approve, exists := approveMapper.GetApprove(txApproveCreate.From, txApproveCreate.To)
	require.True(t, exists)
	require.True(t, txApproveCreate.Approve.Plus(txApproveIncrease.Approve).Equals(approve))
}

func TestTxApproveDecrease_ValidateData(t *testing.T) {
	approveTestInit()

	ctx := txApproveTestContext()
	require.False(t, txApproveDecrease.ValidateData(ctx))

	approveMapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
	err := approveMapper.SaveApprove(txApproveCreate.Approve)
	require.Nil(t, err)

	require.True(t, txApproveDecrease.ValidateData(ctx))

	txApproveDecrease.Qos = btypes.NewInt(100)
	require.True(t, txApproveDecrease.ValidateData(ctx))

	txApproveDecrease.Qos = btypes.NewInt(110)
	require.False(t, txApproveDecrease.ValidateData(ctx))
}

func TestTxApproveDecrease_Exec(t *testing.T) {
	approveTestInit()

	ctx := txApproveTestContext()

	approveMapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
	err := approveMapper.SaveApprove(txApproveCreate.Approve)
	require.Nil(t, err)

	result, cross := txApproveDecrease.Exec(ctx)
	require.Nil(t, cross)
	require.Equal(t, result.Code, btypes.ABCICodeOK)

	approve, exists := approveMapper.GetApprove(txApproveCreate.From, txApproveCreate.To)
	require.True(t, exists)
	require.True(t, txApproveCreate.Approve.Minus(txApproveDecrease.Approve).Equals(approve))
}

func TestTxApproveUse_ValidateData(t *testing.T) {
	approveTestInit()

	ctx := txApproveTestContext()
	require.False(t, txApproveUse.ValidateData(ctx))

	approveMapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
	err := approveMapper.SaveApprove(txApproveCreate.Approve)
	require.Nil(t, err)
	require.False(t, txApproveUse.ValidateData(ctx))

	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	accountMapper.SetAccount(fromAccount)
	accountMapper.SetAccount(toAccount)

	require.True(t, txApproveUse.ValidateData(ctx))

	txApproveUse.Qos = btypes.NewInt(110)
	require.False(t, txApproveUse.ValidateData(ctx))

}

func TestTxApproveUse_GetSigner(t *testing.T) {
	approveTestInit()
	require.Equal(t, txApproveUse.GetSigner(), []btypes.Address{txApproveUse.To})
}

func TestTxApproveUse_GetGasPayer(t *testing.T) {
	approveTestInit()
	require.Equal(t, txApproveUse.GetGasPayer(), txApproveUse.To)
}

func TestTxApproveUse_Exec(t *testing.T) {
	approveTestInit()

	ctx := txApproveTestContext()
	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	accountMapper.SetAccount(fromAccount)
	accountMapper.SetAccount(toAccount)

	result, cross := txApproveUse.Exec(ctx)
	require.Nil(t, cross)
	require.NotEqual(t, result.Code, btypes.ABCICodeOK)

	txApproveCreate.Qos = btypes.NewInt(1)
	approveMapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
	err := approveMapper.SaveApprove(txApproveCreate.Approve)
	require.Nil(t, err)

	result, cross = txApproveUse.Exec(ctx)
	require.Nil(t, cross)
	require.NotEqual(t, result.Code, btypes.ABCICodeOK)

	txApproveCreate.Qos = btypes.NewInt(100)
	err = approveMapper.SaveApprove(txApproveCreate.Approve)
	require.Nil(t, err)

	result, cross = txApproveUse.Exec(ctx)
	require.Nil(t, cross)
	require.Equal(t, result.Code, btypes.ABCICodeOK)

	approve, exists := approveMapper.GetApprove(txApproveUse.From, txApproveUse.To)
	require.True(t, exists)
	require.True(t, txApproveCreate.Minus(txApproveUse.Approve).Equals(approve))

}

func TestTxApproveCancel_ValidateData(t *testing.T) {
	approveTestInit()

	ctx := txApproveTestContext()
	require.False(t, txApproveCancel.ValidateData(ctx))

	mapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
	err := mapper.SaveApprove(txApproveCreate.Approve)
	require.Nil(t, err)

	require.True(t, txApproveCancel.ValidateData(ctx))
}

func TestTxApproveCancel_Exec(t *testing.T) {
	approveTestInit()

	ctx := txApproveTestContext()
	result, cross := txApproveCancel.Exec(ctx)
	require.Nil(t, cross)
	require.NotEqual(t, result.Code, btypes.ABCICodeOK)

	mapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
	err := mapper.SaveApprove(txApproveCreate.Approve)
	require.Nil(t, err)

	result, _ = txApproveCancel.Exec(ctx)
	require.Equal(t, result.Code, btypes.ABCICodeOK)

}
