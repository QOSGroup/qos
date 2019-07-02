# Guardian

QOS加入了由基金会控制的特权系统用户

- 作为`TxTaxUsage`提议目的地址，接收社区费池提取出来的QOS。
- 发送交易无需支付Gas费。

## Struct

### Guardian
```go
const (
	Genesis  GuardianType = 0x01    // Guardian in genesis.json
	Ordinary GuardianType = 0x02    // Guardian created by genesis guardian
)

type Guardian struct {
	Description  string         `json:"description"`
	GuardianType GuardianType   `json:"guardian_type"`  
	Address      btypes.Address `json:"address"`
	Creator      btypes.Address `json:"creator"` 
}
```

字段说明：
- Description 描述信息
- GuardianType 类型，区别为是否在`genesis.json`中
- Address 账户地址
- Creator 创建Guardian地址

## Txs

### TxAddGuardian

添加特权账户
```go
type TxAddGuardian struct {
	Description string         `json:"description"`
	Address     btypes.Address `json:"address"`
	Creator     btypes.Address `json:"creator"`
}
```

字段说明：
- Description 描述信息
- Address 账户地址
- Creator 创建Guardian地址

操作指令：[添加特权账户](../command/qoscli.md#添加特权账户)

### TxDeleteGuardian

删除特权账户
```go
type TxDeleteGuardian struct {
	Address   btypes.Address `json:"address"`   
	DeletedBy btypes.Address `json:"deleted_by"`
}
```

字段说明：
- Address 账户地址
- DeletedBy 执行删除操作账户地址，只能是Genesis Guardian

操作指令：[删除特权账户](../command/qoscli.md#删除特权账户)