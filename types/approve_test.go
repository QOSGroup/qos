package types

import (
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
)

func defaultContext() context.Context {
	ctx := context.NewContext(nil, abci.Header{}, false, log.NewNopLogger(), nil)
	return ctx
}

func TestApprove_ValidateData(t *testing.T) {
	qos := btypes.NewInt(100)
	approve := NewApprove(
		btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		&qos,
		QSCs{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	)
	ctx := defaultContext()
	require.True(t, approve.ValidateData(ctx))

	from := approve.From
	to := approve.To

	approve.From = nil
	require.False(t, approve.ValidateData(ctx))
	approve.To = nil
	require.False(t, approve.ValidateData(ctx))
	approve.From = from
	require.False(t, approve.ValidateData(ctx))

	approve.To = to
	approve.QOS = btypes.NewInt(0)
	require.True(t, approve.ValidateData(ctx))

	approve = NewApprove(
		btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		&qos,
		QSCs{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	)
	require.False(t, approve.ValidateData(ctx))

	approve = NewApprove(
		btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		&qos,
		QSCs{
			{
				Name:   "qos",
				Amount: btypes.NewInt(100),
			},
		},
	)
	require.False(t, approve.ValidateData(ctx))
}

func TestApprove_GetSigner(t *testing.T) {
	qos := btypes.NewInt(100)
	approve := NewApprove(btypes.Address(ed25519.GenPrivKey().PubKey().Address()), btypes.Address(ed25519.GenPrivKey().PubKey().Address()), &qos, QSCs{})
	require.Equal(t, approve.GetSigner(), []btypes.Address{approve.From})
}

func TestApprove_GetGasPayer(t *testing.T) {
	qos := btypes.NewInt(100)
	approve := NewApprove(btypes.Address(ed25519.GenPrivKey().PubKey().Address()), btypes.Address(ed25519.GenPrivKey().PubKey().Address()), &qos, QSCs{})
	require.Equal(t, approve.GetGasPayer(), approve.From)
}

func TestApprove_CalcGas(t *testing.T) {
	qos := btypes.NewInt(100)
	approve := NewApprove(btypes.Address(ed25519.GenPrivKey().PubKey().Address()), btypes.Address(ed25519.GenPrivKey().PubKey().Address()), &qos, QSCs{})
	require.Equal(t, approve.CalcGas(), btypes.NewInt(0))
}

func TestApprove_GetSignData(t *testing.T) {
	qos := btypes.NewInt(100)
	approve := NewApprove(btypes.Address(ed25519.GenPrivKey().PubKey().Address()), btypes.Address(ed25519.GenPrivKey().PubKey().Address()), &qos, QSCs{})
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
	qos := btypes.NewInt(100)
	approve := NewApprove(
		btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		&qos,
		QSCs{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	)
	require.True(t, approve.IsPositive())

	approve.QOS = btypes.NewInt(0)
	require.True(t, approve.IsPositive())

	approve.QSCs[0].Amount = btypes.NewInt(-1)
	require.False(t, approve.IsPositive())
}

func TestApprove_IsNotNegative(t *testing.T) {
	qos := btypes.NewInt(0)
	approve := NewApprove(
		btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		&qos,
		QSCs{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	)
	require.True(t, approve.IsNotNegative())

	approve.QOS = btypes.NewInt(-1)
	require.False(t, approve.IsNotNegative())

	approve.QOS = btypes.NewInt(0)
	approve.QSCs[0].Amount = btypes.NewInt(0)
	require.True(t, approve.IsNotNegative())
}

func TestApprove_Negative(t *testing.T) {
	qos := btypes.NewInt(100)
	approve := NewApprove(
		btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		&qos,
		QSCs{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	)
	negative := approve.Negative()
	require.True(t, negative.QOS.String() == "-100")

	require.Equal(t, approve, negative.Negative())
}

func TestApprove_Plus(t *testing.T) {
	qos1 := btypes.NewInt(100)
	approve1 := NewApprove(
		btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		&qos1,
		QSCs{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	)
	qos2 := btypes.NewInt(100)
	a := approve1.Plus(qos2, QSCs{})
	require.Equal(t, a.QOS.String(), btypes.NewInt(200).String())
	require.Equal(t, a.QSCs[0].Amount, btypes.NewInt(100))
}

func TestApprove_Minus(t *testing.T) {
	qos1 := btypes.NewInt(10)
	approve1 := NewApprove(
		btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		&qos1,
		QSCs{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	)
	qos2 := btypes.NewInt(100)
	a := approve1.Minus(qos2, QSCs{})
	require.Equal(t, a.QOS.String(), btypes.NewInt(-90).String())
	require.Equal(t, a.QSCs[0].Amount, btypes.NewInt(100))
}

func TestApprove_IsGTE(t *testing.T) {
	qos1 := btypes.NewInt(100)
	approve1 := NewApprove(
		btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		&qos1,
		QSCs{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	)
	qos2 := btypes.NewInt(100)
	require.True(t, approve1.IsGTE(qos2, QSCs{}))

	qsc2 := QSCs{
		{
			Name:   "qstar",
			Amount: btypes.NewInt(100),
		},
	}
	require.True(t, approve1.IsGTE(qos2, qsc2))

	qos2 = btypes.NewInt(200)
	require.False(t, approve1.IsGTE(qos2, qsc2))
}

func TestApprove_IsGT(t *testing.T) {
	qos1 := btypes.NewInt(100)
	approve1 := NewApprove(
		btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		&qos1,
		QSCs{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	)
	qos2 := btypes.NewInt(100)
	qsc2 := QSCs{}
	require.True(t, approve1.IsGT(qos2, qsc2))

	qos2 = btypes.NewInt(200)
	require.False(t, approve1.IsGT(qos2, qsc2))

	qos2 = btypes.NewInt(100)
	qsc2 = append(qsc2, &QSC{
		Name:   "qstar",
		Amount: btypes.NewInt(100),
	})
	require.False(t, approve1.IsGT(qos2, qsc2))
}

func TestApprove_Equals(t *testing.T) {
	from := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
	to := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
	qos1 := btypes.NewInt(100)
	approve1 := NewApprove(
		from, to, &qos1,
		QSCs{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	)
	qos2 := btypes.NewInt(100)
	approve2 := NewApprove(
		from, to, &qos2,
		QSCs{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	)
	require.True(t, approve1.Equals(approve2))
}

func TestApproveCancel_ValidateData(t *testing.T) {
	approveCancel := ApproveCancel{
		From: btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		To:   btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
	}

	ctx := defaultContext()
	require.True(t, approveCancel.ValidateData(ctx))

	from := approveCancel.From

	approveCancel.From = nil
	require.False(t, approveCancel.ValidateData(ctx))
	approveCancel.To = nil
	require.False(t, approveCancel.ValidateData(ctx))
	approveCancel.From = from
	require.False(t, approveCancel.ValidateData(ctx))
}

func TestApproveCancel_GetSigner(t *testing.T) {
	approveCancel := ApproveCancel{
		From: btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		To:   btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
	}
	require.Equal(t, approveCancel.GetSigner(), []btypes.Address{approveCancel.From})
}

func TestApproveCancel_GetGasPayer(t *testing.T) {
	approveCancel := ApproveCancel{
		From: btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		To:   btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
	}
	require.Equal(t, approveCancel.GetGasPayer(), approveCancel.From)
}

func TestApproveCancel_CalcGas(t *testing.T) {
	approveCancel := ApproveCancel{
		From: btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		To:   btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
	}
	require.Equal(t, approveCancel.CalcGas(), btypes.NewInt(0))
}

func TestApproveCancel_GetSignData(t *testing.T) {
	approveCancel := ApproveCancel{
		From: btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		To:   btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
	}
	ret := []byte{}
	ret = append(ret, approveCancel.From...)
	ret = append(ret, approveCancel.To...)
	require.Equal(t, approveCancel.GetSignData(), ret)
}