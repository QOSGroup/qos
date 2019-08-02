package mapper

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qos/module/stake/types"
	qtypes "github.com/QOSGroup/qos/types"
)

type Mapper struct {
	*mapper.BaseMapper
	hooks types.Hooks
}

var _ mapper.IMapper = (*Mapper)(nil)

func NewMapper() *Mapper {
	var validatorMapper = Mapper{}
	validatorMapper.BaseMapper = mapper.NewBaseMapper(nil, types.MapperName)
	return &validatorMapper
}

func GetMapper(ctx context.Context) *Mapper {
	return ctx.Mapper(types.MapperName).(*Mapper)
}

func (mapper *Mapper) Copy() mapper.IMapper {
	validatorMapper := &Mapper{}
	validatorMapper.BaseMapper = mapper.BaseMapper.Copy()
	validatorMapper.hooks = mapper.hooks
	return validatorMapper
}

func (mapper *Mapper) SetHooks(sh qtypes.Hooks) {
	if mapper.hooks != nil {
		panic("cannot set validator hooks twice")
	}
	mapper.hooks = sh.(types.Hooks)
}
