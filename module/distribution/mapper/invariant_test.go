package mapper

import (
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	bmapper "github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/distribution/types"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"testing"
)

func defaultContext() context.Context {
	mapperMap := make(map[string]bmapper.IMapper)

	cdc := baseabci.MakeQBaseCodec()
	qtypes.RegisterCodec(cdc)

	accountMapper := bacc.NewAccountMapper(nil, qtypes.ProtoQOSAccount)
	accountMapper.SetCodec(cdc)
	accountKey := accountMapper.GetStoreKey()
	mapperMap[bacc.AccountMapperName] = accountMapper

	distributionMapper := NewMapper()
	distributionMapper.SetCodec(cdc)
	distributionKey := distributionMapper.GetStoreKey()
	mapperMap[types.MapperName] = distributionMapper

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(accountKey, btypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(distributionKey, btypes.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := context.NewContext(cms, abci.Header{}, false, log.NewNopLogger(), mapperMap)
	return ctx
}

func TestFeepoolInvariant(t *testing.T) {
	ctx := defaultContext()

	dm := GetMapper(ctx)
	dm.SetCommunityFeePool(btypes.NewInt(100))
	_, coins, broken := FeepoolInvariant("gov")(ctx)
	assert.False(t, broken)
	assert.Equal(t, coins.AmountOf(qtypes.QOSCoinName), btypes.NewInt(100))

	dm.SetCommunityFeePool(btypes.NewInt(-100))
	_, coins, broken = FeepoolInvariant("gov")(ctx)
	assert.True(t, broken)
	assert.Equal(t, coins.AmountOf(qtypes.QOSCoinName), btypes.NewInt(-100))
}
