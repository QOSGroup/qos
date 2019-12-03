package types

import (
	btypes "github.com/QOSGroup/qbase/types"
)

type KeyValuePair struct {
	Key   []byte
	Value interface{}
}

type KeyValuePairs []KeyValuePair

// 参数接口，不同模块参数结构必须实现该接口并注册到params模块才能正常使用params提供的参数管理方法。
type ParamSet interface {
	// 参数键值对
	KeyValuePairs() KeyValuePairs

	// 单个参数类型和值校验
	ValidateKeyValue(key string, value string) (interface{}, btypes.Error)

	// 参数所属模块名
	GetParamSpace() string

	// 参数校验
	Validate() btypes.Error

	// 设置单个参数
	SetKeyValue(key string, value interface{}) btypes.Error
}