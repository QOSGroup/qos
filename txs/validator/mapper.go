package validator

import (
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/types"
)

const (
	ValidatorMapperName = "validator"
	validatorKey        = "validator:"
	validatorChangedKey = "validator_changed"
)

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

func BuildChangedKey() []byte {
	return []byte(validatorChangedKey)
}

func BuildValidatorKey(consAddress btypes.Address) []byte {
	return append([]byte(validatorKey), consAddress...)
}

func (mapper *ValidatorMapper) IsValidatorChanged() bool {
	v, _ := mapper.GetBool(BuildChangedKey())
	return v
}

func (mapper *ValidatorMapper) SetValidatorChanged() {
	mapper.Set(BuildChangedKey(), true)
}

func (mapper *ValidatorMapper) SetValidatorUnChanged() {
	mapper.Set(BuildChangedKey(), false)
}

// 保存validator
func (mapper *ValidatorMapper) SaveValidator(validator types.Validator) {
	mapper.Set(BuildValidatorKey(validator.ConsPubKey.Address().Bytes()), validator)
}

// 是否已经存在
func (mapper *ValidatorMapper) Exists(consAddress btypes.Address) bool {
	return mapper.Get(BuildValidatorKey(consAddress), &(types.Validator{}))
}

func (mapper *ValidatorMapper) GetByConsAddress(consAddress btypes.Address) (types.Validator, bool) {
	var val types.Validator
	exsits := mapper.Get(BuildValidatorKey(consAddress), &val)
	return val, exsits
}

// 获取所有validator
func (mapper *ValidatorMapper) GetValidators() (validators []types.Validator) {
	iterator := mapper.GetStore().Iterator([]byte(validatorKey), store.PrefixEndBytes([]byte(validatorKey)))
	for ; iterator.Valid(); iterator.Next() {
		validator := types.Validator{}
		mapper.GetCodec().UnmarshalBinaryBare(iterator.Value(), &validator)
		validators = append(validators, validator)
	}
	iterator.Close()
	return
}
