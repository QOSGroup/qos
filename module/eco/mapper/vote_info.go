package mapper

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"

	"github.com/QOSGroup/qos/module/eco/types"

	btypes "github.com/QOSGroup/qbase/types"
)

type VoteInfoMapper struct {
	*mapper.BaseMapper
}

var _ mapper.IMapper = (*VoteInfoMapper)(nil)

func NewVoteInfoMapper() *VoteInfoMapper {
	var VoteInfoMapper = VoteInfoMapper{}
	VoteInfoMapper.BaseMapper = mapper.NewBaseMapper(nil, types.VoteInfoMapperName)
	return &VoteInfoMapper
}

func GetVoteInfoMapper(ctx context.Context) *VoteInfoMapper {
	return ctx.Mapper(types.VoteInfoMapperName).(*VoteInfoMapper)
}

func (mapper *VoteInfoMapper) Copy() mapper.IMapper {
	VoteInfoMapper := &VoteInfoMapper{}
	VoteInfoMapper.BaseMapper = mapper.BaseMapper.Copy()
	return VoteInfoMapper
}

func (mapper *VoteInfoMapper) GetValidatorVoteInfo(valAddr btypes.Address) (info types.ValidatorVoteInfo, exsits bool) {
	key := types.BuildValidatorVoteInfoKey(valAddr)
	exsits = mapper.Get(key, &info)
	return
}

func (mapper *VoteInfoMapper) SetValidatorVoteInfo(valAddr btypes.Address, info types.ValidatorVoteInfo) {
	key := types.BuildValidatorVoteInfoKey(valAddr)
	mapper.Set(key, info)
}

func (mapper *VoteInfoMapper) ResetValidatorVoteInfo(valAddr btypes.Address, info types.ValidatorVoteInfo) {
	key := types.BuildValidatorVoteInfoKey(valAddr)
	mapper.ClearValidatorVoteInfoInWindow(valAddr)
	mapper.Del(key)
}

func (mapper *VoteInfoMapper) DelValidatorVoteInfo(valAddr btypes.Address) {
	key := types.BuildValidatorVoteInfoKey(valAddr)
	mapper.Del(key)
}

func (mapper *VoteInfoMapper) GetVoteInfoInWindow(valAddr btypes.Address, index uint64) (vote bool) {
	key := types.BuildValidatorVoteInfoInWindowKey(index, valAddr)
	vote, exsits := mapper.GetBool(key)

	if !exsits {
		return true
	}

	return vote
}

func (mapper *VoteInfoMapper) SetVoteInfoInWindow(valAddr btypes.Address, index uint64, vote bool) {
	key := types.BuildValidatorVoteInfoInWindowKey(index, valAddr)
	mapper.Set(key, vote)
}

func (mapper *VoteInfoMapper) ClearValidatorVoteInfoInWindow(valAddr btypes.Address) {
	prefixKey := append(types.GetValidatorVoteInfoInWindowKey(), valAddr...)
	endKey := store.PrefixEndBytes(prefixKey)
	iter := mapper.GetStore().Iterator(prefixKey, endKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		mapper.Del(iter.Key())
	}
}

//-------------------------genesis export

func (mapper *VoteInfoMapper) IterateVoteInfos(fn func(btypes.Address, types.ValidatorVoteInfo)) {
	iter := store.KVStorePrefixIterator(mapper.GetStore(), types.GetValidatorVoteInfoKey())
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		valAddr := types.GetValidatorVoteInfoAddr(key)
		var info types.ValidatorVoteInfo
		mapper.DecodeObject(iter.Value(), &info)
		fn(valAddr, info)
	}
}

func (mapper *VoteInfoMapper) IterateVoteInWindowsInfos(fn func(uint64, btypes.Address, bool)) {
	iter := store.KVStorePrefixIterator(mapper.GetStore(), types.GetValidatorVoteInfoInWindowKey())
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		index, addr := types.GetValidatorVoteInfoInWindowIndexAddr(key)
		var vote bool
		mapper.DecodeObject(iter.Value(), &vote)
		fn(index, addr, vote)
	}
}
