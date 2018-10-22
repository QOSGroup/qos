package mapper

import (
	bmapper "github.com/QOSGroup/qbase/mapper"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
	"testing"
)

var cdc *amino.Codec
var approve types.Approve
var approveCancel types.ApproveCancel

func init() {
	cdc = amino.NewCodec()
	cryptoAmino.RegisterAmino(cdc)

	fromPub := ed25519.GenPrivKey().PubKey()
	fromAddr := btypes.Address(fromPub.Address())
	toPub := ed25519.GenPrivKey().PubKey()
	toAddr := btypes.Address(toPub.Address())
	approve = types.Approve{
		From: fromAddr,
		To:   toAddr,
		Qos:  btypes.NewInt(100),
		QscList: []*types.QSC{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(100),
			},
		},
	}
	approveCancel = types.ApproveCancel{
		From: fromAddr,
		To:   toAddr,
	}
}

func TestSaveApprove(t *testing.T) {
	approveMapper := NewApproveMapper()
	approveMapper.SetCodec(cdc)
	storeKey := approveMapper.GetStoreKey()
	mapper := make(map[string]bmapper.IMapper)
	mapper[approveMapper.Name()] = approveMapper
	ctx := defaultContext(storeKey, mapper)
	approveMapper, _ = ctx.Mapper(approveMapper.Name()).(*ApproveMapper)

	err := approveMapper.SaveApprove(approve)
	require.Nil(t, err)
	recover, exists := approveMapper.GetApprove(approve.From, approve.To)
	require.True(t, exists)
	require.True(t, approve.Equals(recover))
}

func TestDeleteApprove(t *testing.T) {
	approveMapper := NewApproveMapper()
	approveMapper.SetCodec(cdc)
	storeKey := approveMapper.GetStoreKey()
	mapper := make(map[string]bmapper.IMapper)
	mapper[approveMapper.Name()] = approveMapper
	ctx := defaultContext(storeKey, mapper)
	approveMapper, _ = ctx.Mapper(approveMapper.Name()).(*ApproveMapper)

	err := approveMapper.SaveApprove(approve)
	require.Nil(t, err)
	recover, exists := approveMapper.GetApprove(approve.From, approve.To)
	require.True(t, exists)
	require.True(t, approve.Equals(recover))

	err = approveMapper.DeleteApprove(approveCancel)
	require.Nil(t, err)

	_, exists = approveMapper.GetApprove(approve.From, approve.To)
	require.False(t, exists)

}
