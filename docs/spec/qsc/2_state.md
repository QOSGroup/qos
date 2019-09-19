# 存储

操作QSC的`MapperName`为`qsc`。

## 根证书

QOS网络启动前需要在`genesis.json`中配置好CA中心用于签发QSC证书的公钥信息。网络启动后会如下保存：

- rootca: `rootca -> amino(pubKey)` 

## QSC信息

QSC代币存储结构如下：
```go
type QSCInfo struct {
	Name         string            `json:"name"`          //币名
	ChainId      string            `json:"chain_id"`      //证书可用链
	ExchangeRate string            `json:"exchange_rate"` //qcs:qos汇率
	Description  string            `json:"description"`   //描述信息
	Banker       btypes.AccAddress `json:"banker"`        //Banker PubKey
	TotalAmount  btypes.BigInt     `json:"total_amount"`  //发行总量
}
```

- approve: `qsc name -> amino(qscInfo)`
