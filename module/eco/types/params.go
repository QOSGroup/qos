package types

import (
	qtypes "github.com/QOSGroup/qos/types"
)

type DistributionParams struct {
	ProposerRewardRate           qtypes.Fraction `json:"proposer_reward_rate"`
	CommunityRewardRate          qtypes.Fraction `json:"community_reward_rate"`
	ValidatorCommissionRate      qtypes.Fraction `json:"validator_commission_rate"`
	DelegatorsIncomePeriodHeight uint64          `json:"delegator_income_period_height"`
}

type StakeParams struct {
	MaxValidatorCnt            uint32 `json:"max_validator_cnt"`
	ValidatorVotingStatusLen   uint32 `json:"voting_status_len"`
	ValidatorVotingStatusLeast uint32 `json:"voting_status_least"`
	ValidatorSurvivalSecs      uint32 `json:"survival_secs"`
}

func DefaultDistributionParams() DistributionParams {
	return DistributionParams{
		//todo
	}
}

func NewStakeParams(maxValidatorCnt uint32, validatorVotingStatusLen uint32, validatorVotingStatusLeast uint32, validatorSurvivalSecs uint32) StakeParams {

	return StakeParams{
		MaxValidatorCnt:            maxValidatorCnt,
		ValidatorVotingStatusLen:   validatorVotingStatusLen,
		ValidatorVotingStatusLeast: validatorVotingStatusLeast,
		ValidatorSurvivalSecs:      validatorSurvivalSecs,
	}
}

func DefaultStakeParams() StakeParams {

	return NewStakeParams(10, 100, 50, 600)
}
