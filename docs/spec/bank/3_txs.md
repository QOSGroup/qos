# 交易

Bank 模块包含交易类型如下：

## 转账

QOS支持一次多账户多币种转账

### 结构
[发送转账交易](../../command/qoscli.md#转账)， 转账交易体结构：
```go
type TxTransfer struct {
	Senders   types.TransItems `json:"senders"`   // 发送集合
	Receivers types.TransItems `json:"receivers"` // 接收集合
}

type TransItems []TransItem

type TransItem struct {
	Address btypes.AccAddress `json:"addr"` // 账户地址
	QOS     btypes.BigInt     `json:"qos"`  // QOS
	QSCs    types.QSCs        `json:"qscs"` // QSCs
}
```

### 验证

转账交易必须通过以下校验才会执行：
- 发送列表/接收列表均不能为空，每个列表中均不能有重复地址，币种币值对应相等
- 发送列表+接收列表账户数小于`MaxTransLen`（默认：500）
- 发送/接收QOS，QSC代币均为正
- 发送账户余额足够本次转账

### 签名

所有发送账户

### 交易费

`0.018QOS`

## 数据检查

检查QOS网络中所有状态数据。

### 结构
[发送数据检查交易](../../command/qoscli.md#数据检查)，其结构如下：
```go
type TxInvariantCheck struct {
	Sender btypes.AccAddress `json:"sender"` // 发送交易账户地址
}
```

### 验证

- `sender`不能为空.

### 签名

`sender`

### 交易费

`200000QOS`

::: warning Note 
为防止开发者随意发送数据校验操作影响全网健康运行，此处设置了较大的交易费，如果数据校验表明全网数据无异常，交易费正常扣除。否则全网宕机，交易费不会扣除。QOS鼓励所有持币用户监控QOS网络的正常运转，及时报告异常情况。
:::