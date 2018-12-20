package staking

import (
	"encoding/binary"
	"fmt"
	"github.com/QOSGroup/qos/types"
	"time"

	"github.com/QOSGroup/qbase/mapper"

	btypes "github.com/QOSGroup/qbase/types"
)

const (
	ValidatorMapperName = "stakingvalidator"
)

var (
	//keys see docs/spec/staking.md
	validatorKey            = []byte{0x01} // 保存Validator信息. key: ValidatorAddress
	validatorByOwnerKey     = []byte{0x02} // 保存Owner与Validator的映射关系. key: OwnerAddress + ValidatorAddress
	validatorByInActiveKey  = []byte{0x03} // 保存处于`inactive`状态的Validator. key: ValidatorInActiveTime + ValidatorAddress
	validatorByVotePowerKey = []byte{0x04} // 按VotePower排序的Validator地址,不包含`pending`状态的Validator. key: VotePower + ValidatorAddress
)

func BuildValidatorStoreQueryPath() []byte {
	return []byte(fmt.Sprintf("/store/%s/key", ValidatorMapperName))
}

func BuildValidatorKey(valAddress btypes.Address) []byte {
	return append(validatorKey, valAddress...)
}

func BuildOwnerWithValidatorKey(ownerAddress btypes.Address, valAddress btypes.Address) []byte {

	lenz := 1 + len(ownerAddress) + len(valAddress)
	bz := make([]byte, lenz)

	copy(bz[0:1], validatorByOwnerKey)
	copy(bz[1:len(ownerAddress)+1], ownerAddress)
	copy(bz[1+len(ownerAddress):1+len(ownerAddress)+len(valAddress)], valAddress)

	return bz
}

func BuildInactiveValidatorKey(inactiveTime time.Time, valAddress btypes.Address) []byte {

	lenz := 1 + 8 + len(valAddress)
	bz := make([]byte, lenz)

	sec := inactiveTime.UTC().Unix()
	secBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(secBytes, uint64(sec))

	copy(bz[0:1], validatorByInActiveKey)
	copy(bz[1:9], secBytes)
	copy(bz[9:len(valAddress)+9], valAddress)

	return bz
}

func BuildValidatorByVotePower(votePower uint64, valAddress btypes.Address) []byte {
	lenz := 1 + 8 + len(valAddress)
	bz := make([]byte, lenz)

	votePowerBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(votePowerBytes, votePower)

	copy(bz[0:1], validatorByVotePowerKey)
	copy(bz[1:9], votePowerBytes)
	copy(bz[9:len(valAddress)+9], valAddress)

	return bz
}

type ValidatorMapper struct {
	*mapper.BaseMapper
}

var _ mapper.IMapper = (*ValidatorMapper)(nil)

func NewValidatorMapper() *ValidatorMapper {
	var validatorMapper = ValidatorMapper{}
	validatorMapper.BaseMapper = mapper.NewBaseMapper(nil, ValidatorMapperName)
	return &validatorMapper
}

func (mapper *ValidatorMapper) Copy() mapper.IMapper {
	validatorMapper := &ValidatorMapper{}
	validatorMapper.BaseMapper = mapper.BaseMapper.Copy()
	return validatorMapper
}

func (mapper *ValidatorMapper) SaveValidator(validator types.Validator) {
	mapper.Set(BuildValidatorKey(validator.ValidatorPubKey.Address().Bytes()), validator)
	mapper.Set(BuildOwnerWithValidatorKey(validator.Owner.Bytes(), validator.ValidatorPubKey.Address().Bytes()), validator)
	mapper.Set(BuildValidatorByVotePower(validator.BondTokens, validator.ValidatorPubKey.Address().Bytes()), validator)
}

func (mapper *ValidatorMapper) Exists(consAddress btypes.Address) bool {
	return mapper.Get(BuildValidatorKey(consAddress), &(types.Validator{}))
}
