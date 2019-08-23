package mint

import (
	"github.com/QOSGroup/qos/module/mint/mapper"
	"github.com/QOSGroup/qos/module/mint/types"
)

var (
	ModuleName      = "mint"
	Cdc             = types.Cdc
	RegisterCodec   = types.RegisterCodec
	NewGenesisState = types.NewGenesisState
	DefaultGenesis  = types.DefaultGenesisState
	ValidateGenesis = types.ValidateGenesis

	MapperName = types.MapperName
	NewMapper  = mapper.NewMapper
	GetMapper  = mapper.GetMapper
)

type (
	GenesisState = types.GenesisState

	InflationPhrase  = types.InflationPhrase
	InflationPhrases = types.InflationPhrases
)
