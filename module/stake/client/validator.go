package client

import (
	"fmt"
	"path/filepath"

	qcliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/context"
	qclitx "github.com/QOSGroup/qbase/client/tx"
	btxs "github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/stake/txs"
	"github.com/QOSGroup/qos/module/stake/types"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/privval"
)

const (
	flagCreator    = "creator"
	flagValidator  = "validator"
	flagBondTokens = "tokens"

	flagCompound = "compound"
	flagNodeHome = "home-node"

	// flags for validator's description
	flagMoniker = "moniker"
	flagLogo    = "logo"
	flagWebsite = "website"
	flagDetails = "details"

	flagCommissionRate          = "commission-rate"
	flagCommissionMaxRate       = "commission-max-rate"
	flagCommissionMaxChangeRate = "commission-max-change-rate"

	DefaultCommissionRate          = "0.1"
	DefaultCommissionMaxRate       = "0.2"
	DefaultCommissionMaxChangeRate = "0.01"
)

func CreateValidatorCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-validator",
		Short: "create new validator initialized with a self-delegation to it",
		Long: `
owner is a keystore name or account address.

example:

	 qoscli tx create-validator --moniker validatorName --creator ownerName --tokens 100

		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return qclitx.BroadcastTxAndPrintResult(cdc, TxCreateValidatorBuilder)
		},
	}

	cmd.Flags().String(flagCreator, "", "keystore name or account creator address")
	cmd.Flags().Int64(flagBondTokens, 0, "bond tokens amount")
	cmd.Flags().Bool(flagCompound, false, "as a self-delegator, whether the income is calculated as compound interest")
	cmd.Flags().String(flagNodeHome, qtypes.DefaultNodeHome, "path of node's config and data files, default: $HOME/.qosd")
	cmd.Flags().String(flagMoniker, "", "The validator's name")
	cmd.Flags().String(flagLogo, "", "The optional logo link")
	cmd.Flags().String(flagWebsite, "", "The validator's (optional) website")
	cmd.Flags().String(flagDetails, "", "The validator's (optional) details")
	cmd.Flags().String(flagCommissionRate, DefaultCommissionRate, "The initial commission rate percentage")
	cmd.Flags().String(flagCommissionMaxRate, DefaultCommissionMaxRate, "The maximum commission rate percentage")
	cmd.Flags().String(flagCommissionMaxChangeRate, DefaultCommissionMaxChangeRate, "The maximum commission change rate percentage (per day)")

	cmd.MarkFlagRequired(flagMoniker)
	cmd.MarkFlagRequired(flagCreator)
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

				owner, err := qcliacc.GetAddrFromFlag(ctx, flagValidator)
				if err != nil {
					return nil, err
				}

				var newRate *qtypes.Dec
				commissionRate := viper.GetString(flagCommissionRate)
				if commissionRate != "" {
					rate, err := qtypes.NewDecFromStr(commissionRate)
					if err != nil {
						return nil, fmt.Errorf("invalid new commission rate: %v", err)
					}

					newRate = &rate
				}

				return txs.NewModifyValidatorTx(btypes.ValAddress(owner), desc, newRate), nil
			})
		},
	}

	cmd.Flags().String(flagMoniker, "", "The validator's name")
	cmd.Flags().String(flagValidator, "", "keystore name or account of validator address")
	cmd.Flags().String(flagLogo, "", "The optional logo link")
	cmd.Flags().String(flagWebsite, "", "The validator's (optional) website")
	cmd.Flags().String(flagDetails, "", "The validator's (optional) details")
	cmd.Flags().String(flagCommissionRate, "", "The initial commission rate percentage")

	cmd.MarkFlagRequired(flagValidator)

	return cmd
}

func RevokeValidatorCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke-validator",
		Short: "Revoke validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			return qclitx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (btxs.ITx, error) {
				owner, err := qcliacc.GetAddrFromFlag(ctx, flagValidator)
				if err != nil {
					return nil, err
				}

				return txs.NewRevokeValidatorTx(btypes.ValAddress(owner)), nil
			})

		},
	}

	cmd.Flags().String(flagValidator, "", "keystore name or account of validator address")

	cmd.MarkFlagRequired(flagValidator)

	return cmd
}

func ActiveValidatorCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "active-validator",
		Short: "Active validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			return qclitx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (btxs.ITx, error) {
				owner, err := qcliacc.GetAddrFromFlag(ctx, flagValidator)
				if err != nil {
					return nil, err
				}

				tokens := viper.GetInt64(flagBondTokens)
				if tokens < 0 {
					return nil, errors.New("tokens lt zero")
				}

				return txs.NewActiveValidatorTx(btypes.ValAddress(owner), uint64(tokens)), nil
			})

		},
	}

	cmd.Flags().String(flagValidator, "", "keystore name or account of validator address")
	cmd.Flags().Int64(flagBondTokens, 0, "bond tokens amount to increase")

	cmd.MarkFlagRequired(flagValidator)

	return cmd
}

func TxCreateValidatorBuilder(ctx context.CLIContext) (btxs.ITx, error) {
	name := viper.GetString(flagMoniker)
	if len(name) == 0 {
		return nil, errors.New("moniker is empty")
	}

	tokens := viper.GetInt64(flagBondTokens)
	if tokens <= 0 {
		return nil, errors.New("tokens lte zero")
	}
	logo := viper.GetString(flagLogo)
	website := viper.GetString(flagWebsite)
	details := viper.GetString(flagDetails)
	desc := types.Description{
		name, logo, website, details,
	}

	commission, err := BuildCommissionRates()
	if err != nil {
		return nil, err
	}

	privValidator := privval.LoadOrGenFilePV(filepath.Join(viper.GetString(flagNodeHome), cfg.DefaultConfig().PrivValidatorKeyFile()),
		filepath.Join(viper.GetString(flagNodeHome), cfg.DefaultConfig().PrivValidatorKeyFile()))

	operator, err := qcliacc.GetAddrFromFlag(ctx, flagCreator)
	if err != nil {
		return nil, err
	}

	isCompound := viper.GetBool(flagCompound)
	return txs.NewCreateValidatorTx(operator, privValidator.GetPubKey(), uint64(tokens), isCompound, desc, *commission), nil
}

func BuildCommissionRates() (*types.CommissionRates, error) {
	rateStr := viper.GetString(flagCommissionRate)
	maxRateStr := viper.GetString(flagCommissionMaxRate)
	maxChangeRateStr := viper.GetString(flagCommissionMaxChangeRate)
	if rateStr == "" || maxRateStr == "" || maxChangeRateStr == "" {
		return nil, errors.New("must specify all validator commission parameters")
	}
	rate, err := qtypes.NewDecFromStr(rateStr)
	if err != nil {
		return nil, err
	}
	maxRate, err := qtypes.NewDecFromStr(maxRateStr)
	if err != nil {
		return nil, err
	}
	maxChangeRate, err := qtypes.NewDecFromStr(maxChangeRateStr)
	if err != nil {
		return nil, err
	}
	commission := types.NewCommissionRates(rate, maxRate, maxChangeRate)

	return &commission, nil
}
