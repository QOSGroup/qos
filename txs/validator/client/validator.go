package validator

import (
	"encoding/base64"
	"fmt"
	"github.com/QOSGroup/qbase/client/context"
	qclitx "github.com/QOSGroup/qbase/client/tx"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qos/txs/validator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

const (
	flagName       = "name"
	flagConsPubKey = "cons-pubkey"
)

func ValidatorCreateCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-validator",
		Short: "Create validator.",
		Long: `
cons-pubkey is a tendermint validator pubkey. the public key of the validator used in
Tendermint consensus.

name is a Operator keybase name.

ex: pubkey: {"type":"tendermint/PubKeyEd25519","value":"VOn2rPx+t7Njdgi+eLb+jBuF175T1b7LAcHElsmIuXA="}

example:

	 ./bin create-validator --name main --cons-pubkey

		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return qclitx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (txs.ITx, error) {
				name := viper.GetString(flagName)
				consPubkey := viper.GetString(flagConsPubKey)
				creatorAddr, err := qclitx.GetAddrFromFlag(ctx, flagName)
				if err != nil {
					return nil, err
				}

				bz, err := base64.StdEncoding.DecodeString(consPubkey)
				if err != nil {
					return nil, fmt.Errorf("cons-pubkey parse error: %s", err.Error())
				}
				var cKey ed25519.PubKeyEd25519
				copy(cKey[:], bz)

				return validator.NewCreateValidatorTx(name, cKey, creatorAddr), nil
			})

		},
	}

	cmd.Flags().String(flagName, "", "operator keys name")
	cmd.Flags().String(flagConsPubKey, "", "tendermint consensus validator public key")

	cmd.MarkFlagRequired(flagName)
	cmd.MarkFlagRequired(flagConsPubKey)

	return cmd
}
