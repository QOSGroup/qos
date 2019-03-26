package types

import (
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/pkg/errors"
	"strings"
)

type Guardian struct {
	Description  string         `json:"description"`
	GuardianType GuardianType   `json:"guardian_type"`
	Address      btypes.Address `json:"address"`
	Creator      btypes.Address `json:"creator"`
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

func AccountTypeFromString(str string) (GuardianType, error) {
	switch strings.ToLower(str) {
	case "genesis":
		return Genesis, nil
	case "ordinary":
		return Ordinary, nil
	default:
		return GuardianType(0xff), errors.Errorf("'%s' is not a valid account type", str)
	}
}
