package mapper

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	bmapper "github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/mint/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"testing"
)

var cdc = baseabci.MakeQBaseCodec()

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

	mapper := GetMapper(defaultMintContext())

	amount := mapper.GetAllTotalMintQOSAmount().NilToZero()
	require.Equal(t, amount, btypes.ZeroInt())

	mapper.SetAllTotalMintQOSAmount(btypes.NewInt(100))

	amount = mapper.GetAllTotalMintQOSAmount()
	require.Equal(t, amount, btypes.NewInt(100))

	mapper.AddAllTotalMintQOSAmount(btypes.NewInt(100))

	amount = mapper.GetAllTotalMintQOSAmount()
	require.Equal(t, amount, btypes.NewInt(200))

	mapper.DelAllTotalMintQOSAmount()

	amount = mapper.GetAllTotalMintQOSAmount().NilToZero()
	require.Equal(t, amount, btypes.ZeroInt())
}
