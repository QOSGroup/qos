package main

import (
	bcli "github.com/QOSGroup/qbase/client"
	bctypes "github.com/QOSGroup/qbase/client/types"
	"github.com/QOSGroup/qos/app"
	"github.com/QOSGroup/qos/txs/approve/client"
	"github.com/QOSGroup/qos/txs/qcp/client"
	"github.com/QOSGroup/qos/txs/qsc/client"
	"github.com/QOSGroup/qos/txs/staking/client"
	"github.com/QOSGroup/qos/txs/transfer/client"
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
