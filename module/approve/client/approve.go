package client

import (
	"errors"

	qcliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/context"
	qclitx "github.com/QOSGroup/qbase/client/tx"
	"github.com/QOSGroup/qbase/txs"
	atxs "github.com/QOSGroup/qos/module/approve/txs"
	approvetypes "github.com/QOSGroup/qos/module/approve/types"
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

	flagFrom  = "from"
	flagTo    = "to"
	flagCoins = "coins"
)

// 查询预授权命令
func QueryApproveCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve",
		Short: "Query approve by from and to",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			// 查询路径，TODO 统一通过app CustomQueryHandler处理
			queryPath := "store/approve/key"
			// 解析授权账户地址
			fromAddr, err := qcliacc.GetAddrFromFlag(cliCtx, flagFrom)
			if err != nil {
				return err
			}
			// 解析被授权账户地址
			toAddr, err := qcliacc.GetAddrFromFlag(cliCtx, flagTo)
			if err != nil {
				return err
			}
			// 获取查询结果
			output, err := cliCtx.Query(queryPath, approvetypes.BuildApproveKey(fromAddr, toAddr))
			if err != nil {
				return err
			}
			if output == nil {
				return errors.New("approve does not exist")
			}

			// 反序列化查询结果
			approve := approvetypes.Approve{}
			cdc.MustUnmarshalBinaryBare(output, &approve)

			return cliCtx.PrintResult(approve)
		},
	}

	cmd.Flags().String(flagFrom, "", "Keybase name or address of approve creator")
	cmd.Flags().String(flagTo, "", "Keybase name or address of approve receiver")
	cmd.MarkFlagRequired(flagFrom)
	cmd.MarkFlagRequired(flagTo)

	return cmd
}

// 创建预授权命令
func CreateApproveCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-approve",
		Short: "Create approve",
		RunE: func(cmd *cobra.Command, args []string) error {
			return handleApproveOperation(cdc, createType)
		},
	}

	cmd.Flags().String(flagFrom, "", "Keybase name or address of approve creator")
	cmd.Flags().String(flagTo, "", "Keybase name or address of approve receiver")
	cmd.Flags().String(flagCoins, "", "Coins to approve. ex: 10qos,100qstars,50qsc")
	cmd.MarkFlagRequired(flagFrom)
	cmd.MarkFlagRequired(flagTo)
	cmd.MarkFlagRequired(flagCoins)

	return cmd
}

// 增加预授权命令
func IncreaseApproveCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "increase-approve",
		Short: "Increase approve",
		RunE: func(cmd *cobra.Command, args []string) error {
			return handleApproveOperation(cdc, increaseType)
		},
	}

	cmd.Flags().String(flagFrom, "", "Keybase name or address of approve creator")
	cmd.Flags().String(flagTo, "", "Keybase name or address of approve receiver")
	cmd.Flags().String(flagCoins, "", "Coins to approve. ex: 10qos,100qstars,50qsc")
	cmd.MarkFlagRequired(flagFrom)
	cmd.MarkFlagRequired(flagTo)
	cmd.MarkFlagRequired(flagCoins)

	return cmd
}

// 减少预授权命令
func DecreaseApproveCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decrease-approve",
		Short: "Decrease approve",
		RunE: func(cmd *cobra.Command, args []string) error {
			return handleApproveOperation(cdc, decreaseType)
		},
	}

	cmd.Flags().String(flagFrom, "", "Keybase name or address of approve creator")
	cmd.Flags().String(flagTo, "", "Keybase name or address of approve receiver")
	cmd.Flags().String(flagCoins, "", "Coins to approve. ex: 10qos,100qstars,50qsc")
	cmd.MarkFlagRequired(flagFrom)
	cmd.MarkFlagRequired(flagTo)
	cmd.MarkFlagRequired(flagCoins)

	return cmd
}

// 使用预授权命令
func UseApproveCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use-approve",
		Short: "Use approve",
		RunE: func(cmd *cobra.Command, args []string) error {
			return handleApproveOperation(cdc, useType)
		},
	}

	cmd.Flags().String(flagFrom, "", "Keybase name or address of approve creator")
	cmd.Flags().String(flagTo, "", "Keybase name or address of approve receiver")
	cmd.Flags().String(flagCoins, "", "Coins to approve. ex: 10qos,100qstars,50qsc")
	cmd.MarkFlagRequired(flagFrom)
	cmd.MarkFlagRequired(flagTo)
	cmd.MarkFlagRequired(flagCoins)

	return cmd
}

// 取消预授权命令
func CancelApproveCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-approve",
		Short: "Cancel approve",
		RunE: func(cmd *cobra.Command, args []string) error {
			return handleApproveOperation(cdc, cancleType)
		},
	}

	cmd.Flags().String(flagFrom, "", "Keybase name or address of approve creator")
	cmd.Flags().String(flagTo, "", "Keybase name or address of approve receiver")
	cmd.MarkFlagRequired(flagFrom)
	cmd.MarkFlagRequired(flagTo)

	return cmd
}

// 创建/增加/减少/使用/取消预授权统一处理方法
func handleApproveOperation(cdc *amino.Codec, operType operateType) error {
	iTxBuilder := func(ctx context.CLIContext) (txs.ITx, error) {
		// 解析授权账户地址
		fromAddr, err := qcliacc.GetAddrFromFlag(ctx, flagFrom)
		if err != nil {
			return nil, err
		}
		// 解析被授权账户地址
		toAddr, err := qcliacc.GetAddrFromFlag(ctx, flagTo)
		if err != nil {
			return nil, err
		}
		// 授权和被授权账户不能相同
		if fromAddr.Equals(toAddr) {
			return nil, errors.New("from and to cannot be the same")
		}

		// 取消预授权不包含币种币值信息，特殊处理
		if operType == cancleType {
			return atxs.TxCancelApprove{
				From: fromAddr,
				To:   toAddr,
			}, nil
		}

		// 解析币种币值信息
		qos, qscs, err := types.ParseCoins(viper.GetString(flagCoins))
		if err != nil {
			return nil, err
		}

		// 构造交易体
		appr := approvetypes.NewApprove(fromAddr, toAddr, qos, qscs)
		switch operType {
		case createType:
			return atxs.TxCreateApprove{Approve: appr}, nil
		case increaseType:
			return atxs.TxIncreaseApprove{Approve: appr}, nil
		case decreaseType:
			return atxs.TxDecreaseApprove{Approve: appr}, nil
		case useType:
			return atxs.TxUseApprove{Approve: appr}, nil
		default:
			return nil, errors.New("operation type invalid")
		}
	}

	return qclitx.BroadcastTxAndPrintResult(cdc, iTxBuilder)
}
