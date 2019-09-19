package types

import (
	btypes "github.com/QOSGroup/qbase/types"
)

// QCP errors reserve 400 ~ 499.
const (
	DefaultCodeSpace btypes.CodespaceType = "qcp"

	CodeInvalidInput       btypes.CodeType = 401 // 信息有误
	CodeInvalidQCPCA       btypes.CodeType = 402 // 无效证书
	CodeWrongQCPCA         btypes.CodeType = 403 // 证书有误
	CodeCreatorNotExists   btypes.CodeType = 404 // 创建账户不存在
	CodeQCPExists          btypes.CodeType = 405 // QCP已存在
	CodeEmptyCreator       btypes.CodeType = 406 // 创建账户为空
	CodeRootCANotConfigure btypes.CodeType = 407 // 没有配置root ca
)

func ErrInvalidInput(msg string) btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeInvalidInput, msg)
}

func ErrInvalidQCPCA() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeInvalidQCPCA, "invalid qcp ca")
}

func ErrWrongQCPCA() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeWrongQCPCA, "wrong qcp ca")
}

func ErrCreatorNotExists() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeCreatorNotExists, "creator not exists")
}

func ErrQCPExists() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeQCPExists, "qcp exists")
}

func ErrEmptyCreator() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeEmptyCreator, "empty creator")
}

func ErrRootCANotConfigure() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeRootCANotConfigure, "no root ca public key initialized")
}
