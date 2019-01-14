package approve

import (
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	bmapper "github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	approvetype "github.com/QOSGroup/qos/module/approve/types"
	"github.com/QOSGroup/qos/module/qsc"
	qsctype "github.com/QOSGroup/qos/module/qsc/types"
	"github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
)

var testFromAddr = btypes.Address(ed25519.GenPrivKey().PubKey().Address())
var testToAddr = btypes.Address(ed25519.GenPrivKey().PubKey().Address())

func genTestApprove() approvetype.Approve {
	return approvetype.Approve{
		From: testFromAddr,
		To:   testToAddr,
		QOS:  btypes.NewInt(100),
		QSCs: types.QSCs{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	}
}

func TestValidateData(t *testing.T) {
	ctx := defaultContext()
	approve := genTestApprove()
	require.NotNil(t, validateData(ctx, approve))

	saveQSCInfo(ctx, "qstar")
	require.Nil(t, validateData(ctx, approve))

	approve.QSCs = append(approve.QSCs, &types.QSC{
		Name:   "qstar",
		Amount: btypes.NewInt(100),
	})
	require.NotNil(t, validateData(ctx, approve))

	approve.QSCs = types.QSCs{
		{
			Name:   "qos",
			Amount: btypes.NewInt(100),
		},
	}
	require.NotNil(t, validateData(ctx, approve))
}

func genTestAccount(addr btypes.Address) *types.QOSAccount {
	return &types.QOSAccount{
		BaseAccount: bacc.BaseAccount{
			AccountAddress: addr,
			Publickey:      nil,
			Nonce:          0,
		},
		QOS: btypes.NewInt(100),
		QSCs: types.QSCs{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	}
}

func genApproveCancelTx() TxCancelApprove {
	return TxCancelApprove{
		From: btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		To:   btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
	}
}

func defaultContext() context.Context {
	mapperMap := make(map[string]bmapper.IMapper)

	approveMapper := NewApproveMapper()
	approveMapper.SetCodec(cdc)
	approveKey := approveMapper.GetStoreKey()
	mapperMap[ApproveMapperName] = approveMapper

	accountMapper := bacc.NewAccountMapper(nil, types.ProtoQOSAccount)
	accountMapper.SetCodec(cdc)
	acountKey := accountMapper.GetStoreKey()
	mapperMap[bacc.AccountMapperName] = accountMapper

	qscMapper := qsc.NewQSCMapper()
	qscMapper.SetCodec(cdc)
	qscKey := qscMapper.GetStoreKey()
	mapperMap[qsc.QSCMapperName] = qscMapper

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(approveKey, store.StoreTypeIAVL, db)
	cms.MountStoreWithDB(acountKey, store.StoreTypeIAVL, db)
	cms.MountStoreWithDB(qscKey, store.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := context.NewContext(cms, abci.Header{}, false, log.NewNopLogger(), mapperMap)
	return ctx
}

func saveQSCInfo(ctx context.Context, qscName string) {
	qscMapper := ctx.Mapper(qsc.QSCMapperName).(*qsc.QSCMapper)
	qscMapper.SaveQsc(&qsctype.QSCInfo{
		Name: qscName,
	})
}

func defaultContextWithQSC() context.Context {
	ctx := defaultContext()
	saveQSCInfo(ctx, "qstar")
	return ctx
}

func TestTxApproveCreate_ValidateData(t *testing.T) {
	ctx := defaultContextWithQSC()

	tx := TxCreateApprove{
		genTestApprove(),
	}
	require.Nil(t, tx.ValidateData(ctx))

	approveMapper := ctx.Mapper(ApproveMapperName).(*ApproveMapper)
	approveMapper.SaveApprove(tx.Approve)

	require.NotNil(t, tx.ValidateData(ctx))
}

func TestTxApproveCreate_Exec(t *testing.T) {
	ctx := defaultContext()

	tx := TxCreateApprove{
		genTestApprove(),
	}
	result, cross := tx.Exec(ctx)
	require.Nil(t, cross)
	require.Equal(t, result.Code, btypes.CodeOK)

	approveMapper := ctx.Mapper(ApproveMapperName).(*ApproveMapper)
	approve, exists := approveMapper.GetApprove(tx.From, tx.To)
	require.True(t, exists)
	require.True(t, tx.Approve.Equals(approve))
}

func TestTxApproveIncrease_ValidateData(t *testing.T) {
	ctx := defaultContextWithQSC()

	createTx := TxCreateApprove{
		genTestApprove(),
	}
	increaseTx := TxIncreaseApprove{
		genTestApprove(),
	}
	require.NotNil(t, increaseTx.ValidateData(ctx))

	approveMapper := ctx.Mapper(ApproveMapperName).(*ApproveMapper)
	approveMapper.SaveApprove(createTx.Approve)

	require.Nil(t, increaseTx.ValidateData(ctx))
}

func TestTxApproveIncrease_Exec(t *testing.T) {
	ctx := defaultContext()

	createTx := TxCreateApprove{
		genTestApprove(),
	}
	increaseTx := TxIncreaseApprove{
		genTestApprove(),
	}

	approveMapper := ctx.Mapper(ApproveMapperName).(*ApproveMapper)
	approveMapper.SaveApprove(createTx.Approve)

	result, cross := increaseTx.Exec(ctx)
	require.Nil(t, cross)
	require.Equal(t, result.Code, btypes.CodeOK)

	approve, exists := approveMapper.GetApprove(createTx.From, createTx.To)
	require.True(t, exists)
	require.True(t, createTx.Approve.Plus(increaseTx.QOS, increaseTx.QSCs).Equals(approve))
}

func TestTxApproveDecrease_ValidateData(t *testing.T) {
	ctx := defaultContextWithQSC()

	createTx := TxCreateApprove{
		genTestApprove(),
	}
	decreaseTx := TxDecreaseApprove{
		genTestApprove(),
	}
	require.NotNil(t, decreaseTx.ValidateData(ctx))

	approveMapper := ctx.Mapper(ApproveMapperName).(*ApproveMapper)
	approveMapper.SaveApprove(createTx.Approve)

	require.Nil(t, decreaseTx.ValidateData(ctx))

	decreaseTx.QOS = btypes.NewInt(100)
	require.Nil(t, decreaseTx.ValidateData(ctx))

	decreaseTx.QOS = btypes.NewInt(110)
	require.NotNil(t, decreaseTx.ValidateData(ctx))
}

func TestTxApproveDecrease_Exec(t *testing.T) {
	ctx := defaultContext()

	createTx := TxCreateApprove{
		genTestApprove(),
	}
	decreaseTx := TxDecreaseApprove{
		genTestApprove(),
	}
	approveMapper := ctx.Mapper(ApproveMapperName).(*ApproveMapper)
	approveMapper.SaveApprove(createTx.Approve)

	result, cross := decreaseTx.Exec(ctx)
	require.Nil(t, cross)
	require.Equal(t, result.Code, btypes.CodeOK)

	approve, exists := approveMapper.GetApprove(createTx.From, createTx.To)
	require.True(t, exists)
	require.True(t, createTx.Approve.Minus(decreaseTx.QOS, decreaseTx.QSCs).Equals(approve))
}

func TestTxApproveUse_ValidateData(t *testing.T) {
	ctx := defaultContextWithQSC()

	createTx := TxCreateApprove{
		genTestApprove(),
	}
	useTx := TxUseApprove{
		genTestApprove(),
	}
	require.NotNil(t, useTx.ValidateData(ctx))

	approveMapper := ctx.Mapper(ApproveMapperName).(*ApproveMapper)
	approveMapper.SaveApprove(createTx.Approve)
	require.NotNil(t, useTx.ValidateData(ctx))

	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	accountMapper.SetAccount(genTestAccount(btypes.Address(useTx.From)))
	accountMapper.SetAccount(genTestAccount(btypes.Address(useTx.To)))

	require.Nil(t, useTx.ValidateData(ctx))

	useTx.QOS = btypes.NewInt(110)
	require.NotNil(t, useTx.ValidateData(ctx))

}

func TestTxApproveUse_GetSigner(t *testing.T) {
	useTx := TxUseApprove{
		genTestApprove(),
	}
	require.Equal(t, useTx.GetSigner(), []btypes.Address{useTx.To})
}

func TestTxApproveUse_GetGasPayer(t *testing.T) {
	useTx := TxUseApprove{
		genTestApprove(),
	}
	require.Equal(t, useTx.GetGasPayer(), useTx.To)
}

func TestTxApproveUse_Exec(t *testing.T) {
	ctx := defaultContext()

	createTx := TxCreateApprove{
		genTestApprove(),
	}
	useTx := TxUseApprove{
		genTestApprove(),
	}
	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	accountMapper.SetAccount(genTestAccount(btypes.Address(useTx.From)))
	accountMapper.SetAccount(genTestAccount(btypes.Address(useTx.To)))

	approveMapper := ctx.Mapper(ApproveMapperName).(*ApproveMapper)
	approveMapper.SaveApprove(createTx.Approve)

	result, cross := useTx.Exec(ctx)
	require.Nil(t, cross)
	require.Equal(t, result.Code, btypes.CodeOK)

	approve, exists := approveMapper.GetApprove(useTx.From, useTx.To)
	require.True(t, exists)
	require.True(t, createTx.Minus(useTx.QOS, useTx.QSCs).Equals(approve))

}

func TestTxApproveCancel_ValidateData(t *testing.T) {
	ctx := defaultContext()
	createTx := TxCreateApprove{
		genTestApprove(),
	}
	cancelTx := TxCancelApprove{
		createTx.From,
		createTx.To,
	}
	require.NotNil(t, cancelTx.ValidateData(ctx))

	mapper := ctx.Mapper(ApproveMapperName).(*ApproveMapper)
	mapper.SaveApprove(createTx.Approve)

	require.Nil(t, cancelTx.ValidateData(ctx))
}

func TestTxApproveCancel_Exec(t *testing.T) {
	ctx := defaultContext()
	createTx := TxCreateApprove{
		genTestApprove(),
	}
	cancelTx := TxCancelApprove{
		createTx.From,
		createTx.To,
	}

	mapper := ctx.Mapper(ApproveMapperName).(*ApproveMapper)
	mapper.SaveApprove(createTx.Approve)

	result, _ := cancelTx.Exec(ctx)
	require.Equal(t, result.Code, btypes.CodeOK)

}

func TestTxApproveCancel_GetSigner(t *testing.T) {
	cancelTx := genApproveCancelTx()
	require.Equal(t, cancelTx.GetSigner(), []btypes.Address{cancelTx.From})
}

func TestTxApproveCancel_GetGasPayer(t *testing.T) {
	cancelTx := genApproveCancelTx()
	require.Equal(t, cancelTx.GetGasPayer(), cancelTx.From)
}

func TestTxApproveCancel_CalcGas(t *testing.T) {
	cancelTx := genApproveCancelTx()
	require.Equal(t, cancelTx.CalcGas(), btypes.NewInt(0))
}

func TestTxApproveCancel_GetSignData(t *testing.T) {
	cancelTx := genApproveCancelTx()
	ret := []byte{}
	ret = append(ret, cancelTx.From...)
	ret = append(ret, cancelTx.To...)
	require.Equal(t, cancelTx.GetSignData(), ret)
}
