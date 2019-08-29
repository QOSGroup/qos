package types

import (
	"github.com/QOSGroup/qbase/mapper"
)

// mapper for db
type MapperWithHooks struct {
	Mapper mapper.IMapper
	Hooks  Hooks
}

func NewMapperWithHooks(mapper mapper.IMapper, hooks Hooks) MapperWithHooks {
	return MapperWithHooks{mapper, hooks}
}

func (mh *MapperWithHooks) IsNil() bool {
	return mh.Mapper == nil
}

type HooksMapperRegistry interface {
	RegisterHooksMapper(map[string]MapperWithHooks)
}

type HooksMapper interface {
	SetHooks(hooks Hooks)
}

type ParamsInitializer interface {
	RegisterParamSet(ps ...ParamSet)
}