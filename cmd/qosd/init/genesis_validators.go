package init

import (
	"encoding/base64"
	"fmt"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"strings"

	"github.com/spf13/viper"

	qbtypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/app"
	"github.com/QOSGroup/qos/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/privval"
)

const (
	flagConsPubKey = "consPubkey"
	flagOperator   = "operator"
	flagPower      = "power"
)

func AddGenesisValidator(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-validator",
		Short: "Add genesis validator to genesis.json",
		RunE: func(_ *cobra.Command, args []string) error {

			home := viper.GetString(cli.HomeFlag)
			genFile := strings.Join([]string{home, "config", "genesis.json"}, "/")

			if !common.FileExists(genFile) {
				return fmt.Errorf("%s does not exist, run `qosd init` first", genFile)
			}

			//TODO
			_, err := qbtypes.GetAddrFromBech32(viper.GetString(flagOperator))
			if err != nil {
				return err
			}

			var consPubkey ed25519.PubKeyEd25519
			consPubKeyStr := viper.GetString(flagConsPubKey)

			if consPubKeyStr == "" {
				//load from priv_validator.json
				priFile := strings.Join([]string{home, "config", "priv_validator.json"}, "/")
				if !common.FileExists(priFile) {
					return fmt.Errorf("%s does not exist, run `qosd init` first", priFile)
				}

				privValidator := privval.LoadFilePV(priFile)
				consPubkey = (privValidator.GetPubKey()).(ed25519.PubKeyEd25519)
			} else {
				bz, err := base64.StdEncoding.DecodeString(consPubKeyStr)
				if err != nil {
					return err
				}
				copy(consPubkey[:], bz)
			}

			name := consPubkey.Address().String()

			genDoc, err := loadGenesisDoc(cdc, genFile)
			if err != nil {
				return err
			}

			var appState app.GenesisState
			if err = cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {
				return err
			}

			for _, v := range appState.Validators {
				if v.Name == name {
					return fmt.Errorf("validator name: %s has already exsits", name)
				}
			}
			//TODO

			// val := types.Validator{
			// 	Name:        name,
			// 	ConsPubKey:  consPubkey,
			// 	Operator:    operatorAddr,
			// 	VotingPower: viper.GetInt64(flagPower),
			// 	Height:      1,
			// }

			// appState.Validators = append(appState.Validators, val)
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

	cmd.Flags().String(flagOperator, "", "operator address")
	cmd.Flags().String(flagConsPubKey, "", "validator's ed25519 consPubkey ")
	cmd.Flags().Int64(flagPower, 10, "validator's voting power. default is 10")
	cmd.Flags().String(cli.HomeFlag, types.DefaultNodeHome, "node's home directory")

	cmd.MarkFlagRequired(flagOperator)
	return cmd
}
