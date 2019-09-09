# Transactions

Guardian module contains  transactions as follow:

## TxAddGuardian

[Sending TxAddGuardian](../../command/qoscli.md#转账) to add a new guardian int the QOS network.

### Struct

```go
type TxAddGuardian struct {
	Description string            `json:"description"` // description
	Address     btypes.AccAddress `json:"address"`     // address
	Creator     btypes.AccAddress `json:"creator"`     // address of creator
}
```

### Validations

This tx is expected to fail if:
- `len(description)` must lte `MaxDescriptionLen`(default 1000)
- `address` can not be empty, and no guardian with this address exists
- `creator`can not be empty, and `Genesis` guardian with this address exists

### Singer

`creator`

### Tx Gas

`0`

## TxDeleteGuardian

[Sending TxDeleteGuardian](../../command/qoscli.md#数据检查) to remove `Oridinary` guardian.

### 结构
```go
type TxDeleteGuardian struct {
	Address   btypes.AccAddress `json:"address"`    // address of guardian to deleted
	DeletedBy btypes.AccAddress `json:"deleted_by"` // address of guardian who execute this transaction
}
```

### Validations

- `address` can not be empty, and `Ordinary` guardian with this address exists
- `deleted_by` can not be empty, and `Genesis` guardian with this address exists

### Singer

`deleted_by`

### Tx Gas

`0`

## TxHaltNetwork

When the network is under attack or a major bug occurs, QOS gives guardian the ability to stop the network in an emergency in order to avoid the loss of the currency account.

[Sending TxHaltNetwork](../../command/qoscli.md#停止网络) will stop the QOS network.

### Struct
```go
type TxHaltNetwork struct {
	Guardian btypes.AccAddress `json:"guardian"` // 
	Reason   string            `json:"reason"`   // reason for halting the network
}
```

### Validations

- `reason` can not be empty, `len(reason)` must lte `MaxDescriptionLen`(default 1000)
- `guardian` can not be empty, and guardian with this address exists

### Singer

`guardian`

### Tx Gas

`0`