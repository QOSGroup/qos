# Events

This module emits the following events:

## Transactions

### TxAddGuardian

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| add-guardian         | creator          | {creator}            |
| add-guardian         | guardian         | {address}            |
| message              | module           | guardian             |
| message              | action           | add-guardian         |
| message              | gas.payer        | {creator}            |

### TxDeleteGuardian

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| delete-guardian      | delete-by        | {deleted_by}         |
| delete-guardian      | guardian         | {address}            |
| message              | module           | guardian             |
| message              | action           | delete-guardian      |
| message              | gas.payer        | {deleted_by}         |

### TxHaltNetwork

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| halt-network         | guardian         | {guardian}           |
| halt-network         | reason           | {reason}             |
| message              | module           | guardian             |
| message              | action           | halt-network         |
| message              | gas.payer        | {guardian}           |