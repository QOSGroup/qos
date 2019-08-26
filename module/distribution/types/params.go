package types

import (
	"fmt"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/params"
	qtypes "github.com/QOSGroup/qos/types"
	"strconv"
)

var (
	ParamSpace = "distribution"

	// keys for distribution p
	KeyProposerRewardRate           = []byte("proposer_reward_rate")
	KeyCommunityRewardRate          = []byte("community_reward_rate")
	KeyDelegatorsIncomePeriodHeight = []byte("delegator_income_period_height")
	KeyGasPerUnitCost               = []byte("gas_per_unit_cost")
)

type Params struct {
	ProposerRewardRate           qtypes.Fraction `json:"proposer_reward_rate"`
	CommunityRewardRate          qtypes.Fraction `json:"community_reward_rate"`
	DelegatorsIncomePeriodHeight uint64          `json:"delegator_income_period_height"`
	GasPerUnitCost               uint64          `json:"gas_per_unit_cost"` // how much gas = 1 QOS
}

func DefaultParams() Params {
	return Params{
		ProposerRewardRate:           qtypes.NewFraction(int64(1), int64(100)),      // 1%
		CommunityRewardRate:          qtypes.NewFraction(int64(2), int64(100)),      // 2%
		DelegatorsIncomePeriodHeight: uint64(60 * 60 / qtypes.DefaultBlockInterval), // 1 hour
		GasPerUnitCost:               qtypes.GasPerUnitCost,
	}
}

func (p *Params) KeyValuePairs() qtypes.KeyValuePairs {
	return qtypes.KeyValuePairs{
		{KeyProposerRewardRate, &p.ProposerRewardRate},
		{KeyCommunityRewardRate, &p.CommunityRewardRate},
		{KeyDelegatorsIncomePeriodHeight, &p.DelegatorsIncomePeriodHeight},
		{KeyGasPerUnitCost, &p.GasPerUnitCost},
	}
}

func (p *Params) Validate(key string, value string) (interface{}, btypes.Error) {
	switch key {
	case string(KeyProposerRewardRate), string(KeyCommunityRewardRate):
		rate, err := qtypes.NewDecFromStr(value)
		if err != nil || rate.GTE(qtypes.OneDec()) || rate.LTE(qtypes.ZeroDec()) {
			return nil, params.ErrInvalidParam(fmt.Sprintf("%s invalid", key))
		}
		return rate, nil
	case string(KeyDelegatorsIncomePeriodHeight), string(KeyGasPerUnitCost):
		height, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, params.ErrInvalidParam(fmt.Sprintf("%s invalid", key))
		}
		return height, nil
	default:
		return nil, params.ErrInvalidParam(fmt.Sprintf("%s not exists", key))
	}
}

func (p *Params) GetParamSpace() string {
	return ParamSpace
}
