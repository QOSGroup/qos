package qsc

import (
	"fmt"
	bacc "github.com/QOSGroup/qbase/account"
	cliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/context"
	"github.com/QOSGroup/qbase/client/keys"
	btx "github.com/QOSGroup/qbase/client/tx"
	btxs "github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/txs/qsc"
	"github.com/QOSGroup/qos/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/common"
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

			var caQsc qsc.Certificate
			err = cdc.UnmarshalBinaryBare(common.MustReadFile(pathqsc), &caQsc)
			if err != nil {
				return err
			}

			var caBanker qsc.Certificate
			err = cdc.UnmarshalBinaryBare(common.MustReadFile(pathbank), &caBanker)
			if err != nil {
				return err
			}

			chainId, err := types.GetDefaultChainId()
			if err != nil {
				return err
			}

			var acs []account.QOSAccount
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
					acc := account.QOSAccount{
						BaseAccount: bacc.BaseAccount{
							info.GetAddress(),
							nil,
							0,
						},
						QOS: btypes.ZeroInt(),
						QSCs: types.QSCs{
							{
								caQsc.CSR.Subj.CN,
								btypes.NewInt(amount),
							},
						},
					}
					acs = append(acs, acc)
				}
			}

			tx := btxs.NewTxStd(
				qsc.TxCreateQSC{
					ChainID:     chainId,
					Creator:     creator.GetAddress(),
					Extrate:     extrate,
					QSCCA:       caQsc,
					BankerCA:    caBanker,
					Description: description,
					Accounts:    acs,
				},
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
		Use:   "query [qsc]",
		Short: "query qsc info by name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			result, err := cliCtx.Client.ABCIQuery("store/qsc/key", qsc.BuildQSCKey(args[0]))
			if err != nil {
				return err
			}

			var info qsc.QSCInfo
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
			qscName := viper.GetString(flagQscname)

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
				qsc.TxIssueQSC{
					qscName,
					btypes.NewInt(amount),
					banker.GetAddress(),
				},
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
