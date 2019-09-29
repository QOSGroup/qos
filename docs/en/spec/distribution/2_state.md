# State

`MapperName` is `distribution`

## Community fee pool

- `0x01 -> amino(BigInt)`

## Proposer address

- `0x02 -> amino(ConsAddress)`

## QOS to be assigned

Inflation plus transaction fee to be distributed:

- `0x04 -> amino(BigInt)`

## Delegator income

### `delegator` income

struct:
```go
type DelegatorEarningsStartInfo struct {
	PreviousPeriod        int64         `json:"previous_period"`       // pre calculation point
	BondToken             btypes.BigInt `json:"bond_token"`            // tokens
	CurrentStartingHeight int64         `json:"earns_starting_height"` // current calculation cycle start height
	FirstDelegateHeight   int64         `json:"first_delegate_height"` // first delegate height
	HistoricalRewardFees  btypes.BigInt `json:"historical_rewards"`    // accumulated no awards
	LastIncomeCalHeight   int64         `json:"last_income_calHeight"` // latest income calculation height
	LastIncomeCalFees     btypes.BigInt `json:"last_income_calFees"`   // latest income
}
```

- `0x12 + validatorAddress + delegatorAddress -> amino(DelegatorEarningsStartInfo)`

### whether a certain high income is paid

- `0x31 + blockheight + validatorAddress + delegatorAddress -> amino(true/false)`


## Validator

### validator historical income point

- `0x13 + validatorAddress + period -> amino(Fraction)`

### validator current income point information 

struct:
```go
type ValidatorCurrentPeriodSummary struct {
	Fees   btypes.BigInt `json:"fees"`
	Period int64         `json:"period"`
}
```

- `0x14 + validatorAddress -> amino(ValidatorCurrentPeriodSummary)`

### validator income

struct:
```go
type ValidatorEcoFeePool struct {
	ProposerTotalRewardFee      btypes.BigInt `json:"proposerTotalRewardFee"`      // rewards for proposer
	CommissionTotalRewardFee    btypes.BigInt `json:"commissionTotalRewardFee"`    // rewards for commissions
	PreDistributeTotalRewardFee btypes.BigInt `json:"preDistributeTotalRewardFee"` // rewards for voting
	PreDistributeRemainTotalFee btypes.BigInt `json:"preDistributeRemainTotalFee"` // reward to be distributed
}
```

- `0x15 + validatorAddress -> amino(ValidatorEcoFeePool)`