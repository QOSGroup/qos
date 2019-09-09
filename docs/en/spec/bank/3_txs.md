# Transactions

Bank module contains transfer, invariant check transactions.

## TxTransfer

QOS supports multi-account multi-currency transfers.

### Struct
[Sending TxTransfer](../../command/qoscli.md#transfer) is simple and convenient.
```go
type TxTransfer struct {
	Senders   types.TransItems `json:"senders"`   // trans items of sender
	Receivers types.TransItems `json:"receivers"` // trans items of receiver
}

type TransItems []TransItem

type TransItem struct {
	Address btypes.AccAddress `json:"addr"` // address
	QOS     btypes.BigInt     `json:"qos"`  // QOS
	QSCs    types.QSCs        `json:"qscs"` // QSC tokens
}
```

### Validations

This tx is expected to fail if:
- The send list/receive list cannot be empty. There can be no duplicate addresses in each list. The currency values ​​are equivalent.
- Send list + receive list account number is less than `MaxTransLen` (default: 500).
- Send/receive QOS, QSC tokens are positive.
- Send account balance enough for this transfer.

### Signer

All sending accounts

### Tx Gas

`0.018QOS`

### Gas Payer

The first send account

## TxInvariantCheck

Check all data in the QOS network.

### Struct
[Sending TxInvariantCheck](../../command/qoscli.md#invariant-check) will send a specific event for checking all data in teh QOS network.
```go
type TxInvariantCheck struct {
	Sender btypes.AccAddress `json:"sender"` // address of tx sender
}
```

### Validations

- `sender` can not be empty.

### Signer

`sender`

### Tx Gas

`200000QOS`

::: warning Note 
In order to prevent developers from arbitrary sending invariant verification operations affecting the healthy operation of the whole network, a large transaction fee is set here. If the invariant verification indicates that there is no abnormality in the whole network data, the transaction fee is normally deducted. Otherwise, the whole network will be down and the transaction fee will not be deducted. QOS encourages all currency users to monitor the normal operation of the QOS network and report abnormal conditions in a timely manner.
:::