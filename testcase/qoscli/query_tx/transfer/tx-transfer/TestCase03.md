# Description
```
指定密钥库中已存在的`--receivers`或`--senders`进行转账
```
# Input
- 指定的`--senders`余额不足
```
$ qoscli tx transfer --senders test,50000000000QOS  --receivers test01,50000000000QOS --indent
```
- 指定的`--senders`余额足够
```
$ qoscli tx transfer --senders test,50000QOS  --receivers test01,50000QOS --indent
```
# Output
- 指定的`--senders`余额不足
```
$ qoscli tx transfer --senders test,50000000000QOS  --receivers test01,50000000000QOS --indent
Password to sign with 'test':<输入密码>
{
  "check_tx": {
    "code": 1,
    "log": "{\"codespace\":\"sdk\",\"code\":1,\"message\":\"TxStd's ITx ValidateData error:  ERROR:\\nCodespace: transfer\\nCode: 203\\nMessage: \\\"sender account has no enough coins\\\"\\n\"}",
    "gasWanted": "100000",
    "gasUsed": "1246",
    "events": []
  },
  "deliver_tx": {},
  "hash": "E0C50F93C7E53F4585779D9EB7E578B6723FFA74A12B60CD004B30F2FCE3AE1C",
  "height": "0"
}
ERROR: {"codespace":"sdk","code":1,"message":"TxStd's ITx ValidateData error:  ERROR:\nCodespace: transfer\nCode: 203\nMessage: \"sender account has no enough coins\"\n"}
```
- 指定的`--senders`余额足够
```
$ qoscli tx transfer --senders test,50000QOS  --receivers test01,50000QOS --indent
Password to sign with 'test':<输入密码>
{
  "check_tx": {
    "gasWanted": "100000",
    "gasUsed": "6952",
    "events": []
  },
  "deliver_tx": {
    "gasWanted": "100000",
    "gasUsed": "16800",
    "events": [
      {
        "type": "message",
        "attributes": [
          {
            "key": "bW9kdWxl",
            "value": "dHJhbnNmZXI="
          },
          {
            "key": "Z2FzLnBheWVy",
            "value": "YWRkcmVzczFodzQzcHdodHNjZWFsdnU5NzNyNjZ2azgzZ3VzOG15cDQwZnk1Ng=="
          }
        ]
      },
      {
        "type": "send",
        "attributes": [
          {
            "key": "YWRkcmVzcw==",
            "value": "YWRkcmVzczFodzQzcHdodHNjZWFsdnU5NzNyNjZ2azgzZ3VzOG15cDQwZnk1Ng=="
          },
          {
            "key": "cW9z",
            "value": "NTAwMDA="
          },
          {
            "key": "cXNjcw==",
            "value": ""
          }
        ]
      },
      {
        "type": "receive",
        "attributes": [
          {
            "key": "YWRkcmVzcw==",
            "value": "YWRkcmVzczFxbmhhazNwaDB5cXB4YXIzcnJrenVhc2duemx6Zm1xNHB5bjczbQ=="
          },
          {
            "key": "cW9z",
            "value": "NTAwMDA="
          },
          {
            "key": "cXNjcw==",
            "value": ""
          }
        ]
      }
    ]
  },
  "hash": "0CEF67AAED1ED02AB0BAD0FA4DBAB6B4806BD525BDEC0EC94C724AE32CCD7930",
  "height": "3454"
}
```