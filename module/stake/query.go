package stake

import (
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	ecomapper "github.com/QOSGroup/qos/module/eco/mapper"
	ecotypes "github.com/QOSGroup/qos/module/eco/types"
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

	if len(route) < 3 {
		return nil, btypes.ErrInternal("custom query miss parameters")
	}

	var data []byte
	var e error

	if route[0] == ecotypes.Delegation {
		deleAddr, _ := btypes.GetAddrFromBech32(route[1])
		ownerAddr, _ := btypes.GetAddrFromBech32(route[2])
		data, e = getDelegationByOwnerAndDelegator(ctx, ownerAddr, deleAddr)

	} else if route[0] == ecotypes.Delegations && route[1] == ecotypes.Owner {
		ownerAddr, _ := btypes.GetAddrFromBech32(route[2])
		data, e = getDelegationsByOwner(ctx, ownerAddr)

	} else if route[0] == ecotypes.Delegations && route[1] == ecotypes.Delegator {
		deleAddr, _ := btypes.GetAddrFromBech32(route[2])
		data, e = getDelegationsByDelegator(ctx, deleAddr)

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
	validatorMapper := ecomapper.GetValidatorMapper(ctx)
	delegationMapper := ecomapper.GetDelegationMapper(ctx)

	validator, exsits := validatorMapper.GetValidatorByOwner(owner)
	if !exsits {
		return nil, fmt.Errorf("validator not exsits. owner: %s", owner.String())
	}

	info, exsits := delegationMapper.GetDelegationInfo(delegator, validator.GetValidatorAddress())
	if !exsits {
		return nil, fmt.Errorf("delegationInfo not exsits. owner: %s , deleAddr: %s", owner.String(), delegator.String())
	}

	result := infoToDelegationQueryResult(validator, info)
	return validatorMapper.GetCodec().MarshalJSON(result)
}

func getDelegationsByOwner(ctx context.Context, owner btypes.Address) ([]byte, error) {
	validatorMapper := ecomapper.GetValidatorMapper(ctx)
	delegationMapper := ecomapper.GetDelegationMapper(ctx)

	validator, exsits := validatorMapper.GetValidatorByOwner(owner)
	if !exsits {
		return nil, fmt.Errorf("validator not exsits. owner: %s", owner.String())
	}

	var result []DelegationQueryResult
	delegationMapper.IterateDelegationsValDeleAddr(validator.GetValidatorAddress(), func(valAddr, deleAddr btypes.Address) {
		info, _ := delegationMapper.GetDelegationInfo(deleAddr, valAddr)
		result = append(result, infoToDelegationQueryResult(validator, info))
	})

	return validatorMapper.GetCodec().MarshalJSON(result)
}

func getDelegationsByDelegator(ctx context.Context, delegator btypes.Address) ([]byte, error) {
	validatorMapper := ecomapper.GetValidatorMapper(ctx)
	delegationMapper := ecomapper.GetDelegationMapper(ctx)

	var result []DelegationQueryResult
	delegationMapper.IterateDelegationsInfo(delegator, func(info ecotypes.DelegationInfo) {
		validator, _ := validatorMapper.GetValidator(info.ValidatorAddr)
		result = append(result, infoToDelegationQueryResult(validator, info))
	})

	return validatorMapper.GetCodec().MarshalJSON(result)
}

func infoToDelegationQueryResult(validator ecotypes.Validator, info ecotypes.DelegationInfo) DelegationQueryResult {
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
