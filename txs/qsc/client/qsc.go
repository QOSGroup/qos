package qsc

import (
	"fmt"
	cliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/context"
	"github.com/QOSGroup/qbase/client/keys"
	btx "github.com/QOSGroup/qbase/client/tx"
	btxs "github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/mapper"
	"github.com/QOSGroup/qos/txs/qsc"
	"github.com/QOSGroup/qos/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"strconv"
	"strings"
)

const (
	flagQscname     = "qsc-name"
	flagCreator     = "creator"
	flagBanker      = "banker"
	flagExtrate     = "extrate"
	flagPathqsc     = "path-qsc"
	flagPathbank    = "path-bank"
	flagAccounts    = "accounts"
	flagAmount      = "amount"
	flagDescription = "desc"
)

func CreateQSCCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-qsc",
		Short: "create qsc",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			extrate := viper.GetString(flagExtrate)
			pathqsc := viper.GetString(flagPathqsc)
			pathbank := viper.GetString(flagPathbank)
			accountStr := viper.GetString(flagAccounts)
			description := viper.GetString(flagDescription)

			creatorName := viper.GetString(flagCreator)
			creatorInfo, err := keys.GetKeyInfo(cliCtx, creatorName)
			if err != nil {
				return err
			}
			creator, err := cliacc.GetAccount(cliCtx, creatorInfo.GetAddress())
			if err != nil {
				return err
			}

			caQsc := qsc.FetchCA(pathqsc)
			caBanker := qsc.FetchCA(pathbank)

			chainId, err := types.GetDefaultChainId()
			if err != nil {
				return err
			}

			var acs []qsc.AddrCoin
			if len(accountStr) > 0 {
				accArrs := strings.Split(accountStr, ";")
				for _, accArrStr := range accArrs {
					accArr := strings.Split(accArrStr, ",")
					info, err := keys.GetKeyInfo(cliCtx, accArr[0])
					if err != nil {
						return err
					}
					amount, err := strconv.ParseInt(strings.TrimSpace(accArr[1]), 10, 64)
					if err != nil {
						return err
					}
					acs = append(acs, qsc.AddrCoin{info.GetAddress(), btypes.NewInt(amount)})
				}
			}

			tx := btxs.NewTxStd(
				qsc.NewCreateQsc(cdc, caQsc, caBanker, chainId, creator.GetAddress(), &acs, extrate, description),
				chainId,
				btypes.ZeroInt())
			tx, err = btx.SignStdTx(cliCtx, creatorName, creator.GetNonce()+1, tx)

			result, err := cliCtx.BroadcastTx(cdc.MustMarshalBinaryBare(tx))
			if err != nil {
				return err
			}

			msg, _ := cdc.MarshalJSON(result)
			fmt.Println(string(msg))

			return err
		},
	}

	cmd.Flags().String(flagCreator, "", "name of banker")
	cmd.Flags().String(flagExtrate, "1:280.0000", "extrate: qos:qscxxx")
	cmd.Flags().String(flagPathqsc, "", "path of CA(qsc)")
	cmd.Flags().String(flagPathbank, "", "path of CA(banker)")
	cmd.Flags().String(flagAccounts, "", "init accounts: Sansa,100;Lisa,100")
	cmd.Flags().String(flagDescription, "", "description")
	cmd.MarkFlagRequired(flagCreator)
	cmd.MarkFlagRequired(flagPathqsc)
	cmd.MarkFlagRequired(flagPathbank)

	return cmd
}

func QueryQscCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "qsc [qsc-name]",
		Short: "query qsc info by name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			key := fmt.Sprintf("qsc/[%s]", args[0])
			result, err := cliCtx.Client.ABCIQuery("store/base/key", []byte(key))
			if err != nil {
				return err
			}

			var info mapper.QscInfo
			err = cdc.UnmarshalBinaryBare(result.Response.GetValue(), &info)
			if err != nil {
				return err
			}
			msg, _ := cdc.MarshalJSON(info)
			fmt.Println(string(msg))

			return nil
		},
	}

	return cmd
}

func IssueQSCCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-qsc",
		Short: "issue qsc",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			amount := viper.GetInt64(flagAmount)
			qscname := viper.GetString(flagQscname)

			chainId, err := types.GetDefaultChainId()
			if err != nil {
				return err
			}

			bankerName := viper.GetString(flagBanker)
			bankerInfo, err := keys.GetKeyInfo(cliCtx, bankerName)
			if err != nil {
				return err
			}
			banker, err := cliacc.GetAccount(cliCtx, bankerInfo.GetAddress())
			if err != nil {
				return err
			}

			tx := btxs.NewTxStd(
				qsc.NewTxIssueQsc(qscname, btypes.NewInt(amount), banker.GetAddress()),
				chainId,
				btypes.ZeroInt())
			tx, err = btx.SignStdTx(cliCtx, bankerName, banker.GetNonce()+1, tx)

			result, err := cliCtx.BroadcastTx(cdc.MustMarshalBinaryBare(tx))
			if err != nil {
				return err
			}

			msg, _ := cdc.MarshalJSON(result)
			fmt.Println(string(msg))

			return err
		},
	}

	cmd.Flags().Int64(flagAmount, 100000, "coin amount send to banker")
	cmd.Flags().String(flagQscname, "", "qsc name")
	cmd.Flags().String(flagBanker, "", "name of banker")
	cmd.MarkFlagRequired(flagAmount)
	cmd.MarkFlagRequired(flagQscname)
	cmd.MarkFlagRequired(flagBanker)

	return cmd
}
