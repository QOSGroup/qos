# 事件

系统账户模块会发出如下事件：

## 交易

### 添加系统账户

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| add-guardian         | creator          | {creator}            |
| add-guardian         | guardian         | {address}            |
| message              | module           | guardian             |
| message              | action           | add-guardian         |
| message              | gas.payer        | {creator}            |

### 删除系统账户

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| delete-guardian      | delete-by        | {deleted_by}         |
| delete-guardian      | guardian         | {address}            |
| message              | module           | guardian             |
| message              | action           | delete-guardian      |
| message              | gas.payer        | {deleted_by}         |

### 停止网络

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| halt-network         | guardian         | {guardian}           |
| halt-network         | reason           | {reason}             |
| message              | module           | guardian             |
| message              | action           | halt-network         |
| message              | gas.payer        | {guardian}           |