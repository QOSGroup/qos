package mint

import (
	"github.com/QOSGroup/qos/module/mint/mapper"
	"github.com/QOSGroup/qos/module/mint/types"
)

var (
	ModuleName          = "mint"
	RegisterCodec       = types.RegisterCodec
	DefaultGenesisState = types.DefaultGenesisState

	MapperName = types.MapperName
	NewMapper  = mapper.NewMapper
	GetMapper  = mapper.GetMapper
)

type (
	GenesisState = types.GenesisState
)
