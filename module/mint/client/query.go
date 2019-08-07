package client

import (
	"errors"
	"github.com/QOSGroup/qbase/client/context"
	"github.com/QOSGroup/qos/module/mint/mapper"
	"github.com/QOSGroup/qos/module/mint/types"
	"github.com/spf13/cobra"
	go_amino "github.com/tendermint/go-amino"
)

func queryInflationPhrases(cdc *go_amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inflation-phrases",
		Short: "Query inflation phrases",
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			path := mapper.BuildQueryPhrasesPath()
			res, err := cliCtx.Query(path, []byte{})

			if len(res) == 0 {
				return errors.New("no result found")
			}

			var result []types.InflationPhrase
			err = cdc.UnmarshalJSON(res, &result)
			if err != nil {
				return err
			}

			return cliCtx.PrintResult(result)
		},
	}

	return cmd
}
