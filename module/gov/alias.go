package gov

import (
	"github.com/QOSGroup/qos/module/gov/mapper"
	"github.com/QOSGroup/qos/module/gov/txs"
	"github.com/QOSGroup/qos/module/gov/types"
)

var (
	ModuleName          = "gov"
	RegisterCodec       = txs.RegisterCodec
	DefaultGenesisState = types.DefaultGenesisState

	MapperName = mapper.MapperName
	NewMapper  = mapper.NewMapper
	GetMapper  = mapper.GetMapper
	Query      = mapper.Query

	ParamsSpace = types.ParamSpace

	StatusDepositPeriod = types.StatusDepositPeriod
	StatusVotingPeriod  = types.StatusVotingPeriod
	StatusPassed        = types.StatusPassed
	StatusRejected      = types.StatusRejected
)

type (
	GenesisState = types.GenesisState
	Params = types.Params
)
