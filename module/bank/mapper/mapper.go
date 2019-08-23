package mapper

import (
	"encoding/binary"
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	qtypes "github.com/QOSGroup/qos/types"
)

var (
	InvariantCheckKey = []byte("0x10")
)

func GetMapper(ctx context.Context) *account.AccountMapper {
	return baseabci.GetAccountMapper(ctx)
}

func GetAccount(ctx context.Context, addr btypes.Address) *qtypes.QOSAccount {
	account := baseabci.GetAccountMapper(ctx).GetAccount(addr)
	if account != nil {
		return account.(*qtypes.QOSAccount)
	} else {
		return nil
	}
}

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
