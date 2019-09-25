# State

`MapperName` is `qsc`

## Root CA

Before the QOS network starts, we configure the public key information of the CA center for issuing QSC certificates in `genesis.json`. 
After the network starts, it will be saved as follows:

- rootca: `rootca -> amino(pubKey)` 

## QSC information

struct:
```go
type QSCInfo struct {
	Name         string            `json:"name"`          //QSC token name
	ChainId      string            `json:"chain_id"`      //chain id of mainnet for holding this token
	ExchangeRate string            `json:"exchange_rate"` //qcs:qos exchange rate
	Description  string            `json:"description"`   //description
	Banker       btypes.AccAddress `json:"banker"`        //banker public key
	TotalAmount  btypes.BigInt     `json:"total_amount"`  //total issue amount
}
```

- approve: `qsc name -> amino(qscInfo)`
