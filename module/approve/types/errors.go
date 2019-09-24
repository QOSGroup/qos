package types

import (
	btypes "github.com/QOSGroup/qbase/types"
)

// Approve errors reserve 100 ~ 199.
const (
	DefaultCodeSpace btypes.CodespaceType = "approve"

	CodeInvalidInput              btypes.CodeType = 101 // 基础数据输入有误
	CodeQSCNotExists              btypes.CodeType = 102 // 联盟币不存在
	CodeApproveExists             btypes.CodeType = 103 // 预授权已存在
	CodeApproveNotExists          btypes.CodeType = 104 // 预授权不存在
	CodeFromAccountNotExists      btypes.CodeType = 105 // 授权账户不存在
	CodeApproveNotEnough          btypes.CodeType = 106 // 授权不足
	CodeFromAccountCoinsNotEnough btypes.CodeType = 107 // 授权账户余额不足
)

func ErrInvalidInput(msg string) btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeInvalidInput, msg)
}

func ErrQSCNotExists() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeQSCNotExists, "approve contains qsc that not exists")
}

func ErrApproveExists() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeApproveExists, "approve exists")
}

func ErrApproveNotExists() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeApproveNotExists, "approve not exists")
}

func ErrFromAccountNotExists() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeFromAccountNotExists, "from account not exists")
}

func ErrApproveNotEnough() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeApproveNotEnough, "approve not enough")
}

func ErrFromAccountCoinsNotEnough() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeFromAccountCoinsNotEnough, "from account has no enough coins")
}
