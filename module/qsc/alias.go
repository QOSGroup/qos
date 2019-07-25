package qsc

import (
	"github.com/QOSGroup/qos/module/qsc/mapper"
	"github.com/QOSGroup/qos/module/qsc/txs"
	"github.com/QOSGroup/qos/module/qsc/types"
)

var (
	ModuleName      = "qsc"
	RegisterCodec   = txs.RegisterCodec
	NewGenesisState = types.NewGenesisState

	MapperName = mapper.MapperName
	GetMapper  = mapper.GetMapper
	NewMapper  = mapper.NewMapper
)

type (
	GenesisState = types.GenesisState
	Info = types.Info
)
