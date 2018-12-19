# QCP

* Struct
```
// init QCP
type TxInitQCP struct {
	Creator btypes.Address       `json:"creator"` //创建账户
	QCPCA   *cert.Certificate    `json:"ca_qcp"`  //CA信息
}
```

* Store
```
QcpMapperName = "qcp"               //需要输出到"chainId"的qcp tx最大序号
outSequenceKey = "sequence/out/%s"  //需要输出到"chainId"的每个qcp tx
outSequenceTxKey = "tx/out/%s/%d"   //已经接受到来自"chainId"的qcp 的合法公钥tx最大序号
inSequenceKey = "sequence/in/%s"    //接受来自"chainId"
inPubkeyKey = "pubkey/in/%s"
```

读写使用QCPMapper，QCPMapper在[qbase]("https://www.github.com/QOSGroup/qbase")中定义。

* valid
1. creator账户存在
2. CA信息正确，与公链保存的RootCA验证通过，未重复使用

* signer
Creator账户