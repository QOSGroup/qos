# Events

This module emits the following events:

## Transactions

### Proposal

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| submit-proposal      | proposal-id      | {proposal_id}        |
| submit-proposal      | proposer         | {proposer}           |
| submit-proposal      | depositor        | {depositor}          |
| submit-proposal      | proposal-type    | {proposal_type}      |
| message              | module           | gov                  |
| message              | action           | submit-proposal      |
| message              | gas.payer        | {proposer}           |

### Deposit

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| deposit-proposal     | proposal-id      | {proposal_id}        |
| deposit-proposal     | depositor        | {depositor}          |
| message              | action           | deposit-proposal     |
| message              | module           | gov                  |
| message              | gas.payer        | {depositor}          |

### Vote

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| vote-proposal        | proposal-id      | {proposal_id}        |
| vote-proposal        | voter            | {voter}              |
| message              | action           | vote-proposal        |
| message              | module           | gov                  |
| message              | gas.payer        | {voter}              |

## EndBlocker

### Proposal dropped

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| inactive-proposal    | proposal-id      | {proposal_id}        |
| inactive-proposal    | proposal-result  | proposal-dropped     |


### Proposal completed

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| active-proposal      | proposal-id      | {proposal_id}        |
| active-proposal      | proposal-result  | {proposal-result}    |