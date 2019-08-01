package stake

import (
	"errors"
	"fmt"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	btxs "github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/bank"
	"github.com/QOSGroup/qos/module/stake/mapper"
	"github.com/QOSGroup/qos/module/stake/txs"
	"github.com/QOSGroup/qos/module/stake/types"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func InitGenesis(ctx context.Context, bapp *baseabci.BaseApp, data types.GenesisState) []abci.ValidatorUpdate {
	validatorMapper := mapper.GetMapper(ctx)

	if len(data.CurrentValidators) > 0 {
		validatorMapper.Set(types.BuildCurrentValidatorsAddressKey(), data.CurrentValidators)
	}

	initValidators(ctx, data.Validators)
	initParams(ctx, data.Params)
	initValidatorsVotesInfo(ctx, data.ValidatorsVoteInfo, data.ValidatorsVoteInWindow)
	initDelegatorsInfo(ctx, data.DelegatorsInfo, data.DelegatorsUnbondInfo, data.ReDelegationsInfo)

	if len(data.GenTxs) > 0 || ctx.BlockHeight() == 0 {
		return initGentxs(ctx, bapp, data.GenTxs)
	} else {
		return GetUpdatedValidators(ctx, uint64(validatorMapper.GetParams(ctx).MaxValidatorCnt))
	}
}

func initGentxs(ctx context.Context, bapp *baseabci.BaseApp, gentxs []btxs.TxStd) []abci.ValidatorUpdate {
	for _, genTx := range gentxs {
		bz := Cdc.MustMarshalBinaryBare(genTx)
		reqDeliverTx := abci.RequestDeliverTx{Tx: bz}
		res := bapp.DeliverTx(reqDeliverTx)
		if !res.IsOK() {
			panic(res.Log)
		}
	}

	validatorSet := []abci.ValidatorUpdate{}
	sm := GetMapper(ctx)
	iterator := sm.IteratorValidatorByVoterPower(false)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		valAddr := btypes.Address(key[9:])
		validator, _ := sm.GetValidator(valAddr)
		validatorSet = append(validatorSet, validator.ToABCIValidatorUpdate(false))
	}

	return validatorSet
}

func initValidators(ctx context.Context, validators []types.Validator) {
	validatorMapper := mapper.GetMapper(ctx)
	for _, v := range validators {

		if validatorMapper.Exists(v.ValidatorPubKey.Address().Bytes()) {
			panic(fmt.Errorf("validator %s already exists", v.ValidatorPubKey.Address()))
		}
		if validatorMapper.ExistsWithOwner(v.Owner) {
			panic(fmt.Errorf("owner %s already bind a validator", v.Owner))
		}
		validatorMapper.CreateValidator(v)
		if !v.IsActive() {
			validatorMapper.MakeValidatorInactive(v.GetValidatorAddress(), v.InactiveHeight, v.InactiveTime, v.InactiveCode)
		}
	}
}

func initValidatorsVotesInfo(ctx context.Context, voteInfos []types.ValidatorVoteInfoState, voteWindowInfos []types.ValidatorVoteInWindowInfoState) {
	sm := mapper.GetMapper(ctx)
	for _, voteInfo := range voteInfos {
		sm.SetValidatorVoteInfo(btypes.Address(voteInfo.ValidatorPubKey.Address()), voteInfo.VoteInfo)
	}

	for _, voteWindowInfo := range voteWindowInfos {
		sm.SetVoteInfoInWindow(btypes.Address(voteWindowInfo.ValidatorPubKey.Address()), voteWindowInfo.Index, voteWindowInfo.Vote)
	}
}

func initDelegatorsInfo(ctx context.Context, delegatorsInfo []types.DelegationInfoState, delegatorsUnbondInfo []types.UnbondingDelegationInfo, redelegationInfo []types.RedelegationInfo) {
	sm := mapper.GetMapper(ctx)

	for _, info := range delegatorsInfo {
		sm.SetDelegationInfo(types.DelegationInfo{
			DelegatorAddr: info.DelegatorAddr,
			ValidatorAddr: btypes.Address(info.ValidatorPubKey.Address()),
			Amount:        info.Amount,
			IsCompound:    info.IsCompound,
		})
	}

	sm.AddUnbondingDelegations(delegatorsUnbondInfo)

	sm.AddRedelegations(redelegationInfo)
}

