package init

import (
	"fmt"
	"github.com/QOSGroup/qbase/server"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/app"
	"github.com/QOSGroup/qos/module/guardian"
	"github.com/QOSGroup/qos/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"
	tmtypes "github.com/tendermint/tendermint/types"
)

func AddGuardian(ctx *server.Context, cdc *amino.Codec) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "add-guardian",
		Short: "Add guardian to genesis",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			address := viper.GetString(flagAddress)
			addr, err := btypes.GetAddrFromBech32(address)
			if err != nil {
				return err
			}

			description := viper.GetString(flagDescription)

			genDoc, err := tmtypes.GenesisDocFromFile(config.GenesisFile())
			if err != nil {
				return err
			}

			var appState app.GenesisState
			if err = cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {
				return err
			}

			for _, v := range appState.GuardianData.Guardians {
				if v.Address.EqualsTo(addr) {
					return fmt.Errorf("guardian: %s has already exists", v.Address.String())
				}
			}

			guardian := guardian.NewGuardian(description, guardian.Genesis, addr, nil)
			appState.GuardianData.Guardians = append(appState.GuardianData.Guardians, *guardian)

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

	cmd.Flags().String(cli.HomeFlag, types.DefaultNodeHome, "node's home directory")
	cmd.Flags().String(flagAddress, "", "address of guardian")
	cmd.Flags().String(flagDescription, "", "description")
	cmd.MarkFlagRequired(flagAddress)

	return cmd
}
