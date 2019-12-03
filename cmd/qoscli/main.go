package main

import (
	bcli "github.com/QOSGroup/qbase/client"
	"github.com/QOSGroup/qbase/client/block"
	"github.com/QOSGroup/qbase/client/config"
	"github.com/QOSGroup/qbase/client/rpc"
	"github.com/QOSGroup/qbase/client/sign"
	"github.com/QOSGroup/qbase/client/tx"
	bctypes "github.com/QOSGroup/qbase/client/types"
	"github.com/QOSGroup/qos/app"
	"github.com/QOSGroup/qos/types"
	"github.com/QOSGroup/qos/version"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
	"net/http"

	_ "github.com/QOSGroup/qos/swagger/statik"
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

	rootCmd.AddCommand(rpc.ServerCommand(cdc, func(server *rpc.RestServer) {
		app.ModuleBasics.RegisterRoutes(server.CliCtx, server.Mux)
		registerSwaggerUI(server)
	}))
	rootCmd.AddCommand(sign.SignCommand(cdc))
	rootCmd.AddCommand(tx.BroadcastCmd(cdc))

	executor := cli.PrepareMainCmd(rootCmd, "qos", types.DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func registerSwaggerUI(server *rpc.RestServer) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}
	staticServer := http.FileServer(statikFS)
	server.Mux.PathPrefix("/swagger-ui/").Handler(http.StripPrefix("/swagger-ui/", staticServer))
}
