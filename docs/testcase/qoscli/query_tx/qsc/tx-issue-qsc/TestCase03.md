# Description
```
正常发放QSC联盟币
```
# Input
```
$ qoscli tx issue-qsc --banker starBanker --qsc-name STAR --amount 100000 --indent
```
# Output
```
$ qoscli tx issue-qsc --banker starBanker --qsc-name STAR --amount 100000 --indent
Password to sign with 'starBanker':
{
  "check_tx": {
    "gasWanted": "100000",
    "gasUsed": "6715"
  },
  "deliver_tx": {
    "gasWanted": "100000",
    "gasUsed": "12760",
    "tags": [
      {
        "key": "YWN0aW9u",
        "value": "aXNzdWUtcXNj"
      },
      {
        "key": "cXNj",
        "value": "U1RBUg=="
      },
      {
        "key": "YmFua2Vy",
        "value": "YWRkcmVzczFhcDZteXYyNDhlbGwwZ2t3bmh6YXBqbnl2dDM4NTg0ZGszZTB0ZQ=="
      }
    ]
  },
  "hash": "1ED60FAB3184EAE7E635616820955AAE39E98208820597F9CFE099DC52246AD2",
  "height": "606613"
}
```
完成后查询`starBanker`账户余额如下:
```
$ qoscli query account starBanker --indent
{
  "type": "qos/types/QOSAccount",
  "value": {
    "base_account": {
      "account_address": "address1ap6myv248ell0gkwnhzapjnyvt38584dk3e0te",
      "public_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "Im+aYoPhzdfmFTw3Sgh3HLRK2s1gpCuhFg+cP3K1OFQ="
      },
      "nonce": "2"
    },
    "qos": "9997137",
    "qscs": [
      {
        "coin_name": "STAR",
        "amount": "100000"
      }
    ]
  }
}
```