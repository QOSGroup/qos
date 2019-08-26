package params

import (
	"github.com/QOSGroup/qos/module/params/mapper"
	"github.com/QOSGroup/qos/module/params/types"
)

var (
	ModuleName    = "params"
	RegisterCodec = types.RegisterCodec

	MapperName = mapper.MapperName
	NewMapper  = mapper.NewMapper
	GetMapper  = mapper.GetMapper

	BuildParamKey = mapper.BuildParamKey

	ErrInvalidParam = types.ErrInvalidParam
)

type (
	Mapper = mapper.Mapper
)
