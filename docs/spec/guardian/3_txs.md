# 交易

系统账户模块包含交易类型如下：

## 添加系统账户

### 结构
[发送添加系统账户交易](../../command/qoscli.md#添加系统账户)， 可以在QOS网络中添加新的系统账户，交易结构如下：
```go
type TxAddGuardian struct {
	Description string            `json:"description"` // 描述信息
	Address     btypes.AccAddress `json:"address"`     // 账户地址
	Creator     btypes.AccAddress `json:"creator"`     // 创建账户地址
}
```

### 验证

添加操作必须通过以下校验才会执行：
- `len(description)`不能大于`MaxDescriptionLen`（默认1000）
- `address`不能为空，且不存在`address`为地址的系统账户
- `creator`不能为空，且存在此地址系统账户，且类型为`Genesis`

### 签名

`creator`

### 交易费

`0`

## 删除系统账户

[发送删除系统账户交易](../../command/qoscli.md#删除系统账户)可删除QOS网络中非`Genesis`类型的系统账户。

### 结构

```go
type TxDeleteGuardian struct {
	Address   btypes.AccAddress `json:"address"`    // 系统账户地址
	DeletedBy btypes.AccAddress `json:"deleted_by"` // 删除操作账户地址
}
```

### 验证

- `address`不能为空，且存在此地址系统账户，且类型是`Ordinary`
- `deleted_by`不能为空，且存在此地址系统账户，且类型是`Genesis`

### 签名

`deleted_by`

### 交易费

`0`

## 停止网络

为了在网络受到攻击或出现重大bug，避免持币账户扩大损失，QOS赋予系统账户可以紧急停网的能力。

[发送停止网络操作](../../command/qoscli.md#停止网络)可使QOS网络在下一块开始停止打块。

### 结构

```go
type TxHaltNetwork struct {
	Guardian btypes.AccAddress `json:"guardian"` // 系统账户地址
	Reason   string            `json:"reason"`   // 停网原因
}
```

### 验证

- `reason`不能为空，`len(reason)`不能大于`MaxDescriptionLen`（默认1000）
- `guardian`不能为空，且存在此地址系统账户

### 签名

`guardian`

### 交易费

`0`