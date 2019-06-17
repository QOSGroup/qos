package mint

import (
	"errors"
	"github.com/QOSGroup/qbase/client/context"
	mtypes "github.com/QOSGroup/qos/module/eco/types"
	"github.com/QOSGroup/qos/module/mint"
	"github.com/spf13/cobra"
	go_amino "github.com/tendermint/go-amino"
)

func queryInflationPhrases(cdc *go_amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inflation-phrases",
		Short: "Query inflation phrases",
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			path := mint.BuildQueryProposalPath()
			res, err := cliCtx.Query(path, []byte{})

			if len(res) == 0 {
				return errors.New("no result found")
			}

			var result []mtypes.InflationPhrase
			err = cdc.UnmarshalJSON(res, &result)
			if err != nil {
				return err
			}

			return cliCtx.PrintResult(result)
		},
	}

	return cmd
}
