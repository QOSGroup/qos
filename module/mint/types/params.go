package types

import "time"

type Params struct {
	Phrases []InflationPhrase `json:"inflation_phrases"`
}

func NewMintParams(phrases []InflationPhrase) Params {
	return Params{phrases}
}

func DefaultMintParams() Params {
	return NewMintParams(
		[]InflationPhrase{
			{
				time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				2.5e12, //mul(10^4),
				0,
			},
			{
				time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC),
				12.75e12, //mul(10^4),
				0,
			},
			{
				time.Date(2031, 1, 1, 0, 0, 0, 0, time.UTC),
				6.375e12, //mul(10^4),
				0,
			},
			{
				time.Date(2035, 1, 1, 0, 0, 0, 0, time.UTC),
				3.185e12, //mul(10^4),
				0,
			},
		},
	)
}
