# QSC

## Struct
```go
// create QSC
type TxCreateQSC struct {
	Creator     btypes.Address        `json:"creator"`     //QSC创建账户
	Extrate     string                `json:"extrate"`     //qcs:qos汇率(amino不支持binary形式的浮点数序列化，精度同qos erc20 [.0000])
	QSCCA       *cert.Certificate     `json:"qsc_crt"`       //CA信息
	Description string                `json:"description"` //描述信息
	Accounts    []*account.QOSAccount `json:"accounts"`
}

// issue QSC
type TxIssueQSC struct {
	QSCName string         `json:"qsc_name"` //币名
	Amount  btypes.BigInt  `json:"amount"`   //金额
	Banker  btypes.Address `json:"banker"`   //banker地址
}
```

## Store
```go
QSCMapperName = "qsc"       // store
QSCKey        = "qsc/[%s]"  // key，qscName，保存types.QSCInfo
```

读写使用QSCMapper
```go
type QSCMapper struct {
	*mapper.BaseMapper      // qbase BaseMapper封装
}
```
提供保存QSC（SaveQsc）、获取QSC（GetQsc）、判断QSC是否存在（Exists）方法

## Create

创建联盟币，发放联盟币到指定账户。
公链中拥有一定数量QOS的账户，即可发起此Tx.

* valid
1. QSCCA数据完整性，ChainId与公链ChainId一致，与公链保存的RootCA验证通过
2. QSC名不能与现有联盟币重复
3. Creator账户存在
4. CA信息正确，与QOS保存的RootCA验证通过
5. Accounts可为空，仅可包含联盟链代币

* signer
Creator账户

## Issue

向create qsc中Banker发币，可重复发放，表现为联盟币累加。

* valid
1. QscName不能为空
2. Amount大于0
3. QSC存在，且名称与CA一致
4. Banker存在，且地址与CA一致

* signer
Banker账户