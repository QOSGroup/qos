package qsc

import (
	"github.com/QOSGroup/qos/module/qsc/client"
	"github.com/QOSGroup/qos/module/qsc/mapper"
	"github.com/QOSGroup/qos/module/qsc/txs"
	"github.com/QOSGroup/qos/module/qsc/types"
)

var (
	ModuleName      = "qsc"
	Cdc             = txs.Cdc
	RegisterCodec   = txs.RegisterCodec
	NewGenesisState = types.NewGenesisState
	DefaultGenesis  = types.DefaultGenesisState
	ValidateGenesis = types.ValidateGenesis

	MapperName = mapper.MapperName
	GetMapper  = mapper.GetMapper
	NewMapper  = mapper.NewMapper

	QueryCommands = client.QueryCommands
	TxCommands    = client.TxCommands
)

type (
	GenesisState = types.GenesisState
	Info         = types.Info
)
