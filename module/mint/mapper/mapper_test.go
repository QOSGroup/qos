package mapper

import (
	"fmt"
	"github.com/QOSGroup/qos/module/mint/types"
	"testing"
	"time"

	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	bmapper "github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

var cdc = baseabci.MakeQBaseCodec()

func TestSaveParams(t *testing.T) {
	params := types.DefaultMintParams()
	mapper := defaultMintContext().Mapper(types.MapperName).(*Mapper)

	mapper.SetMintParams(params)

	blockSec := uint64(time.Now().UTC().Unix())

	currentInflationPhrase, exist := mapper.GetCurrentInflationPhrase(blockSec)
	require.True(t, exist)

	fmt.Println(currentInflationPhrase.EndTime)

	mapper.addCurrentPhraseAppliedQOSAmount(blockSec, 1999)
	require.Equal(t, mapper.getCurrentPhraseAppliedQOSAmount(blockSec), uint64(1999))

	now := time.Now()
	mapper.AddInflationPhrase(types.InflationPhrase{
		now, //插入当前时间
		1000,
		0,
	})

	mapper.AddInflationPhrase(types.InflationPhrase{
		now.Add(time.Minute * 10), //插入十分钟后的时间
		2000,
		0,
	})

	currentInflationPhrase, exist = mapper.GetCurrentInflationPhrase(blockSec)
	require.True(t, exist)

	fmt.Println(currentInflationPhrase.EndTime)
	require.Equal(t, currentInflationPhrase.TotalAmount, uint64(2000))
}

func defaultMintContext() context.Context {
	mapperMap := make(map[string]bmapper.IMapper)

	mapper := NewMapper()
	mapper.SetCodec(cdc)
	mintKey := mapper.GetStoreKey()
	mapperMap[types.MapperName] = mapper

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)

	cms.MountStoreWithDB(mintKey, btypes.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := context.NewContext(cms, abci.Header{}, false, log.NewNopLogger(), mapperMap)
	return ctx
}

func TestMintMapper_GetAllTotalMintQOSAmount(t *testing.T) {

	mapper := defaultMintContext().Mapper(types.MapperName).(*Mapper)

	amount := mapper.GetAllTotalMintQOSAmount()
	require.Equal(t, amount, uint64(0))

	mapper.SetAllTotalMintQOSAmount(uint64(100))

	amount = mapper.GetAllTotalMintQOSAmount()
	require.Equal(t, amount, uint64(100))

	mapper.addAllTotalMintQOSAmount(uint64(100))

	amount = mapper.GetAllTotalMintQOSAmount()
	require.Equal(t, amount, uint64(200))

	mapper.DelAllTotalMintQOSAmount()

	amount = mapper.GetAllTotalMintQOSAmount()
	require.Equal(t, amount, uint64(0))
}
