package validator

import (
	"errors"
	"fmt"
	"github.com/QOSGroup/qbase/client/context"
	btxs "github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/client"
	"github.com/QOSGroup/qos/txs/validator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

const (
	flagName   = "name"
	flagPubKey = "pub-key"
)

func ValidatorCreateCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-create",
		Short: "Create validator.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			name := viper.GetString(flagName)
			pubkeyStr := viper.GetString(flagPubKey)
			if len(name) == 0 || len(pubkeyStr) == 0 {
				return errors.New("name or pub-key is empty")
			}

			var pubKey ed25519.PubKeyEd25519
			copy(pubKey[:], pubkeyStr)

			validatorTx := validator.NewCreateValidatorTx(name, pubKey)

			chainId, err := client.GetDefaultChainId()
			if err != nil {
				return nil
			}

			stdTx := btxs.NewTxStd(&validatorTx, chainId, btypes.ZeroInt())

			result, err := cliCtx.BroadcastTx(cdc.MustMarshalBinaryBare(stdTx))

			msg, _ := cdc.MarshalJSON(result)
			fmt.Println(string(msg))

			return err
		},
	}

	cmd.Flags().String(flagName, "", "name")
	cmd.Flags().String(flagPubKey, "", "public key，amino serialized json")

	return cmd
}
