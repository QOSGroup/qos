package main

import (
	bcli "github.com/QOSGroup/qbase/client"
	"github.com/QOSGroup/qbase/client/config"
	bctypes "github.com/QOSGroup/qbase/client/types"
	"github.com/QOSGroup/qos/app"
	"github.com/QOSGroup/qos/modules/approve/client"
	"github.com/QOSGroup/qos/modules/qcp/client"
	"github.com/QOSGroup/qos/modules/qsc/client"
	"github.com/QOSGroup/qos/modules/stake/client"
	"github.com/QOSGroup/qos/modules/transfer/client"
	"github.com/QOSGroup/qos/types"
	"github.com/QOSGroup/qos/version"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
)

var (
	rootCmd = &cobra.Command{
		Use:   "qoscli",
		Short: "QOS light-client",
	}
)

func main() {
	cobra.EnableCommandSorting = false

	cdc := app.MakeCodec()

	rootCmd.AddCommand(bctypes.LineBreak)

	// config command
	rootCmd.AddCommand(config.Cmd(types.DefaultCLIHome))

	// query commands
	queryCommands := bcli.QueryCommand(cdc)
	queryCommands.AddCommand(approve.QueryCommands(cdc)...)
	queryCommands.AddCommand(qsc.QueryCommands(cdc)...)
	queryCommands.AddCommand(staking.QueryCommands(cdc)...)

	// txs commands
	txsCommands := bcli.TxCommand()
	txsCommands.AddCommand(qsc.TxCommands(cdc)...)
	txsCommands.AddCommand(bctypes.LineBreak)
	txsCommands.AddCommand(qcp.TxCommands(cdc)...)
	txsCommands.AddCommand(bctypes.LineBreak)
	txsCommands.AddCommand(transfer.TxCommands(cdc)...)
	txsCommands.AddCommand(bctypes.LineBreak)
	txsCommands.AddCommand(approve.TxCommands(cdc)...)
	txsCommands.AddCommand(bctypes.LineBreak)
	txsCommands.AddCommand(staking.TxCommands(cdc)...)

	rootCmd.AddCommand(
		bcli.KeysCommand(cdc),
		queryCommands,
		txsCommands,
		bcli.TendermintCommand(cdc),
		version.VersionCmd,
	)

	executor := cli.PrepareMainCmd(rootCmd, "qos", types.DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
