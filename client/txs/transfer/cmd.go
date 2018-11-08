package transfer

import (
	bcli "github.com/QOSGroup/qbase/client"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
)

func AddCommands(cmd *cobra.Command, cdc *amino.Codec) {
	cmd.AddCommand(bcli.GetCommands(
		TransferCmd(cdc),
	)...)
}
