# Events

This module emits the following events:

## Transactions

### TxInitQCP

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| init-qcp             | chain-id         | {chain-id}           |
| init-qcp             | creator          | {creator}            |
| message              | module           | qcp                  |
| message              | action           | init-qcp             |
| message              | gas.payer        | {creator}            |