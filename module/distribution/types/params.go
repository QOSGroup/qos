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
	ProposerRewardRate           qtypes.Dec `json:"proposer_reward_rate"`
	CommunityRewardRate          qtypes.Dec `json:"community_reward_rate"`
	DelegatorsIncomePeriodHeight int64      `json:"delegator_income_period_height"`
	GasPerUnitCost               int64      `json:"gas_per_unit_cost"` // how much gas = 1 QOS
}

func (p *Params) SetKeyValue(key string, value interface{}) btypes.Error {
	switch key {
	case string(KeyProposerRewardRate):
		p.ProposerRewardRate = value.(qtypes.Dec)
		break
	case string(KeyCommunityRewardRate):
		p.CommunityRewardRate = value.(qtypes.Dec)
		break
	case string(KeyDelegatorsIncomePeriodHeight):
		p.DelegatorsIncomePeriodHeight = value.(int64)
		break
	case string(KeyGasPerUnitCost):
		p.GasPerUnitCost = value.(int64)
		break
	default:
		return params.ErrInvalidParam(fmt.Sprintf("%s not exists", key))
	}

	return nil
}

var _ qtypes.ParamSet = (*Params)(nil)

func DefaultParams() Params {
	return Params{
		ProposerRewardRate:           qtypes.NewDecWithPrec(1, 2),           // 1%
		CommunityRewardRate:          qtypes.NewDecWithPrec(2, 2),           // 2%
		DelegatorsIncomePeriodHeight: 60 * 60 / qtypes.DefaultBlockInterval, // 1 hour
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

func (p *Params) ValidateKeyValue(key string, value string) (interface{}, btypes.Error) {
	switch key {
	case string(KeyProposerRewardRate), string(KeyCommunityRewardRate):
		rate, err := qtypes.NewDecFromStr(value)
		if err != nil || rate.GTE(qtypes.OneDec()) || rate.LTE(qtypes.ZeroDec()) {
			return nil, params.ErrInvalidParam(fmt.Sprintf("%s invalid", key))
		}
		return rate, nil
	case string(KeyDelegatorsIncomePeriodHeight), string(KeyGasPerUnitCost):
		height, err := strconv.ParseInt(value, 10, 64)
		if err != nil || height <= 0 {
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

func (p *Params) Validate() btypes.Error {
	if p.ProposerRewardRate.GT(qtypes.OneDec()) || p.ProposerRewardRate.IsNegative() {
		params.ErrInvalidParam(fmt.Sprintf("%s must gte 0 and lte 1", KeyProposerRewardRate))
	}
	if p.CommunityRewardRate.GT(qtypes.OneDec()) || p.CommunityRewardRate.IsNegative() {
		params.ErrInvalidParam(fmt.Sprintf("%s must gte 0 and lte 1", KeyCommunityRewardRate))
	}
	if p.DelegatorsIncomePeriodHeight <= 0 {
		params.ErrInvalidParam(fmt.Sprintf("%s must gt 0", KeyDelegatorsIncomePeriodHeight))
	}
	if p.GasPerUnitCost <= 0 {
		params.ErrInvalidParam(fmt.Sprintf("%s must gt 0", KeyGasPerUnitCost))
	}

	return nil
}
