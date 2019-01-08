package types

type Params struct {
	TotalAmount uint64 `json:"total_amount"`
	TotalBlock  uint64 `json:"total_block"`
}

func NewParams(totalAmount uint64, totalBlock uint64) Params {

	return Params{
		TotalAmount: totalAmount,
		TotalBlock:  totalBlock,
	}
}

func DefaultParams() Params {

	return NewParams(100e8, 6307200)
}
