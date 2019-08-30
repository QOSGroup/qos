package types

import (
	btypes "github.com/QOSGroup/qbase/types"
)

type Guardian struct {
	Description  string            `json:"description"`
	GuardianType GuardianType      `json:"guardian_type"`
	Address      btypes.AccAddress `json:"address"`
	Creator      btypes.AccAddress `json:"creator"`
}

func (g Guardian) Equals(g1 Guardian) bool {
	return g.Description == g1.Description &&
		g.GuardianType == g1.GuardianType &&
		g.Address.Equals(g1.Address) &&
		g.Creator.Equals(g1.Creator)
}

func NewGuardian(description string, guardianType GuardianType, address, creator btypes.AccAddress) *Guardian {
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
