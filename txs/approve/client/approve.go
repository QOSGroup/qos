package approve

import (
	"errors"
	"github.com/QOSGroup/qbase/client/context"
	qclitx "github.com/QOSGroup/qbase/client/tx"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/txs/approve"
	"github.com/QOSGroup/qos/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
)

type operateType int

const (
	createType operateType = iota
	increaseType
	decreaseType
	useType
	cancleType

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

			fromAddr, err := qclitx.GetAddrFromFlag(cliCtx, flagFrom)
			if err != nil {
				return err
			}

			toAddr, err := qclitx.GetAddrFromFlag(cliCtx, flagTo)
			if err != nil {
				return err
			}

			output, err := cliCtx.Query(queryPath, approve.BuildApproveKey(fromAddr.String(), toAddr.String()))
			if err != nil {
				return err
			}

			approve := approve.Approve{}
			cdc.MustUnmarshalBinaryBare(output, &approve)

			return cliCtx.PrintResult(approve)
		},
	}

	cmd.Flags().String(flagFrom, "", "Name or Address of approve creator")
	cmd.Flags().String(flagTo, "", "Name or Address of approve receiver")
	cmd.MarkFlagRequired(flagFrom)
	cmd.MarkFlagRequired(flagTo)

	return cmd
}

func CreateApproveCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-approve",
		Short: "Create approve",
		RunE: func(cmd *cobra.Command, args []string) error {
			return applyApprove(cdc, createType)
		},
	}

	cmd.Flags().String(flagFrom, "", "Name or Address of approve creator")
	cmd.Flags().String(flagTo, "", "Name or Address of approve receiver")
	cmd.Flags().Int64(flagQOS, 0, "Amount of QOS")
	cmd.Flags().String(flagQSCs, "", "Names and amounts of QSCs. ex: 100qstars,50qsc")
	cmd.MarkFlagRequired(flagFrom)
	cmd.MarkFlagRequired(flagTo)

	return cmd
}

func IncreaseApproveCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "increase-approve",
		Short: "Increase approve",
		RunE: func(cmd *cobra.Command, args []string) error {
			return applyApprove(cdc, increaseType)
		},
	}

	cmd.Flags().String(flagFrom, "", "Name or Address of approve creator")
	cmd.Flags().String(flagTo, "", "Name or Address of approve receiver")
	cmd.Flags().Int64(flagQOS, 0, "Amount of QOS")
	cmd.Flags().String(flagQSCs, "", "Names and amounts of QSCs. ex: 100qstars,50qsc")
	cmd.MarkFlagRequired(flagFrom)
	cmd.MarkFlagRequired(flagTo)

	return cmd
}

func DecreaseApproveCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decrease-approve",
		Short: "Decrease approve",
		RunE: func(cmd *cobra.Command, args []string) error {
			return applyApprove(cdc, decreaseType)
		},
	}

	cmd.Flags().String(flagFrom, "", "Name or Address of approve creator")
	cmd.Flags().String(flagTo, "", "Name or Address of approve receiver")
	cmd.Flags().Int64(flagQOS, 0, "Amount of QOS")
	cmd.Flags().String(flagQSCs, "", "Names and amounts of QSCs. ex: 100qstars,50qsc")
	cmd.MarkFlagRequired(flagFrom)
	cmd.MarkFlagRequired(flagTo)

	return cmd
}

func UseApproveCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use-approve",
		Short: "Use approve",
		RunE: func(cmd *cobra.Command, args []string) error {
			return applyApprove(cdc, useType)
		},
	}

	cmd.Flags().String(flagFrom, "", "Name or Address of approve creator")
	cmd.Flags().String(flagTo, "", "Name or Address of approve receiver")
	cmd.Flags().Int64(flagQOS, 0, "Amount of QOS")
	cmd.Flags().String(flagQSCs, "", "Names and amounts of QSCs. ex: 100qstars,50qsc")
	cmd.MarkFlagRequired(flagFrom)
	cmd.MarkFlagRequired(flagTo)

	return cmd
}

func CancelApproveCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-approve",
		Short: "Cancel approve",
		RunE: func(cmd *cobra.Command, args []string) error {
			return applyApprove(cdc, cancleType)
		},
	}

	cmd.Flags().String(flagFrom, "", "Name or Address of approve creator")
	cmd.Flags().String(flagTo, "", "Name or Address of approve receiver")
	cmd.MarkFlagRequired(flagFrom)
	cmd.MarkFlagRequired(flagTo)

	return cmd
}

func applyApprove(cdc *amino.Codec, operType operateType) error {
	iTxBuilder := func(ctx context.CLIContext) (txs.ITx, error) {
		if err := handleOperateFlag(ctx); err != nil {
			return nil, err
		}

		fromAddr := viper.Get(flagFrom).(btypes.Address)
		toAddr := viper.Get(flagTo).(btypes.Address)

		if operType == cancleType {
			return approve.TxCancelApprove{
				From: fromAddr,
				To:   toAddr,
			}, nil
		}

		qos := viper.GetInt64(flagQOS)
		qscs := viper.Get(flagQSCs).(types.QSCs)
		appr := approve.NewApprove(fromAddr, toAddr, btypes.NewInt(qos), qscs)

		switch operType {
		case createType:
			return approve.TxCreateApprove{Approve: appr}, nil
		case increaseType:
			return approve.TxIncreaseApprove{Approve: appr}, nil
		case decreaseType:
			return approve.TxDecreaseApprove{Approve: appr}, nil
		case useType:
			return approve.TxUseApprove{Approve: appr}, nil
		default:
			return nil, errors.New("operType invalid")
		}
	}

	return qclitx.BroadcastTxAndPrintResult(cdc, iTxBuilder)
}

func handleOperateFlag(ctx context.CLIContext) error {

	fromAddr, err := qclitx.GetAddrFromFlag(ctx, flagFrom)
	if err != nil {
		return err
	}

	toAddr, err := qclitx.GetAddrFromFlag(ctx, flagTo)
	if err != nil {
		return err
	}

	viper.Set(flagFrom, fromAddr)
	viper.Set(flagTo, toAddr)

	if viper.IsSet(flagQSCs) {
		qscsStr := viper.GetString(flagQSCs)
		_, qscs, err := types.ParseCoins(qscsStr)
		if err != nil {
			return err
		}
		viper.Set(flagQSCs, qscs)
	}

	return nil
}
