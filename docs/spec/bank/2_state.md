# 状态

账户操作`Mapper`定义在[qbase](https://github.com/QOSGroup/qbase/blob/master/account/accountmapper.go)中，`MapperName`为`acc`。

## 账户

QOS账户结构设计如下：
```go
type QOSAccount struct {
	account.BaseAccount `json:"base_account"`  
	QOS                 btypes.BigInt         `json:"qos"`  // QOS
	QSCs                QSCs                  `json:"qscs"` // QSC代币
}
```
其中`BaseAccount`定义在qbase中：
```go
type BaseAccount struct {
	AccountAddress types.AccAddress `json:"account_address"` // account address
	Publickey      crypto.PubKey    `json:"public_key"`      // public key
	Nonce          int64            `json:"nonce"`           // identifies tx_status of an account
}
```
`QOSAccount`是保存了账户地址，公钥，交易序列，QOS和QSC代币信息。

涉及存储键值对：
- account: `account:address -> amino(account)`

## 数据检查
涉及存储键值对：
- 数据检查: `0x10 height -> amino(true)`

## 锁定-释放账户
`QOS`先期作为以太坊代币，以太坊合约中包含锁定释放账户信息，锁定QOS一亿个，代币发行一年后按24个月释放完成。

数据结构：
```go
type LockInfo struct {
	LockedAccount   types.AccAddress `json:"locked_account"`   // 锁定账户地址
	Receiver        types.AccAddress `json:"receiver"`         // 接收账户地址
	TotalAmount     types.BigInt     `json:"total_amount"`     // 总锁定QOS
	ReleasedAmount  types.BigInt     `json:"released_amount"`  // 已释放QOS
	ReleaseTime     time.Time        `json:"release_time"'`    // 下一次释放时间
	ReleaseInterval int64            `json:"release_interval"` // 释放间隔，以天为单位
	ReleaseTimes    int64            `json:"release_times"`    // 释放次数
}
```
涉及存储键值对：
- 锁定-释放： `0x11 -> amino(lockinfo)`