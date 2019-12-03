package types

import (
	"fmt"
	"github.com/QOSGroup/qos/module/params"
	"strconv"
	"time"

	btypes "github.com/QOSGroup/qbase/types"
	qtypes "github.com/QOSGroup/qos/types"
)

var (
	ParamSpace = "gov"

	KeyNormalMinDeposit             = []byte("normal_min_deposit")
	KeyNormalMinProposerDepositRate = []byte("normal_min_proposer_deposit_rate")
	KeyNormalMaxDepositPeriod       = []byte("normal_max_deposit_period")
	KeyNormalVotingPeriod           = []byte("normal_voting_period")
	KeyNormalQuorum                 = []byte("normal_quorum")
	KeyNormalThreshold              = []byte("normal_threshold")
	KeyNormalVeto                   = []byte("normal_veto")
	KeyNormalPenalty                = []byte("normal_penalty")
	KeyNormalBurnRate               = []byte("normal_burn_rate")

	KeyImportantMinDeposit             = []byte("important_min_deposit")
	KeyImportantMinProposerDepositRate = []byte("important_min_proposer_deposit_rate")
	KeyImportantMaxDepositPeriod       = []byte("important_max_deposit_period")
	KeyImportantVotingPeriod           = []byte("important_voting_period")
	KeyImportantQuorum                 = []byte("important_quorum")
	KeyImportantThreshold              = []byte("important_threshold")
	KeyImportantVeto                   = []byte("important_veto")
	KeyImportantPenalty                = []byte("important_penalty")
	KeyImportantBurnRate               = []byte("important_burn_rate")

	KeyCriticalMinDeposit             = []byte("critical_min_deposit")
	KeyCriticalMinProposerDepositRate = []byte("critical_min_proposer_deposit_rate")
	KeyCriticalMaxDepositPeriod       = []byte("critical_max_deposit_period")
	KeyCriticalVotingPeriod           = []byte("critical_voting_period")
	KeyCriticalQuorum                 = []byte("critical_quorum")
	KeyCriticalThreshold              = []byte("critical_threshold")
	KeyCriticalVeto                   = []byte("critical_veto")
	KeyCriticalPenalty                = []byte("critical_penalty")
	KeyCriticalBurnRate               = []byte("critical_burn_rate")
)

