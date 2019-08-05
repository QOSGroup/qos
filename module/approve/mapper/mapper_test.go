package mapper

import (
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	bmapper "github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/approve/types"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
)

func defaultContext() context.Context {
	mapperMap := make(map[string]bmapper.IMapper)

	cdc := amino.NewCodec()
	approveMapper := NewApproveMapper()
	approveMapper.SetCodec(cdc)
	approveKey := approveMapper.GetStoreKey()
	mapperMap[types.MapperName] = approveMapper

	accountMapper := bacc.NewAccountMapper(nil, qtypes.ProtoQOSAccount)
	accountMapper.SetCodec(cdc)
	acountKey := accountMapper.GetStoreKey()
	mapperMap[bacc.AccountMapperName] = accountMapper

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(approveKey, btypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(acountKey, btypes.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := context.NewContext(cms, abci.Header{}, false, log.NewNopLogger(), mapperMap)
	return ctx
}

var testFromAddr = btypes.Address(ed25519.GenPrivKey().PubKey().Address())
var testToAddr = btypes.Address(ed25519.GenPrivKey().PubKey().Address())

func genTestApprove(from, to btypes.Address, qos, qsc int64) types.Approve {
	return types.Approve{
		From: from,
		To:   to,
		QOS:  btypes.NewInt(qos),
		QSCs: qtypes.QSCs{
			{
				Name:   "qstar",
				Amount: btypes.NewInt(qsc),
			},
		},
	}
}

func TestSaveApprove(t *testing.T) {
	ctx := defaultContext()
	approveMapper := GetMapper(ctx)

	approve := genTestApprove(testFromAddr, testToAddr, 100, 100)
	approveMapper.SaveApprove(approve)

	recover, exists := approveMapper.GetApprove(approve.From, approve.To)
	require.True(t, exists)
	require.True(t, approve.Equals(recover))
}

func TestDeleteApprove(t *testing.T) {
	ctx := defaultContext()
	approveMapper := GetMapper(ctx)

	approve := genTestApprove(testFromAddr, testToAddr, 100, 100)
	approveMapper.SaveApprove(approve)

	recover, exists := approveMapper.GetApprove(approve.From, approve.To)
	require.True(t, exists)
	require.True(t, approve.Equals(recover))

	approveMapper.DeleteApprove(approve.From, approve.To)

	_, exists = approveMapper.GetApprove(approve.From, approve.To)
	require.False(t, exists)

}
