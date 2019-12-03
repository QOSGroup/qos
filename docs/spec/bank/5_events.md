# 事件

Bank 模块会发出如下事件：

## 交易

### 转账

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

### 数据检查

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| invariant_check      | sender           | {sender}             |
| invariant_check      | height           | {ctx.BlockHeight()}  |
| message              | action           | invariant_check      |
| message              | gas.payer        | {senders[0]}         |

## EndBlocker

### 数据检查

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| invariant_check      |                  |                      |

App EndBlocker 会判断`invariant_check`事件类型从而执行数据检查操作。

### 锁定-释放

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| release              | address          | {receiver}           |
| release              | qos              | {releaseAmount}      |