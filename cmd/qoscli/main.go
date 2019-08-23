package main

import (
	bcli "github.com/QOSGroup/qbase/client"
	"github.com/QOSGroup/qbase/client/block"
	"github.com/QOSGroup/qbase/client/config"
	bctypes "github.com/QOSGroup/qbase/client/types"
	"github.com/QOSGroup/qos/app"
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
	app.ModuleBasics.AddQueryCommands(queryCommands, cdc)
	queryCommands.AddCommand(block.BlockCommand(cdc)...)

	// txs commands
	txsCommands := bcli.TxCommand()
	app.ModuleBasics.AddTxCommands(txsCommands, cdc)

	rootCmd.AddCommand(
		bcli.KeysCommand(cdc),
		queryCommands,
		txsCommands,
		version.VersionCmd(),
	)

	executor := cli.PrepareMainCmd(rootCmd, "qos", types.DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
