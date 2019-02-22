package main

import (
	bcli "github.com/QOSGroup/qbase/client"
	"github.com/QOSGroup/qbase/client/config"
	bctypes "github.com/QOSGroup/qbase/client/types"
	"github.com/QOSGroup/qos/app"
	"github.com/QOSGroup/qos/module/approve/client"
	"github.com/QOSGroup/qos/module/distribution/client"
	"github.com/QOSGroup/qos/module/qcp/client"
	"github.com/QOSGroup/qos/module/qsc/client"
	"github.com/QOSGroup/qos/module/stake/client"
	"github.com/QOSGroup/qos/module/transfer/client"
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
	queryCommands.AddCommand(distribution.QueryCommands(cdc)...)

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
	txsCommands.AddCommand(staking.TxValidatorCommands(cdc)...)
	txsCommands.AddCommand(bctypes.LineBreak)
	txsCommands.AddCommand(staking.TxDelegationCommands(cdc)...)

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
