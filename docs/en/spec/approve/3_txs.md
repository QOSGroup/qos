# Transactions

We introduce approve txs here，these txs affect [State](2_state.md).

## TxCreateApprove

### Struct
[Sending TxCreateApprove](../../command/qoscli.md#create-approve) to create an `approve`. The struct of this tx is:
```go
type TxCreateApprove struct {
	types.Approve
}
```
[Approve](2_state.md#Approve) is the common struct of `TxCreateApprove`,[TxIncreaseApprove](#txincreaseapprove),[TxDecreaseApprove](#txdecreaseapprove),[TxUseApprove](#txuseapprove).

### Validations

This tx is expected to fail if:
- `from`,`to` can not be empty, can not be the same.
- `coins`are all positive and no duplicate.
- account of `from` must exists.
- no approve from `from` to `to` exists.

### Singer

`from`

### Tx Gas

`0`

## TxIncreaseApprove

### Struct
[Sending TxIncreaseApprove](../../command/qoscli.md#increase-approve) to increase an `approve`. The struct of this tx is:
```go
type TxIncreaseApprove struct {
	types.Approve
}
```
[Approve](2_state.md#Approve) is the common struct of [TxCreateApprove](#txcreateapprove),`TxIncreaseApprove`,[TxDecreaseApprove](#txdecreaseapprove),[TxUseApprove](#txuseapprove).

### Validations

This tx is expected to fail if:
- approve from `from` to `to` must exist.

### Singer

`from`

### Tx Gas

`0`

## TxDecreaseApprove

### Struct
[Sending TxDecreaseApprove](../../command/qoscli.md#decrease-approve) to decrease an `approve`. The struct of this tx is:
```go
type TxDecreaseApprove struct {
	types.Approve
}
```
[Approve](2_state.md#Approve) is the common struct of [TxCreateApprove](#txcreateapprove),[TxIncreaseApprove](#txincreaseapprove),`TxDecreaseApprove`,[TxUseApprove](#txuseapprove).

### Validations

This tx is expected to fail if:
- approve from `from` to `to` must exist.
- `coins` must lte the existing approve.

### Singer

`from`

### Tx Gas

`0`

## TxUseApprove

The approve receiver extracts the authorized QOS, QSCs to account, 
deducts the approve creator account balance accordingly, and reduces the approve.

> Approve will be deleted when there is no QOS and QSCs left after this tx done. 

### Struct
[Sending TxUseApprove](../../command/qoscli.md#use-approve) to use an `approve`. The struct of this tx is:
```go
type TxUseApprove struct {
	types.Approve
}
```
[Approve](2_state.md#Approve) is the common struct of [TxCreateApprove](#txcreateapprove),[TxIncreaseApprove](#txincreaseapprove),[TxDecreaseApprove](#txdecreaseapprove),`TxUseApprove`.

### Validations

This tx is expected to fail if:
- approve from `from` to `to` must exist.
- `coins` must lte the existing approve.

### Singer

`from`

### Tx Gas

`0`

## TxCancelApprove

Cancel an existing approve and remove from the database.

### Struct
[Sending TxCancelApprove](../../command/qoscli.md#cancel-approve) to cancel an `approve`. The struct of this tx is:
```go
type TxCancelApprove struct {
    From btypes.AccAddress `json:"from"`
    To   btypes.AccAddress `json:"to"`  
}
```

### Validations

This tx is expected to fail if:
- approve from `from` to `to` must exist.

### Singer

`from`

### Tx Gas

`0`