// Params for governance, there are three levels: normal,important,critical.
type Params struct {
	// params of normal level
	// DepositParams
	NormalMinDeposit             btypes.BigInt `json:"normal_min_deposit"`               //  Minimum deposit for a proposal to enter voting period.
	NormalMinProposerDepositRate qtypes.Dec    `json:"normal_min_proposer_deposit_rate"` //  Minimum deposit rate for proposer to submit a proposal.
	NormalMaxDepositPeriod       time.Duration `json:"normal_max_deposit_period"`        //  Maximum period for Atom holders to deposit on a proposal.
	// VotingParams
	NormalVotingPeriod time.Duration `json:"normal_voting_period"` //  Length of the voting period.
	// TallyParams
	NormalQuorum    qtypes.Dec `json:"normal_quorum"`    //  Minimum percentage of total stake needed to vote for a result to be considered valid
	NormalThreshold qtypes.Dec `json:"normal_threshold"` //  Minimum propotion of Yes votes for proposal to pass.
	NormalVeto      qtypes.Dec `json:"normal_veto"`      //  Minimum value of Veto votes to Total votes ratio for proposal to be vetoed.
	NormalPenalty   qtypes.Dec `json:"normal_penalty"`   //  Penalty if validator does not vote
	// BurnRate
	NormalBurnRate qtypes.Dec `json:"normal_burn_rate"` // Deposit burning rate when proposals pass or reject.

	// params of important level
	// DepositParams
	ImportantMinDeposit             btypes.BigInt `json:"important_min_deposit"`               //  Minimum deposit for a proposal to enter voting period.
	ImportantMinProposerDepositRate qtypes.Dec    `json:"important_min_proposer_deposit_rate"` //  Minimum deposit rate for proposer to submit a proposal.
	ImportantMaxDepositPeriod       time.Duration `json:"important_max_deposit_period"`        //  Maximum period for Atom holders to deposit on a proposal.
	// VotingParams
	ImportantVotingPeriod time.Duration `json:"important_voting_period"` //  Length of the voting period.
	// TallyParams
	ImportantQuorum    qtypes.Dec `json:"important_quorum"`    //  Minimum percentage of total stake needed to vote for a result to be considered valid
	ImportantThreshold qtypes.Dec `json:"important_threshold"` //  Minimum propotion of Yes votes for proposal to pass.
	ImportantVeto      qtypes.Dec `json:"important_veto"`      //  Minimum value of Veto votes to Total votes ratio for proposal to be vetoed.
	ImportantPenalty   qtypes.Dec `json:"important_penalty"`   //  Penalty if validator does not vote
	// BurnRate
	ImportantBurnRate qtypes.Dec `json:"important_burn_rate"` // Deposit burning rate when proposals pass or reject.

	// params of critical level
	// DepositParams
	CriticalMinDeposit             btypes.BigInt `json:"critical_min_deposit"`               //  Minimum deposit for a proposal to enter voting period.
	CriticalMinProposerDepositRate qtypes.Dec    `json:"critical_min_proposer_deposit_rate"` //  Minimum deposit rate for proposer to submit a proposal.
	CriticalMaxDepositPeriod       time.Duration `json:"critical_max_deposit_period"`        //  Maximum period for Atom holders to deposit on a proposal.
	// VotingParams
	CriticalVotingPeriod time.Duration `json:"critical_voting_period"` //  Length of the voting period.
	// TallyParams
	CriticalQuorum    qtypes.Dec `json:"critical_quorum"`    //  Minimum percentage of total stake needed to vote for a result to be considered valid
	CriticalThreshold qtypes.Dec `json:"critical_threshold"` //  Minimum propotion of Yes votes for proposal to pass.
	CriticalVeto      qtypes.Dec `json:"critical_veto"`      //  Minimum value of Veto votes to Total votes ratio for proposal to be vetoed.
	CriticalPenalty   qtypes.Dec `json:"critical_penalty"`   //  Penalty if validator does not vote
	// BurnRate
	CriticalBurnRate qtypes.Dec `json:"critical_burn_rate"` // Deposit burning rate when proposals pass or reject.
}

// Set new value
func (p *Params) SetKeyValue(key string, value interface{}) btypes.Error {
	switch key {
	case string(KeyNormalMinDeposit):
		p.NormalMinDeposit = value.(btypes.BigInt)
		break
	case string(KeyNormalMinProposerDepositRate):
		p.NormalMinProposerDepositRate = value.(qtypes.Dec)
		break
	case string(KeyNormalMaxDepositPeriod):
		p.NormalMaxDepositPeriod = value.(time.Duration)
		break
	case string(KeyNormalVotingPeriod):
		p.NormalVotingPeriod = value.(time.Duration)
		break
	case string(KeyNormalQuorum):
		p.NormalQuorum = value.(qtypes.Dec)
		break
	case string(KeyNormalThreshold):
		p.NormalThreshold = value.(qtypes.Dec)
		break
	case string(KeyNormalVeto):
		p.NormalVeto = value.(qtypes.Dec)
		break
	case string(KeyNormalPenalty):
		p.NormalPenalty = value.(qtypes.Dec)
		break
	case string(KeyNormalBurnRate):
		p.NormalBurnRate = value.(qtypes.Dec)
		break
	case string(KeyImportantMinDeposit):
		p.ImportantMinDeposit = value.(btypes.BigInt)
		break
	case string(KeyImportantMinProposerDepositRate):
		p.ImportantMinProposerDepositRate = value.(qtypes.Dec)
		break
	case string(KeyImportantMaxDepositPeriod):
		p.ImportantMaxDepositPeriod = value.(time.Duration)
		break
	case string(KeyImportantVotingPeriod):
		p.ImportantVotingPeriod = value.(time.Duration)
		break
	case string(KeyImportantQuorum):
		p.ImportantQuorum = value.(qtypes.Dec)
		break
	case string(KeyImportantThreshold):
		p.ImportantThreshold = value.(qtypes.Dec)
		break
	case string(KeyImportantVeto):
		p.ImportantVeto = value.(qtypes.Dec)
		break
	case string(KeyImportantPenalty):
		p.ImportantPenalty = value.(qtypes.Dec)
		break
	case string(KeyImportantBurnRate):
		p.ImportantBurnRate = value.(qtypes.Dec)
		break
	case string(KeyCriticalMinDeposit):
		p.CriticalMinDeposit = value.(btypes.BigInt)
		break
	case string(KeyCriticalMinProposerDepositRate):
		p.CriticalMinProposerDepositRate = value.(qtypes.Dec)
		break
	case string(KeyCriticalMaxDepositPeriod):
		p.CriticalMaxDepositPeriod = value.(time.Duration)
		break
	case string(KeyCriticalVotingPeriod):
		p.CriticalVotingPeriod = value.(time.Duration)
		break
	case string(KeyCriticalQuorum):
		p.CriticalQuorum = value.(qtypes.Dec)
		break
	case string(KeyCriticalThreshold):
		p.CriticalThreshold = value.(qtypes.Dec)
		break
	case string(KeyCriticalVeto):
		p.CriticalVeto = value.(qtypes.Dec)
		break
	case string(KeyCriticalPenalty):
		p.CriticalPenalty = value.(qtypes.Dec)
		break
	case string(KeyCriticalBurnRate):
		p.CriticalBurnRate = value.(qtypes.Dec)
		break
	default:
		return params.ErrInvalidParam(fmt.Sprintf("%s not exists", key))
	}

	return nil
}

