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

var approveCreateTx ApproveCreateTx
var approveIncreaseTx ApproveIncreaseTx
var approveDecreaseTx ApproveDecreaseTx
var approveUseTx ApproveUseTx
var approveCancelTx ApproveCancelTx
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
	qos1 := btypes.NewInt(100)
	approve1 := types.NewApprove(fromAddr, toAddr, &qos1,
		[]*types.QSC{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	)
	qos2 := btypes.NewInt(100)
	approve2 := types.NewApprove(fromAddr, toAddr, &qos2,
		[]*types.QSC{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	)
	approveCancel := types.ApproveCancel{
		From: fromAddr,
		To:   toAddr,
	}
	approveCreateTx = ApproveCreateTx{
		approve1,
	}
	approveIncreaseTx = ApproveIncreaseTx{
		approve2,
	}
	approveDecreaseTx = ApproveDecreaseTx{
		approve2,
	}
	approveUseTx = ApproveUseTx{
		approve2,
	}
	approveCancelTx = ApproveCancelTx{
		approveCancel,
	}
}

func txApproveTestContext() context.Context {
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
	require.True(t, approveCreateTx.ValidateData(ctx))

	approveMapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
	err := approveMapper.SaveApprove(approveCreateTx.Approve)
	require.Nil(t, err)

	require.False(t, approveCreateTx.ValidateData(ctx))
}

func TestTxApproveCreate_Exec(t *testing.T) {
	approveTestInit()

	ctx := txApproveTestContext()
	result, cross := approveCreateTx.Exec(ctx)
	require.Nil(t, cross)
	require.Equal(t, result.Code, btypes.ABCICodeOK)

	approveMapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
	approve, exists := approveMapper.GetApprove(approveCreateTx.From, approveCreateTx.To)
	require.True(t, exists)
	require.True(t, approveCreateTx.Approve.Equals(approve))
}

func TestTxApproveIncrease_ValidateData(t *testing.T) {
	approveTestInit()

	ctx := txApproveTestContext()
	require.False(t, approveIncreaseTx.ValidateData(ctx))

	approveMapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
	err := approveMapper.SaveApprove(approveCreateTx.Approve)
	require.Nil(t, err)

	require.True(t, approveIncreaseTx.ValidateData(ctx))
}

func TestTxApproveIncrease_Exec(t *testing.T) {
	approveTestInit()

	ctx := txApproveTestContext()

	approveMapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
	err := approveMapper.SaveApprove(approveCreateTx.Approve)
	require.Nil(t, err)

	result, cross := approveIncreaseTx.Exec(ctx)
	require.Nil(t, cross)
	require.Equal(t, result.Code, btypes.ABCICodeOK)

	approve, exists := approveMapper.GetApprove(approveCreateTx.From, approveCreateTx.To)
	require.True(t, exists)
	require.True(t, approveCreateTx.Approve.Plus(approveIncreaseTx.Qos, approveIncreaseTx.QscList).Equals(approve))
}

func TestTxApproveDecrease_ValidateData(t *testing.T) {
	approveTestInit()

	ctx := txApproveTestContext()
	require.False(t, approveDecreaseTx.ValidateData(ctx))

	approveMapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
	err := approveMapper.SaveApprove(approveCreateTx.Approve)
	require.Nil(t, err)

	require.True(t, approveDecreaseTx.ValidateData(ctx))

	approveDecreaseTx.Qos = btypes.NewInt(100)
	require.True(t, approveDecreaseTx.ValidateData(ctx))

	approveDecreaseTx.Qos = btypes.NewInt(110)
	require.False(t, approveDecreaseTx.ValidateData(ctx))
}

func TestTxApproveDecrease_Exec(t *testing.T) {
	approveTestInit()

	ctx := txApproveTestContext()

	approveMapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
	err := approveMapper.SaveApprove(approveCreateTx.Approve)
	require.Nil(t, err)

	result, cross := approveDecreaseTx.Exec(ctx)
	require.Nil(t, cross)
	require.Equal(t, result.Code, btypes.ABCICodeOK)

	approve, exists := approveMapper.GetApprove(approveCreateTx.From, approveCreateTx.To)
	require.True(t, exists)
	require.True(t, approveCreateTx.Approve.Minus(approveDecreaseTx.Qos, approveDecreaseTx.QscList).Equals(approve))
}

func TestTxApproveUse_ValidateData(t *testing.T) {
	approveTestInit()

	ctx := txApproveTestContext()
	require.False(t, approveUseTx.ValidateData(ctx))

	approveMapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
	err := approveMapper.SaveApprove(approveCreateTx.Approve)
	require.Nil(t, err)
	require.False(t, approveUseTx.ValidateData(ctx))

	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	accountMapper.SetAccount(fromAccount)
	accountMapper.SetAccount(toAccount)

	require.True(t, approveUseTx.ValidateData(ctx))

	approveUseTx.Qos = btypes.NewInt(110)
	require.False(t, approveUseTx.ValidateData(ctx))

}

func TestTxApproveUse_GetSigner(t *testing.T) {
	approveTestInit()
	require.Equal(t, approveUseTx.GetSigner(), []btypes.Address{approveUseTx.To})
}

func TestTxApproveUse_GetGasPayer(t *testing.T) {
	approveTestInit()
	require.Equal(t, approveUseTx.GetGasPayer(), approveUseTx.To)
}

func TestTxApproveUse_Exec(t *testing.T) {
	approveTestInit()

	ctx := txApproveTestContext()
	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	accountMapper.SetAccount(fromAccount)
	accountMapper.SetAccount(toAccount)

	result, cross := approveUseTx.Exec(ctx)
	require.Nil(t, cross)
	require.NotEqual(t, result.Code, btypes.ABCICodeOK)

	approveCreateTx.Qos = btypes.NewInt(1)
	approveMapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
	err := approveMapper.SaveApprove(approveCreateTx.Approve)
	require.Nil(t, err)

	result, cross = approveUseTx.Exec(ctx)
	require.Nil(t, cross)
	require.NotEqual(t, result.Code, btypes.ABCICodeOK)

	approveCreateTx.Qos = btypes.NewInt(100)
	err = approveMapper.SaveApprove(approveCreateTx.Approve)
	require.Nil(t, err)

	result, cross = approveUseTx.Exec(ctx)
	require.Nil(t, cross)
	require.Equal(t, result.Code, btypes.ABCICodeOK)

	approve, exists := approveMapper.GetApprove(approveUseTx.From, approveUseTx.To)
	require.True(t, exists)
	require.True(t, approveCreateTx.Minus(approveUseTx.Qos, approveUseTx.QscList).Equals(approve))

}

func TestTxApproveCancel_ValidateData(t *testing.T) {
	approveTestInit()

	ctx := txApproveTestContext()
	require.False(t, approveCancelTx.ValidateData(ctx))

	mapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
	err := mapper.SaveApprove(approveCreateTx.Approve)
	require.Nil(t, err)

	require.True(t, approveCancelTx.ValidateData(ctx))
}

func TestTxApproveCancel_Exec(t *testing.T) {
	approveTestInit()

	ctx := txApproveTestContext()
	result, cross := approveCancelTx.Exec(ctx)
	require.Nil(t, cross)
	require.NotEqual(t, result.Code, btypes.ABCICodeOK)

	mapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
	err := mapper.SaveApprove(approveCreateTx.Approve)
	require.Nil(t, err)

	result, _ = approveCancelTx.Exec(ctx)
	require.Equal(t, result.Code, btypes.ABCICodeOK)

}
