package gov

import (
	"fmt"
	btypes "github.com/QOSGroup/qbase/types"
	perr "github.com/QOSGroup/qos/module/params"
	ptypes "github.com/QOSGroup/qos/module/params/types"
	"github.com/QOSGroup/qos/types"
	"strconv"
	"time"
)

var (
	ParamSpace = "gov"

	KeyMinDeposit       = []byte("min_deposit")
	KeyMaxDepositPeriod = []byte("max_deposit_period")
	KeyVotingPeriod     = []byte("voting_period")
	KeyQuorum           = []byte("quorum")
	KeyThreshold        = []byte("threshold")
	KeyVeto             = []byte("veto")
	KeyPenalty          = []byte("penalty")
)

// Params returns all of the governance params
type Params struct {
	// DepositParams
	MinDeposit       uint64        `json:"min_deposit"`        //  Minimum deposit for a proposal to enter voting period.
	MaxDepositPeriod time.Duration `json:"max_deposit_period"` //  Maximum period for Atom holders to deposit on a proposal. Initial value: 2 months

	// VotingParams
	VotingPeriod time.Duration `json:"voting_period"` //  Length of the voting period.

	// TallyParams
	Quorum    types.Dec `json:"quorum"`    //  Minimum percentage of total stake needed to vote for a result to be considered valid
	Threshold types.Dec `json:"threshold"` //  Minimum propotion of Yes votes for proposal to pass. Initial value: 0.5
	Veto      types.Dec `json:"veto"`      //  Minimum value of Veto votes to Total votes ratio for proposal to be vetoed. Initial value: 1/3
	Penalty   types.Dec `json:"penalty"`   //  Penalty if validator does not vote
}

func DefaultParams() Params {
	return Params{
		MinDeposit:       10,
		MaxDepositPeriod: DefaultPeriod,
		VotingPeriod:     DefaultPeriod,
		Quorum:           types.NewDecWithPrec(334, 3),
		Threshold:        types.NewDecWithPrec(5, 1),
		Veto:             types.NewDecWithPrec(334, 3),
		Penalty:          types.ZeroDec(),
	}
}

func (params *Params) KeyValuePairs() ptypes.KeyValuePairs {
	return ptypes.KeyValuePairs{
		{KeyMinDeposit, &params.MinDeposit},
		{KeyMaxDepositPeriod, &params.MaxDepositPeriod},
		{KeyVotingPeriod, &params.VotingPeriod},
		{KeyQuorum, &params.Quorum},
		{KeyThreshold, &params.Threshold},
		{KeyVeto, &params.Veto},
		{KeyPenalty, &params.Penalty},
	}
}

func (params *Params) Validate(key string, value string) (interface{}, btypes.Error) {
	switch key {
	case string(KeyMinDeposit):
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, perr.ErrInvalidParam(fmt.Sprintf("%s invalid", key))
		}
		return v, nil
	case string(KeyMaxDepositPeriod),
		string(KeyVotingPeriod):
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, perr.ErrInvalidParam(fmt.Sprintf("%s invalid", key))
		}
		return time.Duration(v), nil
	case string(KeyQuorum), string(KeyThreshold), string(KeyVeto), string(KeyPenalty):
		v, err := types.NewDecFromStr(value)
		if err != nil {
			return nil, perr.ErrInvalidParam(fmt.Sprintf("%s invalid", key))
		}
		return v, nil
	default:
		return nil, perr.ErrInvalidParam(fmt.Sprintf("%s not exists", key))
	}
}

func (params *Params) GetParamSpace() string {
	return ParamSpace
}
