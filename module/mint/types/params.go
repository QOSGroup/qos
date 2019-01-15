package types

import "time"

type Params struct {
	Phrases []InflationPhrase `jason:"inflation_phrases"`
}

type InflationPhrase struct {
	EndTime       time.Time `jason:"endtime"`
	TotalAmount   uint64    `json:"total_amount"`
	AppliedAmount uint64    `json:"applied_amount"`
}

func NewParams(phrases []InflationPhrase) Params {
	return Params{phrases}
}

func DefaultParams() Params {
	return NewParams(
		[]InflationPhrase{
			InflationPhrase{
				time.Date(2023,1,1,0,0,0,0,time.UTC),
				2.5e8,
				0,
			},
			InflationPhrase{
				time.Date(2027,1,1,0,0,0,0,time.UTC),
				12.75e8,
				0,
			},
			InflationPhrase{
				time.Date(2031,1,1,0,0,0,0,time.UTC),
				6.375e8,
				0,
			},
			InflationPhrase{
				time.Date(2035,1,1,0,0,0,0,time.UTC),
				3.185e8,
				0,
			},
		},
	)
}
