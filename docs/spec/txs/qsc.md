# QSC

## Struct
```
// create QSC
type TxCreateQSC struct {
	ChainID     string               `json:"chain_id"`    //chain-id
	Creator     btypes.Address       `json:"creator"`     //QSC创建账户
	Extrate     string               `json:"extrate"`     //qcs:qos汇率(amino不支持binary形式的浮点数序列化，精度同qos erc20 [.0000])
	QSCCA       *Certificate         `json:"ca_qsc"`      //CA信息
	BankerCA    *Certificate         `json:"ca_banker"`   //CA信息
	Description string               `json:"description"` //描述信息
	Accounts    []account.QOSAccount `json:"accounts"`    //初始化时接受qsc的账户
}

// issue QSC
type TxIssueQSC struct {
	QscName string         `json:"qsc_name"` //币名
	Amount  btypes.BigInt  `json:"amount"`   //金额
	Banker  btypes.Address `json:"banker"`   //banker地址
}
```

## Store
```
QSCMapperName = "qsc"       // store
QSCKey        = "qsc/[%s]"  // key，qscName
```

读写使用QSCMapper
```
type QSCMapper struct {
	*mapper.BaseMapper      // qbase BaseMapper封装
}
```
提供保存QSC（SaveQsc）、获取QSC（GetQsc）、判断QSC是否存在（Exists）方法

## Create

创建联盟链，保存联盟链代币信息，保存、初始化联盟链跨链（QCP）信息，发放联盟币到指定账户。

* valid
1. ChainID不为空，不能与现有联盟链重复
2. QSC名不能与现有联盟币重复
3. Creator账户存在
4. CA信息正确，与QOS保存的RootCA验证通过
5. Accounts可为空，仅可包含联盟链代币

* signer
Creator账户

## Issue

向联盟链的Banker发币，可重复发放，表现为联盟币累加。

* valid
1. QscName不能为空
2. Amount大于0
3. QSC存在，且名称一致
4. Banker存在，且地址一致

* signer
Banker账户