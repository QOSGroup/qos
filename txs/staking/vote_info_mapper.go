package staking

import (
	"encoding/binary"
	"fmt"

	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"

	"github.com/QOSGroup/qos/types"

	btypes "github.com/QOSGroup/qbase/types"
)

const (
	VoteInfoMapperName = "voteInfo"
)

var (
	//keys see docs/spec/staking.md
	validatorVoteInfoKey         = []byte{0x01} // 保存Validator在窗口的统计信息
	validatorVoteInfoInWindowKey = []byte{0x02} // 保存Validator在指定窗口签名信息
)

func BuildVoteInfoStoreQueryPath() []byte {
	return []byte(fmt.Sprintf("/store/%s/key", VoteInfoMapperName))
}

func BuildValidatorVoteInfoKey(valAddress btypes.Address) []byte {
	return append(validatorVoteInfoKey, valAddress...)
}

func BuildValidatorVoteInfoInWindowPrefixKey(valAddress btypes.Address) []byte {
	return append(validatorVoteInfoInWindowKey, valAddress...)
}

func BuildValidatorVoteInfoInWindowKey(index uint64, valAddress btypes.Address) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, index)

	bz := append(validatorVoteInfoInWindowKey, valAddress...)
	bz = append(bz, b...)

	return bz
}

type VoteInfoMapper struct {
	*mapper.BaseMapper
}

var _ mapper.IMapper = (*VoteInfoMapper)(nil)

func NewVoteInfoMapper() *VoteInfoMapper {
	var VoteInfoMapper = VoteInfoMapper{}
	VoteInfoMapper.BaseMapper = mapper.NewBaseMapper(nil, VoteInfoMapperName)
	return &VoteInfoMapper
}

func GetVoteInfoMapper(ctx context.Context) *VoteInfoMapper {
	return ctx.Mapper(VoteInfoMapperName).(*VoteInfoMapper)
}

func (mapper *VoteInfoMapper) Copy() mapper.IMapper {
	VoteInfoMapper := &VoteInfoMapper{}
	VoteInfoMapper.BaseMapper = mapper.BaseMapper.Copy()
	return VoteInfoMapper
}

func (mapper *VoteInfoMapper) GetValidatorVoteInfo(valAddr btypes.Address) (VoteInfo types.ValidatorVoteInfo, exsits bool) {
	key := BuildValidatorVoteInfoKey(valAddr)
	exsits = mapper.Get(key, &VoteInfo)
	return
}

func (mapper *VoteInfoMapper) SetValidatorVoteInfo(valAddr btypes.Address, info types.ValidatorVoteInfo) {
	key := BuildValidatorVoteInfoKey(valAddr)
	mapper.Set(key, info)
}

func (mapper *VoteInfoMapper) ResetValidatorVoteInfo(valAddr btypes.Address, info types.ValidatorVoteInfo) {
	key := BuildValidatorVoteInfoKey(valAddr)
	mapper.ClearValidatorVoteInfoInWindow(valAddr)
	mapper.Del(key)
}

func (mapper *VoteInfoMapper) DelValidatorVoteInfo(valAddr btypes.Address) {
	key := BuildValidatorVoteInfoKey(valAddr)
	mapper.Del(key)
}

func (mapper *VoteInfoMapper) GetVoteInfoInWindow(valAddr btypes.Address, index uint64) (vote bool) {
	key := BuildValidatorVoteInfoInWindowKey(index, valAddr)
	vote, exsits := mapper.GetBool(key)

	if !exsits {
		return true
	}

	return vote
}

func (mapper *VoteInfoMapper) SetVoteInfoInWindow(valAddr btypes.Address, index uint64, vote bool) {
	key := BuildValidatorVoteInfoInWindowKey(index, valAddr)
	mapper.Set(key, vote)
}

func (mapper *VoteInfoMapper) ClearValidatorVoteInfoInWindow(valAddr btypes.Address) {
	prefixKey := append(validatorVoteInfoInWindowKey, valAddr...)
	endKey := store.PrefixEndBytes(prefixKey)
	iter := mapper.GetStore().Iterator(prefixKey, endKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		mapper.Del(iter.Key())
	}
}
