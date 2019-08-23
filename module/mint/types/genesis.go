package types

import (
	"errors"
	qtypes "github.com/QOSGroup/qos/types"
	"time"
)

type GenesisState struct {
	InflationPhrases InflationPhrases `json:"inflation_phrases"`
	FirstBlockTime   int64            `json:"first_block_time"` //UTC().UNIX()
	AppliedQOSAmount uint64           `json:"applied_qos_amount"`
	TotalQOSAmount   uint64           `json:"total_qos_amount"`
}

func NewGenesisState(inflationPhrases InflationPhrases, firstBlockTime int64, appliedQOSAmount uint64, totalQOSAmount uint64) GenesisState {
	return GenesisState{
		InflationPhrases: inflationPhrases,
		FirstBlockTime:   firstBlockTime,
		AppliedQOSAmount: appliedQOSAmount,
		TotalQOSAmount:   totalQOSAmount,
	}
}

func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultInflationPhrases(), time.Now().Unix(), 0, qtypes.TotalQOSAmount)
}

func ValidateGenesis(gs GenesisState) error {
	if gs.TotalQOSAmount == 0 {
		return errors.New("total amount must positive")
	}

	return gs.InflationPhrases.Valid()
}
