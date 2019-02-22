package export

import (
	"encoding/json"
	"fmt"
	"github.com/QOSGroup/qbase/server"
	"github.com/QOSGroup/qos/app"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/log"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"io/ioutil"
	"path"

	tmtypes "github.com/tendermint/tendermint/types"

	dbm "github.com/tendermint/tendermint/libs/db"
)

const (
	flagHeight        = "height"
	flagForZeroHeight = "for-zero-height"
	flagTraceStore    = "trace-store"
)

// ExportCmd dumps app state to JSON.
func ExportCmd(ctx *server.Context, cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export state to JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			home := viper.GetString("home")
			traceWriterFile := viper.GetString(flagTraceStore)
			emptyState, err := isEmptyState(home)
			if err != nil {
				return err
			}

			if emptyState {
				fmt.Println("WARNING: State is not initialized. Returning genesis file.")
				genesisFile := path.Join(home, "config", "genesis.json")
				genesis, err := ioutil.ReadFile(genesisFile)
				if err != nil {
					return err
				}
				fmt.Println(string(genesis))
				return nil
			}

			db, err := openDB(home)
			if err != nil {
				return err
			}
			traceWriter, err := openTraceWriter(traceWriterFile)
			if err != nil {
				return err
			}
			height := viper.GetInt64(flagHeight)
			forZeroHeight := viper.GetBool(flagForZeroHeight)
			appState, err := exportAppState(ctx.Logger, db, traceWriter, height, forZeroHeight)
			if err != nil {
				return errors.Errorf("error exporting state: %v\n", err)
			}

			doc, err := tmtypes.GenesisDocFromFile(ctx.Config.GenesisFile())
			if err != nil {
				return err
			}

			doc.AppState = appState

			encoded, err := cdc.MarshalJSONIndent(doc, "", " ")
			if err != nil {
				return err
			}

			fmt.Println(string(encoded))
			return nil
		},
	}
	cmd.Flags().Int64(flagHeight, -1, "Export state from a particular height (-1 means latest height)")
	cmd.Flags().Bool(flagForZeroHeight, false, "Export state to start at height zero (perform preproccessing)")
	return cmd
}

func isEmptyState(home string) (bool, error) {
	files, err := ioutil.ReadDir(path.Join(home, "data"))
	if err != nil {
		return false, err
	}

	return len(files) == 0, nil
}

func openDB(rootDir string) (dbm.DB, error) {
	dataDir := filepath.Join(rootDir, "data")
	db, err := dbm.NewGoLevelDB("application", dataDir)
	return db, err
}

func openTraceWriter(traceWriterFile string) (w io.Writer, err error) {
	if traceWriterFile != "" {
		w, err = os.OpenFile(
			traceWriterFile,
			os.O_WRONLY|os.O_APPEND|os.O_CREATE,
			0666,
		)
		return
	}
	return
}

func exportAppState(logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool) (json.RawMessage, error) {
	qApp := app.NewApp(logger, db, traceStore)
	qApp.LoadVersion(height)
	return qApp.ExportAppStates(forZeroHeight)
}
