package types

import (
	"fmt"
	btypes "github.com/QOSGroup/qbase/types"
	perr "github.com/QOSGroup/qos/module/params"
	ptypes "github.com/QOSGroup/qos/module/params/types"
	"strconv"
	"time"

	qtypes "github.com/QOSGroup/qos/types"
)

var (
	DistributionParamSpace = "distribution"

	// keys for distribution params
	KeyProposerRewardRate           = []byte("proposer_reward_rate")
	KeyCommunityRewardRate          = []byte("community_reward_rate")
	KeyValidatorCommissionRate      = []byte("validator_commission_rate")
	KeyDelegatorsIncomePeriodHeight = []byte("delegator_income_period_height")
	KeyGasPerUnitCost               = []byte("gas_per_unit_cost")

	StakeParamSpace = "stake"

	// keys for stake params
	KeyMaxValidatorCnt             = []byte("max_validator_cnt")
	KeyValidatorVotingStatusLen    = []byte("voting_status_len")
	KeyValidatorVotingStatusLeast  = []byte("voting_status_least")
	KeyValidatorSurvivalSecs       = []byte("survival_secs")
	KeyDelegatorUnbondReturnHeight = []byte("unbond_return_height")
)

type DistributionParams struct {
	ProposerRewardRate           qtypes.Fraction `json:"proposer_reward_rate"`
	CommunityRewardRate          qtypes.Fraction `json:"community_reward_rate"`
	ValidatorCommissionRate      qtypes.Fraction `json:"validator_commission_rate"`
	DelegatorsIncomePeriodHeight uint64          `json:"delegator_income_period_height"`
	GasPerUnitCost               uint64          `json:"gas_per_unit_cost"` // how much gas = 1 QOS
}

func (params *DistributionParams) KeyValuePairs() ptypes.KeyValuePairs {
	return ptypes.KeyValuePairs{
		{KeyProposerRewardRate, &params.ProposerRewardRate},
		{KeyCommunityRewardRate, &params.CommunityRewardRate},
		{KeyValidatorCommissionRate, &params.ValidatorCommissionRate},
		{KeyDelegatorsIncomePeriodHeight, &params.DelegatorsIncomePeriodHeight},
		{KeyGasPerUnitCost, &params.GasPerUnitCost},
	}
}

func (params *DistributionParams) Validate(key string, value string) (interface{}, btypes.Error) {
	switch key {
	case string(KeyProposerRewardRate), string(KeyCommunityRewardRate), string(KeyValidatorCommissionRate):
		rate, err := qtypes.NewDecFromStr(value)
		if err != nil || rate.GTE(qtypes.OneDec()) || rate.LTE(qtypes.ZeroDec()) {
			return nil, perr.ErrInvalidParam(fmt.Sprintf("%s invalid", key))
		}
		return rate, nil
	case string(KeyDelegatorsIncomePeriodHeight), string(KeyGasPerUnitCost):
		height, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, perr.ErrInvalidParam(fmt.Sprintf("%s invalid", key))
		}
		return height, nil
	default:
		return nil, perr.ErrInvalidParam(fmt.Sprintf("%s not exists", key))
	}
}

func (params *DistributionParams) GetParamSpace() string {
	return DistributionParamSpace
}

type StakeParams struct {
	MaxValidatorCnt             uint32 `json:"max_validator_cnt"`
	ValidatorVotingStatusLen    uint32 `json:"voting_status_len"`
	ValidatorVotingStatusLeast  uint32 `json:"voting_status_least"`
	ValidatorSurvivalSecs       uint32 `json:"survival_secs"`
	DelegatorUnbondReturnHeight uint32 `json:"unbond_return_height"`
}

func (params *StakeParams) KeyValuePairs() ptypes.KeyValuePairs {
	return ptypes.KeyValuePairs{
		{KeyMaxValidatorCnt, &params.MaxValidatorCnt},
		{KeyValidatorVotingStatusLen, &params.ValidatorVotingStatusLen},
		{KeyValidatorVotingStatusLeast, &params.ValidatorVotingStatusLeast},
		{KeyValidatorSurvivalSecs, &params.ValidatorSurvivalSecs},
		{KeyDelegatorUnbondReturnHeight, &params.DelegatorUnbondReturnHeight},
	}
}

func (params *StakeParams) Validate(key string, value string) (interface{}, btypes.Error) {
	switch key {
	case string(KeyMaxValidatorCnt),
		string(KeyValidatorVotingStatusLen),
		string(KeyValidatorVotingStatusLeast),
		string(KeyValidatorSurvivalSecs),
		string(KeyDelegatorUnbondReturnHeight):
		v, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return nil, perr.ErrInvalidParam(fmt.Sprintf("%s invalid", key))
		}
		return v, nil
	default:
		return nil, perr.ErrInvalidParam(fmt.Sprintf("%s not exists", key))
	}
}

func (params *StakeParams) GetParamSpace() string {
	return StakeParamSpace
}

type MintParams struct {
	Phrases []InflationPhrase `json:"inflation_phrases"`
}

type InflationPhrase struct {
	EndTime       time.Time `json:"endtime"`
	TotalAmount   uint64    `json:"total_amount"`
	AppliedAmount uint64    `json:"applied_amount"`
}

func DefaultDistributionParams() DistributionParams {
	return DistributionParams{
		ProposerRewardRate:           qtypes.NewFraction(int64(4), int64(100)), // 4%
		CommunityRewardRate:          qtypes.NewFraction(int64(1), int64(100)), // 1%
		ValidatorCommissionRate:      qtypes.NewFraction(int64(1), int64(100)), // 1%
		DelegatorsIncomePeriodHeight: uint64(10),
		GasPerUnitCost:               uint64(10),
	}
}

func NewStakeParams(maxValidatorCnt, validatorVotingStatusLen, validatorVotingStatusLeast, validatorSurvivalSecs, delegatorUnbondReturnHeight uint32) StakeParams {

	return StakeParams{
		MaxValidatorCnt:             maxValidatorCnt,
		ValidatorVotingStatusLen:    validatorVotingStatusLen,
		ValidatorVotingStatusLeast:  validatorVotingStatusLeast,
		ValidatorSurvivalSecs:       validatorSurvivalSecs,
		DelegatorUnbondReturnHeight: delegatorUnbondReturnHeight,
	}
}

func DefaultStakeParams() StakeParams {
	return NewStakeParams(10, 100, 50, 600, 10)
}

func NewMintParams(phrases []InflationPhrase) MintParams {
	return MintParams{phrases}
}

func DefaultMintParams() MintParams {
	return NewMintParams(
		[]InflationPhrase{
			InflationPhrase{
				time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				2.5e12, //mul(10^4),
				0,
			},
			InflationPhrase{
				time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC),
				12.75e12, //mul(10^4),
				0,
			},
			InflationPhrase{
				time.Date(2031, 1, 1, 0, 0, 0, 0, time.UTC),
				6.375e12, //mul(10^4),
				0,
			},
			InflationPhrase{
				time.Date(2035, 1, 1, 0, 0, 0, 0, time.UTC),
				3.185e12, //mul(10^4),
				0,
			},
		},
	)
}
