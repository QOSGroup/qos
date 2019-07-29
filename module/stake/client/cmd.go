package staking

import (
	bctypes "github.com/QOSGroup/qbase/client/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
)

func TxValidatorCommands(cdc *amino.Codec) []*cobra.Command {
	return bctypes.PostCommands(
		CreateValidatorCmd(cdc),
		ModifyValidatorCmd(cdc),
		RevokeValidatorCmd(cdc),
		ActiveValidatorCmd(cdc),
	)
}

func TxDelegationCommands(cdc *amino.Codec) []*cobra.Command {
	return bctypes.PostCommands(
		CreateDelegationCommand(cdc),
		CreateModifyCompoundCommand(cdc),
		CreateUnbondDelegationCommand(cdc),
		CreateReDelegationCommand(cdc),
	)
}

func QueryCommands(cdc *amino.Codec) []*cobra.Command {
	return bctypes.GetCommands(
		queryAllValidatorsCommand(cdc),
		queryValidatorInfoCommand(cdc),
		queryValidatorMissedVoteInfoCommand(cdc),
		queryDelegationInfoCommand(cdc),
		queryDelegationsCommand(cdc),
		queryDelegationsToCommand(cdc),
		queryUnbondingsCommand(cdc),
		queryRedelegationsCommand(cdc),
	)
}
