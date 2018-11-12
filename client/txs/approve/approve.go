package approve

import (
	"fmt"
	cliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/context"
	"github.com/QOSGroup/qbase/client/keys"
	btx "github.com/QOSGroup/qbase/client/tx"
	btxs "github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/client"
	"github.com/QOSGroup/qos/txs/approve"
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

			fromName := viper.GetString(flagFrom)
			fromInfo, err := keys.GetKeyInfo(cliCtx, fromName)
			if err != nil {
				return err
			}

			toName := viper.GetString(flagTo)
			toInfo, err := keys.GetKeyInfo(cliCtx, toName)
			if err != nil {
				return err
			}

			output, err := cliCtx.Query(queryPath, approve.BuildApproveKey(fromInfo.GetAddress().String(), toInfo.GetAddress().String()))
			if err != nil {
				return err
			}

			if len(toName) > 0 {
				approve := approve.Approve{}
				cdc.MustUnmarshalBinaryBare(output, &approve)
				fmt.Println(cliCtx.ToJSONIndentStr(approve))
			} else {
				approves := new([]approve.Approve)
				cdc.MustUnmarshalBinary(output, approves)
				fmt.Println(cliCtx.ToJSONIndentStr(approves))
			}

			return err
		},
	}

	cmd.Flags().String(flagFrom, "", "Account name to create approve")
	cmd.Flags().String(flagTo, "", "Account name to receive approve")

	return cmd
}

func CreateApproveCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve-create",
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

			toName := viper.GetString(flagTo)
			toInfo, err := keys.GetKeyInfo(cliCtx, toName)
			if err != nil {
				return err
			}

			qos := viper.GetInt64(flagQOS)
			qscsStr := viper.GetString(flagQSCs)
			_, qscs, err := client.ParseCoins(qscsStr)
			if err != nil {
				return err
			}

			chainId, err := client.GetDefaultChainId()
			if err != nil {
				return err
			}
			tx := btxs.NewTxStd(
				approve.ApproveCreateTx{
					Approve: approve.NewApprove(from.GetAddress(),
						toInfo.GetAddress(),
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

	cmd.Flags().String(flagFrom, "", "Account name to create approve")
	cmd.Flags().String(flagTo, "", "Account name to receive approve")
	cmd.Flags().Int64(flagQOS, 0, "Amount of QOS")
	cmd.Flags().String(flagQSCs, "", "Names and amounts of QSCs")

	return cmd
}

func IncreaseApproveCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve-increase",
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

			toName := viper.GetString(flagTo)
			toInfo, err := keys.GetKeyInfo(cliCtx, toName)
			if err != nil {
				return err
			}

			qos := viper.GetInt64(flagQOS)
			qscsStr := viper.GetString(flagQSCs)
			_, qscs, err := client.ParseCoins(qscsStr)
			if err != nil {
				return err
			}

			chainId, err := client.GetDefaultChainId()
			if err != nil {
				return err
			}

			tx := btxs.NewTxStd(
				approve.ApproveIncreaseTx{
					Approve: approve.NewApprove(from.GetAddress(),
						toInfo.GetAddress(),
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

	cmd.Flags().String(flagFrom, "", "Account name to create approve")
	cmd.Flags().String(flagTo, "", "Account name to receive approve")
	cmd.Flags().Int64(flagQOS, 0, "Amount of QOS")
	cmd.Flags().String(flagQSCs, "", "Names and amounts of QSCs")

	return cmd
}

func DecreaseApproveCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve-decrease",
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

			toName := viper.GetString(flagTo)
			toInfo, err := keys.GetKeyInfo(cliCtx, toName)
			if err != nil {
				return err
			}

			qos := viper.GetInt64(flagQOS)
			qscsStr := viper.GetString(flagQSCs)
			_, qscs, err := client.ParseCoins(qscsStr)
			if err != nil {
				return err
			}

			chainId, err := client.GetDefaultChainId()
			if err != nil {
				return err
			}

			tx := btxs.NewTxStd(
				approve.ApproveDecreaseTx{
					Approve: approve.NewApprove(from.GetAddress(),
						toInfo.GetAddress(),
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

	cmd.Flags().String(flagFrom, "", "Account name to create approve")
	cmd.Flags().String(flagTo, "", "Account name to receive approve")
	cmd.Flags().Int64(flagQOS, 0, "Amount of QOS")
	cmd.Flags().String(flagQSCs, "", "Names and amounts of QSCs")

	return cmd
}

func UseApproveCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve-use",
		Short: "Use approve",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc)

			fromName := viper.GetString(flagFrom)
			fromInfo, err := keys.GetKeyInfo(cliCtx, fromName)
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
			_, qscs, err := client.ParseCoins(qscsStr)
			if err != nil {
				return err
			}

			chainId, err := client.GetDefaultChainId()
			if err != nil {
				return err
			}

			tx := btxs.NewTxStd(
				approve.ApproveUseTx{
					Approve: approve.NewApprove(fromInfo.GetAddress(),
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

	cmd.Flags().String(flagFrom, "", "Account name to create approve")
	cmd.Flags().String(flagTo, "", "Account name to receive approve")
	cmd.Flags().Int64(flagQOS, 0, "Amount of QOS")
	cmd.Flags().String(flagQSCs, "", "Names and amounts of QSCs")

	return cmd
}

func CancelApproveCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve-cancel",
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

			toName := viper.GetString(flagTo)
			toInfo, err := keys.GetKeyInfo(cliCtx, toName)
			if err != nil {
				return err
			}

			chainId, err := client.GetDefaultChainId()
			if err != nil {
				return err
			}

			tx := btxs.NewTxStd(
				approve.ApproveCancelTx{
					From: fromInfo.GetAddress(),
					To:   toInfo.GetAddress(),
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

	cmd.Flags().String(flagFrom, "", "Account name to create approve")
	cmd.Flags().String(flagTo, "", "Account name to receive approve")

	return cmd
}
