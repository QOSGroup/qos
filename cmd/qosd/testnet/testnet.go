package testnet

import (
	"fmt"
	"github.com/QOSGroup/qbase/server"
	"github.com/QOSGroup/qos/app"
	qosinit "github.com/QOSGroup/qos/cmd/qosd/init"
	"github.com/QOSGroup/qos/module/distribution"
	staketypes "github.com/QOSGroup/qos/module/eco/types"
	"github.com/QOSGroup/qos/module/gov"
	"github.com/QOSGroup/qos/module/mint"
	"github.com/QOSGroup/qos/module/qcp"
	"github.com/QOSGroup/qos/module/qsc"
	"github.com/QOSGroup/qos/module/stake"
	"github.com/QOSGroup/qos/types"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	btypes "github.com/QOSGroup/qbase/types"
	cfg "github.com/tendermint/tendermint/config"
	cmn "github.com/tendermint/tendermint/libs/common"
	ttypes "github.com/tendermint/tendermint/types"
)

var (
	chainId  string
	compound bool

	nValidators   int
	outputDir     string
	nodeDirPrefix string

	populatePersistentPeers bool
	hostnamePrefix          string
	startingIPAddress       string

	qcpRootCA string
	qscRootCA string
	accounts  string
)

const (
	nodeDirPerm  = 0755
	nodeFilePerm = 0644

	validatorBondTokens   = 1000
	validatorOwnerInitQOS = 1000000
	validatorOperatorFile = "priv_validator_owner.json"
)

