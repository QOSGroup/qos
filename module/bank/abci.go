package bank

import (
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/bank/mapper"
	"github.com/QOSGroup/qos/module/bank/types"
	qtypes "github.com/QOSGroup/qos/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"time"
)

func EndBlocker(ctx context.Context, req abci.RequestEndBlock) {
	// 存在数据检查请求时向Event中添加EventTypeInvariantCheck事件
	if NeedInvariantCheck(ctx) {
		ctx.EventManager().EmitEvent(btypes.NewEvent(qtypes.EventTypeInvariantCheck))
	}

	// 锁定账户释放信息
	if lockInfo, exists := mapper.GetLockInfo(ctx); exists {
		ReleaseLockedAccount(ctx, lockInfo)
	}

	return
}

func ReleaseLockedAccount(ctx context.Context, lockInfo LockInfo) {
	if lockInfo.ReleaseTime.Before(ctx.BlockHeader().Time.UTC()) {
		releaseAmount := btypes.ZeroInt()
		if lockInfo.ReleaseTimes != 1 {
			releaseAmount = lockInfo.TotalAmount.Sub(lockInfo.ReleasedAmount).DivRaw(lockInfo.ReleaseTimes)
		} else {
			releaseAmount = lockInfo.TotalAmount.Sub(lockInfo.ReleasedAmount)
		}
		if releaseAmount.GT(btypes.ZeroInt()) {
			// 更新lockinfo
			lockInfo.ReleasedAmount = lockInfo.ReleasedAmount.Add(releaseAmount)
			lockInfo.ReleaseTimes -= 1
			// lockInfo.ReleaseTime = lockInfo.ReleaseTime.Add(time.Hour * 24 * time.Duration(lockInfo.ReleaseInterval))
			lockInfo.ReleaseTime = lockInfo.ReleaseTime.Add(time.Hour * time.Duration(lockInfo.ReleaseInterval))
			mapper.SetLockInfo(ctx, lockInfo)
			// 更新锁定账户
			lockedAccount := mapper.GetAccount(ctx, lockInfo.LockedAccount)
			if lockedAccount == nil {
				panic("LockAccount not exists")
			}
			lockedAccount.MustMinusQOS(releaseAmount)
			mapper.GetMapper(ctx).SetAccount(lockedAccount)
			// 更新接收账户
			receiver := mapper.GetAccount(ctx, lockInfo.Receiver)
			if receiver == nil {
				receiver = qtypes.NewQOSAccountWithAddress(lockInfo.Receiver)
			}
			receiver.MustPlusQOS(releaseAmount)
			mapper.GetMapper(ctx).SetAccount(receiver)

			// 发送事件
			ctx.EventManager().EmitEvent(
				btypes.NewEvent(types.EventTypeRelease,
					btypes.NewAttribute(types.AttributeKeyAddress, receiver.AccountAddress.String()),
					btypes.NewAttribute(types.AttributeKeyQOS, releaseAmount.String()),
				),
			)
		}

		// 释放完成删除锁定信息
		if lockInfo.TotalAmount.Equal(lockInfo.ReleasedAmount) {
			mapper.DelLockInfo(ctx)
		}
	}
}
