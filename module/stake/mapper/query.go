package mapper

import (
	"errors"
	"fmt"
	"github.com/QOSGroup/qos/module/stake/types"
	"runtime/debug"

	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
)

/*

custom path:
/custom/stake/$query path

 query path:
	/delegation/:delegatorAddr/:ownerAddr : 根据delegator和owner查询委托信息(first: delegator)
	/delegations/owner/:ownerAddr : 查询owner下的所有委托信息
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
		deleAddr, _ := btypes.GetAddrFromBech32(route[1])
		ownerAddr, _ := btypes.GetAddrFromBech32(route[2])
		data, e = getDelegationByOwnerAndDelegator(ctx, ownerAddr, deleAddr)

	} else if route[0] == types.Delegations && route[1] == types.Owner {
		ownerAddr, _ := btypes.GetAddrFromBech32(route[2])
		data, e = getDelegationsByOwner(ctx, ownerAddr)

	} else if route[0] == types.Delegations && route[1] == types.Delegator {
		deleAddr, _ := btypes.GetAddrFromBech32(route[2])
		data, e = getDelegationsByDelegator(ctx, deleAddr)

	} else if route[0] == types.Unbondings {
		deleAddr, _ := btypes.GetAddrFromBech32(route[1])
		data, e = getUnbondingsByDelegator(ctx, deleAddr)

	} else if route[0] == types.Redelegations {
		deleAddr, _ := btypes.GetAddrFromBech32(route[1])
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

func getDelegationByOwnerAndDelegator(ctx context.Context, owner, delegator btypes.Address) ([]byte, error) {
	sm := GetMapper(ctx)

	validator, exists := sm.GetValidatorByOwner(owner)
	if !exists {
		return nil, fmt.Errorf("validator not exists. owner: %s", owner.String())
	}

	info, exists := sm.GetDelegationInfo(delegator, validator.GetValidatorAddress())
	if !exists {
		return nil, fmt.Errorf("delegationInfo not exists. owner: %s , deleAddr: %s", owner.String(), delegator.String())
	}

	result := infoToDelegationQueryResult(validator, info)
	return sm.GetCodec().MarshalJSON(result)
}

func getDelegationsByOwner(ctx context.Context, owner btypes.Address) ([]byte, error) {
	sm := GetMapper(ctx)

	validator, exists := sm.GetValidatorByOwner(owner)
	if !exists {
		return nil, fmt.Errorf("validator not exists. owner: %s", owner.String())
	}

	var result []DelegationQueryResult
	sm.IterateDelegationsValDeleAddr(validator.GetValidatorAddress(), func(valAddr, deleAddr btypes.Address) {
		info, _ := sm.GetDelegationInfo(deleAddr, valAddr)
		result = append(result, infoToDelegationQueryResult(validator, info))
	})

	return sm.GetCodec().MarshalJSON(result)
}

func getDelegationsByDelegator(ctx context.Context, delegator btypes.Address) ([]byte, error) {
	sm := GetMapper(ctx)

	var result []DelegationQueryResult
	sm.IterateDelegationsInfo(delegator, func(info types.DelegationInfo) {
		validator, _ := sm.GetValidator(info.ValidatorAddr)
		result = append(result, infoToDelegationQueryResult(validator, info))
	})

	return sm.GetCodec().MarshalJSON(result)
}

func infoToDelegationQueryResult(validator types.Validator, info types.DelegationInfo) DelegationQueryResult {
	return NewDelegationQueryResult(info.DelegatorAddr, validator.Owner, validator.ValidatorPubKey, info.Amount, info.IsCompound)
}

type DelegationQueryResult struct {
	DelegatorAddr   btypes.Address `json:"delegator_address"`
	OwnerAddr       btypes.Address `json:"owner_address"`
	ValidatorPubKey crypto.PubKey  `json:"validator_pub_key"`
	Amount          uint64         `json:"delegate_amount"`
	IsCompound      bool           `json:"is_compound"`
}

func NewDelegationQueryResult(deleAddr, ownerAddr btypes.Address, valPubkey crypto.PubKey, amount uint64, compound bool) DelegationQueryResult {
	return DelegationQueryResult{
		DelegatorAddr:   deleAddr,
		OwnerAddr:       ownerAddr,
		ValidatorPubKey: valPubkey,
		Amount:          amount,
		IsCompound:      compound,
	}
}

func getUnbondingsByDelegator(ctx context.Context, delegator btypes.Address) ([]byte, error) {
	sm := GetMapper(ctx)
	result := sm.GetUnbondingDelegationsByDelegator(delegator)

	return sm.GetCodec().MarshalJSON(result)
}

func getRedelegationsByDelegator(ctx context.Context, delegator btypes.Address) ([]byte, error) {
	sm := GetMapper(ctx)
	result := sm.GetRedelegationsByDelegator(delegator)

	return sm.GetCodec().MarshalJSON(result)
}
