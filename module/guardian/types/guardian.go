package types

import (
	btypes "github.com/QOSGroup/qbase/types"
)

type Guardian struct {
	Description  string         `json:"description"`
	GuardianType GuardianType   `json:"guardian_type"`
	Address      btypes.Address `json:"address"`
	Creator      btypes.Address `json:"creator"`
}

func (g Guardian) Equals(g1 Guardian) bool {
	return g.Description == g1.Description &&
		g.GuardianType == g1.GuardianType &&
		g.Address.EqualsTo(g1.Address) &&
		g.Creator.EqualsTo(g1.Creator)
}

func NewGuardian(description string, guardianType GuardianType, address, creator btypes.Address) *Guardian {
	return &Guardian{
		Description:  description,
		GuardianType: guardianType,
		Address:      address,
		Creator:      creator,
	}
}

type GuardianType byte

const (
	Genesis  GuardianType = 0x01
	Ordinary GuardianType = 0x02
)
