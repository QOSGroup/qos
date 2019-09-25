# Concepts

## Inflation rules

struct:

```go
// inflation rules
type InflationPhrases []InflationPhrase

// inflation phrase
type InflationPhrase struct {
	EndTime       time.Time    `json:"end_time"`       // end time
	TotalAmount   types.BigInt `json:"total_amount"`   // total amount
	AppliedAmount types.BigInt `json:"applied_amount"` // total applied
}
```

default inflation rules:
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