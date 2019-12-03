# Description
```
正常使用预授权
```
注意，使用预授权的`--coins`不可以超过预授权的上限， 也不可以超过`--from`的账户余额。

# Input
```
$ qoscli tx use-approve --coins 300000QOS --from test --to test01 --indent
```
# Output
- 正常情况:
```
$ qoscli tx use-approve --coins 300000QOS --from test --to test01 --indent
Password to sign with 'test03':
{
  "check_tx": {
    "gasWanted": "100000",
    "gasUsed": "7787",
    "events": []
  },
  "deliver_tx": {
    "gasWanted": "100000",
    "gasUsed": "23740",
    "events": [
      {
        "type": "use-approve",
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
            "value": "YWRkcmVzczFxbmhhazNwaDB5cXB4YXIzcnJrenVhc2duemx6Zm1xNHB5bjczbQ=="
          }
        ]
      }
    ]
  },
  "hash": "AF031981DEE8147EB7D9C5640BE1975F13A50A6C3C8B78065EC1F559A27A6415",
  "height": "7139"
}
```
- 当超过预授权余额时:
```
$ qoscli tx use-approve --coins 300000QOS --from test --to test01 --indent
Password to sign with 'test03':
{
  "check_tx": {
    "code": 1,
    "log": "{\"codespace\":\"sdk\",\"code\":1,\"message\":\"TxStd's ITx ValidateData error:  ERROR:\\nCodespace: approve\\nCode: 106\\nMessage: \\\"approve not enough\\\"\\n\"}",
    "gasWanted": "100000",
    "gasUsed": "1156",
    "events": []
  },
  "deliver_tx": {},
  "hash": "FDF31F79E2FC13A73B8607FB71EF4213001E9AC0B59F99416238D5DC19DE9D01",
  "height": "0"
}
ERROR: {"codespace":"sdk","code":1,"message":"TxStd's ITx ValidateData error:  ERROR:\nCodespace: approve\nCode: 106\nMessage: \"approve not enough\"\n"}
```
- 当使用的预授权超过`--from`账户余额时: 
```
$ qoscli tx use-approve --from test01 --to test02 --coins 500000QOS --indent
Password to sign with 'test02':
{
  "check_tx": {
    "code": 1,
    "log": "{\"codespace\":\"sdk\",\"code\":1,\"message\":\"TxStd's ITx ValidateData error:  ERROR:\\nCodespace: approve\\nCode: 107\\nMessage: \\\"from account has no enough coins\\\"\\n\"}",
    "gasWanted": "100000",
    "gasUsed": "2390",
    "events": []
  },
  "deliver_tx": {},
  "hash": "16177E8DCE056A31707183E766F3A205627652D5DC6F3E38E6BE6A660F930C62",
  "height": "0"
}
ERROR: {"codespace":"sdk","code":1,"message":"TxStd's ITx ValidateData error:  ERROR:\nCodespace: approve\nCode: 107\nMessage: \"from account has no enough coins\"\n"}
```