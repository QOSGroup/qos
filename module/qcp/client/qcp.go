package client

import (
	"errors"
	"github.com/QOSGroup/kepler/cert"
	qcliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/context"
	qclitx "github.com/QOSGroup/qbase/client/tx"
	"github.com/QOSGroup/qbase/txs"
	qtxs "github.com/QOSGroup/qos/module/qcp/txs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/common"
)

const (
	flagCreator = "creator"
	flagPathqcp = "qcp.crt"
)

func InitQCPCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init-qcp",
		Short: "init qcp",
		RunE: func(cmd *cobra.Command, args []string) error {
			return qclitx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (txs.ITx, error) {
				pathqcp := viper.GetString(flagPathqcp)

				creatorAddr, err := qcliacc.GetAddrFromFlag(ctx, flagCreator)
				if err != nil {
					return nil, err
				}

				var crt = cert.Certificate{}
				err = cdc.UnmarshalJSON(common.MustReadFile(pathqcp), &crt)
				if err != nil {
					return nil, err
				}

				_, ok := crt.CSR.Subj.(cert.QCPSubject)
				if !ok {
					return nil, errors.New("invalid crt file")
				}

				return qtxs.TxInitQCP{creatorAddr, &crt}, nil
			})
		},
	}

	cmd.Flags().String(flagCreator, "", "address or name of creator")
	cmd.Flags().String(flagPathqcp, "", "path of CA(QCP)")
	cmd.MarkFlagRequired(flagCreator)
	cmd.MarkFlagRequired(flagPathqcp)

	return cmd
}
