package types

type SPOConfig struct {
	TotalAmount uint64 `json:"total_amount"`
	TotalBlock  uint64 `json:"total_block"`
}

func DefaultSPOConfig() SPOConfig {
	return SPOConfig{
		100e8,
		6307200,
	}
}

type StakeConfig struct {
	MaxValidatorCnt            uint32 `json:"max_validator_cnt"`
	ValidatorVotingStatusLen   uint32 `json:"voting_status_len"`
	ValidatorVotingStatusLeast uint32 `json:"voting_status_least"`
	ValidatorSurvivalSecs      uint64 `json:"survival_secs"`
}

func DefaultStakeConfig() StakeConfig {
	return StakeConfig{
		10000,
		10000,
		5000,
		86400,
	}
}
