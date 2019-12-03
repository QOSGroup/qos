package client

import (
	"fmt"
	"github.com/QOSGroup/kepler/cert"
	bacc "github.com/QOSGroup/qbase/account"
	qcliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/context"
	qclitx "github.com/QOSGroup/qbase/client/tx"
	btxs "github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/qsc/txs"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/common"
	"strings"
)

const (
	flagQscname      = "qsc-name"
	flagCreator      = "creator"
	flagBanker       = "banker"
	flagExchangeRate = "exchange-rate"
	flagQscCrtFile   = "qsc.crt"
	flagAccounts     = "accounts"
	flagAmount       = "amount"
	flagDescription  = "desc"
)

func CreateQSCCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-qsc",
		Short: "create qsc",
		RunE: func(cmd *cobra.Command, args []string) error {
			return qclitx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (btxs.ITx, error) {
				//flag args
				exchangeRate := viper.GetString(flagExchangeRate)
				qscCrtFile := viper.GetString(flagQscCrtFile)
				accountStr := viper.GetString(flagAccounts)
				description := viper.GetString(flagDescription)

				creatorAddr, err := qcliacc.GetAddrFromFlag(ctx, flagCreator)
				if err != nil {
					return nil, err
				}

				var crt cert.Certificate
				err = cdc.UnmarshalJSON(common.MustReadFile(qscCrtFile), &crt)
				if err != nil {
					return nil, err
				}

				subj, ok := crt.CSR.Subj.(cert.QSCSubject)
				if !ok {
					return nil, errors.New("invalid crt file")
				}

				acs, err := parseAccountStr(accountStr, subj.Name, func(addrStr string) (addr btypes.AccAddress, e error) {
					return qcliacc.GetAddrFromValue(ctx, addrStr)
				})

				tx := txs.TxCreateQSC{
					creatorAddr,
					exchangeRate,
					&crt,
					description,
					acs,
				}
				if err = tx.ValidateInputs(); err != nil {
					return nil, err
				}
				return tx, nil
			})
		},
	}

	cmd.Flags().String(flagCreator, "", "name or address of creator")
	cmd.Flags().String(flagExchangeRate, "1", "extrate: qos:qscxxx")
	cmd.Flags().String(flagQscCrtFile, "", "path of CA(qsc)")
	cmd.Flags().String(flagDescription, "", "description")
	cmd.Flags().String(flagAccounts, "", "init accounts, eg: address1,100;address2,100")
	cmd.MarkFlagRequired(flagCreator)
	cmd.MarkFlagRequired(flagQscCrtFile)

	return cmd
}

func parseAccountStr(accountsStr, qscName string, fn func(string) (btypes.AccAddress, error)) ([]*qtypes.QOSAccount, error) {
	var acs []*qtypes.QOSAccount
	if len(accountsStr) > 0 {
		accArrs := strings.Split(accountsStr, ";")
		for _, accArrStr := range accArrs {
			accArr := strings.Split(accArrStr, ",")
			accountAddr, err := fn(strings.TrimSpace(accArr[0]))
			if err != nil {
				return nil, err
			}
			amount, ok := btypes.NewIntFromString(strings.TrimSpace(accArr[1]))
			if !ok {
				return nil, fmt.Errorf("%s parse error", accArr[1])
			}
			acc := qtypes.QOSAccount{
				BaseAccount: bacc.BaseAccount{
					accountAddr,
					nil,
					0,
				},
				QOS: btypes.ZeroInt(),
				QSCs: qtypes.QSCs{
					{
						qscName,
						amount,
					},
				},
			}
			acs = append(acs, &acc)
		}
	}

	return acs, nil
}

func IssueQSCCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-qsc",
		Short: "issue qsc",
		RunE: func(cmd *cobra.Command, args []string) error {
			return qclitx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (btxs.ITx, error) {
				amount, err := qtypes.GetIntFromFlag(flagAmount, false)
				if err != nil {
					return nil, err
				}
				qscName := viper.GetString(flagQscname)
				bankerAddr, err := qcliacc.GetAddrFromFlag(ctx, flagBanker)
				if err != nil {
					return nil, err
				}
				return txs.TxIssueQSC{qscName, amount, bankerAddr}, nil
			})
		},
	}

	cmd.Flags().String(flagAmount, "0", "coin amount send to banker")
	cmd.Flags().String(flagQscname, "", "qsc name")
	cmd.Flags().String(flagBanker, "", "address or name of banker")
	cmd.MarkFlagRequired(flagAmount)
	cmd.MarkFlagRequired(flagQscname)
	cmd.MarkFlagRequired(flagBanker)

	return cmd
}
