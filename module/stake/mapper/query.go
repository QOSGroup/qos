package mapper

import (
	"errors"
	"fmt"
	"github.com/QOSGroup/qos/module/stake/types"
	"runtime/debug"

	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

/*

custom path:
/custom/stake/$query path

 query path:
	/delegation/:delegatorAddr/:operatorAddr : 根据delegator和owner查询委托信息(first: delegator)
	/delegations/validator/:operatorAddr : 查询owner下的所有委托信息
	/delegations/delegator/:delegatorAddr : 查询delegator的所有委托信息

return:
  json字节数组
*/

func Query(ctx context.Context, route []string, req abci.RequestQuery) (res []byte, err btypes.Error) {

	defer func() {
		if r := recover(); r != nil {
			err = btypes.ErrInternal(string(debug.Stack()))
			return
		}
	}()

	if len(route) < 2 {
		return nil, btypes.ErrInternal("custom query miss parameters")
	}

	var data []byte
	var e error

	if route[0] == types.Delegation {
		deleAddr, _ := btypes.AccAddressFromBech32(route[1])
		validatorAddr, _ := btypes.ValAddressFromBech32(route[2])
		data, e = getDelegationByOwnerAndDelegator(ctx, validatorAddr, deleAddr)

	} else if route[0] == types.Delegations && route[1] == types.ValidatorFlag {
		ownerAddr, _ := btypes.ValAddressFromBech32(route[2])
		data, e = getDelegationsByValidator(ctx, ownerAddr)

	} else if route[0] == types.Delegations && route[1] == types.Delegator {
		deleAddr, _ := btypes.AccAddressFromBech32(route[2])
		data, e = getDelegationsByDelegator(ctx, deleAddr)

	} else if route[0] == types.Unbondings {
		deleAddr, _ := btypes.AccAddressFromBech32(route[1])
		data, e = getUnbondingsByDelegator(ctx, deleAddr)

	} else if route[0] == types.Redelegations {
		deleAddr, _ := btypes.AccAddressFromBech32(route[1])
		data, e = getRedelegationsByDelegator(ctx, deleAddr)

	} else {
		data = nil
		e = errors.New("not found match path")
	}

	if e != nil {
		return nil, btypes.ErrInternal(e.Error())
	}

	return data, nil
}

func getDelegationByOwnerAndDelegator(ctx context.Context, validatorAddr btypes.ValAddress, delegator btypes.AccAddress) ([]byte, error) {
	sm := GetMapper(ctx)

	validator, exists := sm.GetValidator(validatorAddr)
	if !exists {
		return nil, fmt.Errorf("validator not exists. owner: %s", validatorAddr.String())
	}

	info, exists := sm.GetDelegationInfo(delegator, validator.GetValidatorAddress())
	if !exists {
		return nil, fmt.Errorf("delegationInfo not exists. owner: %s , deleAddr: %s", validatorAddr.String(), delegator.String())
	}

	result := infoToDelegationQueryResult(validator, info)
	return sm.GetCodec().MarshalJSON(result)
}

func getDelegationsByValidator(ctx context.Context, validatorAddr btypes.ValAddress) ([]byte, error) {
	sm := GetMapper(ctx)

	validator, exists := sm.GetValidator(validatorAddr)
	if !exists {
		return nil, fmt.Errorf("validator not exists. owner: %s", validatorAddr.String())
	}

	var result []DelegationQueryResult
	sm.IterateDelegationsValDeleAddr(validator.GetValidatorAddress(), func(valAddr btypes.ValAddress, deleAddr btypes.AccAddress) {
		info, _ := sm.GetDelegationInfo(deleAddr, valAddr)
		result = append(result, infoToDelegationQueryResult(validator, info))
	})

	return sm.GetCodec().MarshalJSON(result)
}

func getDelegationsByDelegator(ctx context.Context, delegator btypes.AccAddress) ([]byte, error) {
	sm := GetMapper(ctx)

	var result []DelegationQueryResult
	sm.IterateDelegationsInfo(delegator, func(info types.DelegationInfo) {
		validator, _ := sm.GetValidator(info.ValidatorAddr)
		result = append(result, infoToDelegationQueryResult(validator, info))
	})

	return sm.GetCodec().MarshalJSON(result)
}

func infoToDelegationQueryResult(validator types.Validator, info types.DelegationInfo) DelegationQueryResult {
	consPubKey, _  := btypes.ConsensusPubKeyString(validator.GetConsensusPubKey())
	return NewDelegationQueryResult(info.DelegatorAddr,
		validator.GetValidatorAddress(), consPubKey, info.Amount, info.IsCompound)
}

type DelegationQueryResult struct {
	DelegatorAddr   btypes.AccAddress `json:"delegator_address"`
	ValidatorAddr       btypes.ValAddress `json:"owner_address"`
	ValidatorConsensusPubKey string  `json:"validator_cons_pub_key"`
	Amount          uint64         `json:"delegate_amount"`
	IsCompound      bool           `json:"is_compound"`
}

func NewDelegationQueryResult(deleAddr btypes.AccAddress, ownerAddr btypes.ValAddress, bench32ConPubKey string, amount uint64, compound bool) DelegationQueryResult {
	return DelegationQueryResult{
		DelegatorAddr:   deleAddr,
		ValidatorAddr:       ownerAddr,
		ValidatorConsensusPubKey: bench32ConPubKey,
		Amount:          amount,
		IsCompound:      compound,
	}
}

func getUnbondingsByDelegator(ctx context.Context, delegator btypes.AccAddress) ([]byte, error) {
	sm := GetMapper(ctx)
	result := sm.GetUnbondingDelegationsByDelegator(delegator)

	return sm.GetCodec().MarshalJSON(result)
}

func getRedelegationsByDelegator(ctx context.Context, delegator btypes.AccAddress) ([]byte, error) {
	sm := GetMapper(ctx)
	result := sm.GetRedelegationsByDelegator(delegator)

	return sm.GetCodec().MarshalJSON(result)
}