func initParams(ctx context.Context, params types.Params) {
	mapper := ctx.Mapper(types.MapperName).(*mapper.Mapper)
	mapper.SetParams(ctx, params)
}

func ExportGenesis(ctx context.Context) types.GenesisState {

	validatorMapper := mapper.GetMapper(ctx)
	sm := mapper.GetMapper(ctx)

	var currentValidators []types.Validator
	validatorMapper.Get(types.BuildCurrentValidatorsAddressKey(), &currentValidators)

	params := validatorMapper.GetParams(ctx)

	var validators []types.Validator
	validatorMapper.IterateValidators(func(validator types.Validator) {
		validators = append(validators, validator)
	})

	var validatorsVoteInfo []types.ValidatorVoteInfoState
	sm.IterateVoteInfos(func(valAddr btypes.Address, info types.ValidatorVoteInfo) {

		validator, exists := validatorMapper.GetValidator(valAddr)
		if exists {
			vvis := ValidatorVoteInfoState{
				ValidatorPubKey: validator.GetValidatorPubKey(),
				VoteInfo:        info,
			}
			validatorsVoteInfo = append(validatorsVoteInfo, vvis)
		}
	})

	var validatorsVoteInWindow []types.ValidatorVoteInWindowInfoState
	sm.IterateVoteInWindowsInfos(func(index uint64, valAddr btypes.Address, vote bool) {

		validator, exists := validatorMapper.GetValidator(valAddr)
		if exists {
			validatorsVoteInWindow = append(validatorsVoteInWindow, ValidatorVoteInWindowInfoState{
				ValidatorPubKey: validator.GetValidatorPubKey(),
				Index:           index,
				Vote:            vote,
			})
		}
	})

	var delegatorsInfo []types.DelegationInfoState
	sm.IterateDelegationsInfo(btypes.Address{}, func(info types.DelegationInfo) {

		validator, exists := validatorMapper.GetValidator(info.ValidatorAddr)
		if !exists {
			panic(fmt.Sprintf("validator:%s not exists", info.ValidatorAddr.String()))
		}

		delegatorsInfo = append(delegatorsInfo, DelegationInfoState{
			DelegatorAddr:   info.DelegatorAddr,
			ValidatorPubKey: validator.GetValidatorPubKey(),
			Amount:          info.Amount,
			IsCompound:      info.IsCompound,
		})
	})

	var delegatorsUnbondInfo []types.UnbondingDelegationInfo
	sm.IterateUnbondingDelegations(func(unbondings []types.UnbondingDelegationInfo) {
		delegatorsUnbondInfo = append(delegatorsUnbondInfo, unbondings...)
	})

	var reDelegationsInfo []types.RedelegationInfo
	sm.IterateRedelegationsInfo(func(reDelegations []types.RedelegationInfo) {
		reDelegationsInfo = append(reDelegationsInfo, reDelegations...)
	})

	return GenesisState{
		Params:                 params,
		Validators:             validators,
		ValidatorsVoteInfo:     validatorsVoteInfo,
		ValidatorsVoteInWindow: validatorsVoteInWindow,
		DelegatorsInfo:         delegatorsInfo,
		DelegatorsUnbondInfo:   delegatorsUnbondInfo,
		ReDelegationsInfo:      reDelegationsInfo,
		CurrentValidators:      currentValidators,
	}
}

func CollectStdTxs(cdc *amino.Codec, nodeID string, genTxsDir string, genDoc *tmtypes.GenesisDoc) (
	genTxs []btxs.TxStd, persistentPeers string, err error) {

	var fos []os.FileInfo
	fos, err = ioutil.ReadDir(genTxsDir)
	if err != nil {
		return genTxs, persistentPeers, err
	}

	var appState qtypes.GenesisState
	if err := cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {
		return genTxs, persistentPeers, err
	}
	var bankState bank.GenesisState
	cdc.MustUnmarshalJSON(appState[bank.ModuleName], &bankState)

	addrMap := make(map[string]*qtypes.QOSAccount, len(bankState.Accounts))
	for i := 0; i < len(bankState.Accounts); i++ {
		acc := bankState.Accounts[i]
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
		var txStd btxs.TxStd
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

		txCreateValidator := itxs[0].(*txs.TxCreateValidator)
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
