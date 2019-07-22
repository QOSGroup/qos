package types

import (
	btypes "github.com/QOSGroup/qbase/types"
)

// QCP errors reserve 400 ~ 499.
const (
	DefaultCodeSpace btypes.CodespaceType = "qcp"

	CodeInvalidInput     btypes.CodeType = 401 // 信息有误
	CodeInvalidQCPCA     btypes.CodeType = 402 // 无效证书
	CodeWrongQCPCA       btypes.CodeType = 403 // 证书有误
	CodeCreatorNotExists btypes.CodeType = 404 // 创建账户不存在
	CodeQCPExists        btypes.CodeType = 405 // QCP已存在
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
	case CodeInvalidQCPCA:
		return "invalid qcp ca"
	case CodeWrongQCPCA:
		return "wrong qcp ca"
	case CodeCreatorNotExists:
		return "creator not exists"
	case CodeQCPExists:
		return "qcp exists"
	default:
		return btypes.CodeToDefaultMsg(code)
	}
}

func ErrInvalidInput(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeInvalidInput, msg)
}

func ErrInvalidQCPCA(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeInvalidQCPCA, msg)
}

func ErrWrongQCPCA(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeWrongQCPCA, msg)
}

func ErrCreatorNotExists(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeCreatorNotExists, msg)
}

func ErrQCPExists(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeQCPExists, msg)
}
