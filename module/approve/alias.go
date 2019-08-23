package approve

import (
	"github.com/QOSGroup/qos/module/approve/client"
	"github.com/QOSGroup/qos/module/approve/mapper"
	"github.com/QOSGroup/qos/module/approve/txs"
	"github.com/QOSGroup/qos/module/approve/types"
)

var (
	ModuleName    = "approve"
	Cdc           = txs.Cdc
	RegisterCodec = txs.RegisterCodec
	DefaultGenesis  = types.DefaultGenesisState
	ValidateGenesis = types.ValidateGenesis

	NewMapper  = mapper.NewApproveMapper
	MapperName = types.MapperName

	QueryCommands = client.QueryCommands
	TxCommands    = client.TxCommands
)

type (
	Mapper = mapper.Mapper

	GenesisState = types.GenesisState

	TxCreateApprove   = txs.TxCreateApprove
	TxIncreaseApprove = txs.TxIncreaseApprove
	TxDecreaseApprove = txs.TxDecreaseApprove
	TxUseApprove      = txs.TxUseApprove
	TxCancelApprove   = txs.TxCancelApprove
)
