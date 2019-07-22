package types

type GenesisState struct {
	Params           Params `json:"params"`
	FirstBlockTime   int64  `json:"first_block_time"` //UTC().UNIX()
	AppliedQOSAmount uint64 `json:"applied_qos_amount"`
}

func NewGenesisState(params Params) GenesisState {
	return GenesisState{
		Params: params,
	}
}

func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultMintParams())
}
