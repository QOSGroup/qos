# 转账设计
实现多账户，多币种交易，只需保证发送和接收集合QOS、QSCs总量相等

## Struct
```
type TransItem struct {
	Address btypes.Address `json:"addr"` // 账户地址
	QOS     btypes.BigInt  `json:"qos"`  // QOS
	QSCs    types.QSCs     `json:"qscs"` // QSCs
}

type TransferTx struct {
	Senders   []TransItem `json:"senders"`   // 发送集合
	Receivers []TransItem `json:"receivers"` // 接收集合
}
```

## TX
Arya向Sansa转账1个qos，1个qstar

### Send
```
qoscli transfer --senders=Arya,1qos,1qstar --receivers=Sansa,1qos,1qstar
Password to sign with 'Arya':
{"check_tx":{},"deliver_tx":{},"hash":"E3D3902CD0C91BB1E982243EF23DDDFF646DED88","height":"231"}
```

### Query

* tx

```
qoscli tx E3D3902CD0C91BB1E982243EF23DDDFF646DED88
{
  "hash": "49OQLNDJG7HpgiQ+8j3d/2Rt7Yg=",
  "height": "231",
  "tx": {
    "type": "qbase/txs/stdtx",
    "value": {
      "itx": {
        "type": "qos/txs/TransferTx",
        "value": {
          "senders": [
            {
              "addr": "address1cnfqru6rts4nz224mvrf58ne427uthmcut4kc3",
              "qos": "1",
              "qscs": [
                {
                  "coin_name": "qstar",
                  "amount": "1"
                }
              ]
            }
          ],
          "receivers": [
            {
              "addr": "address1spdn868fzcpah8zd74tjck0e5akacgt2gmccnq",
              "qos": "1",
              "qscs": [
                {
                  "coin_name": "qstar",
                  "amount": "1"
                }
              ]
            }
          ]
        }
      },
      "sigature": [
        {
          "pubkey": {
            "type": "tendermint/PubKeyEd25519",
            "value": "tJXDzIjW1NZp3XiCxWDFBqyiMb2UpyoU8vp240DqsjY="
          },
          "signature": "gB9TPt0hgscbDnD0CIWMBQwsTgSkaTIZps/cVIqs61ivasWm8+sFymNxMSIvqZ8QDap/9ihRxwyu17u9gdjHAA==",
          "nonce": "1"
        }
      ],
      "chainid": "qos-test",
      "maxgas": "0"
    }
  },
  "result": {}
} <nil>
```

* account

Arya：
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
```
Sansa：
```
qoscli account --name=Sansa
{
  "type": "qbase/account/QOSAccount",
  "value": {
    "base_account": {
      "account_address": "address1spdn868fzcpah8zd74tjck0e5akacgt2gmccnq",
      "public_key": null,
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