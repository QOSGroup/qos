# Events

This module emits the following events:

## Transactions

### Validator

#### TxCreateValidator

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| create-validator     | validator        | {validator}          |
| create-validator     | owner            | {owner}              |
| create-validator     | delegator        | {owner}              |
| message              | module           | stake                |
| message              | action           | create-validator     |
| message              | gas.payer        | {owner}              |

#### TxModifyValidator

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| modify-validator     | owner            | {owner}              |
| modify-validator     | validator        | {validator}          |
| message              | module           | stake                |
| message              | action           | modify-validator     |
| message              | gas.payer        | {owner}              |

#### TxRevokeValidator

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| revoke-validator     | validator        | {validator}          |
| revoke-validator     | owner            | {owner}              |
| message              | module           | stake                |
| message              | action           | revoke-validator     |
| message              | gas.payer        | {owner}              |

#### TxActiveValidator

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| active-validator     | validator        | {validator}          |
| active-validator     | owner            | {owner}              |
| message              | module           | stake                |
| message              | action           | active-validator     |
| message              | gas.payer        | {owner}              |

### Delegation

#### TxCreateDelegate

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| create-delegation    | validator        | {validator}          |
| create-delegation    | delegator        | {delegator}          |
| message              | module           | stake                |
| message              | action           | create-delegation    |
| message              | gas.payer        | {delegator}          |

#### TxModifyCompound

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| modify-compound      | approve-from     | {validator}          |
| modify-compound      | approve-to       | {delegator}          |
| message              | module           | stake                |
| message              | action           | modify-compound      |
| message              | gas.payer        | {delegator}          |

#### TxUnbondDelegation

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| unbond-delegation    | validator        | {validator}          |
| unbond-delegatione   | delegator        | {delegator}          |
| message              | module           | stake                |
| message              | action           | unbond-delegation    |
| message              | gas.payer        | {delegator}          |

#### TxCreateReDelegation

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| create-redelegation  | validator        | {from_validator}     |
| create-redelegation  | new-validator    | {to_validator}       |
| create-redelegation  | delegator        | {delegator}          |
| message              | module           | stake                |
| message              | action           | create-redelegation  |
| message              | gas.payer        | {delegator}          |

## BeginBlocker

| Type                 | Attribute Key    | Attribute Value        |
|----------------------|------------------|------------------------|
| missing-vote         | validator        | {validator}            |
| missing-vote         | missed-blocks    | {MissedBlocksCounter}  |
| missing-vote         | height           | {height}               |
| inactive-validator   | validator        | {validator}            |
| inactive-validator   | height           | {height}               |
| slash                | validator        | {validator}            |
| slash                | owner            | {owner}                |
| slash                | reason           | {double_sign/down_time}|