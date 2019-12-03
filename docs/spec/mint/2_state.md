# 存储

`MapperName`为`mint`

## 第一块出块时间

这里存的是第二块的时间：

`first_block_time -> amino(first_block_time)`

## 总流通

`total_mint_qos -> amino(total_mint_qos)`

可通过[查询流通总量](../../command/qoscli.md#流通总量查询)指令查询当前流通总量。

## 发行总量

`total_qos -> amino(total_qos)`

可通过[查询发行总量](../../command/qoscli.md#发行总量查询)指令查询QOS网络发行总量。

## 通胀规则

通胀规则结构：

```go
// 通胀规则
type InflationPhrases []InflationPhrase

// 通胀阶段
type InflationPhrase struct {
	EndTime       time.Time    `json:"end_time"`       // 结束时间
	TotalAmount   types.BigInt `json:"total_amount"`   // 通胀总量
	AppliedAmount types.BigInt `json:"applied_amount"` // 发行总量
}
```

存储：

`phrases -> amoni(InflationPhrases)`

可通过[查询通胀规则](../../command/qoscli.md#通胀规则)指令查询QOS网络通胀规则。