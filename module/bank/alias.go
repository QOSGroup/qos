package bank

import (
	"github.com/QOSGroup/qos/module/bank/client"
	"github.com/QOSGroup/qos/module/bank/mapper"
	"github.com/QOSGroup/qos/module/bank/txs"
	"github.com/QOSGroup/qos/module/bank/types"
)

var (
	ModuleName      = "bank"
	Cdc             = txs.Cdc
	GetMapper       = mapper.GetMapper
	RegisterCodec   = txs.RegisterCodec
	NewGenesisState = types.NewGenesisState
	DefaultGenesis  = types.DefaultGenesisState
	ValidateGenesis = types.ValidateGenesis

	TxCommands = client.TxCommands

	NeedInvariantCheck = mapper.NeedInvariantCheck
)

type (
	GenesisState = types.GenesisState
)
