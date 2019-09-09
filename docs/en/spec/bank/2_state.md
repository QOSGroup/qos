# State

Bank module use `Mapper` defined in [qbase](https://github.com/QOSGroup/qbase/blob/master/account/accountmapper.go),`MapperName` is `acc`.

## Account

`QOSAccount` is to save the account address, public key, transaction sequence, QOS and QSC token information.

struct:
```go
type QOSAccount struct {
	account.BaseAccount `json:"base_account"`  
	QOS                 btypes.BigInt         `json:"qos"`  // QOS
	QSCs                QSCs                  `json:"qscs"` // QSC tokens
}
```
`BaseAccount` is defined in [qbase](https://github.com/QOSGroup/qbase/blob/master/account/account.go):
```go
type BaseAccount struct {
	AccountAddress types.AccAddress `json:"account_address"` // account address
	Publickey      crypto.PubKey    `json:"public_key"`      // public key
	Nonce          int64            `json:"nonce"`           // identifies tx_status of an account
}
```

- account: `account:address -> amino(account)`

## InvariantCheck

This is a space for holding invariant check flag.

- invariant check: `0x10 height -> amino(true)`

## LockInfo

`QOS` as the Ethereum token, the Ethereum contract contained the information of the lock release account, and locked 100 million QOS. 
The token was released in 24 months after one year of `2018-08-06`.

Struct:
```go
type LockInfo struct {
	LockedAccount   types.AccAddress `json:"locked_account"`   // lock address
	Receiver        types.AccAddress `json:"receiver"`         // receive address
	TotalAmount     types.BigInt     `json:"total_amount"`     // total lock amount
	ReleasedAmount  types.BigInt     `json:"released_amount"`  // total released amount
	ReleaseTime     time.Time        `json:"release_time"'`    // next release time
	ReleaseInterval int64            `json:"release_interval"` // release interval, in days
	ReleaseTimes    int64            `json:"release_times"`    // release times
}
```

- lock info `0x11 -> amino(lockinfo)`