package stake

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	"github.com/QOSGroup/qos/types"

	btypes "github.com/QOSGroup/qbase/types"
)

const (
	ValidatorMapperName = "validator"
)

var (
	//keys see docs/spec/staking.md
	validatorKey            = []byte{0x01} // 保存Validator信息. key: ValidatorAddress
	validatorByOwnerKey     = []byte{0x02} // 保存Owner与Validator的映射关系. key: OwnerAddress, value : ValidatorAddress
	validatorByInactiveKey  = []byte{0x03} // 保存处于`inactive`状态的Validator. key: ValidatorInactiveTime + ValidatorAddress
	validatorByVotePowerKey = []byte{0x04} // 按VotePower排序的Validator地址,不包含`pending`状态的Validator. key: VotePower + ValidatorAddress

	currentValidatorAddressKey = []byte("currentValidatorAddressKey")
)

func BuildValidatorStoreQueryPath() []byte {
	return []byte(fmt.Sprintf("/store/%s/key", ValidatorMapperName))
}

func BuildCurrentValidatorAddressKey() []byte {
	return currentValidatorAddressKey
}

func BuildValidatorKey(valAddress btypes.Address) []byte {
	return append(validatorKey, valAddress...)
}

func BulidValidatorPrefixKey() []byte {
	return validatorKey
}

func BuildOwnerWithValidatorKey(ownerAddress btypes.Address) []byte {

	lenz := 1 + len(ownerAddress)
	bz := make([]byte, lenz)

	copy(bz[0:1], validatorByOwnerKey)
	copy(bz[1:len(ownerAddress)+1], ownerAddress)

	return bz
}

func BuildInactiveValidatorKeyByTime(inactiveTime time.Time, valAddress btypes.Address) []byte {
	return BuildInactiveValidatorKey(uint64(inactiveTime.UTC().Unix()), valAddress)
}

func BuildInactiveValidatorKey(sec uint64, valAddress btypes.Address) []byte {
	lenz := 1 + 8 + len(valAddress)
	bz := make([]byte, lenz)

	secBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(secBytes, sec)

	copy(bz[0:1], validatorByInactiveKey)
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

func GetValidatorMapper(ctx context.Context) *ValidatorMapper {
	return ctx.Mapper(ValidatorMapperName).(*ValidatorMapper)
}

func (mapper *ValidatorMapper) Copy() mapper.IMapper {
	validatorMapper := &ValidatorMapper{}
	validatorMapper.BaseMapper = mapper.BaseMapper.Copy()
	return validatorMapper
}

func (mapper *ValidatorMapper) CreateValidator(validator types.Validator) {
	mapper.Set(BuildValidatorKey(validator.ValidatorPubKey.Address().Bytes()), validator)
	mapper.Set(BuildOwnerWithValidatorKey(validator.Owner), validator.ValidatorPubKey.Address().Bytes())
	mapper.Set(BuildValidatorByVotePower(validator.BondTokens, validator.ValidatorPubKey.Address().Bytes()), 1)
}

func (mapper *ValidatorMapper) Exists(valAddress btypes.Address) bool {
	return mapper.Get(BuildValidatorKey(valAddress), &(types.Validator{}))
}

func (mapper *ValidatorMapper) ExistsWithOwner(owner btypes.Address) bool {
	return mapper.Get(BuildOwnerWithValidatorKey(owner), &(btypes.Address{}))
}

func (mapper *ValidatorMapper) GetValidator(valAddress btypes.Address) (validator types.Validator, exsits bool) {
	validatorKey := BuildValidatorKey(valAddress)
	exsits = mapper.Get(validatorKey, &validator)
	return
}

func (mapper *ValidatorMapper) MakeValidatorInactive(valAddress btypes.Address, inactiveHeight uint64, inactiveTime time.Time, code types.InactiveCode) {
	validator, exsits := mapper.GetValidator(valAddress)
	if !exsits {
		return
	}

	validator.Status = types.Inactive
	validator.InactiveCode = code
	validator.InactiveHeight = inactiveHeight
	validator.InactiveTime = inactiveTime.UTC()
	mapper.Set(BuildValidatorKey(valAddress), validator)

	validatorInactiveKey := BuildInactiveValidatorKeyByTime(inactiveTime, valAddress)
	mapper.Set(validatorInactiveKey, inactiveTime.UTC().Unix())

	validatorVotePowerKey := BuildValidatorByVotePower(validator.BondTokens, valAddress)
	mapper.Del(validatorVotePowerKey)
}

func (mapper *ValidatorMapper) KickValidator(valAddress btypes.Address) (validator types.Validator, ok bool) {
	validator, exsits := mapper.GetValidator(valAddress)
	if !exsits {
		return validator, false
	}

	mapper.Del(BuildValidatorKey(valAddress))
	mapper.Del(BuildOwnerWithValidatorKey(validator.Owner))
	mapper.Del(BuildInactiveValidatorKeyByTime(validator.InactiveTime, valAddress))
	mapper.Del(BuildValidatorByVotePower(validator.BondTokens, valAddress))

	return validator, true
}

func (mapper *ValidatorMapper) IteratorInactiveValidator(fromSecond, endSecond uint64) store.Iterator {

	secBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(secBytes, fromSecond)
	startKey := append(validatorByInactiveKey, secBytes...)

	binary.BigEndian.PutUint64(secBytes, endSecond)
	endKey := append(validatorByInactiveKey, secBytes...)

	return mapper.GetStore().Iterator(startKey, endKey)
}

func (mapper *ValidatorMapper) IteratorInactiveValidatorByTime(fromTime, endTime time.Time) store.Iterator {
	return mapper.IteratorInactiveValidator(uint64(fromTime.UTC().Unix()), uint64(endTime.UTC().Unix()))
}

func (mapper *ValidatorMapper) IteratorValidatrorByVoterPower(ascending bool) store.Iterator {
	if ascending {
		return store.KVStorePrefixIterator(mapper.GetStore(), validatorByVotePowerKey)
	}
	return store.KVStoreReversePrefixIterator(mapper.GetStore(), validatorByVotePowerKey)
}

func (mapper *ValidatorMapper) MakeValidatorActive(valAddress btypes.Address) {
	validator, exsits := mapper.GetValidator(valAddress)
	if !exsits {
		return
	}

	validator.Status = types.Active

	mapper.Set(BuildValidatorKey(validator.ValidatorPubKey.Address().Bytes()), validator)
	mapper.Del(BuildInactiveValidatorKey(uint64(validator.InactiveTime.UTC().Unix()), valAddress))
	mapper.Set(BuildValidatorByVotePower(validator.BondTokens, validator.ValidatorPubKey.Address().Bytes()), 1)
}

func (mapper *ValidatorMapper) GetValidatorByOwner(owner btypes.Address) (validator types.Validator, exsits bool) {
	var valAddress btypes.Address
	exsits = mapper.Get(BuildOwnerWithValidatorKey(owner), &valAddress)
	if !exsits {
		return validator, false
	}

	return mapper.GetValidator(valAddress)
}
