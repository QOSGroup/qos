package init

import (
	"fmt"
	"github.com/QOSGroup/qbase/server"
	"github.com/QOSGroup/qos/module/bank"
	"github.com/QOSGroup/qos/module/mint"
	"github.com/spf13/viper"

	"github.com/QOSGroup/qos/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"
	tmtypes "github.com/tendermint/tendermint/types"
)

func AddGenesisAccount(ctx *server.Context, cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-accounts [accounts]",
		Short: "Add genesis accounts to genesis.json",
		Long: `add-genesis-accounts [accounts] will add [accounts] into app_state.
Multiple accounts separated by ';'.

Example:

	qosd add-genesis-accounts "address1lly0audg7yem8jt77x2jc6wtrh7v96hgve8fh8,1000000qos;address1auhqphrnk74jx2c5n80m9pdgl0ln79tyz32xlc,100000qos"
	`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			accounts, err := types.ParseAccounts(args[0], viper.GetString(flagClientHome))

			genDoc, err := tmtypes.GenesisDocFromFile(config.GenesisFile())
			if err != nil {
				return err
			}

			var appState types.GenesisState
			if err = cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {
				return err
			}

			var bankState bank.GenesisState
			cdc.MustUnmarshalJSON(appState[bank.ModuleName], &bankState)

			var mintState mint.GenesisState
			cdc.MustUnmarshalJSON(appState[mint.ModuleName], &mintState)

			for _, v := range bankState.Accounts {
				for _, acc := range accounts {
					if acc.AccountAddress.EqualsTo(v.GetAddress()) {
						return fmt.Errorf("addr: %s has already exists", v.AccountAddress.String())
					}
				}
			}

			bankState.Accounts = append(bankState.Accounts, accounts...)
			for _, acc := range accounts {
				mintState.AppliedQOSAmount = mintState.AppliedQOSAmount + uint64(acc.QOS.Int64())
			}

			appState[bank.ModuleName] = cdc.MustMarshalJSON(bankState)
			appState[mint.ModuleName] = cdc.MustMarshalJSON(mintState)

			rawMessage, _ := cdc.MarshalJSON(appState)
			genDoc.AppState = rawMessage

			err = genDoc.ValidateAndComplete()
			if err != nil {
				return err
			}

			err = genDoc.SaveAs(config.GenesisFile())
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().String(cli.HomeFlag, types.DefaultNodeHome, "directory for node's data and config files")
	cmd.Flags().String(flagClientHome, types.DefaultCLIHome, "directory for keybase")

	return cmd
}
