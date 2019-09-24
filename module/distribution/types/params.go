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
	ProposerRewardRate           qtypes.Dec `json:"proposer_reward_rate"`           // 块提议者奖励比例
	CommunityRewardRate          qtypes.Dec `json:"community_reward_rate"`          // 社区奖励比例
	DelegatorsIncomePeriodHeight int64      `json:"delegator_income_period_height"` // 奖励发放周期
	GasPerUnitCost               int64      `json:"gas_per_unit_cost"`              // 1QOS折算Gas量
}

// 设置单个参数，不同数据类型对应不同处理
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
		GasPerUnitCost:               qtypes.GasPerUnitCost,                 // 1000
	}
}

// 返回参数键值对
func (p *Params) KeyValuePairs() qtypes.KeyValuePairs {
	return qtypes.KeyValuePairs{
		{KeyProposerRewardRate, &p.ProposerRewardRate},
		{KeyCommunityRewardRate, &p.CommunityRewardRate},
		{KeyDelegatorsIncomePeriodHeight, &p.DelegatorsIncomePeriodHeight},
		{KeyGasPerUnitCost, &p.GasPerUnitCost},
	}
}

// 校验单个参数，返回参数值
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

// 参数所属模块名
func (p *Params) GetParamSpace() string {
	return ParamSpace
}

// 参数校验
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
