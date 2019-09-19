package types

import (
	btypes "github.com/QOSGroup/qbase/types"
)

// 系统账户
type Guardian struct {
	Description  string            `json:"description"`   // 描述
	GuardianType GuardianType      `json:"guardian_type"` // 账户类型：Genesis 创世配置 Ordinary 交易创建
	Address      btypes.AccAddress `json:"address"`       // 账户地址
	Creator      btypes.AccAddress `json:"creator"`       // 创建者账户地址
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
	Genesis  GuardianType = 0x01 // 创世配置
	Ordinary GuardianType = 0x02 // 交易创建
)
