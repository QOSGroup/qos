package init

import (
	"errors"
	"fmt"
	"github.com/QOSGroup/qbase/client/context"
	clikeys "github.com/QOSGroup/qbase/client/keys"
	"github.com/QOSGroup/qbase/keys"
	"github.com/QOSGroup/qbase/server"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/app"
	"github.com/QOSGroup/qos/module/stake/client"
	"github.com/QOSGroup/qos/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	tmtypes "github.com/tendermint/tendermint/types"
	"os"
	"path/filepath"
	"strings"
)

func GenTxCmd(ctx *server.Context, cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gentx",
		Short: "Generate a genesis tx carrying a self delegation",
		Args:  cobra.NoArgs,
		Long: `This command is an alias of the 'gaiad tx create-validator' command'.
qosd gentx --name validatorName --owner ownerName --tokens 100
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))
			nodeID, _, err := server.InitializeNodeValidatorFiles(config)
			if err != nil {
				return err
			}
			genDoc, err := tmtypes.GenesisDocFromFile(config.GenesisFile())
			if err != nil {
				return err
			}
			genesisState := app.GenesisState{}
			if err = cdc.UnmarshalJSON(genDoc.AppState, &genesisState); err != nil {
				return err
			}

			owner := viper.GetString(flagOwner)
			var info keys.Info
			if len(owner) == 0 {
				return errors.New("owner is empty")
			}
			keybase, err := clikeys.GetKeyBaseFromDir(cliCtx, viper.GetString(flagClientHome))
			if err != nil {
				return err
			}
			if strings.HasPrefix(owner, btypes.PREF_ADD) {
				addr, err := btypes.GetAddrFromBech32(owner)
				if err != nil {
					return err
				}
				info, err = keybase.GetByAddress(addr)
				if err != nil {
					return err
				}
			} else {
				info, err = keybase.Get(owner)
				if err != nil {
					return err
				}
			}

			tokens := viper.GetInt64(flagBondTokens)
			if tokens <= 0 {
				return errors.New("tokens lte zero")
			}

			validGenesisAccount(genesisState, info.GetAddress(), btypes.NewInt(tokens))

			itx, err := staking.TxCreateValidatorBuilder(cliCtx)
			if err != nil {
				return err
			}
			txStd := txs.NewTxStd(itx, genDoc.ChainID, btypes.NewInt(1000000))
			sigdata := txStd.BuildSignatureBytes(1, "")
			pass, err := clikeys.GetPassphrase(cliCtx, info.GetName())
			if err != nil {
				panic(fmt.Sprintf("Get %s Passphrase error: %s", info.GetAddress(), err.Error()))
			}
			sig, pubkey, err := keybase.Sign(info.GetName(), pass, sigdata)
			if err != nil {
				panic(err.Error())
			}
			txStd.Signature = append(txStd.Signature, txs.Signature{
				Pubkey:    pubkey,
				Signature: sig,
				Nonce:     1,
			})

			if err := writeSignedGenTx(cdc, filepath.Join(config.RootDir, "config", "gentx"), nodeID, viper.GetString(flagIP), txStd); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(flagName, "", "name for validator")
	cmd.Flags().String(flagOwner, "", "keystore name or account address")
	cmd.Flags().Int64(flagBondTokens, 0, "bond tokens amount")
	cmd.Flags().Bool(flagCompound, false, "as a self-delegator, whether the income is calculated as compound interest")
	cmd.Flags().String(flagDescription, "", "description")
	cmd.Flags().String(flagClientHome, types.DefaultCLIHome, "directory for keybase")
	cmd.Flags().String(flagNodeHome, types.DefaultNodeHome, "directory for your node")
	cmd.Flags().String(flagIP, "127.0.0.1", "ip of your node")

	cmd.MarkFlagRequired(flagName)
	cmd.MarkFlagRequired(flagOwner)
	cmd.MarkFlagRequired(flagBondTokens)

	return cmd
}

func validGenesisAccount(genesisState app.GenesisState, address btypes.Address, amount btypes.BigInt) error {
	accountIsInGenesis := false

	for _, acc := range genesisState.Accounts {
		if acc.AccountAddress.EqualsTo(address) {

			if !acc.EnoughOfQOS(amount) {
				return fmt.Errorf(
					"account %v is in genesis, but it only has %v QOS available to stake, not %v",
					address.String(), acc.QOS, amount,
				)
			}
			accountIsInGenesis = true
			break
		}
	}

	if accountIsInGenesis {
		return nil
	}

	return fmt.Errorf("account %s in not in the app_state.accounts array of genesis.json", address)
}

func writeSignedGenTx(cdc *amino.Codec, genTxDir, nodeID, ip string, tx *txs.TxStd) error {
	if err := common.EnsureDir(genTxDir, 0700); err != nil {
		return err
	}
	genTx := filepath.Join(genTxDir, fmt.Sprintf("%s@%s.json", nodeID, ip))
	outputFile, err := os.OpenFile(genTx, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	json, err := cdc.MarshalJSON(tx)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(outputFile, "%s\n", json)
	return err
}
