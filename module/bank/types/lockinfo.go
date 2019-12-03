package types

import (
	"github.com/QOSGroup/qbase/types"
	"time"
)

// 锁定-释放账户信息
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
