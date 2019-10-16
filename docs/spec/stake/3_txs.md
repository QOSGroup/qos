# 交易

这里具体介绍stake模块包含的所有交易类型，这些交易操作会直接影响[存储](2_state.md)。

## 验证节点

### 创建验证节点

[发送创建验证节点交易](../../command/qoscli.md#成为验证节点)创建验证节点。

#### 结构

```go
type TxCreateValidator struct {
	Owner       btypes.AccAddress      `json:"owner"`        // 操作者, self delegator
	ConsPubKey  crypto.PubKey          `json:"cons_pub_key"` // validator公钥
	BondTokens  btypes.BigInt          `json:"bond_tokens"`  // 绑定Token数量
	IsCompound  bool                   `json:"is_compound"`  // 周期收益是否复投
	Description types.Description      `json:"description"`  // 描述信息
	Commission  types.CommissionRates  `json:"commission"`   // 佣金比例
	Delegations []types.DelegationInfo `json:"delegations"`  // 初始委托，仅在iniChainer中执行有效
}
```

#### 校验

创建验证节点需要通过以下校验：
- `moniker`不能为空，且`length(moniker) < 300`
- `owner`不为空，且不存在`owenr`创建的其他验证节点
- `cons_pub_key`不能为空，且不存在相同共识公钥的验证节点
- `description`中`logo`和`website`长度不能超过255
- `bond_tokens`必须为正，`owner`账户余额大于`bond_tokens`
- `commission`数值校验通过

#### 签名

`owner`账户

#### 交易费

`1.8QOS`

### 编辑验证节点

[发送编辑验证节点交易](../../command/qoscli.md#编辑验证节点)编辑验证节点。

#### 结构

```go
type TxModifyValidator struct {
	Owner          btypes.AccAddress `json:"owner"`           // 验证人Owner地址
	ValidatorAddr  btypes.ValAddress `json:"validator_addr"`  // 验证人地址
	Description    types.Description `json:"description"`     // 描述信息
	CommissionRate *qtypes.Dec       `json:"commission_rate"` // 佣金比例
}
```

#### 校验

编辑验证节点需要通过以下校验：
- `owner`不为空
- `validator_addr`不能为空，验证节点存在，且验证节点owner与`owner`一致
- `description`中`logo`和`website`长度不能超过255
- `commission`数值校验通过

#### 签名

`owner`账户

#### 交易费

`0.18QOS`

### 撤销验证节点

[发送撤销验证节点交易](../../command/qoscli.md#撤销验证节点)撤销验证节点。

#### 结构

```go
type TxRevokeValidator struct {
	Owner         btypes.AccAddress `json:"owner"`          // 验证人Owner地址
	ValidatorAddr btypes.ValAddress `json:"validator_addr"` // 验证人地址
}
```

#### 校验

撤销验证节点需要通过以下校验：
- `owner`不为空
- `validator_addr`不能为空，验证节点存在，且验证节点owner与`owner`一致

#### 签名

`owner`账户

#### 交易费

`18QOS`

### 激活验证节点

[发送激活验证节点交易](../../command/qoscli.md#激活验证节点)激活被撤销且还未关闭的验证节点。

#### 结构

```go
type TxActiveValidator struct {
	Owner         btypes.AccAddress `json:"owner"`          // 验证人Owner地址
	ValidatorAddr btypes.ValAddress `json:"validator_addr"` // 验证人地址
	BondTokens    btypes.BigInt     `json:"bond_tokens"`    // 增加绑定Token数量
}
```

#### 校验

激活验证节点需要通过以下校验：
- `owner`不为空
- `validator_addr`不能为空，验证节点存在，且验证节点owner与`owner`一致
- `bond_tokens`为正，且小于`owner`账户QOS余额

#### 签名

`owner`账户

#### 交易费

`0`

## 委托

### 创建委托

[发送创建委托交易](../../command/qoscli.md#委托)创建委托。

#### 结构

```go
type TxCreateDelegation struct {
	Delegator     btypes.AccAddress `json:"delegator"`      // 委托人
	ValidatorAddr btypes.ValAddress `json:"validator_addr"` // 验证人
	Amount        btypes.BigInt     `json:"amount"`         // 委托QOS数量
	IsCompound    bool              `json:"is_compound"`    // 定期收益是否复投
}
```
#### 校验

创建委托需要通过以下校验：
- `delegator`地址不能为空，且账号存在
- `validator_addr`地址不能为空，且验证节点存在
- `amount`必须为正整数
- `delegator`账户中必须有大于`amount`的QOS

通过校验并成功执行交易后，可通过[委托查询](../../command/qoscli.md#委托查询)搜索委托信息。

#### 签名

`Delegator`账户

#### 交易费

`0` 鼓励用户使用委托功能

### 修改收益复投方式

[发送修改收益复投方式交易](../../command/qoscli.md#修改收益复投方式)修改收益复投方式。

#### 结构

```go
type TxModifyCompound struct {
	Delegator     btypes.AccAddress `json:"delegator"`      // 委托人
	ValidatorAddr btypes.ValAddress `json:"validator_addr"` // 验证者
	IsCompound    bool              `json:"is_compound"`    // 周期收益是否复投: 收益发放周期内多次修改,仅最后一次生效
}
```

#### 校验

修改收益复投方式需要通过以下校验：
- `delegator`地址不能为空，且账号存在
- `validator_addr`地址不能为空，且验证节点存在
- 委托关系已经存在
- 新的复投方式不能与原来的相同

#### 签名

`Delegator`账户

#### 交易费

`0`

### 解除委托

[发送解除委托交易](../../command/qoscli.md#解除委托)，解除委托。

#### 结构

```go
type TxUnbondDelegation struct {
	Delegator     btypes.AccAddress `json:"delegator"`      // 委托人
	ValidatorAddr btypes.ValAddress `json:"validator_addr"` // 验证者
	UnbondAmount  btypes.BigInt     `json:"unbond_amount"`  // unbond数量
	UnbondAll     bool              `json:"unbond_all"`     // 是否全部解绑, 为true时覆盖UnbondAmount
}
```

#### 校验

解除委托需要通过以下校验：
- 当`unbond_all`为空时，`unbond_amount`需为正整数
- 当`unbond_all`为空时，`validator_addr`所有bonded tokens需要大于`Amount`
- `delegator`与`validator_addr`的委托关系需存在

#### 签名

`Delegator`账户

#### 交易费

`0.18QOS`

### 变更委托验证节点

[发送变更委托验证节点交易](../../command/qoscli.md#变更委托验证节点)变更委托验证节点。

#### 结构

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

#### 校验

变更委托验证节点需要通过以下校验：
- 当`redelegate_all`为空时，`Amount`需为正整数
- `from_validator_addr`需存在
- `to_validator_addr`需存在
- `delegator`与`from_validator_addr`的委托关系需存在

#### 签名

`Delegator`账户

#### 交易费

`0`
