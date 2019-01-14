package approve

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSaveApprove(t *testing.T) {
	ctx := defaultContext()
	approveMapper, _ := ctx.Mapper(ApproveMapperName).(*ApproveMapper)

	approve := genTestApprove()
	approveMapper.SaveApprove(approve)

	recover, exists := approveMapper.GetApprove(approve.From, approve.To)
	require.True(t, exists)
	require.True(t, approve.Equals(recover))
}

func TestDeleteApprove(t *testing.T) {
	ctx := defaultContext()
	approveMapper, _ := ctx.Mapper(ApproveMapperName).(*ApproveMapper)

	approve := genTestApprove()
	approveMapper.SaveApprove(approve)

	recover, exists := approveMapper.GetApprove(approve.From, approve.To)
	require.True(t, exists)
	require.True(t, approve.Equals(recover))

	approveMapper.DeleteApprove(approve.From, approve.To)

	_, exists = approveMapper.GetApprove(approve.From, approve.To)
	require.False(t, exists)

}
