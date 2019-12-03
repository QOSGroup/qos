# Events

This module emits the following events:

## Transactions

### TxCreateQSC

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| create-qsc           | name             | {name}               |
| create-qsc           | creator          | {creator}            |
| message              | module           | qsc                  |
| message              | action           | create-qsc           |
| message              | gas.payer        | {creator}            |

### TxIssueQSC

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| issue-qsc            | name             | {name}               |
| issue-qsc            | banker           | {banker}             |
| issue-qsc            | tokens           | {tokens}             |
| message              | module           | qsc                  |
| message              | action           | issue-qsc            |
| message              | gas.payer        | {banker}             |