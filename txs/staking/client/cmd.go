package staking

import (
	bctypes "github.com/QOSGroup/qbase/client/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
)

func TxCommands(cdc *amino.Codec) []*cobra.Command {
	return bctypes.PostCommands(
		CreateValidatorCmd(cdc),
		RevokeValidatorCmd(cdc),
		ActiveValidatorCmd(cdc))
}

func QueryCommands(cdc *amino.Codec) []*cobra.Command {
	return bctypes.GetCommands(
		queryAllValidatorsCommand(cdc),
		queryValidatorInfoCommand(cdc),
	)
}
