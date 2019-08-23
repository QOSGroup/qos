package mapper

import (
	"errors"
	"fmt"
	"github.com/QOSGroup/qos/module/stake"
	"runtime/debug"

	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/distribution/types"
	qtypes "github.com/QOSGroup/qos/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
)

/*

custom path:
/custom/distribution/$query path

query path:
	/validatorPeriodInfo/:ownerAddr : 根据validator owner地址查询validator period info
	/delegatorIncomeInfo/:delegatorAddr/:ownerAddr : 查询delegator地址查询收益计算信息

	xxx为bech32 address

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

	if route[0] == types.ValidatorPeriodInfo {
		ownerAddr, _ := btypes.GetAddrFromBech32(route[1])
		data, e = queryValidatorPeriodInfo(ctx, ownerAddr)
	} else if route[0] == types.DelegatorIncomeInfo {
		deleAddr, _ := btypes.GetAddrFromBech32(route[1])
		ownerAddr, _ := btypes.GetAddrFromBech32(route[2])
		data, e = queryDelegatorIncomeInfo(ctx, deleAddr, ownerAddr)
	} else {
		data = nil
		e = errors.New("not found match path")
	}

	if e != nil {
		return nil, btypes.ErrInternal(e.Error())
	}

	return data, nil
}

func queryValidatorPeriodInfo(ctx context.Context, owner btypes.Address) ([]byte, error) {
	dm := GetMapper(ctx)
	sm := stake.GetMapper(ctx)

	validator, exists := sm.GetValidatorByOwner(owner)
	if !exists {
		return nil, fmt.Errorf("validator not exists. owner: %s", owner.String())
	}

	vcps, exists := dm.GetValidatorCurrentPeriodSummary(validator.GetValidatorAddress())
	if !exists {
		return nil, fmt.Errorf("validator current period not exists. owner: %s", owner.String())
	}

	result := ValidatorPeriodInfoQueryResult{
		OwnerAddr:       validator.GetOwner(),
		ValidatorPubKey: validator.GetValidatorPubKey(),
		Fees:            vcps.Fees,
		CurrentTokens:   validator.GetBondTokens(),
		CurrentPeriod:   vcps.Period,
	}

	if vcps.Period >= 1 {
		frac := dm.GetValidatorHistoryPeriodSummary(validator.GetValidatorAddress(), vcps.Period-1)
		result.LastPeriod = vcps.Period - 1
		result.LastPeriodFraction = frac
	}

	return dm.GetCodec().MarshalJSON(result)
}

func queryDelegatorIncomeInfo(ctx context.Context, delegator btypes.Address, owner btypes.Address) ([]byte, error) {
	dm := GetMapper(ctx)
	sm := stake.GetMapper(ctx)

	validator, exists := sm.GetValidatorByOwner(owner)
	if !exists {
		return nil, fmt.Errorf("validator not exists. owner: %s", owner.String())
	}

	info, exists := dm.GetDelegatorEarningStartInfo(validator.GetValidatorAddress(), delegator)
	if !exists {
		return nil, fmt.Errorf("delegator income info not exists. delegator: %s , owner: %s", delegator.String(), owner.String())
	}

	result := DelegatorIncomeInfoQueryResult{
		OwnerAddr:             validator.GetOwner(),
		ValidatorPubKey:       validator.GetValidatorPubKey(),
		PreviousPeriod:        info.PreviousPeriod,
		BondToken:             info.BondToken,
		CurrentStartingHeight: info.CurrentStartingHeight,
		FirstDelegateHeight:   info.FirstDelegateHeight,
		HistoricalRewardFees:  info.HistoricalRewardFees,
		LastIncomeCalHeight:   info.LastIncomeCalHeight,
		LastIncomeCalFees:     info.LastIncomeCalFees,
	}
	return dm.GetCodec().MarshalJSON(result)
}

type ValidatorPeriodInfoQueryResult struct {
	OwnerAddr          btypes.Address  `json:"owner_address"`
	ValidatorPubKey    crypto.PubKey   `json:"validator_pub_key"`
	Fees               btypes.BigInt   `json:"fees"`
	CurrentTokens      uint64          `json:"current_tokens"`
	CurrentPeriod      uint64          `json:"current_period"`
	LastPeriod         uint64          `json:"last_period"`
	LastPeriodFraction qtypes.Fraction `json:"last_period_fraction"`
}

type DelegatorIncomeInfoQueryResult struct {
	OwnerAddr             btypes.Address `json:"owner_address"`
	ValidatorPubKey       crypto.PubKey  `json:"validator_pub_key"`
	PreviousPeriod        uint64         `json:"previous_validator_period"`
	BondToken             uint64         `json:"bond_token"`
	CurrentStartingHeight uint64         `json:"earns_starting_height"`
	FirstDelegateHeight   uint64         `json:"first_delegate_height"`
	HistoricalRewardFees  btypes.BigInt  `json:"historical_rewards"`
	LastIncomeCalHeight   uint64         `json:"last_income_calHeight"`
	LastIncomeCalFees     btypes.BigInt  `json:"last_income_calFees"`
}
