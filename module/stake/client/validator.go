package staking

import (
	qcliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/context"
	qclitx "github.com/QOSGroup/qbase/client/tx"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qos/module/stake"
	"github.com/QOSGroup/qos/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/privval"
	"path/filepath"
)

const (
	flagName        = "name"
	flagOwner       = "owner"
	flagBondTokens  = "tokens"
	flagDescription = "description"
	flagCompound    = "compound"
	flagNodeHome    = "nodeHome"
)

func CreateValidatorCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-validator",
		Short: "create new validator initialized with a self-delegation to it",
		Long: `
owner is a keystore name or account address.

example:

	 qoscli tx create-validator --name validatorName --owner ownerName --tokens 100

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

				privValidator := privval.LoadOrGenFilePV(filepath.Join(viper.GetString(flagNodeHome), cfg.DefaultConfig().PrivValidatorFile()))

				owner, err := qcliacc.GetAddrFromFlag(ctx, flagOwner)
				if err != nil {
					return nil, err
				}

				isCompound := viper.GetBool(flagCompound)
				return stake.NewCreateValidatorTx(name, owner, privValidator.PubKey, tokens, isCompound, desc), nil
			})

		},
	}

	cmd.Flags().String(flagName, "", "name for validator")
	cmd.Flags().String(flagOwner, "", "keystore name or account address")
	cmd.Flags().Int64(flagBondTokens, 0, "bond tokens amount")
	cmd.Flags().Bool(flagCompound, false, "as a self-delegator, whether the income is calculated as compound interest")
	cmd.Flags().String(flagDescription, "", "description")
	cmd.Flags().String(flagNodeHome, types.DefaultNodeHome, "path of node's config and data files, default: $HOME/.qosd")

	cmd.MarkFlagRequired(flagName)
	cmd.MarkFlagRequired(flagOwner)
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

				return stake.NewRevokeValidatorTx(owner), nil
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

				return stake.NewActiveValidatorTx(owner), nil
			})

		},
	}

	cmd.Flags().String(flagOwner, "", "owner keystore or address")

	cmd.MarkFlagRequired(flagOwner)

	return cmd
}
