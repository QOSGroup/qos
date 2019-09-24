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
)

/*

custom path:
/custom/distribution/$query path

query path:
	/validatorPeriodInfo/:valOperatorAddr : 根据validator operator地址查询validator period info
	/delegatorIncomeInfo/:delegatorAddr/:valOperatorAddr : 查询delegator地址查询收益计算信息

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
		valAddr, _ := btypes.ValAddressFromBech32(route[1])
		data, e = queryValidatorPeriodInfo(ctx, valAddr)
	} else if route[0] == types.DelegatorIncomeInfo {
		deleAddr, _ := btypes.AccAddressFromBech32(route[1])
		valAddr, _ := btypes.ValAddressFromBech32(route[2])
		data, e = queryDelegatorIncomeInfo(ctx, deleAddr, valAddr)
	} else {
		data = nil
		e = errors.New("not found match path")
	}

	if e != nil {
		return nil, btypes.ErrInternal(e.Error())
	}

	return data, nil
}

// 查询验证节点收益信息
func queryValidatorPeriodInfo(ctx context.Context, valAddr btypes.ValAddress) ([]byte, error) {
	dm := GetMapper(ctx)
	sm := stake.GetMapper(ctx)

	validator, exists := sm.GetValidator(valAddr)
	if !exists {
		return nil, fmt.Errorf("validator not exists. validator-address: %s", valAddr.String())
	}

	vcps, exists := dm.GetValidatorCurrentPeriodSummary(validator.GetValidatorAddress())
	if !exists {
		return nil, fmt.Errorf("validator current period not exists. validator-address: %s", valAddr.String())
	}

	result := ValidatorPeriodInfoQueryResult{
		OperatorAddress: validator.OperatorAddress,
		ConsPubKey:      btypes.MustConsensusPubKeyString(validator.ConsPubKey),
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

// 查询委托收益信息
func queryDelegatorIncomeInfo(ctx context.Context, delegator btypes.AccAddress, valAddr btypes.ValAddress) ([]byte, error) {
	dm := GetMapper(ctx)
	sm := stake.GetMapper(ctx)

	validator, exists := sm.GetValidator(valAddr)
	if !exists {
		return nil, fmt.Errorf("validator not exists. validator-address: %s", valAddr.String())
	}

	info, exists := dm.GetDelegatorEarningStartInfo(validator.GetValidatorAddress(), delegator)
	if !exists {
		return nil, fmt.Errorf("delegator income info not exists. delegator: %s , validator-address: %s", delegator.String(), valAddr.String())
	}

	result := DelegatorIncomeInfoQueryResult{
		OperatorAddress:       validator.OperatorAddress,
		ConsPubKey:            btypes.MustConsensusPubKeyString(validator.ConsPubKey),
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
	OperatorAddress    btypes.ValAddress `json:"validator_address"`
	ConsPubKey         string            `json:"consensus_pubkey"`
	Fees               btypes.BigInt     `json:"fees"`
	CurrentTokens      btypes.BigInt     `json:"current_tokens"`
	CurrentPeriod      int64             `json:"current_period"`
	LastPeriod         int64             `json:"last_period"`
	LastPeriodFraction qtypes.Fraction   `json:"last_period_fraction"`
}

type DelegatorIncomeInfoQueryResult struct {
	OperatorAddress       btypes.ValAddress `json:"validator_address"`
	ConsPubKey            string            `json:"consensus_pubkey"`
	PreviousPeriod        int64             `json:"previous_validator_period"`
	BondToken             btypes.BigInt     `json:"bond_token"`
	CurrentStartingHeight int64             `json:"earns_starting_height"`
	FirstDelegateHeight   int64             `json:"first_delegate_height"`
	HistoricalRewardFees  btypes.BigInt     `json:"historical_rewards"`
	LastIncomeCalHeight   int64             `json:"last_income_calHeight"`
	LastIncomeCalFees     btypes.BigInt     `json:"last_income_calFees"`
}
