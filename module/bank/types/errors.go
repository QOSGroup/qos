package types

import (
	btypes "github.com/QOSGroup/qbase/types"
)

// Transfer errors reserve 200 ~ 299.
const (
	DefaultCodeSpace btypes.CodespaceType = "bank"

	CodeInvalidInput                btypes.CodeType = 201 // 基础数据输入有误
	CodeSenderAccountNotExists      btypes.CodeType = 202 // 转出账户不存在
	CodeSenderAccountCoinsNotEnough btypes.CodeType = 203 // 转出账户余额不足
)

func ErrInvalidInput(msg string) btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeInvalidInput, msg)
}

func ErrSenderAccountNotExists() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeSenderAccountNotExists, "sender account not exists")
}

func ErrSenderAccountCoinsNotEnough() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeSenderAccountCoinsNotEnough, "sender account has no enough coins")
}
