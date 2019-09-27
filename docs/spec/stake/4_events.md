# 事件

预授权模块会发出以下事件:

## 交易

### 创建委托

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| create-delegation    | validator        | {validator}          |
| create-delegation    | delegator        | {delegator}          |
| message              | module           | stake                |
| message              | action           | create-delegation    |
| message              | gas.payer        | {delegator}          |

### 修改收益复投方式

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| modify-compound      | approve-from     | {validator}          |
| modify-compound      | approve-to       | {delegator}          |
| message              | module           | stake                |
| message              | action           | modify-compound      |
| message              | gas.payer        | {delegator}          |

### 使用预授权

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| unbond-delegation    | validator        | {validator}          |
| unbond-delegatione   | delegator        | {delegator}          |
| message              | module           | approve              |
| message              | action           | unbond-delegation    |
| message              | gas.payer        | {delegator}          |

### 变更委托验证节点

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| create-redelegation  | validator        | {from_validator}     |
| create-redelegation  | new-validator    | {to_validator}       |
| create-redelegation  | delegator        | {delegator}          |
| message              | module           | stake                |
| message              | action           | create-redelegation  |
| message              | gas.payer        | {delegator}          |

