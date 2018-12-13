package types

import (
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

var testFromAddr = btypes.Address(ed25519.GenPrivKey().PubKey().Address())
var testToAddr = btypes.Address(ed25519.GenPrivKey().PubKey().Address())

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

func TestApprove_IsValid(t *testing.T) {
	approve := genTestApprove()
	_, err := approve.IsValid()
	require.Nil(t, err)

	from := approve.From
	to := approve.To

	approve.From = nil
	_, err = approve.IsValid()
	require.NotNil(t, err)
	approve.To = nil
	_, err = approve.IsValid()
	require.NotNil(t, err)
	approve.From = from
	_, err = approve.IsValid()
	require.NotNil(t, err)

	approve.To = to
	approve.QOS = btypes.NewInt(0)
	_, err = approve.IsValid()
	require.Nil(t, err)

	approve.QSCs = append(approve.QSCs, &types.QSC{
		Name:   "qstar",
		Amount: btypes.NewInt(100),
	})
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
