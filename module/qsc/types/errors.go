package types

import (
	btypes "github.com/QOSGroup/qbase/types"
)

// QSC errors reserve 300 ~ 399.
const (
	DefaultCodeSpace btypes.CodespaceType = "qsc"

	CodeInvalidInput        btypes.CodeType = 301 // 信息有误
	CodeInvalidQSCCA        btypes.CodeType = 302 // 无效证书
	CodeWrongQSCCA          btypes.CodeType = 303 // 证书有误
	CodeInvalidInitAccounts btypes.CodeType = 304 // 散币账户币种币值有误
	CodeCreatorNotExists    btypes.CodeType = 305 // 创建账户不存在
	CodeQSCExists           btypes.CodeType = 306 // QSC已存在
	CodeQSCNotExists        btypes.CodeType = 307 // QSC不存在
	CodeBankerNotExists     btypes.CodeType = 308 // Banker账户不存在
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
		return "invalid tx msg"
	case CodeInvalidQSCCA:
		return "invalid qsc ca"
	case CodeWrongQSCCA:
		return "wrong qsc ca"
	case CodeInvalidInitAccounts:
		return "invalid init accounts"
	case CodeCreatorNotExists:
		return "creator not exists"
	case CodeQSCExists:
		return "qsc exists"
	case CodeQSCNotExists:
		return "qsc not exists"
	case CodeBankerNotExists:
		return "banker not exists"
	default:
		return btypes.CodeToDefaultMsg(code)
	}
}

func ErrInvalidInput(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeInvalidInput, msg)
}

func ErrInvalidQSCCA(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeInvalidQSCCA, msg)
}

func ErrWrongQSCCA(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeWrongQSCCA, msg)
}

func ErrInvalidInitAccounts(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeInvalidInitAccounts, msg)
}

func ErrCreatorNotExists(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeCreatorNotExists, msg)
}

func ErrQSCExists(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeQSCExists, msg)
}

func ErrQSCNotExists(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeQSCNotExists, msg)
}

func ErrBankerNotExists(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeBankerNotExists, msg)
}
