# 事件

Bank 模块会发出如下事件：

## 交易

### 提议

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| submit-proposal      | proposal-id      | {proposal_id}        |
| submit-proposal      | proposer         | {proposer}           |
| submit-proposal      | depositor        | {depositor}          |
| submit-proposal      | proposal-type    | {proposal_type}      |
| message              | module           | gov                  |
| message              | action           | submit-proposal      |
| message              | gas.payer        | {proposer}           |

### 质押

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| deposit-proposal     | proposal-id      | {proposal_id}        |
| deposit-proposal     | depositor        | {depositor}          |
| message              | action           | deposit-proposal     |
| message              | module           | gov                  |
| message              | gas.payer        | {depositor}          |

### 投票

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| vote-proposal        | proposal-id      | {proposal_id}        |
| vote-proposal        | voter            | {voter}              |
| message              | action           | vote-proposal        |
| message              | module           | gov                  |
| message              | gas.payer        | {voter}              |

## EndBlocker

### 提议失效（未到投票期）

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| inactive-proposal    | proposal-id      | {proposal_id}        |
| inactive-proposal    | proposal-result  | proposal-dropped     |


### 提议完成

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| active-proposal      | proposal-id      | {proposal_id}        |
| active-proposal      | proposal-result  | {proposal-result}    |