package approve

import (
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	bmapper "github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/account"
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

func genTestAccount(addr btypes.Address) *account.QOSAccount {
	return &account.QOSAccount{
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

func genTestApprove() Approve {
	return Approve{
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
	mapperMap[GetApproveMapperStoreKey()] = approveMapper

	accountMapper := bacc.NewAccountMapper(nil, account.ProtoQOSAccount)
	accountMapper.SetCodec(cdc)
	acountKey := accountMapper.GetStoreKey()
	mapperMap[bacc.AccountMapperName] = accountMapper

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(approveKey, store.StoreTypeIAVL, db)
	cms.MountStoreWithDB(acountKey, store.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := context.NewContext(cms, abci.Header{}, false, log.NewNopLogger(), mapperMap)
	return ctx
}

func TestApprove_ValidateData(t *testing.T) {
	ctx := defaultContext()
	approve := genTestApprove()
	require.Nil(t, approve.ValidateData(ctx))

	from := approve.From
	to := approve.To

	approve.From = nil
	require.NotNil(t, approve.ValidateData(ctx))
	approve.To = nil
	require.NotNil(t, approve.ValidateData(ctx))
	approve.From = from
	require.NotNil(t, approve.ValidateData(ctx))

	approve.To = to
	approve.QOS = btypes.NewInt(0)
	require.Nil(t, approve.ValidateData(ctx))

	approve.QSCs = append(approve.QSCs, &types.QSC{
		Name:   "qstar",
		Amount: btypes.NewInt(100),
	})
	require.NotNil(t, approve.ValidateData(ctx))

	approve.QSCs = types.QSCs{
		{
			Name:   "qos",
			Amount: btypes.NewInt(100),
		},
	}
	require.NotNil(t, approve.ValidateData(ctx))
}

func TestApprove_GetSigner(t *testing.T) {
	approve := genTestApprove()
	require.Equal(t, approve.GetSigner(), []btypes.Address{approve.From})
}

func TestApprove_GetGasPayer(t *testing.T) {
	approve := genTestApprove()
	require.Equal(t, approve.GetGasPayer(), approve.From)
}

func TestApprove_CalcGas(t *testing.T) {
	approve := genTestApprove()
	require.Equal(t, approve.CalcGas(), btypes.NewInt(0))
}

func TestApprove_GetSignData(t *testing.T) {
	approve := genTestApprove()
	ret := []byte{}
	ret = append(ret, approve.From...)
	ret = append(ret, approve.To...)
	ret = append(ret, approve.QOS.String()...)
	for _, coin := range approve.QSCs {
		ret = append(ret, []byte(coin.Name)...)
		ret = append(ret, []byte(coin.Amount.String())...)
	}
	require.Equal(t, approve.GetSignData(), ret)
}

func TestApprove_IsPositive(t *testing.T) {
	approve := genTestApprove()
	require.True(t, approve.IsPositive())

	approve.QOS = btypes.NewInt(0)
	require.True(t, approve.IsPositive())

	approve.QSCs[0].Amount = btypes.NewInt(-1)
	require.False(t, approve.IsPositive())
}

func TestApprove_IsNotNegative(t *testing.T) {
	approve := genTestApprove()
	require.True(t, approve.IsNotNegative())

	approve.QOS = btypes.NewInt(-1)
	require.False(t, approve.IsNotNegative())

	approve.QOS = btypes.NewInt(0)
	approve.QSCs[0].Amount = btypes.NewInt(0)
	require.True(t, approve.IsNotNegative())
}

func TestApprove_Negative(t *testing.T) {
	approve := genTestApprove()
	negative := approve.Negative()
	require.True(t, negative.QOS.String() == "-100")

	require.Equal(t, approve, negative.Negative())
}

func TestApprove_Plus(t *testing.T) {
	approve := genTestApprove()
	qos := btypes.NewInt(100)
	a := approve.Plus(qos, types.QSCs{})
	require.Equal(t, a.QOS.String(), btypes.NewInt(200).String())
	require.Equal(t, a.QSCs[0].Amount, btypes.NewInt(100))
}

func TestApprove_Minus(t *testing.T) {
	approve := genTestApprove()
	qos := btypes.NewInt(100)
	a := approve.Minus(qos, types.QSCs{})
	require.Equal(t, a.QOS.String(), btypes.NewInt(0).String())
	require.Equal(t, a.QSCs[0].Amount, btypes.NewInt(100))
}

func TestApprove_IsGTE(t *testing.T) {
	approve := genTestApprove()
	qos := btypes.NewInt(100)
	require.True(t, approve.IsGTE(qos, types.QSCs{}))

	qsc := types.QSCs{
		{
			Name:   "qstar",
			Amount: btypes.NewInt(100),
		},
	}
	require.True(t, approve.IsGTE(qos, qsc))

	qos = btypes.NewInt(200)
	require.False(t, approve.IsGTE(qos, qsc))
}

func TestApprove_IsGT(t *testing.T) {
	approve := genTestApprove()
	qos := btypes.NewInt(100)
	qsc := types.QSCs{}
	require.True(t, approve.IsGT(qos, qsc))

	qos = btypes.NewInt(200)
	require.False(t, approve.IsGT(qos, qsc))

	qos = btypes.NewInt(100)
	qsc = append(qsc, &types.QSC{
		Name:   "qstar",
		Amount: btypes.NewInt(100),
	})
	require.False(t, approve.IsGT(qos, qsc))
}

func TestApprove_Equals(t *testing.T) {
	approve1 := genTestApprove()
	approve2 := genTestApprove()
	require.True(t, approve1.Equals(approve2))
}

func TestTxApproveCreate_ValidateData(t *testing.T) {
	ctx := defaultContext()

	tx := TxCreateApprove{
		genTestApprove(),
	}
	require.Nil(t, tx.ValidateData(ctx))

	approveMapper := ctx.Mapper(GetApproveMapperStoreKey()).(*ApproveMapper)
	err := approveMapper.SaveApprove(tx.Approve)
	require.Nil(t, err)

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

	approveMapper := ctx.Mapper(GetApproveMapperStoreKey()).(*ApproveMapper)
	approve, exists := approveMapper.GetApprove(tx.From, tx.To)
	require.True(t, exists)
	require.True(t, tx.Approve.Equals(approve))
}

func TestTxApproveIncrease_ValidateData(t *testing.T) {
	ctx := defaultContext()

	createTx := TxCreateApprove{
		genTestApprove(),
	}
	increaseTx := TxIncreaseApprove{
		genTestApprove(),
	}
	require.NotNil(t, increaseTx.ValidateData(ctx))

	approveMapper := ctx.Mapper(GetApproveMapperStoreKey()).(*ApproveMapper)
	err := approveMapper.SaveApprove(createTx.Approve)
	require.Nil(t, err)

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

	approveMapper := ctx.Mapper(GetApproveMapperStoreKey()).(*ApproveMapper)
	err := approveMapper.SaveApprove(createTx.Approve)
	require.Nil(t, err)

	result, cross := increaseTx.Exec(ctx)
	require.Nil(t, cross)
	require.Equal(t, result.Code, btypes.CodeOK)

	approve, exists := approveMapper.GetApprove(createTx.From, createTx.To)
	require.True(t, exists)
	require.True(t, createTx.Approve.Plus(increaseTx.QOS, increaseTx.QSCs).Equals(approve))
}

func TestTxApproveDecrease_ValidateData(t *testing.T) {
	ctx := defaultContext()

	createTx := TxCreateApprove{
		genTestApprove(),
	}
	decreaseTx := TxDecreaseApprove{
		genTestApprove(),
	}
	require.NotNil(t, decreaseTx.ValidateData(ctx))

	approveMapper := ctx.Mapper(GetApproveMapperStoreKey()).(*ApproveMapper)
	err := approveMapper.SaveApprove(createTx.Approve)
	require.Nil(t, err)

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
	approveMapper := ctx.Mapper(GetApproveMapperStoreKey()).(*ApproveMapper)
	err := approveMapper.SaveApprove(createTx.Approve)
	require.Nil(t, err)

	result, cross := decreaseTx.Exec(ctx)
	require.Nil(t, cross)
	require.Equal(t, result.Code, btypes.CodeOK)

	approve, exists := approveMapper.GetApprove(createTx.From, createTx.To)
	require.True(t, exists)
	require.True(t, createTx.Approve.Minus(decreaseTx.QOS, decreaseTx.QSCs).Equals(approve))
}

func TestTxApproveUse_ValidateData(t *testing.T) {
	ctx := defaultContext()

	createTx := TxCreateApprove{
		genTestApprove(),
	}
	useTx := TxUseApprove{
		genTestApprove(),
	}
	require.NotNil(t, useTx.ValidateData(ctx))

	approveMapper := ctx.Mapper(GetApproveMapperStoreKey()).(*ApproveMapper)
	err := approveMapper.SaveApprove(createTx.Approve)
	require.Nil(t, err)
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

	result, cross := useTx.Exec(ctx)
	require.Nil(t, cross)
	require.NotEqual(t, result.Code, btypes.CodeOK)

	createTx.QOS = btypes.NewInt(1)
	approveMapper := ctx.Mapper(GetApproveMapperStoreKey()).(*ApproveMapper)
	err := approveMapper.SaveApprove(createTx.Approve)
	require.Nil(t, err)

	result, cross = useTx.Exec(ctx)
	require.Nil(t, cross)
	require.NotEqual(t, result.Code, btypes.CodeOK)

	createTx.QOS = btypes.NewInt(100)
	err = approveMapper.SaveApprove(createTx.Approve)
	require.Nil(t, err)

	result, cross = useTx.Exec(ctx)
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

	mapper := ctx.Mapper(GetApproveMapperStoreKey()).(*ApproveMapper)
	err := mapper.SaveApprove(createTx.Approve)
	require.Nil(t, err)

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
	result, cross := cancelTx.Exec(ctx)
	require.Nil(t, cross)
	require.NotEqual(t, result.Code, btypes.CodeOK)

	mapper := ctx.Mapper(GetApproveMapperStoreKey()).(*ApproveMapper)
	err := mapper.SaveApprove(createTx.Approve)
	require.Nil(t, err)

	result, _ = cancelTx.Exec(ctx)
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
