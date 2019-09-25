# Description
```
正常减少预授权
```
注意，预授权的`--coins`可以超过`--from`的账户余额。

# Input
```
$ qoscli tx decrease-approve --coins 200000QOS --from test --to test01 --indent
```
# Output
```
$ qoscli tx decrease-approve --coins 200000QOS --from test --to test01 --indent
Password to sign with 'test':
{
  "check_tx": {
    "gasWanted": "100000",
    "gasUsed": "6865",
    "events": []
  },
  "deliver_tx": {
    "gasWanted": "100000",
    "gasUsed": "11580",
    "events": [
      {
        "type": "decrease-approve",
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
  "hash": "9E76472FE899566E3A6EDB306F14B4E87EB5E5F1D811274882B81E2D7D93EF3B",
  "height": "7030"
}
```