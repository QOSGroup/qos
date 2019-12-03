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
	CodeAmountLTZero        btypes.CodeType = 309 // 币量小于等于0
	CodeInvalidBanker       btypes.CodeType = 310 // banker有误
	CodeEmptyCreator        btypes.CodeType = 311 // 创建地址为空
	CodeDescriptionTooLong  btypes.CodeType = 312 // 描述信息太长
	CodeInvalidExchangeRate btypes.CodeType = 313 // 汇率值有误
	CodeRootCANotConfigure  btypes.CodeType = 314 // 没有配置qsc root ca public key
)

func ErrInvalidInput(msg string) btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeInvalidInput, msg)
}

func ErrInvalidQSCCA() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeInvalidQSCCA, "invalid qsc ca")
}

func ErrWrongQSCCA() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeWrongQSCCA, "wrong qsc ca")
}

func ErrInvalidInitAccounts() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeInvalidInitAccounts, "invalid init accounts")
}

func ErrCreatorNotExists() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeCreatorNotExists, "creator not exists")
}

func ErrQSCExists() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeQSCExists, "qsc exists")
}

func ErrQSCNotExists() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeQSCNotExists, "qsc not exists")
}

func ErrBankerNotExists() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeBankerNotExists, "banker not exists")
}

func ErrAmountLTZero() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeAmountLTZero, "amount must be positive")
}

func ErrInvalidBanker() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeInvalidBanker, "invalid banker")
}

func ErrEmptyCreator() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeEmptyCreator, "empty creator")
}

func ErrDescriptionTooLong() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeDescriptionTooLong, "description is too long")
}

func ErrInvalidExchangeRate() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeInvalidExchangeRate, "invalid exchange rate")
}

func ErrRootCANotConfigure() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeRootCANotConfigure, "no qsc root ca public key initialized")
}
