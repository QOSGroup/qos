package stake

import (
	"github.com/QOSGroup/qos/module/stake/client"
	"github.com/QOSGroup/qos/module/stake/mapper"
	"github.com/QOSGroup/qos/module/stake/txs"
	"github.com/QOSGroup/qos/module/stake/types"
)

var (
	ModuleName      = "stake"
	Cdc             = txs.Cdc
	RegisterCodec   = txs.RegisterCodec
	NewGenesis      = types.NewGenesisState
	DefaultGenesis  = types.DefaultGenesisState
	ValidateGenesis = types.ValidateGenesis

	MapperName                             = types.MapperName
	NewMapper                              = mapper.NewMapper
	GetMapper                              = mapper.GetMapper
	Query                                  = mapper.Query
	BuildUnbondingDelegationByHeightPrefix = types.BuildUnbondingDelegationByHeightPrefix
	GetUnbondingDelegationHeightAddress    = types.GetUnbondingDelegationHeightDelegatorValidator
	BuildRedelegationByHeightPrefix        = types.BuildRedelegationByHeightPrefix
	GetRedelegationHeightAddress           = types.GetRedelegationHeightDelegatorFromValidator

	NewCreateValidatorTx = txs.NewCreateValidatorTx

	NewDelegationInfo = types.NewDelegationInfo

	NewUnbondingInfo = types.NewUnbondingDelegationInfo

	ParamSpace    = types.ParamSpace
	DefaultParams = types.DefaultParams

	Active        = types.Active
	Inactive      = types.Inactive
	Revoke        = types.Revoke
	MissVoteBlock = types.MissVoteBlock
	MaxValidator  = types.MaxValidator

	QueryCommands            = client.QueryCommands
	TxCommands               = client.TxCommands
	TxCreateValidatorBuilder = client.TxCreateValidatorBuilder
	BuildCommissionRates     = client.BuildCommissionRates
	NewCommissionRates       = types.NewCommissionRates

	DefaultCommissionRate          = client.DefaultCommissionRate
	DefaultCommissionMaxRate       = client.DefaultCommissionMaxRate
	DefaultCommissionMaxChangeRate = client.DefaultCommissionMaxChangeRate

	BuildCurrentValidatorsAddressKey = types.BuildCurrentValidatorsAddressKey
)

type (
	GenesisState                   = types.GenesisState
	ValidatorVoteInfoState         = types.ValidatorVoteInfoState
	ValidatorVoteInWindowInfoState = types.ValidatorVoteInWindowInfoState
	DelegationInfoState            = types.DelegationInfoState
	Params                         = types.Params
	Validator                      = types.Validator
	Description                    = types.Description
	Delegation                     = types.DelegationInfo
	Hooks                          = types.Hooks

	TxCreateValidator = txs.TxCreateValidator

	UnbondingDelegationInfo = types.UnbondingDelegationInfo
	ReDelegationInfo        = types.RedelegationInfo
)
