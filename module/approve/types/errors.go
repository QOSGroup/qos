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
		return "invalid approve msg"
	case CodeQSCNotExists:
		return "approve contains qsc that not exists"
	case CodeApproveExists:
		return "approve exists"
	case CodeApproveNotExists:
		return "approve not exists"
	case CodeFromAccountNotExists:
		return "from account not exists"
	case CodeApproveNotEnough:
		return "approve not enough"
	case CodeFromAccountCoinsNotEnough:
		return "from account has no enough coins"
	default:
		return btypes.CodeToDefaultMsg(code)
	}
}

func ErrInvalidInput(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeInvalidInput, msg)
}

func ErrQSCNotExists(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeQSCNotExists, msg)
}

func ErrApproveExists(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeApproveExists, msg)
}

func ErrApproveNotExists(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeApproveNotExists, msg)
}

func ErrFromAccountNotExists(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeFromAccountNotExists, msg)
}

func ErrApproveNotEnough(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeApproveNotEnough, msg)
}

func ErrFromAccountCoinsNotEnough(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeFromAccountCoinsNotEnough, msg)
}
