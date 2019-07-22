package types

import (
	btypes "github.com/QOSGroup/qbase/types"
)

// Transfer errors reserve 200 ~ 299.
const (
	DefaultCodeSpace btypes.CodespaceType = "transfer"

	CodeInvalidInput                btypes.CodeType = 201 // 基础数据输入有误
	CodeSenderAccountNotExists      btypes.CodeType = 202 // 转出账户不存在
	CodeSenderAccountCoinsNotEnough btypes.CodeType = 203 // 转出账户余额不足
)

func msgOrDefaultMsg(msg string, code btypes.CodeType) string {
	if msg != "" {
		return msg
	}
	return codeToDefaultMsg(code)
}

func newError(codeSpace btypes.CodespaceType, code btypes.CodeType, msg string) btypes.Error {
	msg = msgOrDefaultMsg(msg, code)
	return btypes.NewError(codeSpace, code, msg)
}

// NOTE: Don't stringer this, we'll put better messages in later.
func codeToDefaultMsg(code btypes.CodeType) string {
	switch code {
	case CodeInvalidInput:
		return "invalid transfer msg"
	case CodeSenderAccountNotExists:
		return "sender account not exists"
	case CodeSenderAccountCoinsNotEnough:
		return "sender account has no enough coins"
	default:
		return btypes.CodeToDefaultMsg(code)
	}
}

func ErrInvalidInput(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeInvalidInput, msg)
}

func ErrSenderAccountNotExists(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeSenderAccountNotExists, msg)
}

func ErrSenderAccountCoinsNotEnough(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeSenderAccountCoinsNotEnough, msg)
}
