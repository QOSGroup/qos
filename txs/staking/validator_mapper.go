package staking

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

<<<<<<< HEAD
func BuildInactiveValidatorKey(inactiveTime time.Time, valAddress btypes.Address) []byte {
=======
func BuildInActiveValidatorKeyByTime(inActiveTime time.Time, valAddress btypes.Address) []byte {
	return BuildInActiveValidatorKey(uint64(inActiveTime.UTC().Unix()), valAddress)
}

func BuildInActiveValidatorKey(sec uint64, valAddress btypes.Address) []byte {
>>>>>>> 1d82e332dbb31187bce816e1449376a4670642ce

	lenz := 1 + 8 + len(valAddress)
	bz := make([]byte, lenz)

<<<<<<< HEAD
	sec := inactiveTime.UTC().Unix()
=======
>>>>>>> 1d82e332dbb31187bce816e1449376a4670642ce
	secBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(secBytes, sec)

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

func GetValidatorMapper(ctx context.Context) *ValidatorMapper {
	return ctx.Mapper(ValidatorMapperName).(*ValidatorMapper)
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

func (mapper *ValidatorMapper) GetValidator(valAddress btypes.Address) (validator types.Validator, exsits bool) {
	validatorKey := BuildValidatorKey(valAddress)
	exsits = mapper.Get(validatorKey, &validator)
	return
}

func (mapper *ValidatorMapper) MakeValidatorInActive(valAddress btypes.Address, inActiveHeight uint64, inActiveTime time.Time, isRevoke bool) {
	validator, exsits := mapper.GetValidator(valAddress)
	if !exsits {
		return
	}

	validator.Status = types.InActive
	validator.IsRevoke = isRevoke
	validator.InActiveHeight = inActiveHeight
	validator.InActiveTime = inActiveTime.UTC()
	mapper.Set(validatorKey, validator)

	validatorInActiveKey := BuildInActiveValidatorKeyByTime(inActiveTime, valAddress)
	mapper.Set(validatorInActiveKey, inActiveTime.UTC().Unix())

	validatorVotePowerKey := BuildValidatorByVotePower(validator.BondTokens, valAddress)
	mapper.Del(validatorVotePowerKey)
}

func (mapper *ValidatorMapper) KickValidator(valAddress btypes.Address) (validator types.Validator, ok bool) {
	validator, exsits := mapper.GetValidator(valAddress)
	if !exsits {
		return validator, false
	}

	mapper.Del(BuildValidatorKey(valAddress))
	mapper.Del(BuildOwnerWithValidatorKey(validator.Owner, valAddress))
	mapper.Del(BuildInActiveValidatorKeyByTime(validator.InActiveTime, valAddress))
	mapper.Del(BuildValidatorByVotePower(validator.BondTokens, valAddress))

	return validator, true
}

func (mapper *ValidatorMapper) IteratorInActiveValidator(fromSecond, endSecond uint64) store.Iterator {

	secBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(secBytes, fromSecond)
	startKey := append(validatorByInActiveKey, secBytes...)

	binary.BigEndian.PutUint64(secBytes, endSecond)
	endKey := append(validatorByInActiveKey, secBytes...)

	return mapper.GetStore().Iterator(startKey, endKey)
}

func (mapper *ValidatorMapper) IteratorInActiveValidatorByTime(fromTime, endTime time.Time) store.Iterator {
	return mapper.IteratorInActiveValidator(uint64(fromTime.UTC().Unix()), uint64(endTime.UTC().Unix()))
}

func (mapper *ValidatorMapper) IteratorValidatrorByVoterPower(ascending bool) store.Iterator {
	if ascending {
		return store.KVStorePrefixIterator(mapper.GetStore(), validatorByVotePowerKey)
	}
	return store.KVStoreReversePrefixIterator(mapper.GetStore(), validatorByVotePowerKey)
}
