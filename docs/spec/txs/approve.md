# 预授权设计

授权、使用授权、增加授权、减少授权、取消授权，暂未涉及Gas逻辑

## Struct
```go
// 授权、增加授权、减少授权、使用授权
type Approve struct {
    From    btypes.Address `json:"from"` // 授权账号，不能为空
    To      btypes.Address `json:"to"`   // 被授权账号，不能为空
    Qos     btypes.BigInt  `json:"qos"`  // qos
    QscList []*QSC         `json:"qsc"`  // qscs，币种不能重复，不能为"qos"（大小写敏感）
}

// 取消授权 Tx
type TxCancelApprove struct {
	From btypes.Address `json:"from"` // 授权账号
	To   btypes.Address `json:"to"`   // 被授权账号
}
```
## Store
```go
approveStoreKey = "approve"             // store
approveKey      = "from:[%s]/to:[%s]"   // key
```

读写使用ApproveMapper
```go
type ApproveMapper struct {
	*mapper.BaseMapper      // qbase BaseMapper封装 
}
```
提供获取授权（GetApprove）、保存授权（SaveApprove）、删除授权（DeleteApprove）方法

## Create

From账户向To账户预授权一定量的QOS和QSCs，预授权创建成功并非转账成功，仅仅是记录，不改变账户状态。所以From/To账户在链上均可不存在，From拥有资产总量可以小于授权资产总量。假设From仅有1QOS，可授权To2QOS。

* valid
1. QOS、QSCs中币种不能重复、币值必须为正
2. 创建前链上不存在From对To的预授权，若存在请执行approve的其他操作。

* signer
  
From账户

## Increase

From账户向To账户增加授权一定量的QOS和QSCs，在已存在预授权基础上增加预授权资产。假设From已对To授权1QOS，增加授权1QOS，完成后From对To的预授权为2QOS。

* valid
1. QOS、QSCs中币种不能重复、币值必须为正
2. 链上存在From对To的预授权，若不存在请执行create操作。

* signer

From账户

## Decrease

From账户向To账户减少授权一定量的QOS和QSCs，在已存在预授权基础上减少预授权资产。假设From已对To授权2QOS，减少授权1QOS，完成后From对To的预授权为1QOS。

* valid
1. QOS、QSCs中币种不能重复、币值必须为正
2. 链上存在From对To的预授权，若不存在请执行create操作。
3. QOS、QSCs总量不能大于已授权币值总量

* signer

From账户

## Use

To账户使用From账户预授权的QOS和QSCs。假设From已授权To 2QOS，执行use 1QOS后，From向To授权变成 1QOS，From账户向To账户转账1QOS。

* valid
1. QOS，QSCs中币种不能重复、币值必须为正
2. 链上不存在From对To的预授权，若存在请执行approve的其他操作。
3. QOS、QSCs总量不能大于已授权币值总量
4. From账户必须存在
3. QOS、QSCs总量不能大于From账户币值总量

* signer
  
To账户

## Cancel

From账户取消对To账户的预授权信息。假设From已授权To 2QOS，执行cancel后，将删除From对To的预授权信息，已使用的授权币种、币值不变。

* valid
1. 链上不存在From对To的预授权。

* signer
  
From账户