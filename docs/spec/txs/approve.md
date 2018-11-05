# 预授权设计

授权、使用授权、增加授权、减少授权、取消授权，暂未涉及Gas逻辑

## Struct
```
// 授权、增加授权、减少授权、使用授权
type Approve struct {
    From    btypes.Address `json:"from"` // 授权账号，不能为空
    To      btypes.Address `json:"to"`   // 被授权账号，不能为空
    Qos     btypes.BigInt  `json:"qos"`  // qos
    QscList []*QSC         `json:"qsc"`  // qscs，币种不能重复，不能为"qos"（大小写敏感）
}

// 取消授权
type ApproveCancel struct {
	From  btypes.Address `json:"from"`  // 授权账号
	To    btypes.Address `json:"to"`    // 被授权账号
}
```
## Store
storeKey:	approve</br>
key:		from:[addr]/to:[addr]</br>

## Query
```
// TODO 实例
```

## Create
授权账户预授权被授权账户指定币种、币值
```
// TODO 实例
```

## Increase
授权账户增加授权被授权账户指定币种、币值
```
// TODO 实例
```

## Decrease
授权账户减少授权被授权账户指定币种、币值
```
// TODO 实例
```

## Use
被授权用户使用预授权指定币种、币值
```
// TODO 实例
```

## Cancel
取消预授权
```
// TODO 实例
```