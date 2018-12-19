package staking

import (
	"encoding/binary"
	"fmt"

	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"

	"github.com/QOSGroup/qos/types"

	btypes "github.com/QOSGroup/qbase/types"
)

const (
	stakingMapperName = "staking"
)

var (
	//keys see docs/spec/staking.md
	validatorSignInfoKey         = []byte{0x01} // 保存Validator在窗口的统计信息
	validatorSignInfoInWindowKey = []byte{0x02} // 保存Validator在指定窗口签名信息
)

func BuildStakingStoreQueryPath() []byte {
	return []byte(fmt.Sprintf("/store/%s/key", stakingMapperName))
}

func BuildValidatorSignInfoKey(valAddress btypes.Address) []byte {
	return append(validatorSignInfoKey, valAddress...)
}

func BuildValidatorSignInfoInWindowKey(index uint64, valAddress btypes.Address) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, index)

	bz := append(validatorSignInfoInWindowKey, valAddress...)
	bz = append(bz, b...)

	return bz
}

type StakingMapper struct {
	*mapper.BaseMapper
}

var _ mapper.IMapper = (*StakingMapper)(nil)

func NewStakingMapper() *StakingMapper {
	var stakingMapper = StakingMapper{}
	stakingMapper.BaseMapper = mapper.NewBaseMapper(nil, stakingMapperName)
	return &stakingMapper
}

func (mapper *StakingMapper) Copy() mapper.IMapper {
	stakingMapper := &StakingMapper{}
	stakingMapper.BaseMapper = mapper.BaseMapper.Copy()
	return stakingMapper
}

func (mapper *StakingMapper) GetValidatorSignInfo(valAddr btypes.Address) (signInfo types.ValidatorSignInfo, exsits bool) {
	key := BuildValidatorSignInfoKey(valAddr)
	exsits = mapper.Get(key, &signInfo)
	return
}

func (mapper *StakingMapper) SetValidatorSignInfo(valAddr btypes.Address, info types.ValidatorSignInfo) {
	key := BuildValidatorSignInfoKey(valAddr)
	mapper.Set(key, info)
}

func (mapper *StakingMapper) DelValidatorSignInfo(valAddr btypes.Address) {
	key := BuildValidatorSignInfoKey(valAddr)
	mapper.Del(key)
}

func (mapper *StakingMapper) GetVoteInWindowIndexOffset(valAddr btypes.Address, index uint64) (vote bool) {
	key := BuildValidatorSignInfoInWindowKey(index, valAddr)
	isVote, exsits := mapper.GetBool(key)

	if !exsits {
		return false
	}

	return isVote
}

func (mapper *StakingMapper) SetVoteInWindowIndexOffset(valAddr btypes.Address, index uint64, vote bool) {
	key := BuildValidatorSignInfoInWindowKey(index, valAddr)
	mapper.Set(key, vote)
}

func (mapper *StakingMapper) ClearValidatorVoteInWindowIndex(valAddr btypes.Address) {
	prefixKey := append(validatorSignInfoInWindowKey, valAddr...)
	endKey := store.PrefixEndBytes(prefixKey)
	iter := mapper.GetStore().Iterator(prefixKey, endKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		mapper.Del(iter.Key())
	}
}
