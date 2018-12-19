package testnet

import (
	"fmt"
	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/app"
	"github.com/QOSGroup/qos/types"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	cfg "github.com/tendermint/tendermint/config"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	ttypes "github.com/tendermint/tendermint/types"
)

var (
	chainId string
	moniker string

	nValidators    int
	nNonValidators int
	outputDir      string
	nodeDirPrefix  string

	populatePersistentPeers bool
	hostnamePrefix          string
	startingIPAddress       string
	p2pPort                 int

	rootCA   string
	accounts string
)

const (
	nodeDirPerm = 0755
)

func TestnetFileCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "testnet",
		Short: "Initialize files for a QOS testnet",
		Long: `testnet will create "v" + "n" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).

Note, strict routability for addresses is turned off in the config file.

Optionally, it will fill in persistent_peers list in config file using either hostnames or IPs.

Example:

	qosd testnet --chain-id=qostest --v=4 --o=./output --starting-ip-address=192.168.1.2 --genesis-accounts=address16lwp3kykkjdc2gdknpjy6u9uhfpa9q4vj78ytd,1000000qos,1000000qstars
	`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config := cfg.DefaultConfig()

			// moniker
			if moniker != "" {
				config.BaseConfig.Moniker = moniker
			}

			genVals := make([]types.Validator, nValidators)

			// accounts
			genesisAccounts := make([]*account.QOSAccount, 0)
			var err error
			if accounts != "" {
				genesisAccounts, err = account.ParseAccounts(accounts)
				if err != nil {
					return err
				}
			}

			// root ca
			var pubKey crypto.PubKey
			if rootCA != "" {
				err := cdc.UnmarshalJSON(cmn.MustReadFile(rootCA), &pubKey)
				if err != nil {
					return err
				}
			}

			// validators
			for i := 0; i < nValidators; i++ {
				nodeDirName := cmn.Fmt("%s%d", nodeDirPrefix, i)
				nodeDir := filepath.Join(outputDir, nodeDirName)
				config.SetRoot(nodeDir)

				err := os.MkdirAll(filepath.Join(nodeDir, "config"), nodeDirPerm)
				if err != nil {
					_ = os.RemoveAll(outputDir)
					return err
				}

				initFilesWithConfig(config)

				//todo
				// pvFile := filepath.Join(nodeDir, config.BaseConfig.PrivValidator)
				// pv := privval.LoadFilePV(pvFile)

				// genVals[i] = types.Validator{
				// 	Name:        nodeDirName,
				// 	ConsPubKey:  pv.GetPubKey(),
				// 	Operator:    pv.GetPubKey().Address().Bytes(),
				// 	VotingPower: 1,
				// 	Height:      1,
				// }
			}

			// non-validators
			for i := 0; i < nNonValidators; i++ {
				nodeDir := filepath.Join(outputDir, cmn.Fmt("%s%d", nodeDirPrefix, i+nValidators))
				config.SetRoot(nodeDir)

				err := os.MkdirAll(filepath.Join(nodeDir, "config"), nodeDirPerm)
				if err != nil {
					_ = os.RemoveAll(outputDir)
					return err
				}

				initFilesWithConfig(config)
			}

			appState := app.GenesisState{
				CAPubKey:   pubKey,
				Validators: genVals,
				Accounts:   genesisAccounts,
			}
			rawState, _ := cdc.MarshalJSON(appState)

			// Generate genesis doc from generated validators
			genDoc := &ttypes.GenesisDoc{
				GenesisTime: time.Now(),
				Validators:  make([]ttypes.GenesisValidator, 0),
				AppState:    rawState,
			}

			// chainId
			if chainId != "" {
				genDoc.ChainID = chainId
			} else {
				genDoc.ChainID = "chain-" + cmn.RandStr(6)
			}

			// Write genesis file.
			for i := 0; i < nValidators+nNonValidators; i++ {
				nodeDir := filepath.Join(outputDir, cmn.Fmt("%s%d", nodeDirPrefix, i))
				if err := genDoc.SaveAs(filepath.Join(nodeDir, config.BaseConfig.Genesis)); err != nil {
					_ = os.RemoveAll(outputDir)
					return err
				}
			}

			if populatePersistentPeers {
				err := populatePersistentPeersInConfigAndWriteIt(config)
				if err != nil {
					_ = os.RemoveAll(outputDir)
					return err
				}
			}

			fmt.Printf("Successfully initialized %v node directories\n", nValidators+nNonValidators)
			return nil
		},
	}

	cmd.Flags().IntVar(&nValidators, "v", 4,
		"Number of validators to initialize the testnet with")
	cmd.Flags().IntVar(&nNonValidators, "n", 0,
		"Number of non-validators to initialize the testnet with")
	cmd.Flags().StringVar(&outputDir, "o", "./mytestnet",
		"Directory to store initialization data for the testnet")
	cmd.Flags().StringVar(&nodeDirPrefix, "node-dir-prefix", "node",
		"Prefix the directory name for each node with (node results in node0, node1, ...)")

	cmd.Flags().BoolVar(&populatePersistentPeers, "populate-persistent-peers", true,
		"Update config of each node with the list of persistent peers build using either hostname-prefix or starting-ip-address")
	cmd.Flags().StringVar(&hostnamePrefix, "hostname-prefix", "node",
		"Hostname prefix (node results in persistent peers list ID0@node0:26656, ID1@node1:26656, ...)")
	cmd.Flags().StringVar(&startingIPAddress, "starting-ip-address", "",
		"Starting IP address (192.168.0.1 results in persistent peers list ID0@192.168.0.1:26656, ID1@192.168.0.2:26656, ...)")
	cmd.Flags().IntVar(&p2pPort, "p2p-port", 26656,
		"P2P Port")
	cmd.Flags().StringVar(&accounts, "genesis-accounts", "",
		"Add genesis accounts to genesis.json, eg: address16lwp3kykkjdc2gdknpjy6u9uhfpa9q4vj78ytd,1000000qos,1000000qstars. Multiple accounts separated by ';'")
	cmd.Flags().StringVar(&rootCA, "root-ca", "", "Config pubKey of root CA")
	cmd.Flags().StringVar(&chainId, "chain-id", "", "Chain ID")
	cmd.Flags().StringVar(&moniker, "moniker", "", "Moniker")

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

