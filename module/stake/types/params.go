package types

type Params struct {
	MaxValidatorCnt            uint32 `json:"max_validator_cnt"`
	ValidatorVotingStatusLen   uint32 `json:"voting_status_len"`
	ValidatorVotingStatusLeast uint32 `json:"voting_status_least"`
	ValidatorSurvivalSecs      uint32 `json:"survival_secs"`
}

func NewParams(maxValidatorCnt uint32, validatorVotingStatusLen uint32, validatorVotingStatusLeast uint32, validatorSurvivalSecs uint32) Params {

	return Params{
		MaxValidatorCnt:            maxValidatorCnt,
		ValidatorVotingStatusLen:   validatorVotingStatusLen,
		ValidatorVotingStatusLeast: validatorVotingStatusLeast,
		ValidatorSurvivalSecs:      validatorSurvivalSecs,
	}
}

func DefaultParams() Params {

	return NewParams(10, 100, 50, 600)
}
