package client

import (
	bctypes "github.com/QOSGroup/qbase/client/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
)

func QueryCommands(cdc *amino.Codec) []*cobra.Command {
	return bctypes.GetCommands(QueryApproveCmd(cdc))
}

func TxCommands(cdc *amino.Codec) []*cobra.Command {
	return bctypes.PostCommands(
		CreateApproveCmd(cdc),
		IncreaseApproveCmd(cdc),
		DecreaseApproveCmd(cdc),
		UseApproveCmd(cdc),
		CancelApproveCmd(cdc),
	)
}
