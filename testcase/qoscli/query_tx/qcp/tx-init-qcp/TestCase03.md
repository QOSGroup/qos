# Description
```
正常初始化QCP联盟链
```
# Input
```
$ qoscli tx init-qcp --creator starBanker --qcp.crt ./qcp-star.crt --indent
```
# Output
```
$ qoscli tx init-qcp --creator starBanker --qcp.crt ./qcp-star.crt --indent
Password to sign with 'starBanker':
{
  "check_tx": {
    "gasWanted": "100000",
    "gasUsed": "8709"
  },
  "deliver_tx": {
    "gasWanted": "100000",
    "gasUsed": "15870",
    "tags": [
      {
        "key": "YWN0aW9u",
        "value": "aW5pdC1xY3A="
      },
      {
        "key": "cWNw",
        "value": "cWNwLXN0YXI="
      },
      {
        "key": "Y3JlYXRvcg==",
        "value": "YWRkcmVzczFhcDZteXYyNDhlbGwwZ2t3bmh6YXBqbnl2dDM4NTg0ZGszZTB0ZQ=="
      }
    ]
  },
  "hash": "6C616140684675A7CB29ED8975B2161D6102955AC8D743B135CFAF16838AAFB5",
  "height": "606260"
}
```