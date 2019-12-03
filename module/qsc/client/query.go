package client

import (
	"errors"
	"github.com/QOSGroup/qbase/client/context"
	"github.com/QOSGroup/qos/module/qsc/mapper"
	"github.com/QOSGroup/qos/module/qsc/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
	"strings"
)

func QueryQSCCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "qsc [qsc]",
		Short: "query qsc info by name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			info, err := queryQscInfo(args[0], cliCtx)
			if err != nil {
				return err
			}
			return cliCtx.PrintResult(info)
		},
	}

	return cmd
}

func queryQscInfo(qsc string, cliCtx context.CLIContext) (types.QSCInfo, error) {
	path := mapper.BuildQueryQSCPath(strings.TrimSpace(qsc))
	res, err := cliCtx.Query(path, []byte{})
	if err != nil {
		return types.QSCInfo{}, err
	}
	if len(res) == 0 {
		return types.QSCInfo{}, context.RecordsNotFoundError
	}

	var info types.QSCInfo
	err = cliCtx.Codec.UnmarshalJSON(res, &info)

	return info, err
}

func QueryQSCsCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "qscs",
		Short: "query qscs list",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			result, err := queryAllQscs(cliCtx)
			if err != nil {
				return err
			}
			return cliCtx.PrintResult(result)
		},
	}

	return cmd
}

func queryAllQscs(cliCtx context.CLIContext) ([]types.QSCInfo, error) {
	path := mapper.BuildQueryQSCsPath()
	res, err := cliCtx.Query(path, []byte{})
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("no result found")
	}

	var infos []types.QSCInfo
	err = cliCtx.Codec.UnmarshalJSON(res, &infos)

	return infos, err
}
