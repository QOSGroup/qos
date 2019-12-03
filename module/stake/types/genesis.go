package types

import (
	"fmt"
	"github.com/QOSGroup/qbase/txs"

	btypes "github.com/QOSGroup/qbase/types"
)

type GenesisState struct {
	GenTxs                 []txs.TxStd                      `json:"gen_txs"`               // signed TxCreateValidator in genesis.json
	Params                 Params                           `json:"params"`                // stake module parameters
	Validators             []Validator                      `json:"validators"`            // validatorKey, validatorByOwnerKey,validatorByInactiveKey,validatorByVotePowerKey
	ValidatorsVoteInfo     []ValidatorVoteInfoState         `json:"val_votes_info"`        // validatorVoteInfoKey
	ValidatorsVoteInWindow []ValidatorVoteInWindowInfoState `json:"val_votes_in_window"`   // validatorVoteInfoInWindowKey
	DelegatorsInfo         []DelegationInfoState            `json:"delegators_info"`       // DelegationByDelValKey, DelegationByValDelKey
	DelegatorsUnbondInfo   []UnbondingDelegationInfo        `json:"delegator_unbond_info"` // UnbondingHeightDelegatorKey
	ReDelegationsInfo      []RedelegationInfo               `json:"redelegations_info"`    // ReDelegationHeightDelegatorKey
	CurrentValidators      []Validator                      `json:"current_validators"`    // currentValidatorsAddressKey
}

func NewGenesisState(params Params,
	validators []Validator,
	validatorsVoteInfo []ValidatorVoteInfoState,
	validatorsVoteInWindow []ValidatorVoteInWindowInfoState,
	delegatorsInfo []DelegationInfoState,
	delegatorsUnbondInfo []UnbondingDelegationInfo,
	reDelegationsInfo []RedelegationInfo,
	currentValidators []Validator) GenesisState {
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

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

func ValidateGenesis(data GenesisState) error {
	// 校验验证节点信息
	err := validateValidators(data.Validators)
	if err != nil {
		return err
	}

	// 验证参数
	err = data.Params.Validate()
	if err != nil {
		return err
	}

	return nil
}

func validateValidators(validators []Validator) (err error) {
	addrMap := make(map[string]bool, len(validators))
	ownerMap := make(map[string]bool, len(validators))
	for i := 0; i < len(validators); i++ {
		val := validators[i]
		strKey := string(val.GetConsensusPubKey().Bytes())
		if _, ok := addrMap[strKey]; ok {
			return fmt.Errorf("duplicate validator in genesis state: Name %v, operator %v", val.Description.Moniker, val.OperatorAddress.String())
		}
		if _, ok := ownerMap[val.OperatorAddress.String()]; ok {
			return fmt.Errorf("duplicate operator in genesis state: Name %v, operator %v", val.Description.Moniker, val.OperatorAddress.String())
		}
		if val.Status != Active {
			return fmt.Errorf("validator is bonded and jailed in genesis state: Name %v, operator %v", val.Description.Moniker, val.OperatorAddress.String())
		}
		addrMap[strKey] = true
		ownerMap[val.OperatorAddress.String()] = true
	}
	return nil
}

type ValidatorVoteInfoState struct {
	ValidatorAddr btypes.ValAddress `json:"validator_addr"`
	VoteInfo      ValidatorVoteInfo `json:"vote_info"`
}

type ValidatorVoteInWindowInfoState struct {
	ValidatorAddr btypes.ValAddress `json:"validator_addr"`
	Index         int64             `json:"index"`
	Vote          bool              `json:"vote"`
}

type DelegationInfoState struct {
	DelegatorAddr btypes.AccAddress `json:"delegator_addr"`
	ValidatorAddr btypes.ValAddress `json:"validator_addr"`
	Amount        btypes.BigInt     `json:"delegate_amount"`
	IsCompound    bool              `json:"is_compound"`
}
