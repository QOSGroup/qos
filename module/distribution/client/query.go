package distribution

import (
	qcliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/distribution"
	ecotypes "github.com/QOSGroup/qos/module/eco/types"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"
)

const (
	flagOwner     = "owner"
	flagDelegator = "delegator"
)

func queryValidatorPeriodCommand(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-period",
		Short: "Query distribution validator period info",
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var owner btypes.Address
			if o, err := qcliacc.GetAddrFromFlag(cliCtx, flagOwner); err == nil {
				owner = o
			}

			path := ecotypes.BuildQueryValidatorPeriodInfoCustomQueryPath(owner)
			res, err := cliCtx.Query(path, []byte(""))
			if err != nil {
				return err
			}

			var result distribution.ValidatorPeriodInfoQueryResult
			cliCtx.Codec.UnmarshalJSON(res, &result)
			return cliCtx.PrintResult(result)
		},
	}

	cmd.Flags().String(flagOwner, "", "validator's owner address")
	cmd.MarkFlagRequired(flagOwner)
	return cmd
}

func queryDelegatorIncomeInfoCommand(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegator-income",
		Short: "Query distribution delegator income info",
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var owner btypes.Address
			var delegator btypes.Address

			if o, err := qcliacc.GetAddrFromFlag(cliCtx, flagOwner); err == nil {
				owner = o
			}

			if d, err := qcliacc.GetAddrFromFlag(cliCtx, flagDelegator); err == nil {
				delegator = d
			}

			path := ecotypes.BuildQueryDelegatorIncomeInfoCustomQueryPath(delegator, owner)
			res, err := cliCtx.Query(path, []byte(""))
			if err != nil {
				return err
			}

			var result distribution.DelegatorIncomeInfoQueryResult
			cliCtx.Codec.UnmarshalJSON(res, &result)
			return cliCtx.PrintResult(result)
		},
	}

	cmd.Flags().String(flagOwner, "", "validator's owner address")
	cmd.Flags().String(flagDelegator, "", "delegator address")

	cmd.MarkFlagRequired(flagDelegator)
	cmd.MarkFlagRequired(flagOwner)
	return cmd
}
