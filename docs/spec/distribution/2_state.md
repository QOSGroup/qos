# 存储

`MapperName`为`distribution`

## 社区收益池

- `0x01 -> amino(BigInt)`

## 上一块`proposer`地址

- `0x02 -> amino(ConsAddress)`

## 待分配的QOS

待分配的通胀加交易费

- `0x04 -> amino(BigInt)`

## `delegator`收益

### `delegator`收益

`delegator`计算收益信息结构：
```go
type DelegatorEarningsStartInfo struct {
	PreviousPeriod        int64         `json:"previous_period"`       // 前收益计算点
	BondToken             btypes.BigInt `json:"bond_token"`            // 绑定tokens
	CurrentStartingHeight int64         `json:"earns_starting_height"` // 当前计算周期起始高度
	FirstDelegateHeight   int64         `json:"first_delegate_height"` // 首次委托高度
	HistoricalRewardFees  btypes.BigInt `json:"historical_rewards"`    // 累计未发放奖励
	LastIncomeCalHeight   int64         `json:"last_income_calHeight"` // 最后收益计算高度
	LastIncomeCalFees     btypes.BigInt `json:"last_income_calFees"`   // 最后一次发放收益
}
```

- `0x12 + validatorAddress + delegatorAddress -> amino(DelegatorEarningsStartInfo)`

### `delegators`某高度收益是否发放

- `0x31 + blockheight + validatorAddress + delegatorAddress -> amino(true/false)`


## `validator`收益

### `validator`历史计费点汇总收益

- `0x13 + validatorAddress + period -> amino(Fraction)`

### `validator`当前计费点收益信息

`validator`当前周期收益信息结构：
```go
type ValidatorCurrentPeriodSummary struct {
	Fees   btypes.BigInt `json:"fees"`
	Period int64         `json:"period"`
}
```

- `0x14 + validatorAddress -> amino(ValidatorCurrentPeriodSummary)`

### `validator`获得收益信息

`validator`收益信息结构：
```go
type ValidatorEcoFeePool struct {
	ProposerTotalRewardFee      btypes.BigInt `json:"proposerTotalRewardFee"`      //validator 通过proposer获取的总收益
	CommissionTotalRewardFee    btypes.BigInt `json:"commissionTotalRewardFee"`    //validator 通过投票获取的佣金总收益
	PreDistributeTotalRewardFee btypes.BigInt `json:"preDistributeTotalRewardFee"` //validator 通过投票获取的待分配金额总收益
	PreDistributeRemainTotalFee btypes.BigInt `json:"preDistributeRemainTotalFee"` //validator 待分配金额中剩余的收益
}
```

- `0x15 + validatorAddress -> amino(ValidatorEcoFeePool)`