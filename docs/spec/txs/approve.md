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
查询授权
```
qoscli approve --from=Arya --to=Sansa
{
  "from": "address1ah9uz0",
  "to": "address1ah9uz0",
  "qos": "0",
  "qscs": null
} <nil>
```
查询账户
```
qoscli account --name=Arya
{
  "type": "qbase/account/QOSAccount",
  "value": {
    "base_account": {
      "account_address": "address1cnfqru6rts4nz224mvrf58ne427uthmcut4kc3",
      "public_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "tJXDzIjW1NZp3XiCxWDFBqyiMb2UpyoU8vp240DqsjY="
      },
      "nonce": "1"
    },
    "qos": "99999999",
    "qscs": [
      {
        "coin_name": "qstar",
        "amount": "99999999"
      }
    ]
  }
} <nil>
qoscli account --name=Sansa
{
  "type": "qbase/account/QOSAccount",
  "value": {
    "base_account": {
      "account_address": "address1spdn868fzcpah8zd74tjck0e5akacgt2gmccnq",
      "public_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "Dc4YBO1JMlVteO5ka21FNbhRBqCeKacNYs62YBHhFcw="
      },
      "nonce": "0"
    },
    "qos": "1",
    "qscs": [
      {
        "coin_name": "qstar",
        "amount": "1"
      }
    ]
  }
} <nil>
```

## Create
Arya向Sansa授权100个qos，100个qstar
```
qoscli approve-create --from=Arya --to=Sansa --qos=100 --qscs=100qstar
Password to sign with 'Arya':
{"check_tx":{},"deliver_tx":{},"hash":"9E33CEC7C2589D3A62E60D3F1D3B0F8FE330B020","height":"370"}
```
查询授权
```
qoscli approve --from=Arya --to=Sansa
{
  "from": "address1cnfqru6rts4nz224mvrf58ne427uthmcut4kc3",
  "to": "address1spdn868fzcpah8zd74tjck0e5akacgt2gmccnq",
  "qos": "100",
  "qscs": [
    {
      "coin_name": "qstar",
      "amount": "100"
    }
  ]
} <nil>
```

## Increase
Arya向Sansa增加授权100个qos，100个qstar
```
qoscli approve-increase --from=Arya --to=Sansa --qos=100 --qscs=100qstar
Password to sign with 'Arya':
{"check_tx":{},"deliver_tx":{},"hash":"3C06676C53A5439D39CB4D0FBA3213C44DC1BA8E","height":"406"}
```
查询授权
```
qoscli approve --from=Arya --to=Sansa
{
  "from": "address1cnfqru6rts4nz224mvrf58ne427uthmcut4kc3",
  "to": "address1spdn868fzcpah8zd74tjck0e5akacgt2gmccnq",
  "qos": "200",
  "qscs": [
    {
      "coin_name": "qstar",
      "amount": "200"
    }
  ]
} <nil>
```

## Decrease
Arya向Sansa减少授权100个qos，100个qstar
```
qoscli approve-decrease --from=Arya --to=Sansa --qos=100 --qscs=100qstar
Password to sign with 'Arya':
{"check_tx":{},"deliver_tx":{},"hash":"9DC18AD3CB0B59FCD354C267D8C22A1CC75E5624","height":"414"}
```
查询授权
```
qoscli approve --from=Arya --to=Sansa
{
  "from": "address1cnfqru6rts4nz224mvrf58ne427uthmcut4kc3",
  "to": "address1spdn868fzcpah8zd74tjck0e5akacgt2gmccnq",
  "qos": "100",
  "qscs": [
    {
      "coin_name": "qstar",
      "amount": "100"
    }
  ]
} <nil>
```

## Use
Sansa使用Arya向自己授权的10个qos，10个qstar
```
qoscli approve-use --from=Arya --to=Sansa --qos=10 --qscs=10qstar
Password to sign with 'Sansa':
{"check_tx":{},"deliver_tx":{},"hash":"0573760D6B316E6695FBB63A56F2A20C0635FCAE","height":"437"}
```
查询授权
```
qoscli approve --from=Arya --to=Sansa
{
  "from": "address1cnfqru6rts4nz224mvrf58ne427uthmcut4kc3",
  "to": "address1spdn868fzcpah8zd74tjck0e5akacgt2gmccnq",
  "qos": "90",
  "qscs": [
    {
      "coin_name": "qstar",
      "amount": "90"
    }
  ]
} <nil>
```
查询账户
```
qoscli account --name=Arya
{
  "type": "qbase/account/QOSAccount",
  "value": {
    "base_account": {
      "account_address": "address1cnfqru6rts4nz224mvrf58ne427uthmcut4kc3",
      "public_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "tJXDzIjW1NZp3XiCxWDFBqyiMb2UpyoU8vp240DqsjY="
      },
      "nonce": "4"
    },
    "qos": "99999989",
    "qscs": [
      {
        "coin_name": "qstar",
        "amount": "99999989"
      }
    ]
  }
} <nil>
qoscli account --name=Sansa
{
  "type": "qbase/account/QOSAccount",
  "value": {
    "base_account": {
      "account_address": "address1spdn868fzcpah8zd74tjck0e5akacgt2gmccnq",
      "public_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "Dc4YBO1JMlVteO5ka21FNbhRBqCeKacNYs62YBHhFcw="
      },
      "nonce": "1"
    },
    "qos": "11",
    "qscs": [
      {
        "coin_name": "qstar",
        "amount": "11"
      }
    ]
  }
} <nil>
```

## Cancel
Arya取消向Sansa授权任何资产
```
qoscli approve-cancel --from=Arya --to=Sansa
Password to sign with 'Arya':
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"484"}
```
查询授权
```
qoscli approve --from=Arya --to=Sansa
{
  "from": "address1ah9uz0",
  "to": "address1ah9uz0",
  "qos": "0",
  "qscs": null
} <nil>
```