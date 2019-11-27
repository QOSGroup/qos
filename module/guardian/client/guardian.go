package client

import (
	qcliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/context"
	qcltx "github.com/QOSGroup/qbase/client/tx"
	"github.com/QOSGroup/qbase/store"
	"github.com/QOSGroup/qbase/txs"
	types2 "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/guardian/mapper"
	gtxs "github.com/QOSGroup/qos/module/guardian/txs"
	"github.com/QOSGroup/qos/module/guardian/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
)

// 添加系统账户
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

				tx := gtxs.NewTxAddGuardian(description, address, creator)
				if err = tx.ValidateInputs(); err != nil {
					return nil, err
				}
				return tx, nil
			})
		},
	}

	cmd.Flags().String(flagAddress, "", "Address of guardian")
	cmd.Flags().String(flagCreator, "", "Address of creator")
	cmd.Flags().String(flagDescription, "", "Description")
	cmd.MarkFlagRequired(flagAddress)
	cmd.MarkFlagRequired(flagCreator)
	cmd.MarkFlagRequired(flagDescription)

	return cmd
}

// 删除系统账户
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

				tx := gtxs.NewTxDeleteGuardian(address, deleteBy)
				if err = tx.ValidateInputs(); err != nil {
					return nil, err
				}
				return tx, nil
			})
		},
	}

	cmd.Flags().String(flagAddress, "", "Address of guardian")
	cmd.Flags().String(flagDeletedBy, "", "Address of deleteBy guardian")
	cmd.MarkFlagRequired(flagAddress)
	cmd.MarkFlagRequired(flagDeletedBy)

	return cmd
}

// 查询系统账户
func QueryGuardianCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "guardian [guardian]",
		Short: "Query guardian",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			address, err := qcliacc.GetAddrFromValue(cliCtx, args[0])
			if err != nil {
				return err
			}

			result, err := getGuardian(cliCtx, address)
			if err != nil {
				return err
			}
			return cliCtx.PrintResult(result)
		},
	}

	return cmd
}

func getGuardian(cliCtx context.CLIContext, guardian types2.AccAddress) (types.Guardian, error) {
	queryPath := "store/guardian/key"
	output, err := cliCtx.Query(queryPath, mapper.KeyGuardian(guardian))
	if err != nil {
		return types.Guardian{}, err
	}

	if output == nil {
		return types.Guardian{}, context.RecordsNotFoundError
	}

	result := types.Guardian{}
	cliCtx.Codec.MustUnmarshalBinaryBare(output, &result)
	return result, err
}

// 系统账户列表
func QueryGuardiansCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "guardians",
		Short: "Query guardian list",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			result, err := queryAllGuardians(cliCtx)
			if err != nil {
				return err
			}

			return cliCtx.PrintResult(result)
		},
	}

	return cmd
}

func queryAllGuardians(cliCtx context.CLIContext) ([]types.Guardian, error) {
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	result, err := node.ABCIQuery("store/guardian/subspace", mapper.KeyGuardiansSubspace())

	if err != nil {
		return nil, err
	}

	if len(result.Response.Value) == 0 {
		return nil, context.RecordsNotFoundError
	}

	var guardians []types.Guardian
	var vKVPair []store.KVPair
	err = cliCtx.Codec.UnmarshalBinaryLengthPrefixed(result.Response.Value, &vKVPair)
	for _, kv := range vKVPair {
		var guardian types.Guardian
		err = cliCtx.Codec.UnmarshalBinaryBare(kv.Value, &guardian)
		guardians = append(guardians, guardian)
	}

	return guardians, err
}

// 停网操作
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

				tx := gtxs.NewTxHaltNetwork(address, description)
				if err = tx.ValidateInputs(); err != nil {
					return nil, err
				}
				return tx, nil
			})
		},
	}

	cmd.Flags().String(flagAddress, "", "Address of guardian")
	cmd.Flags().String(flagDescription, "", "Description for this operation")
	cmd.MarkFlagRequired(flagAddress)
	cmd.MarkFlagRequired(flagDescription)

	return cmd
}
