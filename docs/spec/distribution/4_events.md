# 事件

分配模块会发出以下事件：

## BeginBlocker

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| proposer_reward      | tokens           | {proposerRewards}    |
| proposer_reward      | validator        | {proposerAddr}       |
| community            | validator        | {communityFeePool}   |

## EndBlocker

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| delegator_rewards    | tokens           | {sharedReward}       |
| delegator_rewards    | validator        | {validator}          |
| commission           | tokens           | {commissionReward}   |
| commission           | validator        | {validator}          |
| delegator_reward     | tokens           | {delegatorRewards}   |
| delegator_reward     | validator        | {validator}          |
| delegator_reward     | delegator        | {delegator}          |
| delegate             | tokens           | {tokens}             |
| delegate             | validator        | {validator}          |
| delegate             | delegator        | {delegator}          |