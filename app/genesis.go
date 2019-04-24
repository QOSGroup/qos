package app

import (
	"fmt"
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/approve"
	"github.com/QOSGroup/qos/module/distribution"
	"github.com/QOSGroup/qos/module/gov"
	"github.com/QOSGroup/qos/module/guardian"
	"github.com/QOSGroup/qos/module/mint"
	"github.com/QOSGroup/qos/module/qcp"
	"github.com/QOSGroup/qos/module/qsc"
	"github.com/QOSGroup/qos/module/stake"
	"github.com/QOSGroup/qos/types"
	"github.com/pkg/errors"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// QOS初始状态
type GenesisState struct {
	GenTxs           []txs.TxStd               `json:"gen_txs"`
	Accounts         []*types.QOSAccount       `json:"accounts"`
	MintData         mint.GenesisState         `json:"mint"`
	StakeData        stake.GenesisState        `json:"stake"`
	QCPData          qcp.GenesisState          `json:"qcp"`
	QSCData          qsc.GenesisState          `json:"qsc"`
	ApproveData      approve.GenesisState      `json:"approve"`
	DistributionData distribution.GenesisState `json:"distribution"`
	GovData          gov.GenesisState          `json:"governance"`
	GuardianData     guardian.GenesisState     `json:"guardian"`
}

func NewGenesisState(accounts []*types.QOSAccount,
	mintData mint.GenesisState,
	stakeData stake.GenesisState,
	qcpData qcp.GenesisState,
	qscData qsc.GenesisState,
	approveData approve.GenesisState,
	distributionData distribution.GenesisState,
	govData gov.GenesisState,
	guardianData guardian.GenesisState,
) GenesisState {
	return GenesisState{
		Accounts:         accounts,
		MintData:         mintData,
		StakeData:        stakeData,
		QCPData:          qcpData,
		QSCData:          qscData,
		ApproveData:      approveData,
		DistributionData: distributionData,
		GovData:          govData,
		GuardianData:     guardianData,
	}
}
func NewDefaultGenesisState() GenesisState {
	return GenesisState{
		MintData:         mint.DefaultGenesisState(),
		StakeData:        stake.DefaultGenesisState(),
		DistributionData: distribution.DefaultGenesisState(),
		GovData:          gov.DefaultGenesisState(),
	}
}

func ValidGenesis(state GenesisState) error {
	if err := validateAccounts(state.Accounts); err != nil {
		return err
	}

	if err := stake.ValidateGenesis(state.Accounts, state.StakeData); err != nil {
		return err
	}

	return nil
}

func InitGenesis(ctx context.Context, state GenesisState) []abci.ValidatorUpdate {
	// accounts init should in the first
	initAccounts(ctx, state.Accounts)
	gov.InitGenesis(ctx, state.GovData)
	guardian.InitGenesis(ctx, state.GuardianData)
	mint.InitGenesis(ctx, state.MintData)
	stake.InitGenesis(ctx, state.StakeData)
	qcp.InitGenesis(ctx, state.QCPData)
	qsc.InitGenesis(ctx, state.QSCData)
	approve.InitGenesis(ctx, state.ApproveData)
	distribution.InitGenesis(ctx, state.DistributionData)

	return stake.GetUpdatedValidators(ctx, uint64(state.StakeData.Params.MaxValidatorCnt))
}

func initAccounts(ctx context.Context, accounts []*types.QOSAccount) {
	if len(accounts) == 0 {
		return
	}
	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	for _, acc := range accounts {
		accountMapper.SetAccount(acc)
	}
}

func validateAccounts(accs []*types.QOSAccount) error {
	addrMap := make(map[string]bool, len(accs))
	for i := 0; i < len(accs); i++ {
		acc := accs[i]
		strAddr := string(acc.AccountAddress)
		if _, ok := addrMap[strAddr]; ok {
			return fmt.Errorf("Duplicate account in genesis state: Address %v", acc.AccountAddress)
		}
		addrMap[strAddr] = true
	}
	return nil
}

func CollectStdTxs(cdc *amino.Codec, nodeID string, genTxsDir string, genDoc *tmtypes.GenesisDoc) (
	genTxs []txs.TxStd, persistentPeers string, err error) {

	var fos []os.FileInfo
	fos, err = ioutil.ReadDir(genTxsDir)
	if err != nil {
		return genTxs, persistentPeers, err
	}

	var appState GenesisState
	if err := cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {
		return genTxs, persistentPeers, err
	}

	addrMap := make(map[string]*types.QOSAccount, len(appState.Accounts))
	for i := 0; i < len(appState.Accounts); i++ {
		acc := appState.Accounts[i]
		addrMap[acc.AccountAddress.String()] = acc
	}

	// addresses and IPs (and port) validator server info
	var addressesIPs []string

	var invalidFileNames []string
	var invalidTxFiles []string
	var accsNotInGenesis []string
	var accsNoEnoughQOS []string

	for _, fo := range fos {
		filename := filepath.Join(genTxsDir, fo.Name())
		ext := filepath.Ext(filename)
		simpleName := strings.TrimSuffix(fo.Name(), ext)
		if !fo.IsDir() && (ext != ".json") {
			invalidFileNames = append(invalidFileNames, simpleName)
			continue
		}

		// validate file name, nodeid@ip
		nodeIdAndIp := strings.Split(simpleName, "@")
		if len(nodeIdAndIp) != 2 {
			//TODO valid ip
			invalidFileNames = append(invalidFileNames, simpleName)
			continue
		}
		nodeId := nodeIdAndIp[0]

		// get the genStdTx
		var jsonRawTx []byte
		if jsonRawTx, err = ioutil.ReadFile(filename); err != nil {
			invalidTxFiles = append(invalidTxFiles, simpleName)
			continue
		}
		var txStd txs.TxStd
		if err = cdc.UnmarshalJSON(jsonRawTx, &txStd); err != nil {
			invalidTxFiles = append(invalidTxFiles, simpleName)
			continue
		}
		genTxs = append(genTxs, txStd)

		// genesis transactions must be single-message
		itxs := txStd.ITxs
		if len(itxs) != 1 {
			invalidTxFiles = append(invalidTxFiles, simpleName)
			continue
		}

		txCreateValidator := itxs[0].(*stake.TxCreateValidator)
		// validate delegator and validator addresses and funds against the accounts in the state
		ownerAddr := txCreateValidator.Owner

		delAcc, delOk := addrMap[ownerAddr.String()]

		if !delOk {
			accsNotInGenesis = append(accsNotInGenesis, simpleName+"-"+ownerAddr.String())
			continue
		} else if !delAcc.EnoughOfQOS(btypes.NewInt(int64(txCreateValidator.BondTokens))) {
			accsNoEnoughQOS = append(accsNoEnoughQOS, simpleName+"-"+ownerAddr.String())
			continue
		}

		// exclude itself from persistent peers
		if nodeID != nodeId {
			addressesIPs = append(addressesIPs, fmt.Sprintf("%s:26656", simpleName))
		}
	}

	var errorInfo string
	if len(invalidFileNames) != 0 {
		errorInfo += fmt.Sprintf("file(s) %v name invalid \n", strings.Join(invalidFileNames, " "))
	}
	if len(invalidTxFiles) != 0 {
		errorInfo += fmt.Sprintf("file(s) %v tx invalid \n", strings.Join(invalidTxFiles, " "))
	}
	if len(accsNotInGenesis) != 0 {
		errorInfo += fmt.Sprintf("account(s) %v not in genesis.json \n", strings.Join(accsNotInGenesis, " "))
	}
	if len(accsNoEnoughQOS) != 0 {
		errorInfo += fmt.Sprintf("account(s) %v no enough QOS in genesis.json \n", strings.Join(accsNoEnoughQOS, " "))
	}

	if len(errorInfo) != 0 {
		return genTxs, persistentPeers, errors.New(errorInfo)
	}

	sort.Strings(addressesIPs)
	persistentPeers = strings.Join(addressesIPs, ",")

	return genTxs, persistentPeers, nil
}
