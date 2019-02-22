package mapper

import (
	"encoding/binary"
	"time"

	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	ecotypes "github.com/QOSGroup/qos/module/eco/types"
)

type ValidatorMapper struct {
	*mapper.BaseMapper
}

var _ mapper.IMapper = (*ValidatorMapper)(nil)

func NewValidatorMapper() *ValidatorMapper {
	var validatorMapper = ValidatorMapper{}
	validatorMapper.BaseMapper = mapper.NewBaseMapper(nil, ecotypes.ValidatorMapperName)
	return &validatorMapper
}

func GetValidatorMapper(ctx context.Context) *ValidatorMapper {
	return ctx.Mapper(ecotypes.ValidatorMapperName).(*ValidatorMapper)
}

func (mapper *ValidatorMapper) Copy() mapper.IMapper {
	validatorMapper := &ValidatorMapper{}
	validatorMapper.BaseMapper = mapper.BaseMapper.Copy()
	return validatorMapper
}

func (mapper *ValidatorMapper) CreateValidator(validator ecotypes.Validator) {
	valAddr := validator.GetValidatorAddress()
	mapper.Set(ecotypes.BuildValidatorKey(valAddr), validator)
	mapper.Set(ecotypes.BuildOwnerWithValidatorKey(validator.Owner), valAddr)
	mapper.Set(ecotypes.BuildValidatorByVotePower(validator.BondTokens, valAddr), true)
}

func (mapper *ValidatorMapper) ChangeValidatorBondTokens(validator ecotypes.Validator, updatedTokens uint64) {
	valAddr := validator.GetValidatorAddress()
	mapper.Del(ecotypes.BuildValidatorByVotePower(validator.BondTokens, valAddr))
	validator.BondTokens = updatedTokens
	mapper.CreateValidator(validator)
}

func (mapper *ValidatorMapper) Exists(valAddress btypes.Address) bool {
	return mapper.Get(ecotypes.BuildValidatorKey(valAddress), &(ecotypes.Validator{}))
}

func (mapper *ValidatorMapper) ExistsWithOwner(owner btypes.Address) bool {
	return mapper.Get(ecotypes.BuildOwnerWithValidatorKey(owner), &(btypes.Address{}))
}

func (mapper *ValidatorMapper) GetValidator(valAddress btypes.Address) (validator ecotypes.Validator, exsits bool) {
	validatorKey := ecotypes.BuildValidatorKey(valAddress)
	exsits = mapper.Get(validatorKey, &validator)
	return
}

func (mapper *ValidatorMapper) MakeValidatorInactive(valAddress btypes.Address, inactiveHeight uint64, inactiveTime time.Time, code ecotypes.InactiveCode) {
	validator, exsits := mapper.GetValidator(valAddress)
	if !exsits {
		return
	}

	validator.Status = ecotypes.Inactive
	validator.InactiveCode = code
	validator.InactiveHeight = inactiveHeight
	validator.InactiveTime = inactiveTime.UTC()
	mapper.Set(ecotypes.BuildValidatorKey(valAddress), validator)

	validatorInactiveKey := ecotypes.BuildInactiveValidatorKeyByTime(inactiveTime, valAddress)
	mapper.Set(validatorInactiveKey, inactiveTime.UTC().Unix())

	validatorVotePowerKey := ecotypes.BuildValidatorByVotePower(validator.BondTokens, valAddress)
	mapper.Del(validatorVotePowerKey)
}

func (mapper *ValidatorMapper) KickValidator(valAddress btypes.Address) (validator ecotypes.Validator, ok bool) {
	validator, exsits := mapper.GetValidator(valAddress)
	if !exsits {
		return validator, false
	}

	mapper.Del(ecotypes.BuildValidatorKey(valAddress))
	mapper.Del(ecotypes.BuildOwnerWithValidatorKey(validator.Owner))
	mapper.Del(ecotypes.BuildInactiveValidatorKeyByTime(validator.InactiveTime, valAddress))
	mapper.Del(ecotypes.BuildValidatorByVotePower(validator.BondTokens, valAddress))

	return validator, true
}

func (mapper *ValidatorMapper) IteratorInactiveValidator(fromSecond, endSecond uint64) store.Iterator {

	secBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(secBytes, fromSecond)
	startKey := append(ecotypes.GetValidatorByInactiveKey(), secBytes...)

	binary.BigEndian.PutUint64(secBytes, endSecond)
	endKey := append(ecotypes.GetValidatorByInactiveKey(), secBytes...)

	return mapper.GetStore().Iterator(startKey, endKey)
}

func (mapper *ValidatorMapper) IteratorInactiveValidatorByTime(fromTime, endTime time.Time) store.Iterator {
	return mapper.IteratorInactiveValidator(uint64(fromTime.UTC().Unix()), uint64(endTime.UTC().Unix()))
}

func (mapper *ValidatorMapper) IteratorValidatrorByVoterPower(ascending bool) store.Iterator {
	if ascending {
		return store.KVStorePrefixIterator(mapper.GetStore(), ecotypes.GetValidatorByVotePowerKey())
	}
	return store.KVStoreReversePrefixIterator(mapper.GetStore(), ecotypes.GetValidatorByVotePowerKey())
}

func (mapper *ValidatorMapper) MakeValidatorActive(valAddress btypes.Address) {
	validator, exsits := mapper.GetValidator(valAddress)
	if !exsits {
		return
	}

	validator.Status = ecotypes.Active

	mapper.Set(ecotypes.BuildValidatorKey(validator.ValidatorPubKey.Address().Bytes()), validator)
	mapper.Del(ecotypes.BuildInactiveValidatorKey(uint64(validator.InactiveTime.UTC().Unix()), valAddress))
	mapper.Set(ecotypes.BuildValidatorByVotePower(validator.BondTokens, validator.ValidatorPubKey.Address().Bytes()), 1)
}

func (mapper *ValidatorMapper) GetValidatorByOwner(owner btypes.Address) (validator ecotypes.Validator, exsits bool) {
	var valAddress btypes.Address
	exsits = mapper.Get(ecotypes.BuildOwnerWithValidatorKey(owner), &valAddress)
	if !exsits {
		return validator, false
	}

	return mapper.GetValidator(valAddress)
}

func (mapper *ValidatorMapper) SetParams(params ecotypes.StakeParams) {
	mapper.Set(ecotypes.BuildStakeParamsKey(), params)
}

func (mapper *ValidatorMapper) GetParams() ecotypes.StakeParams {
	params := ecotypes.StakeParams{}
	mapper.Get(ecotypes.BuildStakeParamsKey(), &params)
	return params
}

//-------------------------genesis export

func (mapper *ValidatorMapper) IterateValidators(fn func(ecotypes.Validator)) {

	iter := store.KVStorePrefixIterator(mapper.GetStore(), ecotypes.BulidValidatorPrefixKey())
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var validator ecotypes.Validator
		mapper.DecodeObject(iter.Value(), &validator)
		fn(validator)
	}
}