var _ qtypes.ParamSet = (*Params)(nil)

func DefaultParams() Params {
	return Params{
		// normal level
		NormalMinDeposit:             btypes.NewInt(100000),
		NormalMinProposerDepositRate: qtypes.NewDecWithPrec(334, 3),
		NormalMaxDepositPeriod:       DefaultDepositPeriod,
		NormalVotingPeriod:           DefaultVotingPeriod,
		NormalQuorum:                 qtypes.NewDecWithPrec(334, 3),
		NormalThreshold:              qtypes.NewDecWithPrec(5, 1),
		NormalVeto:                   qtypes.NewDecWithPrec(334, 3),
		NormalPenalty:                qtypes.ZeroDec(),
		NormalBurnRate:               qtypes.NewDecWithPrec(2, 1),
		// important level
		ImportantMinDeposit:             btypes.NewInt(500000),
		ImportantMinProposerDepositRate: qtypes.NewDecWithPrec(334, 3),
		ImportantMaxDepositPeriod:       DefaultDepositPeriod,
		ImportantVotingPeriod:           DefaultVotingPeriod,
		ImportantQuorum:                 qtypes.NewDecWithPrec(334, 3),
		ImportantThreshold:              qtypes.NewDecWithPrec(5, 1),
		ImportantVeto:                   qtypes.NewDecWithPrec(334, 3),
		ImportantPenalty:                qtypes.ZeroDec(),
		ImportantBurnRate:               qtypes.NewDecWithPrec(2, 1),
		// critical level
		CriticalMinDeposit:             btypes.NewInt(1000000),
		CriticalMinProposerDepositRate: qtypes.NewDecWithPrec(334, 3),
		CriticalMaxDepositPeriod:       DefaultDepositPeriod,
		CriticalVotingPeriod:           DefaultVotingPeriod,
		CriticalQuorum:                 qtypes.NewDecWithPrec(334, 3),
		CriticalThreshold:              qtypes.NewDecWithPrec(5, 1),
		CriticalVeto:                   qtypes.NewDecWithPrec(334, 3),
		CriticalPenalty:                qtypes.ZeroDec(),
		CriticalBurnRate:               qtypes.NewDecWithPrec(2, 1),
	}
}

