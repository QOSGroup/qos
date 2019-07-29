package types

import (
	"fmt"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/params"
	"strconv"
)

var (
	ParamSpace = "stake"

	// keys for stake p
	KeyMaxValidatorCnt             = []byte("max_validator_cnt")
	KeyValidatorVotingStatusLen    = []byte("voting_status_len")
	KeyValidatorVotingStatusLeast  = []byte("voting_status_least")
	KeyValidatorSurvivalSecs       = []byte("survival_secs")
	KeyDelegatorUnbondReturnHeight = []byte("unbond_return_height")
	KeyDelegatorRedelegationHeight = []byte("redelegation_height")
)

type Params struct {
	MaxValidatorCnt             uint32 `json:"max_validator_cnt"`
	ValidatorVotingStatusLen    uint32 `json:"voting_status_len"`
	ValidatorVotingStatusLeast  uint32 `json:"voting_status_least"`
	ValidatorSurvivalSecs       uint32 `json:"survival_secs"`
	DelegatorUnbondReturnHeight uint32 `json:"unbond_return_height"`
	DelegatorRedelegationHeight uint32 `json:"redelegation_height"`
}

func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{KeyMaxValidatorCnt, &p.MaxValidatorCnt},
		{KeyValidatorVotingStatusLen, &p.ValidatorVotingStatusLen},
		{KeyValidatorVotingStatusLeast, &p.ValidatorVotingStatusLeast},
		{KeyValidatorSurvivalSecs, &p.ValidatorSurvivalSecs},
		{KeyDelegatorUnbondReturnHeight, &p.DelegatorUnbondReturnHeight},
		{KeyDelegatorRedelegationHeight, &p.DelegatorRedelegationHeight},
	}
}

func (p *Params) Validate(key string, value string) (interface{}, btypes.Error) {
	switch key {
	case string(KeyMaxValidatorCnt),
		string(KeyValidatorVotingStatusLen),
		string(KeyValidatorVotingStatusLeast),
		string(KeyValidatorSurvivalSecs),
		string(KeyDelegatorUnbondReturnHeight),
		string(KeyDelegatorRedelegationHeight):
		v, err := strconv.ParseUint(value, 10, 32)
		if err != nil || v <= 0 {
			return nil, params.ErrInvalidParam(fmt.Sprintf("%s invalid", key))
		}
		return v, nil
	default:
		return nil, params.ErrInvalidParam(fmt.Sprintf("%s not exists", key))
	}
}

func (p *Params) GetParamSpace() string {
	return ParamSpace
}

func NewParams(maxValidatorCnt, validatorVotingStatusLen, validatorVotingStatusLeast, validatorSurvivalSecs, delegatorUnbondReturnHeight uint32, delegatorRedelegationHeight uint32) Params {

	return Params{
		MaxValidatorCnt:             maxValidatorCnt,
		ValidatorVotingStatusLen:    validatorVotingStatusLen,
		ValidatorVotingStatusLeast:  validatorVotingStatusLeast,
		ValidatorSurvivalSecs:       validatorSurvivalSecs,
		DelegatorUnbondReturnHeight: delegatorUnbondReturnHeight,
		DelegatorRedelegationHeight: delegatorRedelegationHeight,
	}
}

func DefaultParams() Params {
	return NewParams(10, 100, 50, 600, 10, 10)
}
