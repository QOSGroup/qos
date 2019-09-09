# 状态

## 参数

参数数据库操作`Mapper`：
```go
type Mapper struct {
	*mapper.BaseMapper

	// 参数表：模块名-参数
	paramSets map[string]qtypes.ParamSet
}
```
`MapperName`为`params`，`paramSets`在QOS App创建时注册到参数模块，初始化时保存到数据库中。

单个参数存储：

param: `moduleName keyName -> keyValue`
