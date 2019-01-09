package main

import (
	"github.com/QOSGroup/qbase/server"
	"github.com/QOSGroup/qos/app"
	"github.com/QOSGroup/qos/cmd/qosd/export"
	qosdinit "github.com/QOSGroup/qos/cmd/qosd/init"
	"github.com/QOSGroup/qos/cmd/qosd/testnet"
	"github.com/QOSGroup/qos/types"
	"github.com/QOSGroup/qos/version"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"io"
)

func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "qosd",
		Short:             "qos Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	// testnet cmd
	rootCmd.AddCommand(testnet.TestnetFileCmd(cdc))

	// version cmd
	rootCmd.AddCommand(version.VersionCmd)

	rootCmd.AddCommand(server.InitCmd(ctx, cdc, qosdinit.GenQOSGenesisDoc, types.DefaultNodeHome))
	rootCmd.AddCommand(export.ExportCmd(ctx, cdc))
	rootCmd.AddCommand(qosdinit.ConfigRootCA(cdc))
	rootCmd.AddCommand(qosdinit.AddGenesisAccount(cdc))
	rootCmd.AddCommand(qosdinit.AddGenesisValidator(cdc))

	server.AddCommands(ctx, cdc, rootCmd, newApp)

	executor := cli.PrepareBaseCmd(rootCmd, "qos", types.DefaultNodeHome)

	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, storeTracer io.Writer) abci.Application {
	return app.NewApp(logger, db, storeTracer)
}
