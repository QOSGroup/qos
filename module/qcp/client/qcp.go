package client

import (
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
	flagCreator    = "creator"
	flagQcpCrtFile = "qcp.crt"
)

func InitQCPCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init-qcp",
		Short: "init qcp",
		RunE: func(cmd *cobra.Command, args []string) error {
			return qclitx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (txs.ITx, error) {
				crtFile := viper.GetString(flagQcpCrtFile)

				creatorAddr, err := qcliacc.GetAddrFromFlag(ctx, flagCreator)
				if err != nil {
					return nil, err
				}

				var crt = cert.Certificate{}
				err = cdc.UnmarshalJSON(common.MustReadFile(crtFile), &crt)
				if err != nil {
					return nil, err
				}

				tx := qtxs.TxInitQCP{creatorAddr, &crt}
				if err = tx.ValidateInputs(); err != nil {
					return nil, err
				}
				return tx, nil
			})
		},
	}

	cmd.Flags().String(flagCreator, "", "address or name of creator")
	cmd.Flags().String(flagQcpCrtFile, "", "path of CA(QCP)")
	cmd.MarkFlagRequired(flagCreator)
	cmd.MarkFlagRequired(flagQcpCrtFile)

	return cmd
}
