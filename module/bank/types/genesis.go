package types

import (
	"errors"
	"fmt"
	qtypes "github.com/QOSGroup/qos/types"
)

// 创世状态
type GenesisState struct {
	Accounts []*qtypes.QOSAccount `json:"accounts"`
	LockInfo *LockInfo            `json:"lock_info"`
}

func NewGenesisState(accounts []*qtypes.QOSAccount, info *LockInfo) GenesisState {
	return GenesisState{
		Accounts: accounts,
		LockInfo: info,
	}
}

// 默认创世状态
func DefaultGenesisState() GenesisState {
	return GenesisState{[]*qtypes.QOSAccount{}, nil}
}

// 创世数据校验
func ValidateGenesis(gs GenesisState) error {
	// 初始账户
	addrMap := make(map[string]bool, len(gs.Accounts))
	for i := 0; i < len(gs.Accounts); i++ {
		acc := gs.Accounts[i]
		strAddr := string(acc.AccountAddress)
		if _, ok := addrMap[strAddr]; ok {
			return fmt.Errorf("duplicate account in genesis state: %v", acc.AccountAddress)
		}
		addrMap[strAddr] = true
	}

	// 初始锁定-释放信息
	if gs.LockInfo != nil {
		if len(gs.LockInfo.LockedAccount) == 0 || len(gs.LockInfo.Receiver) == 0 ||
			!gs.LockInfo.TotalAmount.GT(gs.LockInfo.ReleasedAmount) ||
			gs.LockInfo.ReleaseInterval == 0 || gs.LockInfo.ReleaseTimes == 0 {
			return errors.New("invalid lock info")
		}
	}

	return nil
}
