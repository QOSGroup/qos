package main

import (
	bcli "github.com/QOSGroup/qbase/client"
	"github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/keys"
	"github.com/QOSGroup/qbase/client/tx"
	"github.com/QOSGroup/qos/app"
	"github.com/QOSGroup/qos/client/txs/approve"
	"github.com/QOSGroup/qos/client/txs/transfer"
	"github.com/QOSGroup/qos/version"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
)

var (
	rootCmd = &cobra.Command{
		Use:   "basecli",
		Short: "Basecoin light-client",
	}
)

func main() {
	cobra.EnableCommandSorting = false

	cdc := app.MakeCodec()

	rootCmd.AddCommand(bcli.LineBreak)

	// account
	rootCmd.AddCommand(
		bcli.GetCommands(account.QueryAccountCmd(cdc))...)
	rootCmd.AddCommand(bcli.LineBreak)

	// keys
	rootCmd.AddCommand(
		bcli.GetCommands(keys.Commands(cdc))...)
	rootCmd.AddCommand(bcli.LineBreak)

	// qos txs
	tx.AddCommands(rootCmd, cdc)
	transfer.AddCommands(rootCmd, cdc)
	approve.AddCommands(rootCmd, cdc)
	rootCmd.AddCommand(bcli.LineBreak)

	// version
	rootCmd.AddCommand(
		version.VersionCmd,
	)

	executor := cli.PrepareMainCmd(rootCmd, "qos", app.DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
