# 转账设计

实现多账户，多币种交易，只需保证发送和接收集合QOS、QSCs总量相等

* Struct
```go
type TransItem struct {
	Address btypes.Address `json:"addr"` // 账户地址
	QOS     btypes.BigInt  `json:"qos"`  // QOS
	QSCs    types.QSCs     `json:"qscs"` // QSCs
}

type TransItems []TransItem 

type TxTransfer struct {
	Senders   TransItems `json:"senders"`   // 发送集合
	Receivers TransItems `json:"receivers"` // 接收集合
}
```

* valid

1. Senders、Receivers不为空，地址不重复，币值大于0
2. Senders中账号对应币种、币值足够转出
3. Senders、Receivers 币值总和对应币种相等

* signer

Senders中账户按顺序依次对交易签名