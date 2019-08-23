package mapper

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
	haltKey     = []byte{0x01}
)

func KeyGuardian(address btypes.Address) []byte {
	return append(guardianKey, address...)
}

func KeyGuardiansSubspace() []byte {
	return guardianKey
}

type Mapper struct {
	*mapper.BaseMapper
}

func (mapper *Mapper) Copy() mapper.IMapper {
	govMapper := &Mapper{}
	govMapper.BaseMapper = mapper.BaseMapper.Copy()
	return govMapper
}

var _ mapper.IMapper = (*Mapper)(nil)

func GetMapper(ctx context.Context) *Mapper {
	return ctx.Mapper(MapperName).(*Mapper)
}

func NewMapper() *Mapper {
	var guardianMapper = Mapper{}
	guardianMapper.BaseMapper = mapper.NewBaseMapper(nil, MapperName)
	return &guardianMapper
}

func (mapper Mapper) AddGuardian(guardian types.Guardian) {
	mapper.Set(KeyGuardian(guardian.Address), guardian)
}

func (mapper Mapper) DeleteGuardian(address btypes.Address) {
	mapper.Del(KeyGuardian(address))
}

func (mapper Mapper) GetGuardian(address btypes.Address) (guardian types.Guardian, exists bool) {
	exists = mapper.Get(KeyGuardian(address), &guardian)
	return guardian, exists
}

func (mapper Mapper) GuardiansIterator() store.Iterator {
	return btypes.KVStorePrefixIterator(mapper.GetStore(), KeyGuardiansSubspace())
}

func (mapper Mapper) SetHalt(halt string) {
	mapper.Set(haltKey, halt)
}

func (mapper Mapper) NeedHalt(height uint64) bool {
	var halt string
	exists := mapper.Get(haltKey, &halt)
	return exists
}
