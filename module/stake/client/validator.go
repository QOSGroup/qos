package staking

import (
	qcliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/context"
	qclitx "github.com/QOSGroup/qbase/client/tx"
	btxs "github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qos/module/stake/txs"
	"github.com/QOSGroup/qos/module/stake/types"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/privval"
	"path/filepath"
)

const (
	flagOwner      = "owner"
	flagBondTokens = "tokens"

	flagCompound = "compound"
	flagNodeHome = "home-node"

	// flags for validator's description
	flagMoniker = "moniker"
	flagLogo    = "logo"
	flagWebsite = "website"
	flagDetails = "details"
)

func CreateValidatorCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-validator",
		Short: "create new validator initialized with a self-delegation to it",
		Long: `
owner is a keystore name or account address.

example:

	 qoscli tx create-validator --moniker validatorName --owner ownerName --tokens 100

		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return qclitx.BroadcastTxAndPrintResult(cdc, TxCreateValidatorBuilder)
		},
	}

	cmd.Flags().String(flagOwner, "", "keystore name or account address")
	cmd.Flags().Int64(flagBondTokens, 0, "bond tokens amount")
	cmd.Flags().Bool(flagCompound, false, "as a self-delegator, whether the income is calculated as compound interest")
	cmd.Flags().String(flagNodeHome, qtypes.DefaultNodeHome, "path of node's config and data files, default: $HOME/.qosd")
	cmd.Flags().String(flagMoniker, "", "The validator's name")
	cmd.Flags().String(flagLogo, "", "The optional logo link")
	cmd.Flags().String(flagWebsite, "", "The validator's (optional) website")
	cmd.Flags().String(flagDetails, "", "The validator's (optional) details")

	cmd.MarkFlagRequired(flagMoniker)
	cmd.MarkFlagRequired(flagOwner)
	cmd.MarkFlagRequired(flagBondTokens)

	return cmd
}

func ModifyValidatorCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modify-validator",
		Short: "modify an existing validator account",
		RunE: func(cmd *cobra.Command, args []string) error {
			return qclitx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (tx btxs.ITx, e error) {
				name := viper.GetString(flagMoniker)
				logo := viper.GetString(flagLogo)
				website := viper.GetString(flagWebsite)
				details := viper.GetString(flagDetails)
				desc := types.Description{
					name, logo, website, details,
				}

				owner, err := qcliacc.GetAddrFromFlag(ctx, flagOwner)
				if err != nil {
					return nil, err
				}

				return txs.NewModifyValidatorTx(owner, desc), nil
			})
		},
	}

	cmd.Flags().String(flagMoniker, "", "The validator's name")
	cmd.Flags().String(flagOwner, "", "keystore name or account address")
	cmd.Flags().String(flagLogo, "", "The optional logo link")
	cmd.Flags().String(flagWebsite, "", "The validator's (optional) website")
	cmd.Flags().String(flagDetails, "", "The validator's (optional) details")

	cmd.MarkFlagRequired(flagOwner)

	return cmd
}

func RevokeValidatorCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke-validator",
		Short: "Revoke validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			return qclitx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (btxs.ITx, error) {
				owner, err := qcliacc.GetAddrFromFlag(ctx, flagOwner)
				if err != nil {
					return nil, err
				}

				return txs.NewRevokeValidatorTx(owner), nil
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
			return qclitx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (btxs.ITx, error) {
				owner, err := qcliacc.GetAddrFromFlag(ctx, flagOwner)
				if err != nil {
					return nil, err
				}

				return txs.NewActiveValidatorTx(owner), nil
			})

		},
	}

	cmd.Flags().String(flagOwner, "", "owner keystore or address")

	cmd.MarkFlagRequired(flagOwner)

	return cmd
}

func TxCreateValidatorBuilder(ctx context.CLIContext) (btxs.ITx, error) {
	name := viper.GetString(flagMoniker)
	if len(name) == 0 {
		return nil, errors.New("moniker is empty")
	}
	tokens := uint64(viper.GetInt64(flagBondTokens))
	if tokens <= 0 {
		return nil, errors.New("tokens lte zero")
	}
	logo := viper.GetString(flagLogo)
	website := viper.GetString(flagWebsite)
	details := viper.GetString(flagDetails)
	desc := types.Description{
		name, logo, website, details,
	}

	privValidator := privval.LoadOrGenFilePV(filepath.Join(viper.GetString(flagNodeHome), cfg.DefaultConfig().PrivValidatorKeyFile()),
		filepath.Join(viper.GetString(flagNodeHome), cfg.DefaultConfig().PrivValidatorKeyFile()))

	owner, err := qcliacc.GetAddrFromFlag(ctx, flagOwner)
	if err != nil {
		return nil, err
	}

	isCompound := viper.GetBool(flagCompound)
	return txs.NewCreateValidatorTx(owner, privValidator.GetPubKey(), tokens, isCompound, desc), nil
}
