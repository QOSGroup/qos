package approve

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSaveApprove(t *testing.T) {
	ctx := defaultContext()
	approveMapper, _ := ctx.Mapper(GetApproveMapperStoreKey()).(*ApproveMapper)

	approve := genTestApprove()
	err := approveMapper.SaveApprove(approve)
	require.Nil(t, err)
	recover, exists := approveMapper.GetApprove(approve.From, approve.To)
	require.True(t, exists)
	require.True(t, approve.Equals(recover))
}

func TestDeleteApprove(t *testing.T) {
	ctx := defaultContext()
	approveMapper, _ := ctx.Mapper(GetApproveMapperStoreKey()).(*ApproveMapper)

	approve := genTestApprove()
	err := approveMapper.SaveApprove(approve)
	require.Nil(t, err)
	recover, exists := approveMapper.GetApprove(approve.From, approve.To)
	require.True(t, exists)
	require.True(t, approve.Equals(recover))

	err = approveMapper.DeleteApprove(approve.From, approve.To)
	require.Nil(t, err)

	_, exists = approveMapper.GetApprove(approve.From, approve.To)
	require.False(t, exists)

}
