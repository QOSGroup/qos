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
	AccountAddress types.Address `json:"account_address"` // account address
	Publickey      crypto.PubKey `json:"public_key"`      // public key
	Nonce          int64         `json:"nonce"`           // identifies tx_status of an account
}
```

* Address 

采用ed25519加密，Bech32编码，"address"前缀

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