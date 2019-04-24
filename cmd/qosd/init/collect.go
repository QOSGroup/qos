package init

import (
	"errors"
	"github.com/QOSGroup/qbase/server"
	"github.com/QOSGroup/qos/app"
	"github.com/QOSGroup/qos/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	tmtypes "github.com/tendermint/tendermint/types"
	"path/filepath"
)

const flagGenTxDir string = "gentx-dir"

func CollectGenTxsCmd(ctx *server.Context, cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collect-gentxs",
		Short: "Collect genesis txs and output a genesis.json file",
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))
			nodeID, _, err := server.InitializeNodeValidatorFiles(config)
			if err != nil {
				return err
			}

			return UpdateGenesisStateFromGenTxs(config, cdc, nodeID)
		},
	}

	cmd.Flags().String(cli.HomeFlag, types.DefaultNodeHome, "node's home directory, default: $HOME/.qosd")
	cmd.Flags().String(flagGenTxDir, "", "directory of gentx files, default: $HOME/.qosd/config/gentx/")
	return cmd
}

func UpdateGenesisStateFromGenTxs(config *cfg.Config, cdc *amino.Codec, nodeID string) (err error) {
	genDoc, err := tmtypes.GenesisDocFromFile(config.GenesisFile())
	if err != nil {
		return err
	}

	genTxsDir := viper.GetString(flagGenTxDir)
	if genTxsDir == "" {
		genTxsDir = filepath.Join(config.RootDir, "config", "gentx")
	}

	genTxs, persistentPeers, err := app.CollectStdTxs(cdc, nodeID, genTxsDir, genDoc)
	if err != nil {
		return err
	}
	if len(genTxs) == 0 {
		return errors.New("there must be at least one genesis tx")
	}

	// update config.toml
	config.P2P.PersistentPeers = persistentPeers
	cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)

	// update genesis.json
	var genesisState app.GenesisState
	if err = cdc.UnmarshalJSON(genDoc.AppState, &genesisState); err != nil {
		return err
	}
	genesisState.GenTxs = genTxs
	genDoc.AppState, err = server.MarshalJSONIndent(cdc, genesisState)
	if err != nil {
		return
	}
	if err := genDoc.ValidateAndComplete(); err != nil {
		return err
	}
	err = genDoc.SaveAs(config.GenesisFile())

	return
}
