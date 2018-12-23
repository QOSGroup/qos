package staking

import (
	"encoding/binary"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"

	qosacc "github.com/QOSGroup/qos/account"
	qosmapper "github.com/QOSGroup/qos/mapper"
	qostypes "github.com/QOSGroup/qos/types"

	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/tendermint/tendermint/crypto/ed25519"
)

func TestValidatorMapper(t *testing.T) {

	ctx := defaultContext()
	validatorMapper := GetValidatorMapper(ctx)

	validator := qostypes.Validator{
		Name:            "test",
		Owner:           btypes.Address(ed25519.GenPrivKey().PubKey().Address()),
		ValidatorPubKey: ed25519.GenPrivKey().PubKey(),
		BondTokens:      500,
		Status:          qostypes.Active,
		BondHeight:      1,
	}

	valAddr := btypes.Address(validator.ValidatorPubKey.Address())
	key := BuildValidatorKey(valAddr)
	validatorMapper.Set(key, validator)

	v, exsits := validatorMapper.GetValidator(valAddr)
	require.Equal(t, true, exsits)
	require.Equal(t, uint64(1), v.BondHeight)
	require.Equal(t, true, v.IsActive())

	now := uint64(time.Now().UTC().Unix())
	for i := uint64(0); i <= uint64(100); i++ {
		addr := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
		validatorMapper.Set(BuildInactiveValidatorKey(now+i, addr), i)
	}

	iter := validatorMapper.IteratorInactiveValidator(0, now+20)
	defer iter.Close()

	i := 0
	for ; iter.Valid(); iter.Next() {
		i++
		k := iter.Key()
		cp := binary.BigEndian.Uint64(k[1:9])
		require.Equal(t, true, cp >= now)
		now = cp
	}
	require.Equal(t, 20, i)

	for i := uint64(100); i <= uint64(200); i++ {
		addr := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
		validatorMapper.Set(BuildValidatorByVotePower(i, addr), 1)
	}

	descIter := validatorMapper.IteratorValidatrorByVoterPower(false)
	defer descIter.Close()

	power := uint64(200)
	for ; descIter.Valid(); descIter.Next() {
		p := descIter.Key()
		cp := binary.BigEndian.Uint64(p[1:9])

		require.Equal(t, true, power >= cp)
		power = cp
	}

	ascIter := validatorMapper.IteratorValidatrorByVoterPower(true)
	defer ascIter.Close()

	power = uint64(0)
	for ; ascIter.Valid(); ascIter.Next() {
		p := ascIter.Key()
		cp := binary.BigEndian.Uint64(p[1:9])

		require.Equal(t, true, power <= cp)
		power = cp
	}

}

func TestVoteInfoMapper(t *testing.T) {

	ctx := defaultContext()

	VoteInfoMapper := GetVoteInfoMapper(ctx)

	addr := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
	voteInfo := qostypes.NewValidatorVoteInfo(1, 1, 1)

	VoteInfoMapper.SetValidatorVoteInfo(addr, voteInfo)

	v, exsits := VoteInfoMapper.GetValidatorVoteInfo(addr)
	require.Equal(t, uint64(1), v.StartHeight)

	VoteInfoMapper.DelValidatorVoteInfo(addr)

	v, exsits = VoteInfoMapper.GetValidatorVoteInfo(addr)
	require.Equal(t, false, exsits)

	for i := uint64(0); i <= 10; i++ {
		VoteInfoMapper.SetVoteInfoInWindow(addr, i, false)
	}

	vote := VoteInfoMapper.GetVoteInfoInWindow(addr, 11)
	require.Equal(t, true, vote)

	vote = VoteInfoMapper.GetVoteInfoInWindow(addr, 10)
	require.Equal(t, false, vote)

	VoteInfoMapper.ClearValidatorVoteInfoInWindow(addr)

	vote = VoteInfoMapper.GetVoteInfoInWindow(addr, 10)
	require.Equal(t, true, vote)

}

func defaultContext() context.Context {

	mapperMap := make(map[string]mapper.IMapper)

	mainMapper := qosmapper.NewMainMapper()
	mapperMap[qosmapper.BaseMapperName] = mainMapper

	accountMapper := account.NewAccountMapper(cdc, qosacc.ProtoQOSAccount)
	mapperMap[account.AccountMapperName] = accountMapper

	validatorMapper := NewValidatorMapper()
	validatorMapper.SetCodec(cdc)
	mapperMap[ValidatorMapperName] = validatorMapper

	signInfoMapper := NewVoteInfoMapper()
	signInfoMapper.SetCodec(cdc)
	mapperMap[VoteInfoMapperName] = signInfoMapper

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)

	for _, v := range mapperMap {
		cms.MountStoreWithDB(v.GetStoreKey(), store.StoreTypeIAVL, db)
	}
	cms.LoadLatestVersion()

	ctx := context.NewContext(cms, abci.Header{}, false, log.NewNopLogger(), mapperMap)
	return ctx
}
