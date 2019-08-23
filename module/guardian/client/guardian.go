package client

import (
	"errors"
	qcliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/context"
	qcltx "github.com/QOSGroup/qbase/client/tx"
	"github.com/QOSGroup/qbase/store"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qos/module/guardian/mapper"
	gtxs "github.com/QOSGroup/qos/module/guardian/txs"
	"github.com/QOSGroup/qos/module/guardian/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
)

func AddGuardianCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-guardian",
		Short: "Add guardian",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			return qcltx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (txs.ITx, error) {
				address, err := qcliacc.GetAddrFromFlag(ctx, flagAddress)
				if err != nil {
					return nil, err
				}

				creator, err := qcliacc.GetAddrFromFlag(ctx, flagCreator)
				if err != nil {
					return nil, err
				}

				description := viper.GetString(flagDescription)
				if len(description) < 0 || len(description) > gtxs.MaxDescriptionLen {

				}

				return gtxs.NewTxAddGuardian(description, address, creator), nil
			})
		},
	}

	cmd.Flags().String(flagAddress, "", "address of guardian")
	cmd.Flags().String(flagCreator, "", "address of creator")
	cmd.Flags().String(flagDescription, "", "description")
	cmd.MarkFlagRequired(flagAddress)
	cmd.MarkFlagRequired(flagCreator)
	cmd.MarkFlagRequired(flagDescription)

	return cmd
}

func DeleteGuardianCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-guardian",
		Short: "Delete guardian",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			return qcltx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (txs.ITx, error) {
				address, err := qcliacc.GetAddrFromFlag(ctx, flagAddress)
				if err != nil {
					return nil, err
				}

				deleteBy, err := qcliacc.GetAddrFromFlag(ctx, flagDeletedBy)
				if err != nil {
					return nil, err
				}

				description := viper.GetString(flagDescription)
				if len(description) < 0 || len(description) > gtxs.MaxDescriptionLen {

				}

				return gtxs.NewTxDeleteGuardian(address, deleteBy), nil
			})
		},
	}

	cmd.Flags().String(flagAddress, "", "address of guardian")
	cmd.Flags().String(flagDeletedBy, "", "address of deleteBy guardian")
	cmd.MarkFlagRequired(flagAddress)
	cmd.MarkFlagRequired(flagDeletedBy)

	return cmd
}

func QueryGuardianCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "guardian [guardian]",
		Short: "Query guardian",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			queryPath := "store/guardian/key"

			address, err := qcliacc.GetAddrFromValue(cliCtx, args[0])
			if err != nil {
				return err
			}

			output, err := cliCtx.Query(queryPath, mapper.KeyGuardian(address))
			if err != nil {
				return err
			}

			if output == nil {
				return errors.New("guardian does not exist")
			}

			guardian := types.Guardian{}
			cdc.MustUnmarshalBinaryBare(output, &guardian)

			return cliCtx.PrintResult(guardian)
		},
	}

	return cmd
}

func QueryGuardiansCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "guardians",
		Short: "Query guardian list",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			node, err := cliCtx.GetNode()
			if err != nil {
				return err
			}

			result, err := node.ABCIQuery("store/guardian/subspace", mapper.KeyGuardiansSubspace())

			if err != nil {
				return err
			}

			if len(result.Response.Value) == 0 {
				return errors.New("no guardian")
			}

			var guardians []types.Guardian
			var vKVPair []store.KVPair
			cdc.UnmarshalBinaryLengthPrefixed(result.Response.Value, &vKVPair)
			for _, kv := range vKVPair {
				var guardian types.Guardian
				cdc.UnmarshalBinaryBare(kv.Value, &guardian)
				guardians = append(guardians, guardian)
			}

			return cliCtx.PrintResult(guardians)
		},
	}

	return cmd
}

func HaltCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "halt-network",
		Short: "Halt the network",
		RunE: func(cmd *cobra.Command, args []string) error {
			return qcltx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (txs.ITx, error) {
				address, err := qcliacc.GetAddrFromFlag(ctx, flagAddress)
				if err != nil {
					return nil, err
				}

				description := viper.GetString(flagDescription)
				if len(description) < 0 || len(description) > gtxs.MaxDescriptionLen {

				}

				return gtxs.NewTxHaltNetwork(address, description), nil
			})
		},
	}

	cmd.Flags().String(flagAddress, "", "address of guardian")
	cmd.Flags().String(flagDescription, "", "description for this operation")
	cmd.MarkFlagRequired(flagAddress)
	cmd.MarkFlagRequired(flagDescription)

	return cmd
}
