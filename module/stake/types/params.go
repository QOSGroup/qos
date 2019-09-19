package types

import (
	"fmt"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/params"
	"github.com/QOSGroup/qos/types"
	qtypes "github.com/QOSGroup/qos/types"
	"strconv"
	"time"
)

var (
	defaultMaxValidatorCnt             = int64(21)
	defaultValidatorVotingStatusLen    = int64(10000)
	defaultValidatorVotingStatusLeast  = int64(7000)
	defaultValidatorSurvivalSecs       = int64(8 * 60 * 60)
	defaultDelegatorUnbondReturnHeight = int64(15 * 24 * 60 * 60 / qtypes.DefaultBlockInterval) // 15 days
	defaultMaxEvidenceAge              = time.Duration(1814400000000000)                        // ~= 21 days
	defaultSlashFractionDoubleSign     = types.NewDecWithPrec(2, 1)                             // 0.2
	defaultSlashFractionDowntime       = types.NewDecWithPrec(1, 4)                             // 0.0001
)

var (
	ParamSpace = "stake"

	// keys for stake p
	KeyMaxValidatorCnt             = []byte("max_validator_cnt")
	KeyValidatorVotingStatusLen    = []byte("voting_status_len")
	KeyValidatorVotingStatusLeast  = []byte("voting_status_least")
	KeyValidatorSurvivalSecs       = []byte("survival_secs")
	KeyDelegatorUnbondFrozenHeight = []byte("unbond_frozen_height")
	KeyMaxEvidenceAge              = []byte("max_evidence_age")
	KeySlashFractionDoubleSign     = []byte("slash_fraction_double_sign")
	KeySlashFractionDowntime       = []byte("slash_fraction_downtime")
)

type Params struct {
	MaxValidatorCnt             int64         `json:"max_validator_cnt"`
	ValidatorVotingStatusLen    int64         `json:"voting_status_len"`
	ValidatorVotingStatusLeast  int64         `json:"voting_status_least"`
	ValidatorSurvivalSecs       int64         `json:"survival_secs"`
	DelegatorUnbondFrozenHeight int64         `json:"unbond_frozen_height"`
	MaxEvidenceAge              time.Duration `json:"max_evidence_age"`
	SlashFractionDoubleSign     types.Dec     `json:"slash_fraction_double_sign"`
	SlashFractionDowntime       types.Dec     `json:"slash_fraction_downtime"`
}

func (p *Params) SetKeyValue(key string, value interface{}) btypes.Error {
	switch key {
	case string(KeyMaxValidatorCnt):
		p.MaxValidatorCnt = value.(int64)
		break
	case string(KeyValidatorVotingStatusLen):
		p.ValidatorVotingStatusLen = value.(int64)
		break
	case string(KeyValidatorVotingStatusLeast):
		p.ValidatorVotingStatusLeast = value.(int64)
		break
	case string(KeyValidatorSurvivalSecs):
		p.ValidatorSurvivalSecs = value.(int64)
		break
	case string(KeyDelegatorUnbondFrozenHeight):
		p.DelegatorUnbondFrozenHeight = value.(int64)
		break
	case string(KeyMaxEvidenceAge):
		p.MaxEvidenceAge = value.(time.Duration)
		break
	case string(KeySlashFractionDoubleSign):
		p.SlashFractionDoubleSign = value.(qtypes.Dec)
		break
	case string(KeySlashFractionDowntime):
		p.SlashFractionDowntime = value.(qtypes.Dec)
		break
	default:
		return params.ErrInvalidParam(fmt.Sprintf("%s not exists", key))
	}

	return nil
}

var _ qtypes.ParamSet = (*Params)(nil)

func (p *Params) KeyValuePairs() qtypes.KeyValuePairs {
	return qtypes.KeyValuePairs{
		{KeyMaxValidatorCnt, &p.MaxValidatorCnt},
		{KeyValidatorVotingStatusLen, &p.ValidatorVotingStatusLen},
		{KeyValidatorVotingStatusLeast, &p.ValidatorVotingStatusLeast},
		{KeyValidatorSurvivalSecs, &p.ValidatorSurvivalSecs},
		{KeyDelegatorUnbondFrozenHeight, &p.DelegatorUnbondFrozenHeight},
		{KeyMaxEvidenceAge, &p.MaxEvidenceAge},
		{KeySlashFractionDoubleSign, &p.SlashFractionDoubleSign},
		{KeySlashFractionDowntime, &p.SlashFractionDowntime},
	}
}

