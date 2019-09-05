# 存储

## 预授权

预授权存储结构如下：
```go
type Approve struct {
	From btypes.AccAddress `json:"from"` // 授权账号
	To   btypes.AccAddress `json:"to"`   // 被授权账号
	QOS  btypes.BigInt     `json:"qos"`  // QOS
	QSCs types.QSCs        `json:"qscs"` // QSCs
}
```
操作预授权的`MapperName`为`approve`，`Mapper`中的具体存储如下：
- approve: `0x01 from to -> amino(approve)`
