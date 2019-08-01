package bank

import (
	"github.com/QOSGroup/qos/module/bank/client"
	"github.com/QOSGroup/qos/module/bank/txs"
	"github.com/QOSGroup/qos/module/bank/types"
)

var (
	ModuleName      = "bank"
	Cdc             = txs.Cdc
	RegisterCodec   = txs.RegisterCodec
	NewGenesisState = types.NewGenesisState
	DefaultGenesis  = types.DefaultGenesisState
	ValidateGenesis = types.ValidateGenesis

	TxCommands = client.TxCommands
)

type (
	GenesisState = types.GenesisState
)
