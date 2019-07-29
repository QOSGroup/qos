package stake

import (
	"github.com/QOSGroup/qos/module/stake/mapper"
	"github.com/QOSGroup/qos/module/stake/txs"
	"github.com/QOSGroup/qos/module/stake/types"
)

var (
	ModuleName          = "stake"
	RegisterCodec       = txs.RegisterCodec
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis

	MapperName                             = types.MapperName
	NewMapper                              = mapper.NewMapper
	GetMapper                              = mapper.GetMapper
	Query                                  = mapper.Query
	BuildUnbondingDelegationByHeightPrefix = types.BuildUnbondingDelegationByHeightPrefix
	GetUnbondingDelegationHeightAddress    = types.GetUnbondingDelegationHeightAddress
	BuildRedelegationByHeightPrefix        = types.BuildRedelegationByHeightPrefix
	GetRedelegationHeightAddress           = types.GetRedelegationHeightAddress

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
)

type (
	GenesisState = types.GenesisState
	ValidatorVoteInfoState = types.ValidatorVoteInfoState
	ValidatorVoteInWindowInfoState = types.ValidatorVoteInWindowInfoState
	DelegationInfoState = types.DelegationInfoState
	Params = types.Params
	Validator = types.Validator
	Description = types.Description
	Delegation = types.DelegationInfo
	Hooks = types.Hooks

	TxCreateValidator = txs.TxCreateValidator

	UnbondingDelegationInfo = types.UnbondingDelegationInfo
	ReDelegationInfo = types.RedelegationInfo
)
