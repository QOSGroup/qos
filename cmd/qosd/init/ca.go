package init

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto"
	"strings"

	"github.com/spf13/viper"

	"github.com/QOSGroup/qos/app"
	"github.com/QOSGroup/qos/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
)

const (
	flagQCP = "qcp"
	flagQSC = "qsc"
)

func ConfigRootCA(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config-root-ca",
		Short: "Config pubKey of root CA for QCP and QSC",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {

			home := viper.GetString(cli.HomeFlag)
			genFile := strings.Join([]string{home, "config", "genesis.json"}, "/")

			if !common.FileExists(genFile) {
				return fmt.Errorf("%s does not exist, run `qosd init` first", genFile)
			}

			qcpFile := viper.GetString(flagQCP)
			qscFile := viper.GetString(flagQSC)
			if qcpFile == "" && qscFile == "" {
				return errors.New("empty input")
			}

			var qcpPubKey crypto.PubKey
			if qcpFile != "" {
				err := cdc.UnmarshalJSON(common.MustReadFile(qcpFile), &qcpPubKey)
				if err != nil {
					return err
				}
			}

			var qscPubKey crypto.PubKey
			if qscFile != "" {
				err := cdc.UnmarshalJSON(common.MustReadFile(qscFile), &qscPubKey)
				if err != nil {
					return err
				}
			}

			genDoc, err := loadGenesisDoc(cdc, genFile)
			if err != nil {
				return err
			}

			var appState app.GenesisState
			if err = cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {
				return err
			}

			appState.QSCData.RootPubKey = qcpPubKey
			appState.QCPData.RootPubKey = qscPubKey

			rawMessage, _ := cdc.MarshalJSON(appState)
			genDoc.AppState = rawMessage

			err = genDoc.ValidateAndComplete()
			if err != nil {
				return err
			}

			err = genDoc.SaveAs(genFile)
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().String(cli.HomeFlag, types.DefaultNodeHome, "node's home directory")
	cmd.Flags().String(flagQCP, "", "directory of QCP root.pub")
	cmd.Flags().String(flagQSC, "", "directory of QSC root.pub")

	return cmd
}