func (p *Params) KeyValuePairs() qtypes.KeyValuePairs {
	return qtypes.KeyValuePairs{
		{KeyNormalMinDeposit, &p.NormalMinDeposit},
		{KeyNormalMinProposerDepositRate, &p.NormalMinProposerDepositRate},
		{KeyNormalMaxDepositPeriod, &p.NormalMaxDepositPeriod},
		{KeyNormalVotingPeriod, &p.NormalVotingPeriod},
		{KeyNormalQuorum, &p.NormalQuorum},
		{KeyNormalThreshold, &p.NormalThreshold},
		{KeyNormalVeto, &p.NormalVeto},
		{KeyNormalPenalty, &p.NormalPenalty},
		{KeyNormalBurnRate, &p.NormalBurnRate},

		{KeyImportantMinDeposit, &p.ImportantMinDeposit},
		{KeyImportantMinProposerDepositRate, &p.ImportantMinProposerDepositRate},
		{KeyImportantMaxDepositPeriod, &p.ImportantMaxDepositPeriod},
		{KeyImportantVotingPeriod, &p.ImportantVotingPeriod},
		{KeyImportantQuorum, &p.ImportantQuorum},
		{KeyImportantThreshold, &p.ImportantThreshold},
		{KeyImportantVeto, &p.ImportantVeto},
		{KeyImportantPenalty, &p.ImportantPenalty},
		{KeyImportantBurnRate, &p.ImportantBurnRate},

		{KeyCriticalMinDeposit, &p.CriticalMinDeposit},
		{KeyCriticalMinProposerDepositRate, &p.CriticalMinProposerDepositRate},
		{KeyCriticalMaxDepositPeriod, &p.CriticalMaxDepositPeriod},
		{KeyCriticalVotingPeriod, &p.CriticalVotingPeriod},
		{KeyCriticalQuorum, &p.CriticalQuorum},
		{KeyCriticalThreshold, &p.CriticalThreshold},
		{KeyCriticalVeto, &p.CriticalVeto},
		{KeyCriticalPenalty, &p.CriticalPenalty},
		{KeyCriticalBurnRate, &p.CriticalBurnRate},
	}
}

// Valid single parameter, return the value
func (p *Params) ValidateKeyValue(key string, value string) (interface{}, btypes.Error) {
	switch key {
	case string(KeyNormalMinDeposit), string(KeyImportantMinDeposit), string(KeyCriticalMinDeposit):
		v, ok := btypes.NewIntFromString(value)
		if !ok || !v.GT(btypes.ZeroInt()) {
			return nil, params.ErrInvalidParam(fmt.Sprintf("%s invalid", key))
		}
		return v, nil
	case string(KeyNormalMaxDepositPeriod), string(KeyImportantMaxDepositPeriod), string(KeyCriticalMaxDepositPeriod),
		string(KeyNormalVotingPeriod), string(KeyImportantVotingPeriod), string(KeyCriticalVotingPeriod):
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil || v <= 0 {
			return nil, params.ErrInvalidParam(fmt.Sprintf("%s invalid", key))
		}
		return time.Duration(v), nil
	case string(KeyNormalMinProposerDepositRate), string(KeyNormalQuorum), string(KeyNormalThreshold), string(KeyNormalVeto), string(KeyNormalPenalty), string(KeyNormalBurnRate),
		string(KeyImportantMinProposerDepositRate), string(KeyImportantQuorum), string(KeyImportantThreshold), string(KeyImportantVeto), string(KeyImportantPenalty), string(KeyImportantBurnRate),
		string(KeyCriticalMinProposerDepositRate), string(KeyCriticalQuorum), string(KeyCriticalThreshold), string(KeyCriticalVeto), string(KeyCriticalPenalty), string(KeyCriticalBurnRate):
		v, err := qtypes.NewDecFromStr(value)
		if err != nil || v.LT(qtypes.ZeroDec()) {
			return nil, params.ErrInvalidParam(fmt.Sprintf("%s invalid", key))
		}
		return v, nil
	default:
		return nil, params.ErrInvalidParam(fmt.Sprintf("%s not exists", key))
	}
}

// Param space, the same as the module name.
func (p *Params) GetParamSpace() string {
	return ParamSpace
}

// Params struct for different level
type LevelParams struct {
	MinDeposit             btypes.BigInt
	MinProposerDepositRate qtypes.Dec
	MaxDepositPeriod       time.Duration
	VotingPeriod           time.Duration
	Quorum                 qtypes.Dec
	Threshold              qtypes.Dec
	Veto                   qtypes.Dec
	Penalty                qtypes.Dec
	BurnRate               qtypes.Dec
}

