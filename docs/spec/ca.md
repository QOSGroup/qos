# CA

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

## QSC
```go
type QSCSubject struct {
	ChainId string        `json:"chain_id"` // 证书可用链
	Name    string        `json:"name"`     // 联盟币名称，大写字母、长度不超过8
	Banker  crypto.PubKey `json:"banker"`   // banker公钥，ed25519
}
```

## QCP
```go
type QCPSubject struct {
	ChainId  string `json:"chain_id"`      // 证书可用链
	QCPChain string `json:"qcp_chain"`     // 联盟链
}
```

## 申请方式
目前仅支持邮件申请方式，请将以上结构中包含的数据以及必要的公司/个人资质证明打包发送至ca@tokenxy.cn。
主题务必以“QSC证书申请“或”QCP证书申请，我们审核后将在一周内将证书文件发送至申请邮箱。