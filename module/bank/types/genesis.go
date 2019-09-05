package types

import (
	"errors"
	"fmt"
	"github.com/QOSGroup/qbase/types"
	qtypes "github.com/QOSGroup/qos/types"
	"time"
)

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

func DefaultGenesisState() GenesisState {
	return GenesisState{[]*qtypes.QOSAccount{}, nil}
}

func ValidateGenesis(gs GenesisState) error {
	addrMap := make(map[string]bool, len(gs.Accounts))
	for i := 0; i < len(gs.Accounts); i++ {
		acc := gs.Accounts[i]
		strAddr := string(acc.AccountAddress)
		if _, ok := addrMap[strAddr]; ok {
			return fmt.Errorf("duplicate account in genesis state: Address %v", acc.AccountAddress)
		}
		addrMap[strAddr] = true
	}
	if gs.LockInfo != nil {
		if len(gs.LockInfo.LockedAccount) == 0 || len(gs.LockInfo.Receiver) == 0 ||
			!gs.LockInfo.TotalAmount.GT(gs.LockInfo.ReleasedAmount) ||
			gs.LockInfo.ReleaseInterval == 0 || gs.LockInfo.ReleaseTimes == 0 {
			return errors.New("invalid lock account")
		}
	}
	return nil
}

type LockInfo struct {
	LockedAccount   types.AccAddress `json:"locked_account"`   // 锁定账户地址
	Receiver        types.AccAddress `json:"receiver"`         // 接收账户地址
	TotalAmount     types.BigInt     `json:"total_amount"`     // 总锁定QOS
	ReleasedAmount  types.BigInt     `json:"released_amount"`  // 已释放QOS
	ReleaseTime     time.Time        `json:"release_time"'`    // 下一次释放时间
	ReleaseInterval int64            `json:"release_interval"` // 释放间隔，以天为单位
	ReleaseTimes    int64            `json:"release_times"`    // 释放次数
}

func NewLockInfo(lockedAccount, receiver types.AccAddress, totalAmount, releasedAmount types.BigInt, releaseTime time.Time, releaseInterval, releaseTimes int64) LockInfo {
	return LockInfo{
		LockedAccount:   lockedAccount,
		Receiver:        receiver,
		TotalAmount:     totalAmount,
		ReleasedAmount:  releasedAmount,
		ReleaseTime:     releaseTime,
		ReleaseInterval: releaseInterval,
		ReleaseTimes:    releaseTimes,
	}
}
