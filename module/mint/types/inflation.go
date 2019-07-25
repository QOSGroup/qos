package types

import "time"

type InflationPhrase struct {
	EndTime       time.Time `json:"end_time"`
	TotalAmount   uint64    `json:"total_amount"`
	AppliedAmount uint64    `json:"applied_amount"`
}
