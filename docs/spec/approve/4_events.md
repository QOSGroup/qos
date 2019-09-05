# 事件

预授权模块会发出以下事件:

## 交易

### 创建预授权

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| create-approve       | approve-from     | {from}               |
| create-approve       | approve-to       | {to}                 |
| message              | module           | approve              |
| message              | action           | create-approve       |
| message              | gas.payer        | {from}               |

### 增加预授权

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| increase-approve     | approve-from     | {from}               |
| increase-approve     | approve-to       | {to}                 |
| message              | module           | approve              |
| message              | action           | increase-approve     |
| message              | gas.payer        | {from}               |

### 减少预授权

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| decrease-approve     | approve-from     | {from}               |
| decrease-approve     | approve-to       | {to}                 |
| message              | module           | approve              |
| message              | action           | decrease-approve     |
| message              | gas.payer        | {from}               |

### 使用预授权

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| use-approve          | approve-from     | {from}               |
| use-approve          | approve-to       | {to}                 |
| message              | module           | approve              |
| message              | action           | use-approve          |
| message              | gas.payer        | {to}                 |

### 取消预授权

| Type                 | Attribute Key    | Attribute Value      |
|----------------------|------------------|----------------------|
| cancel-approve       | approve-from     | {from}               |
| cancel-approve       | approve-to       | {to}                 |
| message              | module           | approve              |
| message              | action           | cancel-approve       |
| message              | gas.payer        | {from}               |
