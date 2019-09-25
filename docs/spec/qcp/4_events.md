# 事件

QCP模块会发出以下事件:

## 交易

### 初始化联盟链

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| init-qcp             | chain-id         | {chain-id}           |
| init-qcp             | creator          | {creator}            |
| message              | module           | qcp                  |
| message              | action           | init-qcp             |
| message              | gas.payer        | {creator}            |