package bank

import (
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	bmapper "github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/bank/mapper"
	"github.com/QOSGroup/qos/module/bank/types"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
	"time"
)

func testContext() context.Context {
	mapperMap := make(map[string]bmapper.IMapper)
	accountMapper := bacc.NewAccountMapper(nil, qtypes.ProtoQOSAccount)
	accountMapper.SetCodec(Cdc)
	acountKey := accountMapper.GetStoreKey()
	mapperMap[bacc.AccountMapperName] = accountMapper

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(acountKey, btypes.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := context.NewContext(cms, abci.Header{Time: time.Now().UTC()}, false, log.NewNopLogger(), mapperMap)
	return ctx
}

func TestReleaseLockedAccount(t *testing.T) {
	ctx := testContext()
	am := GetMapper(ctx)

	_, exists := mapper.GetLockInfo(ctx)
	require.False(t, exists)

	// init lock info
	lockAddress := btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	receiverAddress := btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	totalAmount := uint64(10000000000000)
	releasedAmountInit := uint64(8333332)
	releaseTime := ctx.BlockHeader().Time.Add(-time.Hour)
	releaseInterval := uint(30)
	releaseTimesInit := uint(22)
	lockInfo := types.NewLockInfo(lockAddress, receiverAddress, totalAmount, releasedAmountInit, releaseTime, releaseInterval, releaseTimesInit)
	am.SetAccount(qtypes.NewQOSAccount(lockAddress, btypes.NewInt(int64(totalAmount-releasedAmountInit)), nil))
	lockAccount := mapper.GetAccount(ctx, lockAddress)
	require.Equal(t, lockAccount.QOS, btypes.NewInt(int64(totalAmount-releasedAmountInit)))
	mapper.SetLockInfo(ctx, lockInfo)
	lockInfo, exists = mapper.GetLockInfo(ctx)
	require.True(t, exists)

	// first time release
	ReleaseLockedAccount(ctx, lockInfo)
	lockInfo, exists = mapper.GetLockInfo(ctx)
	require.True(t, exists)
	require.Equal(t, lockInfo.ReleaseTimes, releaseTimesInit-1)
	require.Equal(t, lockInfo.ReleasedAmount, releasedAmountInit+(totalAmount-releasedAmountInit)/uint64(releaseTimesInit))
	ReleaseLockedAccount(ctx, lockInfo)
	lockInfo, exists = mapper.GetLockInfo(ctx)
	require.True(t, exists)
	require.Equal(t, lockInfo.ReleaseTimes, releaseTimesInit-1)
	require.Equal(t, lockInfo.ReleasedAmount, releasedAmountInit+(totalAmount-releasedAmountInit)/uint64(releaseTimesInit))
	require.Equal(t, lockInfo.ReleaseTime, releaseTime.Add(time.Hour*24*time.Duration(releaseInterval)))

	headerTime := ctx.BlockHeader().Time
	releasedAmount := lockInfo.ReleasedAmount
	releaseTimes := lockInfo.ReleaseTimes
	for i := uint(1); i < releaseTimesInit; i++ {

		// before release time
		ctx = ctx.WithBlockHeader(abci.Header{Time: releaseTime.Add(time.Hour * 24 * time.Duration(i*releaseInterval))})
		ReleaseLockedAccount(ctx, lockInfo)
		lockInfo, exists = mapper.GetLockInfo(ctx)
		require.True(t, exists)
		require.Equal(t, lockInfo.ReleaseTimes, releaseTimes)
		require.Equal(t, lockInfo.ReleasedAmount, releasedAmount)
		require.Equal(t, lockInfo.ReleaseTime, releaseTime.Add(time.Hour*24*time.Duration(i*releaseInterval)))

		// after release time
		ctx = ctx.WithBlockHeader(abci.Header{Time: headerTime.Add(time.Hour * 24 * time.Duration(i*releaseInterval))})
		ReleaseLockedAccount(ctx, lockInfo)
		lockInfo, exists = mapper.GetLockInfo(ctx)
		if i != releaseTimesInit-1 {
			require.True(t, exists)
			require.Equal(t, lockInfo.ReleaseTimes, releaseTimes-1)
			require.Equal(t, lockInfo.ReleasedAmount, releasedAmount+(totalAmount-releasedAmount)/uint64(releaseTimes))
			require.Equal(t, lockInfo.ReleaseTime, releaseTime.Add(time.Hour*24*time.Duration((i+1)*releaseInterval)))

			releasedAmount = lockInfo.ReleasedAmount
			releaseTimes = lockInfo.ReleaseTimes
		} else {
			require.False(t, exists)
		}
	}

	lockAccount = mapper.GetAccount(ctx, lockAddress)
	require.Equal(t, lockAccount.QOS, btypes.ZeroInt())
	receiver := mapper.GetAccount(ctx, receiverAddress)
	require.Equal(t, receiver.QOS, btypes.NewInt(int64(totalAmount-releasedAmountInit)))
}
