# Events

This module emits the following events:

## Transactions

### TxTransfer

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| send                 | address          | {addr}               |
| send                 | qos              | {qos}                |
| send                 | qscs             | {qscs}               |
| receive              | address          | {addr}               |
| receive              | qos              | {qos}                |
| receive              | qscs             | {qscs}               |
| message              | action           | transfer             |
| message              | gas.payer        | {senders[0]}         |

### TxInvariantCheck

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| invariant_check      | sender           | {sender}             |
| invariant_check      | height           | {ctx.BlockHeight()}  |
| message              | action           | invariant_check      |
| message              | gas.payer        | {senders[0]}         |

## EndBlocker

### Check Invariant

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| invariant_check      |                  |                      |

App EndBlocker will check `invariant_check` event for invariant checking.

### Release Lock

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| release              | address          | {receiver}           |
| release              | qos              | {releaseAmount}      |