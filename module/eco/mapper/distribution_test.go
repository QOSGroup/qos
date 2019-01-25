package mapper

import (
	"testing"

	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/eco/types"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/require"
	go_amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func defaultContext(key store.StoreKey, mapperMap map[string]mapper.IMapper) context.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, store.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := context.NewContext(cms, abci.Header{}, false, log.NewNopLogger(), mapperMap)
	return ctx
}

func getDistributionMapper() *DistributionMapper {

	cdc := go_amino.NewCodec()

	seedMapper := NewDistributionMapper()
	seedMapper.SetCodec(cdc)

	mapperMap := make(map[string]mapper.IMapper)
	mapperMap[seedMapper.MapperName()] = seedMapper

	ctx := defaultContext(seedMapper.GetStoreKey(), mapperMap)

	return GetDistributionMapper(ctx)
}

func TestDistributionMapper_GetCommunityFeePool(t *testing.T) {

	mapper := getDistributionMapper()
	fp := mapper.GetCommunityFeePool()

	require.Equal(t, btypes.ZeroInt(), fp)

	fp = fp.Add(btypes.OneInt())
	mapper.Set(types.BuildCommunityFeePoolKey(), fp)

	fp = mapper.GetCommunityFeePool()
	require.Equal(t, btypes.OneInt(), fp)
}

func TestDistributionMapper_GetValidatorHistoryPeriodSummary(t *testing.T) {
	mapper := getDistributionMapper()
	addr := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
	period := uint64(100)

	frac := mapper.GetValidatorHistoryPeriodSummary(addr, period)
	require.Equal(t, qtypes.ZeroFraction(), frac)

	frac = frac.Add(qtypes.OneFraction())
	mapper.Set(types.BuildValidatorHistoryPeriodSummaryKey(addr, period), frac)

	frac = mapper.GetValidatorHistoryPeriodSummary(addr, period)
	require.Equal(t, qtypes.OneFraction(), frac)

}

func TestDistributionMapper_GetPreDistributionQOS(t *testing.T) {
	mapper := getDistributionMapper()

	a := mapper.GetPreDistributionQOS()
	require.Equal(t, btypes.ZeroInt(), a)

	add := btypes.NewInt(10)
	mapper.AddPreDistributionQOS(add)

	a = mapper.GetPreDistributionQOS()
	require.Equal(t, add, a)

	mapper.ClearPreDistributionQOS()

	a = mapper.GetPreDistributionQOS()
	require.Equal(t, btypes.ZeroInt(), a)

	mapper.AddPreDistributionQOS(add)
	mapper.AddPreDistributionQOS(add)
	mapper.AddPreDistributionQOS(add)

	a = mapper.GetPreDistributionQOS()
	require.Equal(t, btypes.NewInt(30), a)

}

func TestDistributionMapper_ClearValidatorPeriodSummaryInfo(t *testing.T) {

	mapper := getDistributionMapper()
	pubkey := ed25519.GenPrivKey().PubKey()
	addr := btypes.Address(pubkey.Address())

	validator := types.Validator{
		ValidatorPubKey: pubkey,
		BondTokens:      uint64(200),
	}

	mapper.InitValidatorPeriodSummaryInfo(addr)
	mapper.IncrementValidatorPeriod(validator)
	mapper.IncrementValidatorPeriod(validator)
	mapper.IncrementValidatorPeriod(validator)
	mapper.IncrementValidatorPeriod(validator)

	exsits := mapper.GetStore().Has(types.BuildValidatorCurrentPeriodSummaryKey(addr))
	require.Equal(t, true, exsits)

	for i := uint64(0); i <= 4; i++ {
		exsits = mapper.GetStore().Has(types.BuildValidatorHistoryPeriodSummaryKey(addr, uint64(i)))
		require.Equal(t, true, exsits)
	}

	mapper.DeleteValidatorPeriodSummaryInfo(addr)

	exsits = mapper.GetStore().Has(types.BuildValidatorCurrentPeriodSummaryKey(addr))
	require.Equal(t, false, exsits)

	for i := uint64(0); i <= 4; i++ {
		exsits = mapper.GetStore().Has(types.BuildValidatorHistoryPeriodSummaryKey(addr, uint64(i)))
		require.Equal(t, false, exsits)
	}

}

func TestDistributionMapper_calculateRewardsBetweenPeriod(t *testing.T) {
	mapper := getDistributionMapper()

	pubkey := ed25519.GenPrivKey().PubKey()
	addr := btypes.Address(pubkey.Address())

	a := mapper.CalculateRewardsBetweenPeriod(addr, uint64(10), uint64(100), uint64(0))
	require.Equal(t, btypes.ZeroInt(), a)

}
