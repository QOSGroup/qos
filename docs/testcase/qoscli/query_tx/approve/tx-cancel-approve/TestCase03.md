# Description
```
正常取消预授权
```
# Input
```
$ qoscli tx cancel-approve --from test --to test01 --indent
```
# Output
- 第一次调用，此时预授权存在
```
$ qoscli tx cancel-approve --from test --to test01 --indent
Password to sign with 'test':
{
  "check_tx": {
    "gasWanted": "100000",
    "gasUsed": "6877",
    "events": []
  },
  "deliver_tx": {
    "gasWanted": "100000",
    "gasUsed": "7870",
    "events": [
      {
        "type": "cancel-approve",
        "attributes": [
          {
            "key": "YXBwcm92ZS1mcm9t",
            "value": "YWRkcmVzczFodzQzcHdodHNjZWFsdnU5NzNyNjZ2azgzZ3VzOG15cDQwZnk1Ng=="
          },
          {
            "key": "YXBwcm92ZS10bw==",
            "value": "YWRkcmVzczFxbmhhazNwaDB5cXB4YXIzcnJrenVhc2duemx6Zm1xNHB5bjczbQ=="
          }
        ]
      },
      {
        "type": "message",
        "attributes": [
          {
            "key": "bW9kdWxl",
            "value": "YXBwcm92ZQ=="
          },
          {
            "key": "Z2FzLnBheWVy",
            "value": "YWRkcmVzczFodzQzcHdodHNjZWFsdnU5NzNyNjZ2azgzZ3VzOG15cDQwZnk1Ng=="
          }
        ]
      }
    ]
  },
  "hash": "611C0A405A895F5F8CD1939BAFCC0B6FB798E2B8C9B4CADD808367F19F6A9A56",
  "height": "4060"
}
```
- 第二次调用，此时预授权不存在
```
$ qoscli tx cancel-approve --from test --to test01 --indent
Password to sign with 'test':
{
  "check_tx": {
    "code": 1,
    "log": "{\"codespace\":\"sdk\",\"code\":1,\"message\":\"TxStd's ITx ValidateData error:  ERROR:\\nCodespace: approve\\nCode: 104\\nMessage: \\\"approve not exists\\\"\\n\"}",
    "gasWanted": "100000",
    "gasUsed": "1000",
    "events": []
  },
  "deliver_tx": {},
  "hash": "9CF96526C4EB97ABBF6749AD0BBAEDCCB4D488F94D9984E18CC7A68A9CEACB66",
  "height": "0"
}
ERROR: {"codespace":"sdk","code":1,"message":"TxStd's ITx ValidateData error:  ERROR:\nCodespace: approve\nCode: 104\nMessage: \"approve not exists\"\n"}
```