// Return levelParams for the specific level
func (p *Params) GetLevelParams(level ProposalLevel) LevelParams {
	switch level {
	case LevelNormal:
		return LevelParams{
			MinDeposit:             p.NormalMinDeposit,
			MinProposerDepositRate: p.NormalMinProposerDepositRate,
			MaxDepositPeriod:       p.NormalMaxDepositPeriod,
			VotingPeriod:           p.NormalVotingPeriod,
			Quorum:                 p.NormalQuorum,
			Threshold:              p.NormalThreshold,
			Veto:                   p.NormalVeto,
			Penalty:                p.NormalPenalty,
			BurnRate:               p.NormalBurnRate,
		}
	case LevelImportant:
		return LevelParams{
			MinDeposit:             p.ImportantMinDeposit,
			MinProposerDepositRate: p.ImportantMinProposerDepositRate,
			MaxDepositPeriod:       p.ImportantMaxDepositPeriod,
			VotingPeriod:           p.ImportantVotingPeriod,
			Quorum:                 p.ImportantQuorum,
			Threshold:              p.ImportantThreshold,
			Veto:                   p.ImportantVeto,
			Penalty:                p.ImportantPenalty,
			BurnRate:               p.ImportantBurnRate,
		}
	case LevelCritical:
		return LevelParams{
			MinDeposit:             p.CriticalMinDeposit,
			MinProposerDepositRate: p.CriticalMinProposerDepositRate,
			MaxDepositPeriod:       p.CriticalMaxDepositPeriod,
			VotingPeriod:           p.CriticalVotingPeriod,
			Quorum:                 p.CriticalQuorum,
			Threshold:              p.CriticalThreshold,
			Veto:                   p.CriticalVeto,
			Penalty:                p.CriticalPenalty,
			BurnRate:               p.CriticalBurnRate,
		}
	}

	return LevelParams{}
}

// Validate the params as a whole
func (p *Params) Validate() btypes.Error {
	for _, level := range ProposalLevels {
		levelParams := p.GetLevelParams(level)
		if !levelParams.MinDeposit.GT(btypes.ZeroInt()) {
			params.ErrInvalidParam(fmt.Sprintf("%s_min_deposit must gt 0", level))
		}
		if levelParams.MinProposerDepositRate.GT(qtypes.OneDec()) || levelParams.MinProposerDepositRate.IsNegative() {
			params.ErrInvalidParam(fmt.Sprintf("%s_min_proposer_deposit_rate must gte 0 and lte 1", level))
		}
		if levelParams.MaxDepositPeriod.Nanoseconds() <= 0 {
			params.ErrInvalidParam(fmt.Sprintf("%s_max_deposit_period must gt 0", level))
		}
		if levelParams.VotingPeriod.Nanoseconds() <= 0 {
			params.ErrInvalidParam(fmt.Sprintf("%s_voting_period must gt 0", level))
		}
		if levelParams.Quorum.GT(qtypes.OneDec()) || levelParams.Quorum.IsNegative() {
			params.ErrInvalidParam(fmt.Sprintf("%s_quorum must gte 0 and lte 1", level))
		}
		if levelParams.Threshold.GT(qtypes.OneDec()) || levelParams.Threshold.IsNegative() {
			params.ErrInvalidParam(fmt.Sprintf("%s_threshold must gte 0 and lte 1", level))
		}
		if levelParams.Veto.GT(qtypes.OneDec()) || levelParams.Veto.IsNegative() {
			params.ErrInvalidParam(fmt.Sprintf("%s_veto must gte 0 and lte 1", level))
		}
		if levelParams.Penalty.GT(qtypes.OneDec()) || levelParams.Penalty.IsNegative() {
			params.ErrInvalidParam(fmt.Sprintf("%s_penalty must gte 0 and lte 1", level))
		}
		if levelParams.BurnRate.GT(qtypes.OneDec()) || levelParams.BurnRate.IsNegative() {
			params.ErrInvalidParam(fmt.Sprintf("%s_burn_rate must gte 0 and lte 1", level))
		}
	}

	return nil
}
