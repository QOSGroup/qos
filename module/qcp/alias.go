package qcp

import (
	"github.com/QOSGroup/qbase/qcp"
	"github.com/QOSGroup/qos/module/qcp/client"
	"github.com/QOSGroup/qos/module/qcp/mapper"
	"github.com/QOSGroup/qos/module/qcp/txs"
	"github.com/QOSGroup/qos/module/qcp/types"
)

var (
	ModuleName      = "qcp"
	Cdc             = txs.Cdc
	RegisterCodec   = txs.RegisterCodec
	NewGenesis      = types.NewGenesisState
	DefaultGenesis  = types.DefaultGenesisState
	ValidateGenesis = types.ValidateGenesis

	MapperName = qcp.QcpMapperName
	GetMapper  = mapper.GetMapper
	NewMapper  = qcp.NewQcpMapper

	QueryCommands = client.QueryCommands
	TxCommands    = client.TxCommands
)

type (
	GenesisState = types.GenesisState
)
