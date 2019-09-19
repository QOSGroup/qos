package client

import (
	"errors"
	"github.com/QOSGroup/qbase/client/context"
	"github.com/QOSGroup/qos/module/qsc/mapper"
	"github.com/QOSGroup/qos/module/qsc/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
)

func QueryQSCCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "qsc [qsc]",
		Short: "query qsc info by name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			path := mapper.BuildQueryQSCPath(args[0])
			res, err := cliCtx.Query(path, []byte{})
			if err != nil {
				return nil
			}
			if len(res) == 0 {
				return errors.New("no result found")
			}

			var info types.QSCInfo
			cdc.UnmarshalJSON(res, &info)

			return cliCtx.PrintResult(info)
		},
	}

	return cmd
}

func QueryQSCsCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "qscs",
		Short: "query qscs list",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			path := mapper.BuildQueryQSCsPath()
			res, err := cliCtx.Query(path, []byte{})
			if err != nil {
				return nil
			}
			if len(res) == 0 {
				return errors.New("no result found")
			}

			var infos []types.QSCInfo
			cdc.UnmarshalJSON(res, &infos)

			return cliCtx.PrintResult(infos)
		},
	}

	return cmd
}
