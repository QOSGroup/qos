package types

// 创世状态
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

// 校验创世状态
func ValidateGenesis(gs GenesisState) error {
	for _, approve := range gs.Approves {
		if err := approve.Valid(); err != nil {
			return err
		}
	}

	return nil
}
