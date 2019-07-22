package qcp

import (
	"github.com/QOSGroup/qbase/qcp"
	"github.com/QOSGroup/qos/module/qcp/mapper"
	"github.com/QOSGroup/qos/module/qcp/txs"
	"github.com/QOSGroup/qos/module/qcp/types"
)

var (
	ModuleName      = "qcp"
	RegisterCodec   = txs.RegisterCodec
	NewGenesisState = types.NewGenesisState

	MapperName = qcp.QcpMapperName
	GetMapper  = mapper.GetMapper
	NewMapper  = qcp.NewQcpMapper
)

type (
	GenesisState = types.GenesisState
)
