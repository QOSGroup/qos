# QCP

QOS跨链协议（QCP），初始化新的联盟链

* Struct
```go
// init QCP
type TxInitQCP struct {
	Creator btypes.Address       `json:"creator"` //创建账户
	QCPCA   *cert.Certificate    `json:"ca_qcp"`  //CA信息
}
```

需要提供创建账户和QCP证书，创建账户需在对应QOS网络中存在，证书申请参照[QCP证书](../ca.md#QCP)。

* Store
```go
QcpMapperName = "qcp"               
outSequenceKey = "sequence/out/%s"  //需要输出到"chainId"的qcp tx最大序号
outSequenceTxKey = "tx/out/%s/%d"   //需要输出到"chainId"的每个qcp tx
inSequenceKey = "sequence/in/%s"    //已经接受到来自"chainId"的qcp 的合法公钥tx最大序号
inPubkeyKey = "pubkey/in/%s"        //接受来自"chainId"
```

读写使用QCPMapper，QCPMapper在[qbase]("https://www.github.com/QOSGroup/qbase")中定义。

* valid
1. creator账户存在
2. CA信息正确，与公链保存的RootCA验证通过，未重复使用

* signer
Creator账户