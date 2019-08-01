package guardian

import (
	"github.com/QOSGroup/qos/module/guardian/client"
	"github.com/QOSGroup/qos/module/guardian/mapper"
	"github.com/QOSGroup/qos/module/guardian/txs"
	"github.com/QOSGroup/qos/module/guardian/types"
)

var (
	ModuleName      = "guardian"
	Cdc             = txs.Cdc
	RegisterCodec   = txs.RegisterCodec
	NewGenesisState = types.NewGenesisState
	DefaultGenesis  = types.DefaultGenesisState
	ValidateGenesis = types.ValidateGenesis

	MapperName = mapper.MapperName
	GetMapper  = mapper.GetMapper
	NewMapper  = mapper.NewMapper

	NewGuardian = types.NewGuardian
	Genesis     = types.Genesis

	QueryCommands = client.QueryCommands
	TxCommands    = client.TxCommands
)

type (
	GenesisState = types.GenesisState
	Guardian     = types.Guardian
)
