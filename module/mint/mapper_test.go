package mint

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	minttypes "github.com/QOSGroup/qos/module/mint/types"
	bmapper "github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/store"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/QOSGroup/qbase/baseabci"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

var cdc = baseabci.MakeQBaseCodec()

func TestSaveParams(t *testing.T) {
	params := minttypes.DefaultParams()
	mintMapper := defaultContext().Mapper(MintMapperName).(*MintMapper)

	mintMapper.SetParams(params)

	currentInflationPhrase, exist := mintMapper.GetCurrentInflationPhrase()
	require.True(t, exist)

	fmt.Println(currentInflationPhrase.EndTime)

	mintMapper.AddAppliedQOSAmount(1999)
	require.Equal(t, mintMapper.GetAppliedQOSAmount(), uint64(1999))

	now := time.Now()
	mintMapper.AddInflationPhrase(minttypes.InflationPhrase{
		now, //插入当前时间
		1000,
		0,
	},)

	mintMapper.AddInflationPhrase(minttypes.InflationPhrase{
		now.Add(time.Minute*10), //插入十分钟后的时间
		2000,
		0,
	},)

	currentInflationPhrase, exist = mintMapper.GetCurrentInflationPhrase()
	require.True(t, exist)

	fmt.Println(currentInflationPhrase.EndTime)
	require.Equal(t, currentInflationPhrase.TotalAmount, uint64(2000))
}

func defaultContext() context.Context {
	mapperMap := make(map[string]bmapper.IMapper)

	mintMapper := NewMintMapper()
	mintMapper.SetCodec(cdc)
	mintKey := mintMapper.GetStoreKey()
	mapperMap[MintMapperName] = mintMapper

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)

	cms.MountStoreWithDB(mintKey, store.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := context.NewContext(cms, abci.Header{}, false, log.NewNopLogger(), mapperMap)
	return ctx
}