# Description
```
正常创建QSC联盟币
```
# Input
```
$ qoscli tx create-qsc --creator starCreator --qsc.crt ./qsc-star.crt --indent
```
# Output
```
$ qoscli tx create-qsc --creator starCreator --qsc.crt ./qsc-star.crt --indent
Password to sign with 'starCreator':
{
  "check_tx": {
    "gasWanted": "100000",
    "gasUsed": "8709"
  },
  "deliver_tx": {
    "gasWanted": "100000",
    "gasUsed": "13350",
    "tags": [
      {
        "key": "YWN0aW9u",
        "value": "Y3JlYXRlLXFzYw=="
      },
      {
        "key": "cXNj",
        "value": "U1RBUg=="
      },
      {
        "key": "Y3JlYXRvcg==",
        "value": "YWRkcmVzczFwYWhuZ3lyZWNqOGpoZWRrZWM4YzNqdXg4dnV4eWcyMnFzZ216cQ=="
      }
    ]
  },
  "hash": "DC7ABEB34786856276F3B7C3843DC4EB8B4B0F2053443856297787118D2E1098",
  "height": "606426"
}
```