func (p *Params) ValidateKeyValue(key string, value string) (interface{}, btypes.Error) {
	switch key {
	case string(KeyMaxValidatorCnt),
		string(KeyValidatorVotingStatusLen),
		string(KeyValidatorVotingStatusLeast),
		string(KeyValidatorSurvivalSecs),
		string(KeyDelegatorUnbondFrozenHeight):
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil || v <= 0 {
			return nil, params.ErrInvalidParam(fmt.Sprintf("%s invalid", key))
		}
		return v, nil
	case string(KeySlashFractionDoubleSign),
		string(KeySlashFractionDowntime):
		v, err := qtypes.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidParam(fmt.Sprintf("%s invalid", key))
		}
		return v, nil
	case string(KeyMaxEvidenceAge):
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil || v < 0 {
			return nil, params.ErrInvalidParam(fmt.Sprintf("%s invalid", key))
		}
		return time.Duration(v), nil
	default:
		return nil, params.ErrInvalidParam(fmt.Sprintf("%s not exists", key))
	}
}

func (p *Params) GetParamSpace() string {
	return ParamSpace
}

func NewParams(maxValidatorCnt, validatorVotingStatusLen, validatorVotingStatusLeast, validatorSurvivalSecs,
	delegatorUnbondFrozenHeight int64, maxEvidenceAge time.Duration,
	slashFractionDoubleSign types.Dec, slashFractionDowntime types.Dec) Params {

	return Params{
		MaxValidatorCnt:             maxValidatorCnt,
		ValidatorVotingStatusLen:    validatorVotingStatusLen,
		ValidatorVotingStatusLeast:  validatorVotingStatusLeast,
		ValidatorSurvivalSecs:       validatorSurvivalSecs,
		DelegatorUnbondFrozenHeight: delegatorUnbondFrozenHeight,
		MaxEvidenceAge:              maxEvidenceAge,
		SlashFractionDoubleSign:     slashFractionDoubleSign,
		SlashFractionDowntime:       slashFractionDowntime,
	}
}

func DefaultParams() Params {
	return NewParams(defaultMaxValidatorCnt, defaultValidatorVotingStatusLen, defaultValidatorVotingStatusLeast, defaultValidatorSurvivalSecs, defaultDelegatorUnbondReturnHeight, defaultMaxEvidenceAge, defaultSlashFractionDoubleSign, defaultSlashFractionDowntime)
}

func (p *Params) Validate() btypes.Error {
	if p.MaxValidatorCnt <= 0 {
		return params.ErrInvalidParam(fmt.Sprintf("%s must gt 0", KeyMaxValidatorCnt))
	}
	if p.ValidatorVotingStatusLen <= 0 {
		return params.ErrInvalidParam(fmt.Sprintf("%s must gt 0", KeyValidatorVotingStatusLen))
	}
	if p.ValidatorVotingStatusLeast <= 0 {
		return params.ErrInvalidParam(fmt.Sprintf("%s must gt 0", KeyValidatorVotingStatusLeast))
	}
	if p.ValidatorVotingStatusLeast > p.ValidatorVotingStatusLen {
		return params.ErrInvalidParam(fmt.Sprintf("%s must gte %s", KeyValidatorVotingStatusLen, KeyValidatorVotingStatusLeast))
	}
	if p.ValidatorSurvivalSecs <= 0 {
		return params.ErrInvalidParam(fmt.Sprintf("%s must gt 0", KeyValidatorSurvivalSecs))
	}
	if p.DelegatorUnbondFrozenHeight <= 0 {
		return params.ErrInvalidParam(fmt.Sprintf("%s must gt 0", KeyDelegatorUnbondFrozenHeight))
	}
	if p.SlashFractionDoubleSign.IsNegative() || p.SlashFractionDoubleSign.GT(qtypes.OneDec()) {
		return params.ErrInvalidParam(fmt.Sprintf("%s must gte 0 and lte 1", KeySlashFractionDoubleSign))
	}
	if p.SlashFractionDowntime.IsNegative() || p.SlashFractionDowntime.GT(qtypes.OneDec()) {
		return params.ErrInvalidParam(fmt.Sprintf("%s must gte 0 and lte 1", KeySlashFractionDowntime))
	}
	if p.MaxEvidenceAge < 0 {
		return params.ErrInvalidParam(fmt.Sprintf("%s must gt 0", KeyMaxEvidenceAge))
	}

	return nil
}
