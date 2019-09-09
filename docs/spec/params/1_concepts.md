# 概念

## ParamSet

其他模块参数结构实现以下规范：

```go
type ParamSet interface {
	// 参数键值对
	KeyValuePairs() KeyValuePairs

	// 单个参数类型和值校验
	ValidateKeyValue(key string, value string) (interface{}, btypes.Error)

	// 参数所属模块名
	GetParamSpace() string

	// 参数校验
	Validate() btypes.Error
}
```

`ParamSet`通过`KeyValuePairs()`返回模块中所有参数键值对:
```go
// 参数键值对
type KeyValuePair struct {
	Key   []byte
	Value interface{}
}

type KeyValuePairs []KeyValuePair
```

QOS `stake`, `gov`, `distribution` 模块涉及到参数配置。

## 参数管理流程

### 注册mappers

注册mappers和hooks时，初始化参数配置：
```go
// app/app.go
app.mm.RegisterMapperAndHooks(app, params.ModuleName, &stake.Params{}, &distribution.Params{}, &gov.Params{})

// module/params/mapper/mapper.go
// 参数注册
func (mapper Mapper) RegisterParamSet(ps ...qtypes.ParamSet) {
	for _, ps := range ps {
		if ps != nil {
			// 禁止重复注册
			if _, ok := mapper.paramSets[ps.GetParamSpace()]; ok {
				panic(fmt.Sprintf("<%s> already registered ", ps.GetParamSpace()))
			}
			mapper.paramSets[ps.GetParamSpace()] = ps
		}
	}
}
```

### 初始化

在有参数的模块初始化时, 保存参数信息：
```go
func InitGenesis(ctx context.Context, bapp *baseabci.BaseApp, data types.GenesisState) []abci.ValidatorUpdate{
    ...
    mapper.SetParams(ctx, params)
    ...
}
```

### 更新

在治理模块，提交参数修改提议更新参数。
