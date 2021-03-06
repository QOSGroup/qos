package mapper

import (
	"encoding/binary"
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/bank/types"
	qtypes "github.com/QOSGroup/qos/types"
)

var (
	InvariantCheckKey = []byte("0x10") // 数据检查
	LockInfoKey       = []byte("0x11") // 锁定-释放信息
)

// 基于qbase中定义的AccountMapper，扩展qos所需操作
func GetMapper(ctx context.Context) *account.AccountMapper {
	return baseabci.GetAccountMapper(ctx)
}

// 获取QOS账户信息
func GetAccount(ctx context.Context, addr btypes.AccAddress) *qtypes.QOSAccount {
	account := baseabci.GetAccountMapper(ctx).GetAccount(addr)
	if account != nil {
		return account.(*qtypes.QOSAccount)
	} else {
		return nil
	}
}

// 获取QOS账户列表
func GetAccounts(ctx context.Context) []*qtypes.QOSAccount {
	accounts := []*qtypes.QOSAccount{}
	baseabci.GetAccountMapper(ctx).IterateAccounts(func(acc account.Account) (stop bool) {
		accounts = append(accounts, acc.(*qtypes.QOSAccount))
		return false
	})

	return accounts
}

// ----------------
// invariant check

func BuildInvariantCheckKey(height uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, height)
	return append(InvariantCheckKey, b...)
}

// 保存某以高度执行检查信息
func SetInvariantCheck(ctx context.Context) {
	GetMapper(ctx).Set(BuildInvariantCheckKey(uint64(ctx.BlockHeight())), true)
}

// 某一高度是否需要执行数据检查
func NeedInvariantCheck(ctx context.Context) bool {
	b := false
	exists := GetMapper(ctx).Get(BuildInvariantCheckKey(uint64(ctx.BlockHeight())), &b)
	return exists && b
}

// 清空数据检查
func ClearInvariantCheck(ctx context.Context) {
	GetMapper(ctx).IteratorWithKV(InvariantCheckKey, func(key []byte, value []byte) (stop bool) {
		GetMapper(ctx).Del(key)
		return false
	})
}

// ----------------
// lock account info

// 保存锁定账户信息
func SetLockInfo(ctx context.Context, info types.LockInfo) {
	GetMapper(ctx).Set(LockInfoKey, info)
}

// 获取锁定账户信息
func GetLockInfo(ctx context.Context) (info types.LockInfo, exists bool) {
	exists = GetMapper(ctx).Get(LockInfoKey, &info)

	return
}

// 删除锁定账户信息
func DelLockInfo(ctx context.Context) {
	GetMapper(ctx).Del(LockInfoKey)
}
