# Description
```
正常增加预授权
```
注意，预授权的`--coins`可以超过`--from`的账户余额。

# Input
```
$ qoscli tx increase-approve --coins 500000QOS --from test --to test01 --indent
```
# Output
```
$ qoscli tx increase-approve --coins 500000QOS --from test --to test01 --indent
Password to sign with 'test':
{
  "check_tx": {
    "gasWanted": "100000",
    "gasUsed": "6862",
    "events": []
  },
  "deliver_tx": {
    "gasWanted": "100000",
    "gasUsed": "11600",
    "events": [
      {
        "type": "increase-approve",
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
  "hash": "652626F7C40D46FCCD00A57ACB33246714484BC369253D994F4FE57840C16ED1",
  "height": "6983"
}
```