func TestnetFileCmd(ctx *server.Context, cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "testnet",
		Short: "Initialize files for a QOS testnet",
		Long: `testnet will create "v" + "n" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).

Note, strict routability for addresses is turned off in the config file.

Example:

	qosd testnet --chain-id=qostest --v=4 --o=./output --starting-ip-address=192.168.1.2 --genesis-accounts=address16lwp3kykkjdc2gdknpjy6u9uhfpa9q4vj78ytd,1000000qos,1000000qstars
	`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config

			var nodeDirs []string
			var persistentPeers []string

			// accounts
			genesisAccounts := make([]*types.QOSAccount, 0)
			var err error
			if accounts != "" {
				genesisAccounts, err = types.ParseAccounts(accounts)
				if err != nil {
					return err
				}
			}

			// root ca
			var qcpPubKey crypto.PubKey
			if qcpRootCA != "" {
				err := cdc.UnmarshalJSON(cmn.MustReadFile(qcpRootCA), &qcpPubKey)
				if err != nil {
					return err
				}
			}
			var qscPubKey crypto.PubKey
			if qscRootCA != "" {
				err := cdc.UnmarshalJSON(cmn.MustReadFile(qscRootCA), &qscPubKey)
				if err != nil {
					return err
				}
			}

			appState := app.GenesisState{
				Accounts:         genesisAccounts,
				MintData:         mint.DefaultGenesisState(),
				StakeData:        stake.NewGenesisState(staketypes.DefaultStakeParams(), nil, nil, nil, nil, nil, nil),
				QCPData:          qcp.NewGenesisState(qcpPubKey, nil),
				QSCData:          qsc.NewGenesisState(qscPubKey, nil),
				DistributionData: distribution.DefaultGenesisState(),
				GovData:          gov.DefaultGenesisState(),
			}

			// validators
			genVals := make([]staketypes.Validator, nValidators)
			for i := 0; i < nValidators; i++ {
				nodeDirName := fmt.Sprintf("%s%d", nodeDirPrefix, i)
				nodeDir := filepath.Join(outputDir, nodeDirName)
				nodeDirs = append(nodeDirs, nodeDir)
				config.SetRoot(nodeDir)
				config.Moniker = nodeDirName

				err := os.MkdirAll(filepath.Join(nodeDir, "config"), nodeDirPerm)
				if err != nil {
					_ = os.RemoveAll(outputDir)
					return err
				}

				err = os.MkdirAll(filepath.Join(nodeDir, "data"), nodeDirPerm)
				if err != nil {
					_ = os.RemoveAll(outputDir)
					return err
				}

				nodeId, valPubKey, err := server.InitializeNodeValidatorFiles(config)
				if err != nil {
					_ = os.RemoveAll(outputDir)
					return err
				}
				persistentPeers = append(persistentPeers, fmt.Sprintf("%s@%s:%d", nodeId, hostnameOrIP(i), 26656))

				// create validator
				owner := ed25519.GenPrivKey()
				genVals[i] = staketypes.Validator{
					Name:            nodeDirName,
					ValidatorPubKey: valPubKey,
					Owner:           btypes.Address(owner.PubKey().Address()),
					Status:          staketypes.Active,
					BondTokens:      validatorBondTokens,
					BondHeight:      1,
				}
				genesisAccounts = append(genesisAccounts, types.NewQOSAccount(owner.PubKey().Address().Bytes(), btypes.NewInt(validatorOwnerInitQOS), nil))
				appState.Accounts = genesisAccounts
				qosinit.AddValidator(&appState, genVals[i], compound)

				// write private key of validator owner
				ownerFile := filepath.Join(nodeDir, "config", validatorOperatorFile)
				ownerBz, _ := cdc.MarshalJSON(owner)
				cmn.MustWriteFile(ownerFile, ownerBz, nodeFilePerm)
			}

			rawState, _ := cdc.MarshalJSON(appState)
			genDoc := &ttypes.GenesisDoc{
				GenesisTime: time.Now(),
				AppState:    rawState,
			}

			// chainId
			if chainId != "" {
				genDoc.ChainID = chainId
			} else {
				genDoc.ChainID = "test-chain-" + cmn.RandStr(6)
			}

			// Write genesis file.
			for i := 0; i < nValidators; i++ {
				if err := genDoc.SaveAs(filepath.Join(nodeDirs[i], config.Genesis)); err != nil {
					_ = os.RemoveAll(outputDir)
					return err
				}
				config.P2P.PersistentPeers = strings.Join(persistentPeers, ",")
				config.P2P.AddrBookStrict = false
				cfg.WriteConfigFile(filepath.Join(nodeDirs[i], "config", "config.toml"), config)
			}


			fmt.Printf("Successfully initialized %v node directories\n", nValidators)
			return nil
		},
	}

	cmd.Flags().IntVar(&nValidators, "v", 4,
		"Number of validators to initialize the testnet with")
	cmd.Flags().StringVar(&outputDir, "o", "./mytestnet",
		"Directory to store initialization data for the testnet")
	cmd.Flags().StringVar(&nodeDirPrefix, "node-dir-prefix", "node",
		"Prefix the directory name for each node with (node results in node0, node1, ...)")
	cmd.Flags().StringVar(&hostnamePrefix, "hostname-prefix", "node",
		"Hostname prefix (node results in persistent peers list ID0@node0:26656, ID1@node1:26656, ...)")
	cmd.Flags().StringVar(&startingIPAddress, "starting-ip-address", "",
		"Starting IP address (192.168.0.1 results in persistent peers list ID0@192.168.0.1:26656, ID1@192.168.0.2:26656, ...)")
	cmd.Flags().StringVar(&accounts, "genesis-accounts", "",
		"Add genesis accounts to genesis.json, eg: address16lwp3kykkjdc2gdknpjy6u9uhfpa9q4vj78ytd,1000000qos,1000000qstars. Multiple accounts separated by ';'")
	cmd.Flags().StringVar(&qcpRootCA, "qsc-root-ca", "", "Config pubKey of root CA for QCP")
	cmd.Flags().StringVar(&qscRootCA, "qcp-root-ca", "", "Config pubKey of root CA for QSC")
	cmd.Flags().StringVar(&chainId, "chain-id", "", "Chain ID")
	cmd.Flags().BoolVar(&compound, "compound", true, "whether the validator's income is calculated as compound interest, default: true")

	return cmd
}

func hostnameOrIP(i int) string {
	if startingIPAddress != "" {
		ip := net.ParseIP(startingIPAddress)
		ip = ip.To4()
		if ip == nil {
			fmt.Printf("%v: non ipv4 address\n", startingIPAddress)
			os.Exit(1)
		}

		for j := 0; j < i; j++ {
			ip[3]++
		}
		return ip.String()
	}

	return fmt.Sprintf("%s%d", hostnamePrefix, i)
}
