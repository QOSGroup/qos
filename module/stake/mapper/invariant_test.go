package mapper

import (
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	bmapper "github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/stake/types"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
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

	stakeMapper := NewMapper()
	stakeMapper.SetCodec(cdc)
	distributionKey := stakeMapper.GetStoreKey()
	mapperMap[types.MapperName] = stakeMapper

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(accountKey, btypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(distributionKey, btypes.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := context.NewContext(cms, abci.Header{}, false, log.NewNopLogger(), mapperMap)
	return ctx
}

func TestUnbondingInvariant(t *testing.T) {
	ctx := defaultContext()
	sm := GetMapper(ctx)

	val := btypes.ValAddress(ed25519.GenPrivKey().PubKey().Address())
	del := btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	unbonding := types.NewUnbondingDelegationInfo(del, val, 10, 100, 100)
	sm.SetUnbondingDelegation(unbonding)
	_, coins, broken := UnbondingInvariant("stake")(ctx)
	assert.False(t, broken)
	assert.Equal(t, coins.AmountOf(qtypes.QOSCoinName), btypes.NewInt(100))

	unbonding = types.NewUnbondingDelegationInfo(del, val, 10, 100, 100)
	sm.SetUnbondingDelegation(unbonding)
	_, coins, broken = UnbondingInvariant("stake")(ctx)
	assert.False(t, broken)
	assert.Equal(t, coins.AmountOf(qtypes.QOSCoinName), btypes.NewInt(100))
}

func TestRedelegationInvariant(t *testing.T) {
	ctx := defaultContext()
	sm := GetMapper(ctx)

	val := btypes.ValAddress(ed25519.GenPrivKey().PubKey().Address())
	del := btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	redelegation := types.NewRedelegateInfo(del, val, val, 100, 10, 100, false)
	sm.SetRedelegation(redelegation)
	_, coins, broken := RedelegationInvariant("stake")(ctx)
	assert.False(t, broken)
	assert.Equal(t, coins.AmountOf(qtypes.QOSCoinName), btypes.NewInt(100))

	redelegation = types.NewRedelegateInfo(del, val, val, 100, 10, 100, false)
	sm.AddRedelegation(redelegation)
	_, coins, broken = RedelegationInvariant("stake")(ctx)
	assert.False(t, broken)
	assert.Equal(t, coins.AmountOf(qtypes.QOSCoinName), btypes.NewInt(200))
}

func TestDelegationInvariant(t *testing.T) {
	ctx := defaultContext()
	sm := GetMapper(ctx)
	val := ed25519.GenPrivKey().PubKey()
	owner := btypes.ValAddress(ed25519.GenPrivKey().PubKey().Address())
	del1 := btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	del2 := btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	sm.CreateValidator(types.Validator{OperatorAddress: owner, ConsPubKey: val, BondTokens: 100})
	sm.SetDelegationInfo(types.NewDelegationInfo(del1, btypes.ValAddress(val.Address()), 100, false))
	_, coins, broken := DelegationInvariant("stake")(ctx)
	assert.True(t, broken)
	assert.Equal(t, coins.AmountOf(qtypes.QOSCoinName), btypes.NewInt(100))

	sm.SetDelegationInfo(types.NewDelegationInfo(del2, btypes.ValAddress(val.Address()), 100, false))
	_, coins, broken = DelegationInvariant("stake")(ctx)
	assert.True(t, broken)
	assert.Equal(t, coins.AmountOf(qtypes.QOSCoinName), btypes.NewInt(100))
}
