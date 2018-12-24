package staking

import (
	"encoding/base64"
	"fmt"
	qcliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/context"
	qclitx "github.com/QOSGroup/qbase/client/tx"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qos/txs/staking"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

const (
	flagName        = "name"
	flagOwner       = "owner"
	flagPubKey      = "pubkey"
	flagBondTokens  = "tokens"
	flagDescription = "description"
)

func CreateValidatorCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-validator",
		Short: "Create validator",
		Long: `
pubkey is a tendermint validator pubkey. the public key of the validator used in
Tendermint consensus.

owner is a keystore name or account address.

ex: pubkey: {"type":"tendermint/PubKeyEd25519","value":"VOn2rPx+t7Njdgi+eLb+jBuF175T1b7LAcHElsmIuXA="}

example:

	 qoscli create-validator --name validatorName --owner ownerName --pubkey "VOn2rPx+t7Njdgi+eLb+jBuF175T1b7LAcHElsmIuXA=" --tokens 100

		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return qclitx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (txs.ITx, error) {
				name := viper.GetString(flagName)
				if len(name) == 0 {
					return nil, errors.New("name is empty")
				}
				tokens := uint64(viper.GetInt64(flagBondTokens))
				if tokens <= 0 {
					return nil, errors.New("tokens lte zero")
				}
				desc := viper.GetString(flagDescription)
				valPubkey := viper.GetString(flagPubKey)
				owner, err := qcliacc.GetAddrFromFlag(ctx, flagOwner)
				if err != nil {
					return nil, err
				}

				bz, err := base64.StdEncoding.DecodeString(valPubkey)
				if err != nil {
					return nil, fmt.Errorf("pubkey parse error: %s", err.Error())
				}
				var cKey ed25519.PubKeyEd25519
				copy(cKey[:], bz)

				return staking.NewCreateValidatorTx(name, owner, cKey, tokens, desc), nil
			})

		},
	}

	cmd.Flags().String(flagName, "", "name for validator")
	cmd.Flags().String(flagOwner, "", "keystore name or account address")
	cmd.Flags().String(flagPubKey, "", "tendermint consensus validator public key")
	cmd.Flags().Int64(flagBondTokens, 0, "bond tokens amount")
	cmd.Flags().String(flagDescription, "", "description")

	cmd.MarkFlagRequired(flagName)
	cmd.MarkFlagRequired(flagOwner)
	cmd.MarkFlagRequired(flagPubKey)
	cmd.MarkFlagRequired(flagBondTokens)

	return cmd
}

func RevokeValidatorCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke-validator",
		Short: "Revoke validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			return qclitx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (txs.ITx, error) {
				owner, err := qcliacc.GetAddrFromFlag(ctx, flagOwner)
				if err != nil {
					return nil, err
				}

				return staking.NewRevokeValidatorTx(owner), nil
			})

		},
	}

	cmd.Flags().String(flagOwner, "", "owner keystore name or address")

	cmd.MarkFlagRequired(flagOwner)

	return cmd
}

func ActiveValidatorCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "active-validator",
		Short: "Active validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			return qclitx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (txs.ITx, error) {
				owner, err := qcliacc.GetAddrFromFlag(ctx, flagOwner)
				if err != nil {
					return nil, err
				}

				return staking.NewActiveValidatorTx(owner), nil
			})

		},
	}

	cmd.Flags().String(flagOwner, "", "owner keystore or address")

	cmd.MarkFlagRequired(flagOwner)

	return cmd
}
