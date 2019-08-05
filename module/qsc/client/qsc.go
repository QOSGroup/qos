package client

import (
	"fmt"
	"github.com/QOSGroup/kepler/cert"
	bacc "github.com/QOSGroup/qbase/account"
	qcliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/context"
	"github.com/QOSGroup/qbase/client/keys"
	qclitx "github.com/QOSGroup/qbase/client/tx"
	btxs "github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/qsc/mapper"
	"github.com/QOSGroup/qos/module/qsc/txs"
	"github.com/QOSGroup/qos/module/qsc/types"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/pkg/errors"
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
	flagPathqsc     = "qsc.crt"
	flagAccounts    = "accounts"
	flagAmount      = "amount"
	flagDescription = "desc"
)

func CreateQSCCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-qsc",
		Short: "create qsc",
		RunE: func(cmd *cobra.Command, args []string) error {
			return qclitx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (btxs.ITx, error) {
				//flag args
				extrate := viper.GetString(flagExtrate)
				pathqsc := viper.GetString(flagPathqsc)
				accountStr := viper.GetString(flagAccounts)
				description := viper.GetString(flagDescription)

				creatorAddr, err := qcliacc.GetAddrFromFlag(ctx, flagCreator)
				if err != nil {
					return nil, err
				}

				var crt cert.Certificate
				err = cdc.UnmarshalJSON(common.MustReadFile(pathqsc), &crt)
				if err != nil {
					return nil, err
				}

				subj, ok := crt.CSR.Subj.(cert.QSCSubject)
				if !ok {
					return nil, errors.New("invalid crt file")
				}

				var acs []*qtypes.QOSAccount
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
						acc := qtypes.QOSAccount{
							BaseAccount: bacc.BaseAccount{
								info.GetAddress(),
								nil,
								0,
							},
							QOS: btypes.ZeroInt(),
							QSCs: qtypes.QSCs{
								{
									subj.Name,
									btypes.NewInt(amount),
								},
							},
						}
						acs = append(acs, &acc)
					}
				}

				return txs.TxCreateQSC{
					creatorAddr,
					extrate,
					&crt,
					description,
					acs,
				}, nil

			})
		},
	}

	cmd.Flags().String(flagCreator, "", "name or address of creator")
	cmd.Flags().String(flagExtrate, "1", "extrate: qos:qscxxx")
	cmd.Flags().String(flagPathqsc, "", "path of CA(qsc)")
	cmd.Flags().String(flagDescription, "", "description")
	cmd.Flags().String(flagAccounts, "", "init accounts, eg: address1,100;address2,100")
	cmd.MarkFlagRequired(flagCreator)
	cmd.MarkFlagRequired(flagPathqsc)

	return cmd
}

func QueryQscCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "qsc [qsc]",
		Short: "query qsc info by name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			result, err := cliCtx.Client.ABCIQuery("store/qsc/key", mapper.BuildQSCKey(args[0]))
			if err != nil {
				return err
			}

			if len(result.Response.GetValue()) == 0 {
				return fmt.Errorf("%s not exists.", args[0])
			}

			var info types.Info
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
			return qclitx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (btxs.ITx, error) {
				amount := viper.GetInt64(flagAmount)
				qscName := viper.GetString(flagQscname)
				bankerAddr, err := qcliacc.GetAddrFromFlag(ctx, flagBanker)
				if err != nil {
					return nil, err
				}
				return txs.TxIssueQSC{qscName, btypes.NewInt(amount), bankerAddr}, nil
			})
		},
	}

	cmd.Flags().Int64(flagAmount, 100000, "coin amount send to banker")
	cmd.Flags().String(flagQscname, "", "qsc name")
	cmd.Flags().String(flagBanker, "", "address or name of banker")
	cmd.MarkFlagRequired(flagAmount)
	cmd.MarkFlagRequired(flagQscname)
	cmd.MarkFlagRequired(flagBanker)

	return cmd
}
