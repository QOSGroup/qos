## Validator Staking

Validator状态图:

```

      TxCreateValidator->           ->TxRevokeValidator          TimePeriod
none---------------------->Active <------------------->Inactive------------->none
                                    <-TxActiveValidator
```

### Transcation

#### TxCreateValidator

```go
type TxCreateValidator struct {
  Name  string
  Owner Address //操作者
  ValidatorPubkey crypto.Pubkey //validator公钥
  BondTokens  uint64 //绑定Token数量
  Description string
}


type Validator struct {
  Name string
  Owner Address
  ValidatorPubkey crypto.Pubkey
  BondTokens uint64
  Description string

  Status  enum// ACTIVE/INACTIVE
  IsRevoke bool
  InactiveTime time.Time
  InactiveHeight uint64

  BondHeight uint64
}
```
创建validator,validator的`VotingPower`由`BondTokens`决定. `BondTokens`与`QOS`等价.

validator总数量不应超过`maxValidatorCnt`配置


#### TxRevokeValidator

```go
type TxRevokeValidator struct {
  Owner Address //操作者
  ValidatorPubkey crypto.Pubkey //validator地址
}
```

撤销validator后,validator状态变更为pending状态

#### TxActiveValidator

```go
type TxActiveValidator struct {
  Owner Address //操作者
  ValidatorPubkey crypto.Pubkey //validator地址
}

```

激活在`pending`状态的`validator`

由Revoke操作转入pending状态的`validator`不能被激活.

validator总数量不应超过`maxValidatorCnt`配置

### 数据存储

1. Validator Store

path: /store/validator/key

|Index| Prefix Key | Key     | Value | 备注|
|:--|:----       | :-------| :---- | :----|
|a| []byte{0x01} | ValidatorAddress | Validator | 保存Validator信息|
|b| []byte{0x02} | OwnerAddress-ValidatorAddress | 1 | Owner与Validator映射关系 |
|c| []byte{0x03} | ValidatorInactiveTime-ValidatorAddress |InactiveTime| 处于inactive状态的Validator|
|d| []byte{0x04} | VotePower-ValidatorAddress|1| 按VotePower排序的Validator地址,不包含`inactive`状态的Validator|

2. Staking Store

path: /store/staking/key

| Prefix Key | Key     | Value | 备注|
|:----       | :-------| :---- | :----|
| []byte{0x01} | ValidatorAddress | ValidatorSignInfo | 保存Validator在窗口内的签名信息|
| []byte{0x02} | ValidatorAddress-index | true/false | validator在指定窗口偏移量位置是否签名 |


3. Main Store

path: /store/main/key

| Prefix Key | Key     | Value | 备注|
|:----       | :-------| :---- | :----|
| | APPLIED_QOS_AMOUNT|uint64|已分配QOS总量|



```go

type ValidatorSignInfo struct { //签名窗口信息
	StartHeight         uint64 //开始统计高度
	IndexOffset         uint64 //偏移量
	MissedBlocksCounter uint64 //未打块数量
}

```


### 初始化参数

* SPOTotalAmount: 增发总数
* SPOTotalBlock: 总增发块数
* MaxValidatorCnt: validator最大数目,默认10000
* ValidatorVotingStatusLen: 投票窗口高度
* ValidatorVotingStatusLeast: 投票窗口高度内最小投票数
* ValidatorSurvivalSecs: 处于Inactive状态的vaidator生存时间(s)


### BeginBlocker处理步骤

1. 记录Validator窗口投票信息

2. 将不活跃的Validator转为`pending`状态

3. 挖矿奖励

   按增发总数`SPOAmount`在`SPOTotalBlock`块中增发完毕计算每块的奖励数量

   根据`VotePower`计算出每个Validator的`Owner`应奖励的`QOS`数量, 若当前Validator已被撤销或转入`pending`状态,则不发放奖励.

### EndBlocker处理步骤

1. pending状态的validator到期删除

2. 统计新的validator




