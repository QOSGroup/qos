# 概念

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

默认通胀规则：
```go
[
    {
        time.Date(2023, 10, 20, 0, 0, 0, 0, time.UTC),
        types.NewInt(25.5e12),
        types.ZeroInt(),
    },
    {
        time.Date(2027, 10, 20, 0, 0, 0, 0, time.UTC),
        types.NewInt(12.75e12),
        types.ZeroInt(),
    },
    {
        time.Date(2031, 10, 20, 0, 0, 0, 0, time.UTC),
        types.NewInt(6.375e12),
        types.ZeroInt(),
    },
    {
        time.Date(2035, 10, 20, 0, 0, 0, 0, time.UTC),
        types.NewInt(3.1875e12),
        types.ZeroInt(),
    },
    {
        time.Date(2039, 10, 20, 0, 0, 0, 0, time.UTC),
        types.NewInt(1.59375e12),
        types.ZeroInt(),
    },
    {
        time.Date(2043, 10, 20, 0, 0, 0, 0, time.UTC),
        types.NewInt(0.796875e12),
        types.ZeroInt(),
    },
    {
        time.Date(2047, 10, 20, 0, 0, 0, 0, time.UTC),
        types.NewInt(0.796875e12),
        types.ZeroInt(),
    }
]
```
每四年一个通胀阶段，共七个通胀阶段。