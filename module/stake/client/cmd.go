package client

import (
	bctypes "github.com/QOSGroup/qbase/client/types"
	"github.com/QOSGroup/qos/module/stake/txs"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
)

func TxCommands(cdc *amino.Codec) []*cobra.Command {
	return bctypes.PostCustomMaxGasCommands([]*cobra.Command{
		CreateValidatorCmd(cdc),
		ModifyValidatorCmd(cdc),
		RevokeValidatorCmd(cdc),
		ActiveValidatorCmd(cdc),
		bctypes.LineBreak,
		CreateDelegationCommand(cdc),
		CreateModifyCompoundCommand(cdc),
		CreateUnbondDelegationCommand(cdc),
		CreateReDelegationCommand(cdc),
	}, []int64{
		txs.GasForCreateValidator + bctypes.DefaultMaxGas,
		txs.GasForModifyValidator + bctypes.DefaultMaxGas,
		txs.GasForRevokeValidator + bctypes.DefaultMaxGas,
		bctypes.DefaultMaxGas,
		bctypes.DefaultMaxGas,
		bctypes.DefaultMaxGas,
		bctypes.DefaultMaxGas,
		txs.GasForUnbond + bctypes.DefaultMaxGas,
		bctypes.DefaultMaxGas,
	})

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
