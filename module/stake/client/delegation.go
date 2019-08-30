package client

import (
	qcliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/context"
	qclitx "github.com/QOSGroup/qbase/client/tx"
	btxs "github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qos/module/stake/txs"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
)

const (
	flagDelegator     = "delegator"
	flagValidator     = "validator"
	flagFromValidator = "from-validator"
	flagToValidator   = "to-validator"
	flagAll           = "all"
)

func CreateDelegationCommand(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegate",
		Short: "delegate QOS to a validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			return qclitx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (btxs.ITx, error) {

				tokens := viper.GetInt64(flagBondTokens)
				if tokens <= 0 {
					return nil, errors.New("delegate QOS amount must gt 0")
				}

				validatorAddr, err := qcliacc.GetValidatorAddrFromFlag(ctx, flagValidator)
				if err != nil {
					return nil, err
				}

				delegator, err := qcliacc.GetAddrFromFlag(ctx, flagDelegator)
				if err != nil {
					return nil, err
				}

				return &txs.TxCreateDelegation{
					Delegator:     delegator,
					ValidatorAddr: validatorAddr,
					Amount:        uint64(tokens),
					IsCompound:    viper.GetBool(flagCompound),
				}, nil
			})
		},
	}

	cmd.Flags().String(flagDelegator, "", "delegator account address")
	cmd.Flags().String(flagValidator, "", "keystore name or account of validator address")
	cmd.Flags().Int64(flagBondTokens, 0, "amount of QOS to delegate")
	cmd.Flags().Bool(flagCompound, false, " whether the income is calculated as compound interest")

	cmd.MarkFlagRequired(flagDelegator)
	cmd.MarkFlagRequired(flagValidator)
	cmd.MarkFlagRequired(flagBondTokens)

	return cmd
}

func CreateModifyCompoundCommand(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modify-compound",
		Short: "modify compound info in a delegation",
		RunE: func(cmd *cobra.Command, args []string) error {
			return qclitx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (btxs.ITx, error) {

				validatorAddr, err := qcliacc.GetValidatorAddrFromFlag(ctx, flagValidator)
				if err != nil {
					return nil, err
				}

				delegator, err := qcliacc.GetAddrFromFlag(ctx, flagDelegator)
				if err != nil {
					return nil, err
				}

				return &txs.TxModifyCompound{
					Delegator:     delegator,
					ValidatorAddr: validatorAddr,
					IsCompound:    viper.GetBool(flagCompound),
				}, nil
			})
		},
	}

	cmd.Flags().String(flagDelegator, "", "delegator account address")
	cmd.Flags().String(flagValidator, "", "keystore name or account of validator address")
	cmd.Flags().Bool(flagCompound, false, " whether the income is calculated as compound interest")

	cmd.MarkFlagRequired(flagDelegator)
	cmd.MarkFlagRequired(flagValidator)

	return cmd
}

func CreateUnbondDelegationCommand(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unbond",
		Short: "unbond QOS from a validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			return qclitx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (btxs.ITx, error) {

				tokens := viper.GetInt64(flagBondTokens)
				isUnbondAll := viper.GetBool(flagAll)

				if isUnbondAll {
					tokens = 0
				}

				if !isUnbondAll && tokens <= 0 {
					return nil, errors.New("unbond QOS amount must gt 0")
				}

				validatorAddr, err := qcliacc.GetValidatorAddrFromFlag(ctx, flagValidator)
				if err != nil {
					return nil, err
				}

				delegator, err := qcliacc.GetAddrFromFlag(ctx, flagDelegator)
				if err != nil {
					return nil, err
				}

				return &txs.TxUnbondDelegation{
					Delegator:     delegator,
					ValidatorAddr: validatorAddr,
					UnbondAmount:  uint64(tokens),
					IsUnbondAll:   isUnbondAll,
				}, nil
			})
		},
	}

	cmd.Flags().String(flagDelegator, "", "delegator account address")
	cmd.Flags().String(flagValidator, "", "keystore name or account of validator address")
	cmd.Flags().Int64(flagBondTokens, 0, "amount of QOS to unbond")
	cmd.Flags().Bool(flagAll, false, "whether unbond all QOS amount. override --tokens if true")

	cmd.MarkFlagRequired(flagDelegator)
	cmd.MarkFlagRequired(flagValidator)

	return cmd
}

func CreateReDelegationCommand(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redelegate",
		Short: "redelegate QOS from a validator to another",
		RunE: func(cmd *cobra.Command, args []string) error {
			return qclitx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (btxs.ITx, error) {

				tokens := viper.GetInt64(flagBondTokens)
				all := viper.GetBool(flagAll)

				if all {
					tokens = 0
				}

				if !all && tokens <= 0 {
					return nil, errors.New("redelegate QOS amount must gt 0")
				}

				delegator, err := qcliacc.GetAddrFromFlag(ctx, flagDelegator)
				if err != nil {
					return nil, err
				}

				fromValidatorAddr, err := qcliacc.GetValidatorAddrFromFlag(ctx, flagFromValidator)
				if err != nil {
					return nil, err
				}

				toValidatorAddr, err := qcliacc.GetValidatorAddrFromFlag(ctx, flagToValidator)
				if err != nil {
					return nil, err
				}

				return &txs.TxCreateReDelegation{
					Delegator:         delegator,
					FromValidatorAddr: fromValidatorAddr,
					ToValidatorAddr:   toValidatorAddr,
					Amount:            uint64(tokens),
					IsCompound:        viper.GetBool(flagCompound),
					IsRedelegateAll:   all,
				}, nil
			})
		},
	}

	cmd.Flags().String(flagDelegator, "", "delegator account address")
	cmd.Flags().String(flagFromValidator, "", "keystore name or account of validator address")
	cmd.Flags().String(flagToValidator, "", "keystore name or account of validator address")
	cmd.Flags().Int64(flagBondTokens, 0, "amount of QOS to redelegate")
	cmd.Flags().Bool(flagAll, false, "whether redelegate all QOS amount. override --tokens if true")
	cmd.Flags().Bool(flagCompound, false, "whether the income is calculated as compound interest")

	cmd.MarkFlagRequired(flagDelegator)
	cmd.MarkFlagRequired(flagFromValidator)
	cmd.MarkFlagRequired(flagToValidator)

	return cmd
}
