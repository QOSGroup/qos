package mapper

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestApproveInvariant(t *testing.T) {

	ctx := defaultContext()
	approveMapper := GetMapper(ctx)

	approve := genTestApprove(testFromAddr, testToAddr, 100, 100)
	approveMapper.SaveApprove(approve)
	_, coins, broken := ApproveInvariant("approve")(ctx)
	require.False(t, broken)
	require.True(t, coins.IsZero())

	approve = genTestApprove(testFromAddr, testToAddr, -100, 100)
	approveMapper.SaveApprove(approve)
	_, coins, broken = ApproveInvariant("approve")(ctx)
	require.True(t, broken)
	require.True(t, coins.IsZero())

	approve = genTestApprove(testFromAddr, testToAddr, 100, -100)
	approveMapper.SaveApprove(approve)
	_, coins, broken = ApproveInvariant("approve")(ctx)
	require.True(t, broken)
	require.True(t, coins.IsZero())

	approve = genTestApprove(testFromAddr, testToAddr, -100, -100)
	approveMapper.SaveApprove(approve)
	_, coins, broken = ApproveInvariant("approve")(ctx)
	require.True(t, broken)
	require.True(t, coins.IsZero())

	approve = genTestApprove(testFromAddr, testToAddr, 100, 100)
	approveMapper.SaveApprove(approve)
	approve = genTestApprove(testToAddr, testFromAddr, 100, 100)
	approveMapper.SaveApprove(approve)
	_, coins, broken = ApproveInvariant("approve")(ctx)
	require.False(t, broken)
	require.True(t, coins.IsZero())

	approve = genTestApprove(testFromAddr, testToAddr, 100, 100)
	approveMapper.SaveApprove(approve)
	approve = genTestApprove(testToAddr, testFromAddr, 100, -100)
	approveMapper.SaveApprove(approve)
	_, coins, broken = ApproveInvariant("approve")(ctx)
	require.True(t, broken)
	require.True(t, coins.IsZero())

}
