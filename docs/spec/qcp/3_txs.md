# 交易

这里具体介绍QCP模块包含的所有交易类型。

## 初始化联盟链

[发送初始化联盟链交易](../../command/qoscli.md#初始化联盟链)在QOS网络中初始化联盟链信息。

### 结构

```go
type TxInitQCP struct {
	Creator btypes.AccAddress `json:"creator"` //创建账户
	QCPCA   *cert.Certificate `json:"ca_qcp"`  //CA信息
}
```

### 校验

初始化联盟链需要通过以下校验：
- `creator`地址均不能为空，且账户必须存在
- 证书文件校验通过
- 不存在同名的已初始化的联盟链信息

通过校验并成功执行交易后，可通过[查询联盟链](../../command/qoscli.md#查询联盟链)查询此联盟链相关信息。

### 签名

`creator`账户

### 交易费

`1.8QOS`