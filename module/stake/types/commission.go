package types

import (
	"fmt"
	btypes "github.com/QOSGroup/qbase/types"
	qtypes "github.com/QOSGroup/qos/types"
	"time"
)

type (
	// Commission defines a commission parameters for a given validator.
	Commission struct {
		CommissionRates `json:"commission_rates" yaml:"commission_rates"`
		UpdateTime      time.Time `json:"update_time"` // the last time the commission rate was changed
	}

	// CommissionRates defines the initial commission rates to be used for creating a
	// validator.
	CommissionRates struct {
		Rate          qtypes.Dec `json:"rate"`            // the commission rate charged to delegators, as a fraction
		MaxRate       qtypes.Dec `json:"max_rate"`        // maximum commission rate which validator can ever charge, as a fraction
		MaxChangeRate qtypes.Dec `json:"max_change_rate"` // maximum daily increase of the validator commission, as a fraction
	}
)

// NewCommissionRates returns an initialized validator commission rates.
func NewCommissionRates(rate, maxRate, maxChangeRate qtypes.Dec) CommissionRates {
	return CommissionRates{
		Rate:          rate,
		MaxRate:       maxRate,
		MaxChangeRate: maxChangeRate,
	}
}

// NewCommission returns an initialized validator commission.
func NewCommission(rate, maxRate, maxChangeRate qtypes.Dec) Commission {
	return Commission{
		CommissionRates: NewCommissionRates(rate, maxRate, maxChangeRate),
		UpdateTime:      time.Unix(0, 0).UTC(),
	}
}

// NewCommission returns an initialized validator commission with a specified
// update time which should be the current block BFT time.
func NewCommissionWithTime(rate, maxRate, maxChangeRate qtypes.Dec, updatedAt time.Time) Commission {
	return Commission{
		CommissionRates: NewCommissionRates(rate, maxRate, maxChangeRate),
		UpdateTime:      updatedAt,
	}
}

// Equal checks if the given Commission object is equal to the receiving
// Commission object.
func (c Commission) Equal(c2 Commission) bool {
	return c.Rate.Equal(c2.Rate) &&
		c.MaxRate.Equal(c2.MaxRate) &&
		c.MaxChangeRate.Equal(c2.MaxChangeRate) &&
		c.UpdateTime.Equal(c2.UpdateTime)
}

// String implements the Stringer interface for a Commission.
func (c Commission) String() string {
	return fmt.Sprintf("rate: %s, maxRate: %s, maxChangeRate: %s, updateTime: %s",
		c.Rate, c.MaxRate, c.MaxChangeRate, c.UpdateTime,
	)
}

// Validate performs basic sanity validation checks of initial commission
// parameters. If validation fails, an SDK error is returned.
func (c CommissionRates) Validate() btypes.Error {
	switch {
	case c.MaxRate.LT(qtypes.ZeroDec()):
		// max rate cannot be negative
		return ErrCommissionNegative(DefaultCodeSpace)

	case c.MaxRate.GT(qtypes.OneDec()):
		// max rate cannot be greater than 1
		return ErrCommissionHuge(DefaultCodeSpace)

	case c.Rate.LT(qtypes.ZeroDec()):
		// rate cannot be negative
		return ErrCommissionNegative(DefaultCodeSpace)

	case c.Rate.GT(c.MaxRate):
		// rate cannot be greater than the max rate
		return ErrCommissionGTMaxRate(DefaultCodeSpace)

	case c.MaxChangeRate.LT(qtypes.ZeroDec()):
		// change rate cannot be negative
		return ErrCommissionChangeRateNegative(DefaultCodeSpace)

	case c.MaxChangeRate.GT(c.MaxRate):
		// change rate cannot be greater than the max rate
		return ErrCommissionChangeRateGTMaxRate(DefaultCodeSpace)
	}

	return nil
}

// ValidateNewRate performs basic sanity validation checks of a new commission
// rate. If validation fails, an SDK error is returned.
func (c Commission) ValidateNewRate(newRate qtypes.Dec, blockTime time.Time) btypes.Error {
	switch {
	case blockTime.Sub(c.UpdateTime).Hours() < 24:
		// new rate cannot be changed more than once within 24 hours
		return ErrCommissionUpdateTime(DefaultCodeSpace)

	case newRate.LT(qtypes.ZeroDec()):
		// new rate cannot be negative
		return ErrCommissionNegative(DefaultCodeSpace)

	case newRate.GT(c.MaxRate):
		// new rate cannot be greater than the max rate
		return ErrCommissionGTMaxRate(DefaultCodeSpace)

	case newRate.Sub(c.Rate).GT(c.MaxChangeRate):
		// new rate % points change cannot be greater than the max change rate
		return ErrCommissionGTMaxChangeRate(DefaultCodeSpace)
	}

	return nil
}
