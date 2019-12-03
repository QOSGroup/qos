# State

`MapperName` is `mint`

## First block time

Here we save the time of second block:

`first_block_time -> amino(first_block_time)`

## Total mint

`total_mint_qos -> amino(total_mint_qos)`

Using [query total applied](../../command/qoscli.md#query-total-applied) to see the value of `total_mint_qos`.

## Total issue

`total_qos -> amino(total_qos)`

Using [query total inflation](../../command/qoscli.md#query-total-inflation) to see the value of `total_qos`.

## Inflation rules

struct:

```go
type InflationPhrases []InflationPhrase

type InflationPhrase struct {
	EndTime       time.Time    `json:"end_time"`       // 结束时间
	TotalAmount   types.BigInt `json:"total_amount"`   // 通胀总量
	AppliedAmount types.BigInt `json:"applied_amount"` // 已发行总量
}
```

key-value:

`phrases -> amoni(InflationPhrases)`

Using [query inflation rules](../../command/qoscli.md#query-inflation-rules) to see the inflation rules.