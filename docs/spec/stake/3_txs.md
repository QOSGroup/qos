# 交易

这里具体介绍stake模块包含的所有交易类型，这些交易操作会直接影响[存储](2_state.md)。

## 创建委托

[发送创建委托交易](../../command/qoscli.md#委托)创建委托。

### 结构

```go
type TxCreateDelegation struct {
	Delegator     btypes.AccAddress `json:"delegator"`      // 委托人
	ValidatorAddr btypes.ValAddress `json:"validator_addr"` // 验证人
	Amount        btypes.BigInt     `json:"amount"`         // 委托QOS数量
	IsCompound    bool              `json:"is_compound"`    // 定期收益是否复投
}
```
### 校验

创建委托需要通过以下校验：
- `Delegator`地址不能为空，且账号存在
- `Validator`地址不能为空，且验证节点存在
- `Amount`必须为正整数
- `Delegator`账户中必须有大于`Amount`的QOS

通过校验并成功执行交易后，可通过[委托查询](../../command/qoscli.md#委托查询)搜索委托信息。

### 签名

`Delegator`账户

### 交易费

`0` 鼓励用户使用委托功能

## 修改收益复投方式

[发送修改收益复投方式交易](../../command/qoscli.md#修改收益复投方式)修改收益复投方式。

### 结构

```go
type TxModifyCompound struct {
	Delegator     btypes.AccAddress `json:"delegator"`      // 委托人
	ValidatorAddr btypes.ValAddress `json:"validator_addr"` // 验证者
	IsCompound    bool              `json:"is_compound"`    // 周期收益是否复投: 收益发放周期内多次修改,仅最后一次生效
}
```

### 校验

修改收益复投方式需要通过以下校验：
- `Delegator`地址不能为空，且账号存在
- `Validator`地址不能为空，且验证节点存在
- 委托关系已经存在
- 新的复投方式不能与原来的相同

### 签名

`Delegator`账户

### 交易费

`0` 鼓励用户使用预授权功能

## 解除委托

[发送解除委托交易](../../command/qoscli.md#解除委托)，解除委托。

### 结构

```go
type TxUnbondDelegation struct {
	Delegator     btypes.AccAddress `json:"delegator"`      // 委托人
	ValidatorAddr btypes.ValAddress `json:"validator_addr"` // 验证者
	UnbondAmount  btypes.BigInt     `json:"unbond_amount"`  // unbond数量
	UnbondAll     bool              `json:"unbond_all"`     // 是否全部解绑, 为true时覆盖UnbondAmount
}
```

### 校验

创建预授权需要通过以下校验：
- 当`UnbondAll`为空时，`Amount`需为正整数
- 当`UnbondAll`为空时，`ValidatorAddr`所有bonded tokens需要大于`Amount`
- `Delegator`与`ValidatorAddr`的委托关系需存在

### 签名

`Delegator`账户

### 交易费

`0.18QOS`

## 变更委托验证节点

[发送变更委托验证节点交易](../../command/qoscli.md#变更委托验证节点)变更委托验证节点。

### 结构

```go
type TxCreateReDelegation struct {
	Delegator         btypes.AccAddress `json:"delegator"`           // 委托人
	FromValidatorAddr btypes.ValAddress `json:"from_validator_addr"` // 原委托验证人
	ToValidatorAddr   btypes.ValAddress `json:"to_validator_addr"`   // 现委托验证人
	Amount            btypes.BigInt     `json:"amount"`              // 委托数量
	RedelegateAll     bool              `json:"redelegate_all"`      // 转委托所有
	Compound          bool              `json:"compound"`            // 复投
}
```

### 校验

创建预授权需要通过以下校验：
- 当`RedelegateAll`为空时，`Amount`需为正整数
- `FromValidatorAddr`需存在
- `ToValidatorAddr`需存在
- `Delegator`与`FromValidatorAddr`的委托关系需存在

### 签名

`Delegator`账户

### 交易费

`0` 鼓励用户使用变更委托验证节点功能
