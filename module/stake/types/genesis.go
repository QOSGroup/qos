package types

import (
	"fmt"
	"github.com/QOSGroup/qbase/txs"

	btypes "github.com/QOSGroup/qbase/types"
	"github.com/tendermint/tendermint/crypto"
)

type GenesisState struct {
	GenTxs                 []txs.TxStd                      `json:"gen_txs"`
	Params                 Params                           `json:"params"`
	Validators             []Validator                      `json:"validators"`            //validatorKey, validatorByOwnerKey,validatorByInactiveKey,validatorByVotePowerKey
	ValidatorsVoteInfo     []ValidatorVoteInfoState         `json:"val_votes_info"`        //validatorVoteInfoKey
	ValidatorsVoteInWindow []ValidatorVoteInWindowInfoState `json:"val_votes_in_window"`   //validatorVoteInfoInWindowKey
	DelegatorsInfo         []DelegationInfoState            `json:"delegators_info"`       //DelegationByDelValKey, DelegationByValDelKey
	DelegatorsUnbondInfo   []UnbondingDelegationInfo        `json:"delegator_unbond_info"` //UnbondingHeightDelegatorKey
	ReDelegationsInfo      []RedelegationInfo               `json:"redelegations_info"`    //ReDelegationHeightDelegatorKey
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
	err := validateValidators(data.Validators)
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
		strKey := string(val.ValidatorPubKey.Bytes())
		if _, ok := addrMap[strKey]; ok {
			return fmt.Errorf("duplicate validator in genesis state: Name %v, Owner %v", val.Description.Moniker, val.Owner)
		}
		if _, ok := ownerMap[val.Owner.String()]; ok {
			return fmt.Errorf("duplicate owner in genesis state: Name %v, Owner %v", val.Description.Moniker, val.Owner)
		}
		if val.Status != Active {
			return fmt.Errorf("validator is bonded and jailed in genesis state: Name %v, Owner %v", val.Description.Moniker, val.Owner)
		}
		addrMap[strKey] = true
		ownerMap[val.Owner.String()] = true
	}
	return nil
}

type ValidatorVoteInfoState struct {
	ValidatorPubKey crypto.PubKey     `json:"validator_pub_key"`
	VoteInfo        ValidatorVoteInfo `json:"vote_info"`
}

type ValidatorVoteInWindowInfoState struct {
	ValidatorPubKey crypto.PubKey `json:"validator_pub_key"`
	Index           uint64        `json:"index"`
	Vote            bool          `json:"vote"`
}

type DelegationInfoState struct {
	DelegatorAddr   btypes.Address `json:"delegator_addr"`
	ValidatorPubKey crypto.PubKey  `json:"validator_pub_key"`
	Amount          uint64         `json:"delegate_amount"`
	IsCompound      bool           `json:"is_compound"`
}
