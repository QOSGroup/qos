package approve

import (
	"fmt"
	cliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/context"
	"github.com/QOSGroup/qbase/client/keys"
	btx "github.com/QOSGroup/qbase/client/tx"
	btxs "github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/txs/approve"
	"github.com/QOSGroup/qos/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
)

const (
	flagFrom = "from"
	flagTo   = "to"
	flagQOS  = "qos"
	flagQSCs = "qscs"
)

func QueryApproveCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve",
		Short: "Query approve by from and to",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			queryPath := "store/approve/key"

			fromStr := viper.GetString(flagFrom)
			_, err := btypes.GetAddrFromBech32(fromStr)
			if err != nil {
				return err
			}

			toStr := viper.GetString(flagTo)
			_, err = btypes.GetAddrFromBech32(toStr)
			if err != nil {
				return err
			}

			output, err := cliCtx.Query(queryPath, approve.BuildApproveKey(fromStr, toStr))
			if err != nil {
				return err
			}

			approve := approve.Approve{}
			cdc.MustUnmarshalBinaryBare(output, &approve)
			fmt.Println(cliCtx.ToJSONIndentStr(approve))

			return err
		},
	}

	cmd.Flags().String(flagFrom, "", "Address of approve creator")
	cmd.Flags().String(flagTo, "", "Address of approve receiver")
	cmd.MarkFlagRequired(flagFrom)
	cmd.MarkFlagRequired(flagTo)

	return cmd
}

func CreateApproveCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-approve",
		Short: "Create approve",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc)

			fromName := viper.GetString(flagFrom)
			fromInfo, err := keys.GetKeyInfo(cliCtx, fromName)
			if err != nil {
				return err
			}
			from, err := cliacc.GetAccount(cliCtx, fromInfo.GetAddress())
			if err != nil {
				return err
			}

			toStr := viper.GetString(flagTo)
			toAddr, err := btypes.GetAddrFromBech32(toStr)
			if err != nil {
				return err
			}

			qos := viper.GetInt64(flagQOS)
			qscsStr := viper.GetString(flagQSCs)
			_, qscs, err := types.ParseCoins(qscsStr)
			if err != nil {
				return err
			}

			chainId, err := types.GetDefaultChainId()
			if err != nil {
				return err
			}
			tx := btxs.NewTxStd(
				approve.ApproveCreateTx{
					Approve: approve.NewApprove(from.GetAddress(),
						toAddr,
						btypes.NewInt(qos),
						qscs),
				},
				chainId,
				btypes.NewInt(0))
			tx, err = btx.SignStdTx(cliCtx, fromName, from.GetNonce()+1, tx)
			if err != nil {
				return err
			}

			result, err := cliCtx.BroadcastTx(cdc.MustMarshalBinaryBare(tx))

			msg, _ := cdc.MarshalJSON(result)
			fmt.Println(string(msg))

			return err
		},
	}

	cmd.Flags().String(flagFrom, "", "Name of approve creator")
	cmd.Flags().String(flagTo, "", "Address of approve receiver")
	cmd.Flags().Int64(flagQOS, 0, "Amount of QOS")
	cmd.Flags().String(flagQSCs, "", "Names and amounts of QSCs")
	cmd.MarkFlagRequired(flagFrom)
	cmd.MarkFlagRequired(flagTo)

	return cmd
}

func IncreaseApproveCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "increase-approve",
		Short: "Increase approve",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc)

			fromName := viper.GetString(flagFrom)
			fromInfo, err := keys.GetKeyInfo(cliCtx, fromName)
			if err != nil {
				return err
			}
			from, err := cliacc.GetAccount(cliCtx, fromInfo.GetAddress())
			if err != nil {
				return err
			}

			toStr := viper.GetString(flagTo)
			toAddr, err := btypes.GetAddrFromBech32(toStr)
			if err != nil {
				return err
			}

			qos := viper.GetInt64(flagQOS)
			qscsStr := viper.GetString(flagQSCs)
			_, qscs, err := types.ParseCoins(qscsStr)
			if err != nil {
				return err
			}

			chainId, err := types.GetDefaultChainId()
			if err != nil {
				return err
			}

			tx := btxs.NewTxStd(
				approve.ApproveIncreaseTx{
					Approve: approve.NewApprove(from.GetAddress(),
						toAddr,
						btypes.NewInt(qos),
						qscs),
				},
				chainId,
				btypes.NewInt(0))
			tx, err = btx.SignStdTx(cliCtx, fromName, from.GetNonce()+1, tx)
			if err != nil {
				return err
			}

			result, err := cliCtx.BroadcastTx(cdc.MustMarshalBinaryBare(tx))

			msg, _ := cdc.MarshalJSON(result)
			fmt.Println(string(msg))

			return err
		},
	}

	cmd.Flags().String(flagFrom, "", "Name of approve creator")
	cmd.Flags().String(flagTo, "", "Address of approve receiver")
	cmd.Flags().Int64(flagQOS, 0, "Amount of QOS")
	cmd.Flags().String(flagQSCs, "", "Names and amounts of QSCs")
	cmd.MarkFlagRequired(flagFrom)
	cmd.MarkFlagRequired(flagTo)

	return cmd
}

func DecreaseApproveCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decrease-approve",
		Short: "Decrease approve",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc)

			fromName := viper.GetString(flagFrom)
			fromInfo, err := keys.GetKeyInfo(cliCtx, fromName)
			if err != nil {
				return err
			}
			from, err := cliacc.GetAccount(cliCtx, fromInfo.GetAddress())
			if err != nil {
				return err
			}

			toStr := viper.GetString(flagTo)
			toAddr, err := btypes.GetAddrFromBech32(toStr)
			if err != nil {
				return err
			}

			qos := viper.GetInt64(flagQOS)
			qscsStr := viper.GetString(flagQSCs)
			_, qscs, err := types.ParseCoins(qscsStr)
			if err != nil {
				return err
			}

			chainId, err := types.GetDefaultChainId()
			if err != nil {
				return err
			}

			tx := btxs.NewTxStd(
				approve.ApproveDecreaseTx{
					Approve: approve.NewApprove(from.GetAddress(),
						toAddr,
						btypes.NewInt(qos),
						qscs),
				},
				chainId,
				btypes.NewInt(0))
			tx, err = btx.SignStdTx(cliCtx, fromName, from.GetNonce()+1, tx)
			if err != nil {
				return err
			}

			result, err := cliCtx.BroadcastTx(cdc.MustMarshalBinaryBare(tx))

			msg, _ := cdc.MarshalJSON(result)
			fmt.Println(string(msg))

			return err
		},
	}

	cmd.Flags().String(flagFrom, "", "Name of approve creator")
	cmd.Flags().String(flagTo, "", "Address of approve receiver")
	cmd.Flags().Int64(flagQOS, 0, "Amount of QOS")
	cmd.Flags().String(flagQSCs, "", "Names and amounts of QSCs")
	cmd.MarkFlagRequired(flagFrom)
	cmd.MarkFlagRequired(flagTo)

	return cmd
}

func UseApproveCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use-approve",
		Short: "Use approve",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc)

			fromStr := viper.GetString(flagFrom)
			fromAddr, err := btypes.GetAddrFromBech32(fromStr)
			if err != nil {
				return err
			}

			toName := viper.GetString(flagTo)
			toInfo, err := keys.GetKeyInfo(cliCtx, toName)
			if err != nil {
				return err
			}
			to, err := cliacc.GetAccount(cliCtx, toInfo.GetAddress())
			if err != nil {
				return err
			}

			qos := viper.GetInt64(flagQOS)
			qscsStr := viper.GetString(flagQSCs)
			_, qscs, err := types.ParseCoins(qscsStr)
			if err != nil {
				return err
			}

			chainId, err := types.GetDefaultChainId()
			if err != nil {
				return err
			}

			tx := btxs.NewTxStd(
				approve.ApproveUseTx{
					Approve: approve.NewApprove(fromAddr,
						toInfo.GetAddress(),
						btypes.NewInt(qos),
						qscs),
				},
				chainId,
				btypes.NewInt(0))
			tx, err = btx.SignStdTx(cliCtx, toName, to.GetNonce()+1, tx)
			if err != nil {
				return err
			}

			result, err := cliCtx.BroadcastTx(cdc.MustMarshalBinaryBare(tx))

			msg, _ := cdc.MarshalJSON(result)
			fmt.Println(string(msg))

			return err
		},
	}

	cmd.Flags().String(flagFrom, "", "Address of approve creator")
	cmd.Flags().String(flagTo, "", "Name of approve receiver")
	cmd.Flags().Int64(flagQOS, 0, "Amount of QOS")
	cmd.Flags().String(flagQSCs, "", "Names and amounts of QSCs")
	cmd.MarkFlagRequired(flagFrom)
	cmd.MarkFlagRequired(flagTo)

	return cmd
}

func CancelApproveCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-approve",
		Short: "Cancel approve",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc)

			fromName := viper.GetString(flagFrom)
			fromInfo, err := keys.GetKeyInfo(cliCtx, fromName)
			if err != nil {
				return err
			}
			from, err := cliacc.GetAccount(cliCtx, fromInfo.GetAddress())
			if err != nil {
				return err
			}

			toStr := viper.GetString(flagTo)
			toAddr, err := btypes.GetAddrFromBech32(toStr)
			if err != nil {
				return err
			}

			chainId, err := types.GetDefaultChainId()
			if err != nil {
				return err
			}

			tx := btxs.NewTxStd(
				approve.ApproveCancelTx{
					From: fromInfo.GetAddress(),
					To:   toAddr,
				},
				chainId,
				btypes.NewInt(0))
			tx, err = btx.SignStdTx(cliCtx, fromName, from.GetNonce()+1, tx)
			if err != nil {
				return err
			}

			result, err := cliCtx.BroadcastTx(cdc.MustMarshalBinaryBare(tx))

			msg, _ := cdc.MarshalJSON(result)
			fmt.Println(string(msg))

			return err
		},
	}

	cmd.Flags().String(flagFrom, "", "Name of approve creator")
	cmd.Flags().String(flagTo, "", "Address of approve receiver")
	cmd.MarkFlagRequired(flagFrom)
	cmd.MarkFlagRequired(flagTo)

	return cmd
}
