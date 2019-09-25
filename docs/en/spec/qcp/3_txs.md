# Transactions

## TxInitQCP

[Sending TxInitQCP](../../command/qoscli.md#init-qcp) to initialize the chain information in the QOS network.

### Struct

```go
type TxInitQCP struct {
	Creator btypes.AccAddress `json:"creator"` //creator address
	QCPCA   *cert.Certificate `json:"ca_qcp"`  //QCP CA
}
```

This tx is expected to fail if:
- `creator` can not be empty, and the account must exists
- ca information must be legal
- no QCP chain with the same `chain-id` exists

[Query QCP](../../command/qoscli.md#query-qcp) can find the QCP chain information.

### Signer

`creator`

### Tx Gas

`1.8QOS`