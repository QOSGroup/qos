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
	approve := Approve{
		From: btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		To:   btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		Qos:  btypes.NewInt(100),
		QscList: []*QSC{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	}
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
	approve.Qos = btypes.NewInt(0)
	require.True(t, approve.ValidateData(ctx))

	approve = Approve{
		From: btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		To:   btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		Qos:  btypes.NewInt(100),
		QscList: []*QSC{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	}
	require.False(t, approve.ValidateData(ctx))

	approve = Approve{
		From: btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		To:   btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		Qos:  btypes.NewInt(100),
		QscList: []*QSC{
			{
				Name:   "qos",
				Amount: btypes.NewInt(100),
			},
		},
	}
	require.False(t, approve.ValidateData(ctx))
}

func TestApprove_GetSigner(t *testing.T) {
	approve := Approve{
		From:    btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		To:      btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		Qos:     btypes.NewInt(100),
		QscList: []*QSC{},
	}
	require.Equal(t, approve.GetSigner(), []btypes.Address{approve.From})
}

func TestApprove_GetGasPayer(t *testing.T) {
	approve := Approve{
		From:    btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		To:      btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		Qos:     btypes.NewInt(100),
		QscList: []*QSC{},
	}
	require.Equal(t, approve.GetGasPayer(), approve.From)
}

func TestApprove_CalcGas(t *testing.T) {
	approve := Approve{
		From:    btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		To:      btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		Qos:     btypes.NewInt(100),
		QscList: []*QSC{},
	}
	require.Equal(t, approve.CalcGas(), btypes.NewInt(0))
}

func TestApprove_GetSignData(t *testing.T) {
	approve := Approve{
		From:    btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		To:      btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		Qos:     btypes.NewInt(100),
		QscList: []*QSC{},
	}
	ret := []byte{}
	ret = append(ret, approve.From...)
	ret = append(ret, approve.To...)
	if &approve.Qos == nil {
		ret = append(ret, approve.Qos.String()...)
	}
	for _, coin := range approve.QscList {
		ret = append(ret, []byte(coin.Name)...)
		ret = append(ret, []byte(coin.Amount.String())...)
	}
	require.Equal(t, approve.GetSignData(), ret)
}

func TestApprove_IsPositive(t *testing.T) {
	approve := Approve{
		Qos: btypes.NewInt(100),
		QscList: []*QSC{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	}
	require.True(t, approve.IsPositive())

	approve.Qos = btypes.NewInt(0)
	require.True(t, approve.IsPositive())

	approve = Approve{
		Qos: btypes.NewInt(100),
		QscList: []*QSC{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(-1),
			},
		},
	}
	require.False(t, approve.IsPositive())
}

func TestApprove_IsNotNegative(t *testing.T) {
	approve := Approve{
		Qos: btypes.NewInt(0),
		QscList: []*QSC{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	}
	require.True(t, approve.IsNotNegative())

	approve.Qos = btypes.NewInt(-1)
	require.False(t, approve.IsNotNegative())

	approve = Approve{
		Qos: btypes.NewInt(0),
		QscList: []*QSC{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(0),
			},
		},
	}
	require.True(t, approve.IsNotNegative())
}

func TestApprove_Negative(t *testing.T) {
	approve := Approve{
		Qos: btypes.NewInt(100),
		QscList: []*QSC{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	}
	negative := approve.Negative()
	require.True(t, negative.Qos.String() == "-100")

	require.Equal(t, approve, negative.Negative())
}

func TestApprove_Plus(t *testing.T) {
	approve1 := Approve{
		Qos: btypes.NewInt(100),
		QscList: []*QSC{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	}
	approve2 := Approve{
		Qos:     btypes.NewInt(100),
		QscList: []*QSC{},
	}
	a := approve1.Plus(approve2)
	require.Equal(t, a.Qos.String(), btypes.NewInt(200).String())
	require.Equal(t, a.QscList[0].Amount, btypes.NewInt(100))
}

func TestApprove_Minus(t *testing.T) {
	approve1 := Approve{
		Qos: btypes.NewInt(10),
		QscList: []*QSC{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	}
	approve2 := Approve{
		Qos:     btypes.NewInt(100),
		QscList: []*QSC{},
	}
	a := approve1.Minus(approve2)
	require.Equal(t, a.Qos.String(), btypes.NewInt(-90).String())
	require.Equal(t, a.QscList[0].Amount, btypes.NewInt(100))
}

func TestApprove_IsGTE(t *testing.T) {
	approve1 := Approve{
		Qos: btypes.NewInt(100),
		QscList: []*QSC{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	}
	approve2 := Approve{
		Qos:     btypes.NewInt(100),
		QscList: []*QSC{},
	}
	require.True(t, approve1.IsGTE(approve2))

	approve2 = Approve{
		Qos: btypes.NewInt(100),
		QscList: []*QSC{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	}
	require.True(t, approve1.IsGTE(approve2))

	approve2.Qos = btypes.NewInt(200)
	require.False(t, approve1.IsGTE(approve2))
}

func TestApprove_IsGT(t *testing.T) {
	approve1 := Approve{
		Qos: btypes.NewInt(100),
		QscList: []*QSC{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	}
	approve2 := Approve{
		Qos:     btypes.NewInt(100),
		QscList: []*QSC{},
	}
	require.True(t, approve1.IsGT(approve2))

	approve2.Qos = btypes.NewInt(200)
	require.False(t, approve1.IsGT(approve2))

	approve2 = Approve{
		Qos: btypes.NewInt(100),
		QscList: []*QSC{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	}
	require.False(t, approve1.IsGT(approve2))
}

func TestApprove_Equals(t *testing.T) {
	approve1 := Approve{
		Qos: btypes.NewInt(100),
		QscList: []*QSC{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	}
	approve2 := Approve{
		Qos: btypes.NewInt(100),
		QscList: []*QSC{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	}
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
