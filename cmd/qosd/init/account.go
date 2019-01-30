package init

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"github.com/QOSGroup/qos/app"
	"github.com/QOSGroup/qos/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
)

func AddGenesisAccount(cdc *amino.Codec) *cobra.Command {
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

			home := viper.GetString(cli.HomeFlag)
			genFile := strings.Join([]string{home, "config", "genesis.json"}, "/")

			if !common.FileExists(genFile) {
				return fmt.Errorf("%s does not exist, run `qosd init` first", genFile)
			}

			accounts, err := types.ParseAccounts(args[0])

			genDoc, err := loadGenesisDoc(cdc, genFile)
			if err != nil {
				return err
			}

			var appState app.GenesisState
			if err = cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {
				return err
			}

			for _, v := range appState.Accounts {
				for _, acc := range accounts {
					if acc.AccountAddress.EqualsTo(v.GetAddress()) {
						return fmt.Errorf("addr: %s has already exsits", v.AccountAddress.String())
					}
				}
			}

			appState.Accounts = append(appState.Accounts, accounts...)
			for _, acc := range accounts {
				appState.MintData.AppliedQOSAmount = appState.MintData.AppliedQOSAmount + uint64(acc.QOS.Int64())
			}

			rawMessage, _ := cdc.MarshalJSON(appState)
			genDoc.AppState = rawMessage

			err = genDoc.ValidateAndComplete()
			if err != nil {
				return err
			}

			err = genDoc.SaveAs(genFile)
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().String(cli.HomeFlag, types.DefaultNodeHome, "node's home directory")

	return cmd
}
