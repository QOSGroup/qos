package types

type GenesisState struct {
	Approves []Approve `json:"approves"`
}

func NewGenesisState(approves []Approve) GenesisState {
	return GenesisState{
		approves,
	}
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func ValidateGenesis(_ GenesisState) error { return nil }
