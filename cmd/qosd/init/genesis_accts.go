package init

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	qbtypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/app"
	"github.com/QOSGroup/qos/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
)

const (
	flagAddr  = "addr"
	flagCoins = "coins"
)

func AddGenesisAccount(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-account",
		Short: "Add genesis account to genesis.json",
		RunE: func(_ *cobra.Command, args []string) error {

			home := viper.GetString(cli.HomeFlag)
			genFile := strings.Join([]string{home, "config", "genesis.json"}, "/")

			if !common.FileExists(genFile) {
				return fmt.Errorf("%s does not exist, run `qosd init` first", genFile)
			}

			addr, err := qbtypes.GetAddrFromBech32(viper.GetString(flagAddr))
			if err != nil {
				return err
			}
			qos, qscs, err := types.ParseCoins(viper.GetString(flagCoins))
			if err != nil {
				return err
			}

			genDoc, err := loadGenesisDoc(cdc, genFile)
			if err != nil {
				return err
			}

			var appState app.GenesisState
			if err = cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {
				return err
			}

			for _, v := range appState.Accounts {
				if qbtypes.Address(addr).EqualsTo(v.GetAddress()) {
					return fmt.Errorf("addr: %s has already exsits", viper.GetString(flagAddr))
				}
			}

			internalQOSAccount := account.NewQOSAccount()
			internalQOSAccount.SetAddress(addr)
			internalQOSAccount.SetQOS(qos)

			for _, qsc := range qscs {
				internalQOSAccount.SetQSC(qsc)
			}
			appState.Accounts = append(appState.Accounts, internalQOSAccount)

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

	cmd.Flags().String(flagAddr, "", "default account address")
	cmd.Flags().String(flagCoins, "", "default account's coins")
	cmd.Flags().String(cli.HomeFlag, types.DefaultNodeHome, "node's home directory")

	cmd.MarkFlagRequired(flagAddr)

	return cmd
}
