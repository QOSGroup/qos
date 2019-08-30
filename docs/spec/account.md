# 账户设计
QOS账户包含账户地址，PubKey，Nonce，QOS法定代币以及QSCs联盟币集合

* Struct
```go
type QOSAccount struct {
	account.BaseAccount `json:"base_account"`       // inherits BaseAccount
	QOS                 btypes.BigInt `json:"qos"`  // coins in public chain
	QSCs                types.QSCs    `json:"qscs"` // varied QSCs
}

// BaseAccount in qbase.
type BaseAccount struct {
	AccountAddress types.AccAddress `json:"account_address"` // account address
	Publickey      crypto.PubKey `json:"public_key"`      // public key
	Nonce          int64         `json:"nonce"`           // identifies tx_status of an account
}
```

* AccAddress

采用ed25519加密，Bech32编码，"qosacc"前缀

* QOS地址区分如下

| HRP               | Definition                            |
|-------------------|---------------------------------------|
| qosacc            | QOS账户地址                            |
| qosaccpub         | QOS账户公钥                            |
| qoscons           | QOS验证人共识地址                       |
| qosconspub        | QOS验证人共识公钥                       |
| qosval            | QOS验证人地址                          |
| qosvalpub         | QOS验证人公钥                          |


* QOS

qos独立于其他联盟币

BigInt：-(2^255-1) to 2^255-1

* QSCs

联盟币
```go
type QSCs = types.BaseCoins

type BaseCoins []*BaseCoin      // in qbase
type BaseCoin struct {          // in qbase
	Name   string `json:"coin_name"`
	Amount BigInt `json:"amount"`
}
```
