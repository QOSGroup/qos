# 状态

## 系统账户

系统账户结构设计如下：
```go
type Guardian struct {
	Description  string            `json:"description"`   // 描述
	GuardianType GuardianType      `json:"guardian_type"` // 账户类型：Genesis 创世配置 Ordinary 交易创建
	Address      btypes.AccAddress `json:"address"`       // 账户地址
	Creator      btypes.AccAddress `json:"creator"`       // 创建者账户地址
}
```
`GuardianType`为`Genesis`时，`Creator`为空。

涉及存储键值对：
- guardian: `0x00 address -> amino(guardian)`

## 停止网络

停网标志存储停网原因：
- 停止网络: `0x01 -> amino(reason)`