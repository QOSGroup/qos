package mapper

import (
	"encoding/binary"
	"time"

	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	staketypes "github.com/QOSGroup/qos/module/eco/types"
)

type ValidatorMapper struct {
	*mapper.BaseMapper
}

var _ mapper.IMapper = (*ValidatorMapper)(nil)

func NewValidatorMapper() *ValidatorMapper {
	var validatorMapper = ValidatorMapper{}
	validatorMapper.BaseMapper = mapper.NewBaseMapper(nil, staketypes.ValidatorMapperName)
	return &validatorMapper
}

func GetValidatorMapper(ctx context.Context) *ValidatorMapper {
	return ctx.Mapper(staketypes.ValidatorMapperName).(*ValidatorMapper)
}

func (mapper *ValidatorMapper) Copy() mapper.IMapper {
	validatorMapper := &ValidatorMapper{}
	validatorMapper.BaseMapper = mapper.BaseMapper.Copy()
	return validatorMapper
}

func (mapper *ValidatorMapper) CreateValidator(validator staketypes.Validator) {
	mapper.Set(staketypes.BuildValidatorKey(validator.ValidatorPubKey.Address().Bytes()), validator)
	mapper.Set(staketypes.BuildOwnerWithValidatorKey(validator.Owner), validator.ValidatorPubKey.Address().Bytes())
	mapper.Set(staketypes.BuildValidatorByVotePower(validator.BondTokens, validator.ValidatorPubKey.Address().Bytes()), 1)
}

func (mapper *ValidatorMapper) Exists(valAddress btypes.Address) bool {
	return mapper.Get(staketypes.BuildValidatorKey(valAddress), &(staketypes.Validator{}))
}

func (mapper *ValidatorMapper) ExistsWithOwner(owner btypes.Address) bool {
	return mapper.Get(staketypes.BuildOwnerWithValidatorKey(owner), &(btypes.Address{}))
}

func (mapper *ValidatorMapper) GetValidator(valAddress btypes.Address) (validator staketypes.Validator, exsits bool) {
	validatorKey := staketypes.BuildValidatorKey(valAddress)
	exsits = mapper.Get(validatorKey, &validator)
	return
}

func (mapper *ValidatorMapper) MakeValidatorInactive(valAddress btypes.Address, inactiveHeight uint64, inactiveTime time.Time, code staketypes.InactiveCode) {
	validator, exsits := mapper.GetValidator(valAddress)
	if !exsits {
		return
	}

	validator.Status = staketypes.Inactive
	validator.InactiveCode = code
	validator.InactiveHeight = inactiveHeight
	validator.InactiveTime = inactiveTime.UTC()
	mapper.Set(staketypes.BuildValidatorKey(valAddress), validator)

	validatorInactiveKey := staketypes.BuildInactiveValidatorKeyByTime(inactiveTime, valAddress)
	mapper.Set(validatorInactiveKey, inactiveTime.UTC().Unix())

	validatorVotePowerKey := staketypes.BuildValidatorByVotePower(validator.BondTokens, valAddress)
	mapper.Del(validatorVotePowerKey)
}

func (mapper *ValidatorMapper) KickValidator(valAddress btypes.Address) (validator staketypes.Validator, ok bool) {
	validator, exsits := mapper.GetValidator(valAddress)
	if !exsits {
		return validator, false
	}

	mapper.Del(staketypes.BuildValidatorKey(valAddress))
	mapper.Del(staketypes.BuildOwnerWithValidatorKey(validator.Owner))
	mapper.Del(staketypes.BuildInactiveValidatorKeyByTime(validator.InactiveTime, valAddress))
	mapper.Del(staketypes.BuildValidatorByVotePower(validator.BondTokens, valAddress))

	return validator, true
}

func (mapper *ValidatorMapper) IteratorInactiveValidator(fromSecond, endSecond uint64) store.Iterator {

	secBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(secBytes, fromSecond)
	startKey := append(staketypes.GetValidatorByInactiveKey(), secBytes...)

	binary.BigEndian.PutUint64(secBytes, endSecond)
	endKey := append(staketypes.GetValidatorByInactiveKey(), secBytes...)

	return mapper.GetStore().Iterator(startKey, endKey)
}

func (mapper *ValidatorMapper) IteratorInactiveValidatorByTime(fromTime, endTime time.Time) store.Iterator {
	return mapper.IteratorInactiveValidator(uint64(fromTime.UTC().Unix()), uint64(endTime.UTC().Unix()))
}

func (mapper *ValidatorMapper) IteratorValidatrorByVoterPower(ascending bool) store.Iterator {
	if ascending {
		return store.KVStorePrefixIterator(mapper.GetStore(), staketypes.GetValidatorByVotePowerKey())
	}
	return store.KVStoreReversePrefixIterator(mapper.GetStore(), staketypes.GetValidatorByVotePowerKey())
}

func (mapper *ValidatorMapper) MakeValidatorActive(valAddress btypes.Address) {
	validator, exsits := mapper.GetValidator(valAddress)
	if !exsits {
		return
	}

	validator.Status = staketypes.Active

	mapper.Set(staketypes.BuildValidatorKey(validator.ValidatorPubKey.Address().Bytes()), validator)
	mapper.Del(staketypes.BuildInactiveValidatorKey(uint64(validator.InactiveTime.UTC().Unix()), valAddress))
	mapper.Set(staketypes.BuildValidatorByVotePower(validator.BondTokens, validator.ValidatorPubKey.Address().Bytes()), 1)
}

func (mapper *ValidatorMapper) GetValidatorByOwner(owner btypes.Address) (validator staketypes.Validator, exsits bool) {
	var valAddress btypes.Address
	exsits = mapper.Get(staketypes.BuildOwnerWithValidatorKey(owner), &valAddress)
	if !exsits {
		return validator, false
	}

	return mapper.GetValidator(valAddress)
}

func (mapper *ValidatorMapper) SetParams(params staketypes.StakeParams) {
	mapper.Set(staketypes.BuildStakeParamsKey(), params)
}

func (mapper *ValidatorMapper) GetParams() staketypes.StakeParams {
	params := staketypes.StakeParams{}
	mapper.Get(staketypes.BuildStakeParamsKey(), &params)
	return params
}
