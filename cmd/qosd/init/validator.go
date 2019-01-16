package init

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/app"
	staketypes "github.com/QOSGroup/qos/module/eco/types"
	"github.com/QOSGroup/qos/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	"path/filepath"
	"strings"
)

const (
	flagName        = "name"
	flagOwner       = "owner"
	flagPubKey      = "pubkey"
	flagBondTokens  = "tokens"
	flagDescription = "description"
)

func AddGenesisValidator(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-validator",
		Short: "Add genesis validator to genesis.json",
		Long: `
pubkey is a tendermint validator pubkey. the public key of the validator used in
Tendermint consensus.

home node's home directory.

owner is account address.

ex: pubkey: {"type":"tendermint/PubKeyEd25519","value":"VOn2rPx+t7Njdgi+eLb+jBuF175T1b7LAcHElsmIuXA="}

example:

	 qosd add-genesis-validator --home "/.qosd/" --name validatorName --owner address1vdp54s5za8tl4dmf9dcldfzn62y66m40ursfsa --pubkey "VOn2rPx+t7Njdgi+eLb+jBuF175T1b7LAcHElsmIuXA=" --tokens 100

		`,
		RunE: func(_ *cobra.Command, args []string) error {

			home := viper.GetString(cli.HomeFlag)
			genFile := filepath.Join(home, "config", "genesis.json")
			if !common.FileExists(genFile) {
				return fmt.Errorf("%s does not exist, run `qosd init` first", genFile)
			}

			name := viper.GetString(flagName)
			if len(name) == 0 {
				return errors.New("name is empty")
			}

			ownerStr := viper.GetString(flagOwner)
			if !strings.HasPrefix(ownerStr, "address") {
				return errors.New("owner is invalid")
			}

			owner, err := btypes.GetAddrFromBech32(ownerStr)
			if err != nil {
				return err
			}
			valPubkey := viper.GetString(flagPubKey)
			if len(valPubkey) == 0 {
				return errors.New("pubkey is empty")
			}
			tokens := uint64(viper.GetInt64(flagBondTokens))
			if tokens <= 0 {
				return errors.New("tokens lte zero")
			}
			desc := viper.GetString(flagDescription)

			bz, err := base64.StdEncoding.DecodeString(valPubkey)
			if err != nil {
				return err
			}
			var cKey ed25519.PubKeyEd25519
			copy(cKey[:], bz)

			val := staketypes.Validator{
				Name:            name,
				ValidatorPubKey: cKey,
				Owner:           owner,
				BondTokens:      uint64(tokens),
				Status:          staketypes.Active,
				BondHeight:      1,
				Description:     desc,
			}

			genDoc, err := loadGenesisDoc(cdc, genFile)
			if err != nil {
				return err
			}

			var appState app.GenesisState
			if err = cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {
				return err
			}

			for _, v := range appState.StakeData.Validators {
				if v.ValidatorPubKey.Equals(val.ValidatorPubKey) {
					return errors.New("validator already exists")
				}
				if bytes.Equal(v.Owner, val.Owner) {
					return fmt.Errorf("owner %s already bind a validator", val.Owner)
				}
			}

			appState.StakeData.Validators = append(appState.StakeData.Validators, val)
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

	cmd.Flags().String(flagName, "", "name for validator")
	cmd.Flags().String(flagOwner, "", "account address")
	cmd.Flags().String(flagPubKey, "", "tendermint consensus validator public key")
	cmd.Flags().Int64(flagBondTokens, 0, "bond tokens amount")
	cmd.Flags().String(flagDescription, "", "description")
	cmd.Flags().String(cli.HomeFlag, types.DefaultNodeHome, "node's home directory")

	cmd.MarkFlagRequired(flagName)
	cmd.MarkFlagRequired(flagOwner)
	cmd.MarkFlagRequired(flagPubKey)
	cmd.MarkFlagRequired(flagBondTokens)

	return cmd
}
