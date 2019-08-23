package gov

import (
	"github.com/QOSGroup/qos/module/gov/client"
	"github.com/QOSGroup/qos/module/gov/mapper"
	"github.com/QOSGroup/qos/module/gov/txs"
	"github.com/QOSGroup/qos/module/gov/types"
)

var (
	ModuleName      = "gov"
	Cdc             = txs.Cdc
	RegisterCodec   = txs.RegisterCodec
	DefaultGenesis  = types.DefaultGenesisState
	ValidateGenesis = types.ValidateGenesis

	MapperName = mapper.MapperName
	NewMapper  = mapper.NewMapper
	GetMapper  = mapper.GetMapper
	Query      = mapper.Query

	ParamsSpace = types.ParamSpace

	StatusDepositPeriod = types.StatusDepositPeriod
	StatusVotingPeriod  = types.StatusVotingPeriod
	StatusPassed        = types.StatusPassed
	StatusRejected      = types.StatusRejected

	QueryCommands = client.QueryCommands
	TxCommands    = client.TxCommands
)

type (
	GenesisState = types.GenesisState
	Params       = types.Params
)
