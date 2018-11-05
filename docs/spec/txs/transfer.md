# 转账设计
实现多账户，多币种交易，只需保证发送和接收集合QOS、QSCs总量相等

## Struct
```
type TransItem struct {
	Address btypes.Address `json:"addr"` // 账户地址
	QOS     btypes.BigInt  `json:"qos"`  // QOS
	QSCs    types.QSCs     `json:"qscs"` // QSCs
}

type TransferTx struct {
	Senders   []TransItem `json:"senders"`   // 发送集合
	Receivers []TransItem `json:"receivers"` // 接收集合
}
```

## TX
```
// TODO
```