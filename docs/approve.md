# 预授权设计

授权、使用授权、增加授权、减少授权、取消授权，暂未涉及Gas逻辑

## Tx数据结构
```
// 授权、增加授权、减少授权、使用授权
type Approve struct {
	From  btypes.Address `json:"from"`  // 授权账号
	To    btypes.Address `json:"to"`    // 被授权账号
	Coins []types.Qsc    `json:"coins"` // 授权币种、币值
}

// 取消授权
type ApproveCancel struct {
	From  btypes.Address `json:"from"`  // 授权账号
	To    btypes.Address `json:"to"`    // 被授权账号
}
```
## 存储
storeKey:	approve</br>
key:		from:[addr]/to:[addr]</br>

## 1. 创建授权 TxApproveCreate
授权账户预授权被授权账户指定币种、币值

1.授权账户必须存在、被授权账户可不存在</br>
2.创建时无需校验授权账户币种、币值</br>
3.签名、Gas payer：授权账户</br>

## 2. 增加授权 TxApproveIncrease
授权账户增加授权被授权账户指定币种、币值

1.授权、被授权账户必须都存在</br>
2.预授权必须存在</br>
3.无需校验授权账户新增授权币种、币值</br>
4.新增授权币种不在原授权列表时，预授权币种、币值列表添加新币种</br>
5.签名、Gas payer：授权账户</br>

## 3. 减少授权 TxApproveDecrease
授权账户减少授权被授权账户指定币种、币值

1.授权、被授权账户必须都存在</br>
2.预授权必须存在</br>
3.减少授权币种、币值必须小于或等于已授权对应的币种、币值</br>
4.签名、Gas payer：授权账户</br>

## 4. 使用授权 TxApproveUse
被授权用户使用预授权指定币种、币值

1.授权、被授权账户必须都存在</br>
2.预授权必须存在</br>
3.使用币种、币值必须小于或等于已授权对应的币种、币值</br>
4.授权用户拥有的币种、币值必须大于或等于已将要使用的币种、币值</br>
5.签名、Gas payer：被授权账户</br>

## 5. 取消授权 TxApproveCancel
取消预授权

1.授权、被授权账户必须都存在</br>
2.预授权必须存在</br>
3.签名、Gas payer：授权账户</br>