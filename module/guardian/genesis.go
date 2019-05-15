package guardian

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qos/module/guardian/types"
	"github.com/pkg/errors"
)

type GenesisState struct {
	Guardians []types.Guardian `json:"guardians"`
}

func (gs GenesisState) Equals(gs1 GenesisState) bool {
	if len(gs.Guardians) != len(gs1.Guardians) {
		return false
	}
	for _, g := range gs.Guardians {
		exists := false
		for _, g1 := range gs1.Guardians {
			if g.Equals(g1) {
				exists = true
				break
			}
		}
		if !exists {
			return false
		}
	}

	return true
}

func NewGenesisState(guardians []types.Guardian) GenesisState {
	return GenesisState{
		Guardians: guardians,
	}
}

func InitGenesis(ctx context.Context, data GenesisState) {
	err := ValidateGenesis(data)
	if err != nil {
		panic(err)
	}

	mapper := GetGuardianMapper(ctx)
	for _, guardian := range data.Guardians {
		mapper.AddGuardian(guardian)
	}
}

func ExportGenesis(ctx context.Context) GenesisState {
	mapper := GetGuardianMapper(ctx)
	iterator := mapper.GuardiansIterator()
	defer iterator.Close()
	var guardians []types.Guardian
	for ; iterator.Valid(); iterator.Next() {
		var guardian types.Guardian
		mapper.GetCodec().MustUnmarshalBinaryBare(iterator.Value(), &guardian)
		guardians = append(guardians, guardian)
	}

	return NewGenesisState(guardians)
}

func ValidateGenesis(gs GenesisState) error {
	addrs := map[string]bool{}
	for _, g := range gs.Guardians {
		if g.GuardianType != types.Genesis || (g.Creator != nil && len(g.Creator) != 0) {
			return errors.New("invalid genesis guardian")
		}
		if _, ok := addrs[string(g.Address.String())]; ok {
			return errors.New(fmt.Sprintf("duplicate guardian %s", string(g.Address.String())))
		}
		addrs[string(g.Address.String())] = true
	}

	return nil
}
