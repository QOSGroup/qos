# Transactions

## TxCreateQSC

[Sending TxCreateQSC](../../command/qoscli.md#create-qsc) to create QSC token。

### Struct

```go
type TxCreateQSC struct {
	Creator      btypes.AccAddress    `json:"creator"`       // creator address
	ExchangeRate string               `json:"exchange_rate"` // qcs:qos exchange rate
	QSCCA        *cert.Certificate    `json:"qsc_crt"`       // CA, applied from CA center
	Description  string               `json:"description"`   // description
	Accounts     []*qtypes.QOSAccount `json:"accounts"`      // initial QSC holders
}
```

### Validations

This tx is expected to fail if:
- `creator` can not be empty, and the account must exists
- max `len(description)` lte 1000
- `exchange_rate` must be float value
- ca information must be legal
- `accounts` can only hold token that this tx will creating
- no QSC token info with the same `name` exists

[Query QSC](../../command/qoscli.md#query-qsc) can find QSC tokens that have been created.

### Signer

`creator`

### Tx Gas

`1.8QOS`

## TxIssueQSC

[Sending TxIssueQSC](../../command/qoscli.md#issue-qsc) to issue QSC token.

### Struct

```go
type TxIssueQSC struct {
	QSCName string            `json:"qsc_name"` //QSC name
	Amount  btypes.BigInt     `json:"amount"`   //issue amount
	Banker  btypes.AccAddress `json:"banker"`   //banker address
}

```

### Validations

This tx is expected to fail if:
- `amount` must be positive
- QSC information with `qsc_name` must exists and has the same address with `banker`

### Signer

`banker`

### Tx Gas

`0.18QOS`