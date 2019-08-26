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

	KeyMinDeposit             = []byte("min_deposit")
	keyMinProposerDepositRate = []byte("min_proposer_deposit_rate")
	KeyMaxDepositPeriod       = []byte("max_deposit_period")
	KeyVotingPeriod           = []byte("voting_period")
	KeyQuorum                 = []byte("quorum")
	KeyThreshold              = []byte("threshold")
	KeyVeto                   = []byte("veto")
	KeyPenalty                = []byte("penalty")
	KeyBurnRate               = []byte("burn_rate")
)

// Params returns all of the governance p
type Params struct {
	// DepositParams
	MinDeposit             uint64        `json:"min_deposit"`               //  Minimum deposit for a proposal to enter voting period.
	MinProposerDepositRate qtypes.Dec    `json:"min_proposer_deposit_rate"` //  Minimum deposit rate for proposer to submit a proposal. Initial value: 1/3
	MaxDepositPeriod       time.Duration `json:"max_deposit_period"`        //  Maximum period for Atom holders to deposit on a proposal. Initial value: 2 months

	// VotingParams
	VotingPeriod time.Duration `json:"voting_period"` //  Length of the voting period.

	// TallyParams
	Quorum    qtypes.Dec `json:"quorum"`    //  Minimum percentage of total stake needed to vote for a result to be considered valid
	Threshold qtypes.Dec `json:"threshold"` //  Minimum propotion of Yes votes for proposal to pass. Initial value: 0.5
	Veto      qtypes.Dec `json:"veto"`      //  Minimum value of Veto votes to Total votes ratio for proposal to be vetoed. Initial value: 1/3
	Penalty   qtypes.Dec `json:"penalty"`   //  Penalty if validator does not vote

	BurnRate qtypes.Dec `json:"burn_rate"` // Deposit burning rate when proposals pass or reject. Initial value: 1/2
}

func DefaultParams() Params {
	return Params{
		MinDeposit:             10,
		MinProposerDepositRate: qtypes.NewDecWithPrec(334, 3),
		MaxDepositPeriod:       DefaultPeriod,
		VotingPeriod:           DefaultPeriod,
		Quorum:                 qtypes.NewDecWithPrec(334, 3),
		Threshold:              qtypes.NewDecWithPrec(5, 1),
		Veto:                   qtypes.NewDecWithPrec(334, 3),
		Penalty:                qtypes.ZeroDec(),
		BurnRate:               qtypes.NewDecWithPrec(5, 1),
	}
}

func (p *Params) KeyValuePairs() qtypes.KeyValuePairs {
	return qtypes.KeyValuePairs{
		{KeyMinDeposit, &p.MinDeposit},
		{keyMinProposerDepositRate, &p.MinProposerDepositRate},
		{KeyMaxDepositPeriod, &p.MaxDepositPeriod},
		{KeyVotingPeriod, &p.VotingPeriod},
		{KeyQuorum, &p.Quorum},
		{KeyThreshold, &p.Threshold},
		{KeyVeto, &p.Veto},
		{KeyPenalty, &p.Penalty},
		{KeyBurnRate, &p.BurnRate},
	}
}

func (p *Params) Validate(key string, value string) (interface{}, btypes.Error) {
	switch key {
	case string(KeyMinDeposit):
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, params.ErrInvalidParam(fmt.Sprintf("%s invalid", key))
		}
		return v, nil
	case string(KeyMaxDepositPeriod),
		string(KeyVotingPeriod):
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, params.ErrInvalidParam(fmt.Sprintf("%s invalid", key))
		}
		return time.Duration(v), nil
	case string(keyMinProposerDepositRate), string(KeyQuorum), string(KeyThreshold), string(KeyVeto), string(KeyPenalty):
		v, err := qtypes.NewDecFromStr(value)
		if err != nil {
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
