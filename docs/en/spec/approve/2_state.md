# State

## Approve

The struct of approve we store into the db as follow:
```go
type Approve struct {
	From btypes.AccAddress `json:"from"` // approve creator
	To   btypes.AccAddress `json:"to"`   // approve receiver
	QOS  btypes.BigInt     `json:"qos"`  // QOS
	QSCs types.QSCs        `json:"qscs"` // QSCs
}
```
`MapperName`is `approve`,key-values in `Mapper`:
- approve: `0x01 from to -> amino(approve)`
