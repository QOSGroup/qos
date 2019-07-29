# Description
```
正常查询查询QCP Out Tx
```
# Input
```
$ qoscli query qcp tx qcp-rocket --seq 1 --indent
```
# Output
```
$ qoscli query qcp tx qcp-rocket --seq 1 --indent
{
  "type": "qbase/txs/qcptx",
  "value": {
    "txstd": {
      "itx": [
        {
          "type": "qbase/txs/qcpresult",
          "value": {
            "result": {
              "Code": 0,
              "Codespace": "",
              "Data": null,
              "Log": "",
              "GasWanted": "100000",
              "GasUsed": "26520",
              "FeeAmount": "0",
              "FeeDenom": "",
              "Tags": [
                {
                  "key": "c2VuZGVy",
                  "value": "YWRkcmVzczF1OWt0bnplNXNya3JndGR0YTRndHhqbGh2eXFmMmc3enhnbTR3cQ=="
                },
                {
                  "key": "cmVjZWl2ZXI=",
                  "value": "YWRkcmVzczFnNDc3N2Z4ZWphY2FndHRjbjJtdWxqbXZuN3I1MnJnaGZkZHpscQ=="
                },
                {
                  "key": "cmVjZWl2ZXI=",
                  "value": "YWRkcmVzczFrbDdxczdoNzY1eXpucmNuZ3d0NmVzN2ZhbnV3a2tucnI1dXJuZw=="
                },
                {
                  "key": "cWNwLmZyb20=",
                  "value": "Y2Fwcmljb3JuLTMwMDA="
                },
                {
                  "key": "cWNwLnRv",
                  "value": "cWNwLXJvY2tldA=="
                }
              ]
            },
            "qcporiginalsequence": "1",
            "qcpextends": "",
            "info": ""
          }
        }
      ],
      "sigature": null,
      "chainid": "qcp-rocket",
      "maxgas": "0"
    },
    "from": "capricorn-3000",
    "to": "qcp-rocket",
    "sequence": "1",
    "sig": {
      "pubkey": null,
      "signature": null,
      "nonce": "0"
    },
    "blockheight": "620570",
    "txindex": "0",
    "isresult": true,
    "extends": ""
  }
}
```
