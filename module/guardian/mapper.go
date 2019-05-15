package guardian

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/guardian/types"
)

const (
	MapperName = "guardian"
)

var (
	guardianKey = []byte{0x00}
)

func KeyGuardian(address btypes.Address) []byte {
	return append(guardianKey, address...)
}

func KeyGuardiansSubspace() []byte {
	return guardianKey
}

type GuardianMapper struct {
	*mapper.BaseMapper
}

func (mapper *GuardianMapper) Copy() mapper.IMapper {
	govMapper := &GuardianMapper{}
	govMapper.BaseMapper = mapper.BaseMapper.Copy()
	return govMapper
}

var _ mapper.IMapper = (*GuardianMapper)(nil)

func GetGuardianMapper(ctx context.Context) *GuardianMapper {
	return ctx.Mapper(MapperName).(*GuardianMapper)
}

func NewGuardianMapper() *GuardianMapper {
	var govMapper = GuardianMapper{}
	govMapper.BaseMapper = mapper.NewBaseMapper(nil, MapperName)
	return &govMapper
}

func (mapper GuardianMapper) AddGuardian(guardian types.Guardian) {
	mapper.Set(KeyGuardian(guardian.Address), guardian)
}

func (mapper GuardianMapper) DeleteGuardian(address btypes.Address) {
	mapper.Del(KeyGuardian(address))
}

func (mapper GuardianMapper) GetGuardian(address btypes.Address) (guardian types.Guardian, exists bool) {
	exists = mapper.Get(KeyGuardian(address), &guardian)
	return guardian, exists
}

func (mapper GuardianMapper) GuardiansIterator() store.Iterator {
	return store.KVStorePrefixIterator(mapper.GetStore(), KeyGuardiansSubspace())
}