func populatePersistentPeersInConfigAndWriteIt(config *cfg.Config) error {
	persistentPeers := make([]string, nValidators+nNonValidators)
	for i := 0; i < nValidators+nNonValidators; i++ {
		nodeDir := filepath.Join(outputDir, cmn.Fmt("%s%d", nodeDirPrefix, i))
		config.SetRoot(nodeDir)
		nodeKey, err := p2p.LoadNodeKey(config.NodeKeyFile())
		if err != nil {
			return err
		}
		persistentPeers[i] = p2p.IDAddressString(nodeKey.ID(), fmt.Sprintf("%s:%d", hostnameOrIP(i), p2pPort))
	}
	persistentPeersList := strings.Join(persistentPeers, ",")

	for i := 0; i < nValidators+nNonValidators; i++ {
		nodeDir := filepath.Join(outputDir, cmn.Fmt("%s%d", nodeDirPrefix, i))
		config.SetRoot(nodeDir)
		config.P2P.PersistentPeers = persistentPeersList
		config.P2P.AddrBookStrict = false

		// overwrite default config
		cfg.WriteConfigFile(filepath.Join(nodeDir, "config", "config.toml"), config)
	}

	return nil
}

func initFilesWithConfig(config *cfg.Config) error {
	// private validator
	privValFile := config.PrivValidatorFile()
	var pv *privval.FilePV
	if cmn.FileExists(privValFile) {
		pv = privval.LoadFilePV(privValFile)
	} else {
		pv = privval.GenFilePV(privValFile)
		pv.Save()
	}

	nodeKeyFile := config.NodeKeyFile()
	if !cmn.FileExists(nodeKeyFile) {
		if _, err := p2p.LoadOrGenNodeKey(nodeKeyFile); err != nil {
			return err
		}
	}

	// genesis file
	genFile := config.GenesisFile()
	if !cmn.FileExists(genFile) {
		genDoc := ttypes.GenesisDoc{
			ChainID:         cmn.Fmt("test-chain-%v", cmn.RandStr(6)),
			GenesisTime:     time.Now(),
			ConsensusParams: ttypes.DefaultConsensusParams(),
		}
		genDoc.Validators = []ttypes.GenesisValidator{{
			PubKey: pv.GetPubKey(),
			Power:  10,
		}}

		if err := genDoc.SaveAs(genFile); err != nil {
			return err
		}
	}

	return nil
}
