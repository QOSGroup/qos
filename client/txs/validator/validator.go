package validator

import (


	"encoding/base64"
	"errors"
	"fmt"
	"github.com/QOSGroup/qbase/client/context"
	ctxs "github.com/QOSGroup/qbase/client/tx"
	btxs "github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qbase/client/keys"
	"github.com/QOSGroup/qos/client"
	"github.com/QOSGroup/qos/txs/validator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto"
)

const (
	flagName   = "name"
	flagNonce = "nonce"
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
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			name := viper.GetString(flagName)
			consPubkey := viper.GetString(flagConsPubKey)
			nonce := viper.GetInt64(flagNonce)

			if name == "" {
				return errors.New("missing name flag")
			}

			if consPubkey == "" {
				return errors.New("missing cons-pubkey flag")
			}

			if nonce == 0 {
				return errors.New("operator account nonce not valid")
			}

			bz, err := base64.StdEncoding.DecodeString(consPubkey)
			if err != nil {
				return fmt.Errorf("cons-pubkey parse error: %s" , err.Error())
			}
			var cKey ed25519.PubKeyEd25519
			copy(cKey[:], bz)

			validatorTx, err := buildCreateValidatorTx(cliCtx , name ,cKey)
			if err != nil {
				return err
			}

			chainID, err := client.GetDefaultChainId()
			if err != nil {
				return err
			}

			stdTx := btxs.NewTxStd(validatorTx, chainID, btypes.ZeroInt())

			tx, err:= ctxs.SignStdTx(cliCtx,name,nonce,stdTx)
			if err != nil {
				return err
			}

			result, err := cliCtx.BroadcastTx(cdc.MustMarshalBinaryBare(tx))

			msg, _ := cdc.MarshalJSON(result)
			fmt.Println(string(msg))

			return err
		},
	}

	cmd.Flags().String(flagName, "", "operator keys name")
	cmd.Flags().String(flagConsPubKey, "", "tendermint consensus validator public key")
	cmd.Flags().Int(flagNonce , 0 , "operator account nonce")

	cmd.MarkFlagRequired(flagNonce)
	cmd.MarkFlagRequired(flagName)
	cmd.MarkFlagRequired(flagConsPubKey)

	return cmd
}


func buildCreateValidatorTx(ctx context.CLIContext , name string , consPubkey crypto.PubKey) ( *validator.CreateValidatorTx , error ){

	info , err := keys.GetKeyInfo(ctx , name)
	if err != nil {
		return nil , err
	}

	tx := validator.NewCreateValidatorTx(name , consPubkey, info.GetAddress())
	return tx , nil
}

