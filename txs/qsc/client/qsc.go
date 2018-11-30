package qsc

import (
	"fmt"
	bacc "github.com/QOSGroup/qbase/account"
	qcliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/context"
	"github.com/QOSGroup/qbase/client/keys"
	qclitx "github.com/QOSGroup/qbase/client/tx"
	"github.com/QOSGroup/qbase/txs"
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
	flagQscChainID  = "qsc-chain"
	flagQscname     = "qsc-name"
	flagCreator     = "creator"
	flagBanker      = "banker"
	flagExtrate     = "extrate"
	flagPathqsc     = "qsc.crt"
	flagPathbank    = "banker.crt"
	flagAccounts    = "accounts"
	flagAmount      = "amount"
	flagDescription = "desc"
)

func CreateQSCCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-qsc",
		Short: "create qsc",
		RunE: func(cmd *cobra.Command, args []string) error {
			return qclitx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (txs.ITx, error) {
				//flag args
				qscChainID := viper.GetString(flagQscChainID)
				extrate := viper.GetString(flagExtrate)
				pathqsc := viper.GetString(flagPathqsc)
				pathbank := viper.GetString(flagPathbank)
				accountStr := viper.GetString(flagAccounts)
				description := viper.GetString(flagDescription)

				creatorAddr, err := qcliacc.GetAddrFromFlag(ctx, flagCreator)
				if err != nil {
					return nil, err
				}

				var caQsc *qsc.Certificate
				err = cdc.UnmarshalBinaryBare(common.MustReadFile(pathqsc), caQsc)
				if err != nil {
					return nil, err
				}

				var caBanker *qsc.Certificate
				if pathbank != "" {
					err = cdc.UnmarshalBinaryBare(common.MustReadFile(pathbank), &caBanker)
				}
				if err != nil {
					return nil, err
				}

				var acs []account.QOSAccount
				if len(accountStr) > 0 {
					accArrs := strings.Split(accountStr, ";")
					for _, accArrStr := range accArrs {
						accArr := strings.Split(accArrStr, ",")
						info, err := keys.GetKeyInfo(ctx, accArr[0])
						if err != nil {
							return nil, err
						}
						amount, err := strconv.ParseInt(strings.TrimSpace(accArr[1]), 10, 64)
						if err != nil {
							return nil, err
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

				return qsc.TxCreateQSC{
					ChainID:     qscChainID,
					Creator:     creatorAddr,
					Extrate:     extrate,
					QSCCA:       caQsc,
					BankerCA:    caBanker,
					Description: description,
					Accounts:    acs,
				}, nil

			})
		},
	}

	cmd.Flags().String(flagQscChainID, "", "chainID for the qsc corresponding to")
	cmd.Flags().String(flagCreator, "", "name of banker")

	cmd.Flags().String(flagExtrate, "1:280.0000", "extrate: qos:qscxxx")
	cmd.Flags().String(flagPathqsc, "", "path of CA(qsc)")
	cmd.Flags().String(flagPathbank, "", "path of CA(banker)")
	cmd.Flags().String(flagAccounts, "", "init accounts: Sansa,100;Lisa,100")
	cmd.Flags().String(flagDescription, "", "description")

	cmd.MarkFlagRequired(flagQscChainID)
	cmd.MarkFlagRequired(flagCreator)
	cmd.MarkFlagRequired(flagPathqsc)
	//cmd.MarkFlagRequired(flagPathbank)

	return cmd
}

func QueryQscCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "qsc [qsc]",
		Short: "query qsc info by name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			result, err := cliCtx.Client.ABCIQuery("store/qsc/key", qsc.BuildQSCKey(args[0]))
			if err != nil {
				return err
			}

			if len(result.Response.GetValue()) == 0 {
				return fmt.Errorf("%s not exists.", args[0])
			}

			var info qsc.QSCInfo
			err = cdc.UnmarshalBinaryBare(result.Response.GetValue(), &info)
			if err != nil {
				return err
			}

			return cliCtx.PrintResult(info)
		},
	}

	return cmd
}

func IssueQSCCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-qsc",
		Short: "issue qsc",
		RunE: func(cmd *cobra.Command, args []string) error {
			return qclitx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (txs.ITx, error) {
				amount := viper.GetInt64(flagAmount)
				qscName := viper.GetString(flagQscname)
				bankerAddr, err := qcliacc.GetAddrFromFlag(ctx, flagBanker)
				if err != nil {
					return nil, err
				}
				return qsc.TxIssueQSC{qscName, btypes.NewInt(amount), bankerAddr}, nil
			})
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
