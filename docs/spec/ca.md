# CA

QOS 证书源码参照[kepler](https://github.com/QOSGroup/kepler.git)，可编译可执行命令工具`kepler`。

## 概念

QOS创建QSC、初始化QCP需要申请CA，获得`qsc.crt`、`qcp.crt`等证书文件

`qoscli tx create-qsc`、`qoscli tx init-qcp`需要传入`qsc.crt`、`qcp.crt`，对应数据结构：

```go
type CertificateSigningRequest struct {
	Subj      Subject       `json:"subj"`
	IsCa      bool          `json:"is_ca"`
	NotBefore time.Time     `json:"not_before"`
	NotAfter  time.Time     `json:"not_after"`
	PublicKey crypto.PubKey `json:"public_key"`
}
```
不同证书，Subject不一样。

### QSC
```go
type QSCSubject struct {
	ChainId string        `json:"chain_id"` // 证书可用链
	Name    string        `json:"name"`     // 联盟币名称，大写字母、长度不超过8
	Banker  crypto.PubKey `json:"banker"`   // banker公钥，ed25519
}
```

### QCP
```go
type QCPSubject struct {
	ChainId  string `json:"chain_id"`      // 证书可用链
	QCPChain string `json:"qcp_chain"`     // 联盟链
}
```

## 申请方式

证书中心提供rpc供开发者申请和查询，申请成功后，等待CA管理人员审核，通过后证书会通过邮件的形式发送到申请邮箱。

### QCP

申请：
```shell script
curl -X POST "{baseUrl}/qcp/apply?qcpChainId={QCP chain-id}&qosChainId={QOS mainnet chain-id}&qcpPub={QCP public key}&phone={mobile}&email={email}&info={apply info}" -H "accept: application/json"
```
必填参数：

- qcpChainId 联盟链chain-id
- qosChainId QOS主网chain-id
- qcpPub 联盟链公钥信息，eg.{"type":"tendermint/PubKeyEd25519","value":"VxlG2mFfbAaWSqIOuntAB3wojZwQz+sK070LC4ZBAeg="}
- phone 手机号
- email 邮箱
- info 备注信息

查询：
```shell script
curl -X GET "{baseUrl}/qcp/apply?phone={mobile}&email={email}" -H "accept: application/json"
```

必填参数：

- phone 手机号
- email 邮箱

### QSC

申请：
```shell script
curl -X POST "{baseUrl}/qsc/apply?qscName={QSC name}&qosChainId={QOS mainnet chain-id}&qscPub={QSC public key}&bankerPub={QSC banker public key}&phone={mobile}&email={email}&info={apply info}" -H "accept: application/json"
```
必填参数：

- qscName QSC代币名称，8个字符长度以内，大写
- qosChainId QOS主网chain-id
- qscPub QSC公钥信息，eg.{"type":"tendermint/PubKeyEd25519","value":"VxlG2mFfbAaWSqIOuntAB3wojZwQz+sK070LC4ZBAeg="}
- bankerPub QSC发币账户公钥信息，eg.{"type":"tendermint/PubKeyEd25519","value":"VxlG2mFfbAaWSqIOuntAB3wojZwQz+sK070LC4ZBAeg="}
- phone 手机号
- email 邮箱
- info 备注信息

查询：
```shell script
curl -X GET "{baseUrl}/qsc/apply?phone={mobile}&email={email}" -H "accept: application/json"
```

必填参数：

- phone 手机号
- email 邮箱