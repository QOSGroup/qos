package mint

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

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