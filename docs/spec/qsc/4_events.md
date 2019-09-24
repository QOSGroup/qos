# 事件

QSC模块会发出以下事件:

## 交易

### 创建QSC

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| create-qsc           | name             | {name}               |
| create-qsc           | creator          | {creator}            |
| message              | module           | qsc                  |
| message              | action           | create-qsc           |
| message              | gas.payer        | {creator}            |

### 发行QSC

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| issue-qsc            | name             | {name}               |
| issue-qsc            | banker           | {banker}             |
| issue-qsc            | tokens           | {tokens}             |
| message              | module           | qsc                  |
| message              | action           | issue-qsc            |
| message              | gas.payer        | {banker}             |