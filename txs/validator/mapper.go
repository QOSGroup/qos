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
	validatorUpdatedKey = "updated"
)

type ValidatorMapper struct {
	*mapper.BaseMapper
}

func NewValidatorMapper() *ValidatorMapper {
	var validatorMapper = ValidatorMapper{}
	validatorMapper.BaseMapper = mapper.NewBaseMapper(nil, ValidatorMapperName)
	return &validatorMapper
}

func BuildUpdatedKey() []byte {
	return []byte(validatorUpdatedKey)
}

func BuildValidatorKey(address btypes.Address) []byte {
	return append([]byte(validatorKey), address...)
}

func (mapper *ValidatorMapper) MapperName() string {
	return ValidatorMapperName
}

func (mapper *ValidatorMapper) Copy() mapper.IMapper {
	validatorMapper := &ValidatorMapper{}
	validatorMapper.BaseMapper = mapper.BaseMapper.Copy()
	return validatorMapper
}

// 是否有新增validator
func (mapper *ValidatorMapper) HasNewValidator() bool {
	v, exists := mapper.GetBool(BuildUpdatedKey())
	return exists && !v
}

// 设置 新增validator 标识
func (mapper *ValidatorMapper) SetUpdated(b bool) {
	mapper.Set(BuildUpdatedKey(), b)
}

// 保存validator
func (mapper *ValidatorMapper) SaveValidator(validator types.Validator) {
	mapper.Set(BuildValidatorKey(validator.PubKey.Address().Bytes()), validator)
}

// 是否已经存在
func (mapper *ValidatorMapper) Exists(address btypes.Address) bool {
	return mapper.Get(BuildValidatorKey(address), &(types.Validator{}))
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
