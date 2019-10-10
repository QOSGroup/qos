package client

import (
	"errors"
	"github.com/QOSGroup/qbase/client/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/mint/mapper"
	"github.com/QOSGroup/qos/module/mint/types"
	"github.com/spf13/cobra"
	go_amino "github.com/tendermint/go-amino"
)

// 查询通胀规则
func queryInflationPhrasesCommand(cdc *go_amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inflation-phrases",
		Short: "Query inflation phrases",
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			result, err := queryInflationPhrases(cliCtx)
			if err != nil {
				return err
			}

			return cliCtx.PrintResult(result)
		},
	}

	return cmd
}

func queryInflationPhrases(cliCtx context.CLIContext) ([]types.InflationPhrase, error) {
	path := mapper.BuildQueryPhrasesPath()
	res, err := cliCtx.Query(path, []byte{})

	if len(res) == 0 {
		return nil, errors.New("no result found")
	}

	var result []types.InflationPhrase
	err = cliCtx.Codec.UnmarshalJSON(res, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// 查询QOS发行总量
func queryTotalCommand(cdc *go_amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-inflation",
		Short: "Query total inflation QOS amount, both in genesis and inflation phrases",
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			result, err := queryTotal(cliCtx)
			if err != nil {
				return err
			}

			return cliCtx.PrintResult(result)
		},
	}

	return cmd
}

func queryTotal(cliCtx context.CLIContext) (btypes.BigInt, error) {
	path := mapper.BuildQueryTotalPath()
	res, err := cliCtx.Query(path, []byte{})

	if len(res) == 0 {
		return btypes.BigInt{}, errors.New("no result found")
	}

	var result btypes.BigInt
	err = cliCtx.Codec.UnmarshalJSON(res, &result)
	if err != nil {
		return btypes.BigInt{}, err
	}

	return result, nil
}

// 查询QOS流通总量
func queryAppliedCommand(cdc *go_amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-applied",
		Short: "Query total applied QOS amount, both in genesis and inflation phrases",
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			result, err := queryApplied(cliCtx)
			if err != nil {
				return err
			}

			return cliCtx.PrintResult(result)
		},
	}

	return cmd
}

func queryApplied(cliCtx context.CLIContext) (btypes.BigInt, error) {
	path := mapper.BuildQueryAppliedPath()
	res, err := cliCtx.Query(path, []byte{})

	if len(res) == 0 {
		return btypes.BigInt{}, errors.New("no result found")
	}

	var result btypes.BigInt
	err = cliCtx.Codec.UnmarshalJSON(res, &result)
	if err != nil {
		return btypes.BigInt{}, err
	}

	return result, nil
}
