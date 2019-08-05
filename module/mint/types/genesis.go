package types

import "time"

type GenesisState struct {
	Params           Params `json:"params"`
	FirstBlockTime   int64  `json:"first_block_time"` //UTC().UNIX()
	AppliedQOSAmount uint64 `json:"applied_qos_amount"`
}

func NewGenesisState(params Params, firstBlockTime int64, appliedQOSAmount uint64) GenesisState {
	return GenesisState{
		Params:           params,
		FirstBlockTime:   firstBlockTime,
		AppliedQOSAmount: appliedQOSAmount,
	}
}

func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultMintParams(), time.Now().Unix(), 0)
}

func ValidateGenesis(gs GenesisState) error {
	return nil
}
