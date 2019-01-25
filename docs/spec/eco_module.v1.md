## 经济模型

### 原则

* 不管是验证人还是委托人，币有相同的收益权
* 不管是验证人还是委托人，币有相同的治理权
* 验证人能接收多少委托人的币：系统上不控制，把验证人当前数据公布，由委托人自己决定；
* 准验证人排名依据：自己抵押币+委托人币；
* 新晋验证人竞争规则，倒数第一名验证人自己抵押10币+委托人10币，等待队伍排名第一的准验证人自己抵押11币，是否允许准验证人把验证人挤掉：不允许；
* 惩罚验证人时，除了验证人自己的币，委托人的币是否受惩罚：惩罚，但惩罚力度可以比验证人小；
* 验证人挖矿收益的到帐方式：系统定期打款到帐，并可选择自动复投；

### 模块

#### 通胀

#### 惩罚

#### 委托

* TxCreateDelegate

创建delegate, 创建之后,根据`分配周期`参数,该委托人获取收益的周期即已确定. unbond操作不影响delegator分配周期.

可以追加delegate,追加后,delegator收益发放周期不变

``` go

type TxDelegate struct {
	delegator
    validator
    amount
    isCompound
}

```

* TxUnbondDelegate
取回部分代理

```

     |                x                             y          |
     |  --------------------------------|----------------------|
     |                                 unbond -----------------|-------------------------|

上次收益发放                                             下次收益发放                unbond解绑

```

unbond后,对应的validator将会增加一个计费点,unbond金额将在`unbond周期`之后返还至delegator账户.
unbond操作立即生效, 先统计出当前收益,并追加到下次收益发放总额中.

下次收益发放时, 发放金额为 x + y

金额与token的换算??


```go
type TxUnDelegate struct {
	delegator
    validator
    amount
}
```

* TxModifyCompound

修改绑定delegate是否为复利, 在一个周期内,仅最后一次操作生效

```go
type TxModifyCompound struct {
	delegator
    validator
    isCompound
}
```








 * Store

path: /store/delegator/key

|Index| Prefix Key | Key     | Value | 备注|
|:--|:----       | :-------| :---- | :----|
|a| []byte{0x31} | DelegatorAddress-ValidatorAddress |DelegationInfo| delegator信息|
|b| []byte{0x32} | ValidatorAddress-DelegatorAddress | struct{}{}| validator与delegator映射|
|c| []byte{0x41} | BlockHeight+Delegator|QOS|delegator在指定高度解绑的QOS数量,在此高度将QOS返回至delegator|


#### Fee分配

1.  每块分配的QOS数量由通胀策略决定
       * mint 数量
       * tx 手续费
2.  每块挖出的QOS数量:  `x%`proposer + `y%`validators + `z%`community
       * `x%`proposer: 验证人获得的奖励
       * `y%`validators: 根据每个validator的power占比平均分配
3.  validator奖励数 =  validator佣金 +  平分金额Fee
       * validator佣金奖励: 佣金 = validator奖励数 * `commission rate`
       * 平分金额Fee由validator,delegator根据各自绑定的stake平均分配
4.  validator的proposer奖励,佣金奖励 均按周期发放


* store

path: /store/distribution/key

|Index| Prefix Key | Key     | Value | 备注|
|:--|:----       | :-------| :---- | :----|
|a| []byte{0x01} | communityFeePoolKey | QOS | 存储社区获得的收益|
|b| []byte{0x02} | lastBlockProposerKey | Proposer Address | 存储上一块中proposer地址 |
|c| []byte{0x04} | ValidatorAddress-DelegatorAdddress|DelegatorStaringInfo| delegator开始计算收益信息|
|d| []byte{0x05} | ValidatorAddress-Period | QOS | validator 历史汇总收益|
|e| []byte{0x06} | ValidatorAddress |ValidatorCurrentPeriodRewards | validator当前周期汇总收益信息 |
|g| []byte{0x10} | blockDistributionKey |QOS|待分配QOS总数 = mint + tx fee|
|m| []byte{0x80} | BlockHeight + ValidatorAddress + DelegatorAddress| struct{}{} | 某高度下需要分配收益的delegator|


```go

type DelegatorStartingInfo struct {
	PreviousPeriod uint64  `json:"previous_period"`
	bondToken      bigInt  `json:"bondToken"`
	StartingHeight         uint64  `json:"height"`
  HistoricalReward  bigInt `json:"historical_rewards"`
}

type ValidatorCurrentRewards struct {
	Rewards bigInt       `json:"rewards"`
	Period  uint64       `json:"period"`
}

```


























