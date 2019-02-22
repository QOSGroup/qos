package types

import (
	"testing"

	btypes "github.com/QOSGroup/qbase/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func TestKey(t *testing.T) {

	key := ed25519.GenPrivKey()
	valAddr := btypes.Address(key.PubKey().Address())

	key = ed25519.GenPrivKey()
	deleAddr := btypes.Address(key.PubKey().Address())

	infoKey := BuildDelegatorEarningStartInfoKey(valAddr, deleAddr)
	vAddr, dAddr := GetDelegatorEarningStartInfoAddr(infoKey)

	require.Equal(t, valAddr, vAddr)
	require.Equal(t, deleAddr, dAddr)

	period := uint64(278)
	skey := BuildValidatorHistoryPeriodSummaryKey(valAddr, period)
	vAddr, p := GetValidatorHistoryPeriodSummaryAddrPeriod(skey)
	require.Equal(t, valAddr, vAddr)
	require.Equal(t, period, p)

	skey = BuildValidatorCurrentPeriodSummaryKey(valAddr)
	vAddr = GetValidatorCurrentPeriodSummaryAddr(skey)
	require.Equal(t, valAddr, vAddr)

	height := uint64(10086)
	skey = BuildDelegatorPeriodIncomeKey(valAddr, deleAddr, height)

	vAddr, dAddr, h := GetDelegatorPeriodIncomeHeightAddr(skey)
	require.Equal(t, height, h)
	require.Equal(t, valAddr, vAddr)
	require.Equal(t, deleAddr, dAddr)

	skey = BuildUnbondingDelegationByHeightDelKey(height, deleAddr)
	h, dAddr = GetUnbondingDelegationHeightAddress(skey)
	require.Equal(t, height, h)
	require.Equal(t, deleAddr, dAddr)
}
