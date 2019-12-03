# 交易

这里具体介绍QSC模块包含的所有交易类型。

## 创建QSC

[发送创建QSC交易](../../command/qoscli.md#创建QSC)创建新代币。

### 结构

```go
type TxCreateQSC struct {
	Creator      btypes.AccAddress    `json:"creator"`       // QSC创建账户
	ExchangeRate string               `json:"exchange_rate"` // qcs:qos汇率
	QSCCA        *cert.Certificate    `json:"qsc_crt"`       // CA信息
	Description  string               `json:"description"`   // 描述信息
	Accounts     []*qtypes.QOSAccount `json:"accounts"`      // 初始账户
}
```

### 校验

创建QSC需要通过以下校验：
- `creator`地址均不能为空，且账户必须存在
- `description`最大长度1000
- `exchange_rate`为浮点数
- 证书文件校验通过
- 初始账户校验，只能包含即将初始化的代币
- 不存在同名的已初始化的代币信息

通过校验并成功执行交易后，可通过[查询QSC](../../command/qoscli.md#查询QSC)搜索代币信息。

### 签名

`creator`账户

### 交易费

`1.8QOS`

## 发行QSC

[发送发行QSC交易](../../command/qoscli.md#发行QSC)发行代币。

### 结构

```go
type TxIssueQSC struct {
	QSCName string            `json:"qsc_name"` //币名
	Amount  btypes.BigInt     `json:"amount"`   //币量
	Banker  btypes.AccAddress `json:"banker"`   //banker地址
}

```

### 校验

发行QSC需要通过以下校验：
- `amount`必须为正
- 代币对应的QSC信息存在，且`banker`一致

### 签名

`banker`账户

### 交易费

`0.18QOS`