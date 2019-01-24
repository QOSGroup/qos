package types

import (
	"time"

	qtypes "github.com/QOSGroup/qos/types"
)

type DistributionParams struct {
	ProposerRewardRate           qtypes.Fraction `json:"proposer_reward_rate"`
	CommunityRewardRate          qtypes.Fraction `json:"community_reward_rate"`
	ValidatorCommissionRate      qtypes.Fraction `json:"validator_commission_rate"`
	DelegatorsIncomePeriodHeight uint64          `json:"delegator_income_period_height"`
}

type StakeParams struct {
	MaxValidatorCnt                 uint32 `json:"max_validator_cnt"`
	ValidatorVotingStatusLen        uint32 `json:"voting_status_len"`
	ValidatorVotingStatusLeast      uint32 `json:"voting_status_least"`
	ValidatorSurvivalSecs           uint32 `json:"survival_secs"`
	DelegatorUnbondDistributeHeight uint32 `json:"unbond_distribute_height"`
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
	}
}

func NewStakeParams(maxValidatorCnt, validatorVotingStatusLen, validatorVotingStatusLeast, validatorSurvivalSecs, delegatorUnbondDistributeHeight uint32) StakeParams {

	return StakeParams{
		MaxValidatorCnt:                 maxValidatorCnt,
		ValidatorVotingStatusLen:        validatorVotingStatusLen,
		ValidatorVotingStatusLeast:      validatorVotingStatusLeast,
		ValidatorSurvivalSecs:           validatorSurvivalSecs,
		DelegatorUnbondDistributeHeight: delegatorUnbondDistributeHeight,
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
				2.5e8,
				0,
			},
			InflationPhrase{
				time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC),
				12.75e8,
				0,
			},
			InflationPhrase{
				time.Date(2031, 1, 1, 0, 0, 0, 0, time.UTC),
				6.375e8,
				0,
			},
			InflationPhrase{
				time.Date(2035, 1, 1, 0, 0, 0, 0, time.UTC),
				3.185e8,
				0,
			},
		},
	)
}
