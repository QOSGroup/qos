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

var approve Approve
var approveCancel ApproveCancel

func init() {
	approve = Approve{
		From: btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		To:   btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		Coins: QSCS{
			{
				Name:   "qos",
				Amount: btypes.NewInt(100),
			},
		},
	}
	approveCancel = ApproveCancel{
		From: btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		To:   btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
	}
}

func defaultContext() context.Context {
	ctx := context.NewContext(nil, abci.Header{}, false, log.NewNopLogger(), nil)
	return ctx
}

func TestApprove_ValidateData(t *testing.T) {
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
	approve.Coins[0].Amount = btypes.NewInt(0)
	require.False(t, approve.ValidateData(ctx))
}

func TestApprove_GetSigner(t *testing.T) {
	require.Equal(t, approve.GetSigner(), []btypes.Address{approve.From})
}

func TestApprove_GetGasPayer(t *testing.T) {
	require.Equal(t, approve.GetGasPayer(), approve.From)
}

func TestApprove_CalcGas(t *testing.T) {
	require.Equal(t, approve.CalcGas(), btypes.NewInt(0))
}

func TestApprove_GetSignData(t *testing.T) {
	ret := []byte{}
	ret = append(ret, approve.From...)
	ret = append(ret, approve.To...)
	for _, coin := range approve.Coins {
		ret = append(ret, []byte(coin.Name)...)
		ret = append(ret, []byte(coin.Amount.String())...)
	}
	require.Equal(t, approve.GetSignData(), ret)
}

func TestApproveCancel_ValidateData(t *testing.T) {
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
	require.Equal(t, approveCancel.GetSigner(), []btypes.Address{approveCancel.From})
}

func TestApproveCancel_GetGasPayer(t *testing.T) {
	require.Equal(t, approveCancel.GetGasPayer(), approveCancel.From)
}

func TestApproveCancel_CalcGas(t *testing.T) {
	require.Equal(t, approveCancel.CalcGas(), btypes.NewInt(0))
}

func TestApproveCancel_GetSignData(t *testing.T) {
	ret := []byte{}
	ret = append(ret, approveCancel.From...)
	ret = append(ret, approveCancel.To...)
	require.Equal(t, approveCancel.GetSignData(), ret)
}
