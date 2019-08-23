# Description
```
查询区块链中已存在的[hash]
```
# Input
- 原始输出：
```
$ qoscli query tx 19B5D448A55B99C03B2AB435BBE97F55DEDF921E5A957C4CD03C6A972E10DD29
```
- 格式化输出：
```
$ qoscli query tx 19B5D448A55B99C03B2AB435BBE97F55DEDF921E5A957C4CD03C6A972E10DD29 --indent
```
# Output
- 原始输出：
```
$ qoscli query tx 19B5D448A55B99C03B2AB435BBE97F55DEDF921E5A957C4CD03C6A972E10DD29
{"height":"2876","txhash":"19B5D448A55B99C03B2AB435BBE97F55DEDF921E5A957C4CD03C6A972E10DD29","gas_wanted":"100000","gas_used":"16640","events":[{"type":"message","attributes":[{"key":"module","value":"transfer"},{"key":"gas.payer","value":"address1hw43pwhtscealvu973r66vk83gus8myp40fy56"}]},{"type":"receive","attributes":[{"key":"address","value":"address1qnhak3ph0yqpxar3rrkzuasgnzlzfmq4pyn73m"},{"key":"qos","value":"100"},{"key":"qscs"}]},{"type":"send","attributes":[{"key":"address","value":"address1hw43pwhtscealvu973r66vk83gus8myp40fy56"},{"key":"qos","value":"100"},{"key":"qscs"}]}],"tx":{"type":"qbase/txs/stdtx","value":{"itx":[{"type":"transfer/txs/TxTransfer","value":{"senders":[{"addr":"address1hw43pwhtscealvu973r66vk83gus8myp40fy56","qos":"100","qscs":null}],"receivers":[{"addr":"address1qnhak3ph0yqpxar3rrkzuasgnzlzfmq4pyn73m","qos":"100","qscs":null}]}}],"sigature":[{"pubkey":{"type":"tendermint/PubKeyEd25519","value":"heAy23lzdDVvEDXHpkL8A+huCcslZDkLiFcK2Xk9J/E="},"signature":"McBrtjJskel+BB3U2sPKBCbNU2KoEVxrPf3vRz8uyBbC6YHriWIMTS3yh2rvuZ7EaDnqcCFRC6eJCtkLhQ3ACw==","nonce":"2"}],"chainid":"test-chain","maxgas":"100000"}},"timestamp":"2019-08-06T03:38:09Z"}
```
- 格式化输出：
```
$ qoscli query tx 19B5D448A55B99C03B2AB435BBE97F55DEDF921E5A957C4CD03C6A972E10DD29 --indent
{
  "height": "2876",
  "txhash": "19B5D448A55B99C03B2AB435BBE97F55DEDF921E5A957C4CD03C6A972E10DD29",
  "gas_wanted": "100000",
  "gas_used": "16640",
  "events": [
    {
      "type": "message",
      "attributes": [
        {
          "key": "module",
          "value": "transfer"
        },
        {
          "key": "gas.payer",
          "value": "address1hw43pwhtscealvu973r66vk83gus8myp40fy56"
        }
      ]
    },
    {
      "type": "receive",
      "attributes": [
        {
          "key": "address",
          "value": "address1qnhak3ph0yqpxar3rrkzuasgnzlzfmq4pyn73m"
        },
        {
          "key": "qos",
          "value": "100"
        },
        {
          "key": "qscs"
        }
      ]
    },
    {
      "type": "send",
      "attributes": [
        {
          "key": "address",
          "value": "address1hw43pwhtscealvu973r66vk83gus8myp40fy56"
        },
        {
          "key": "qos",
          "value": "100"
        },
        {
          "key": "qscs"
        }
      ]
    }
  ],
  "tx": {
    "type": "qbase/txs/stdtx",
    "value": {
      "itx": [
        {
          "type": "transfer/txs/TxTransfer",
          "value": {
            "senders": [
              {
                "addr": "address1hw43pwhtscealvu973r66vk83gus8myp40fy56",
                "qos": "100",
                "qscs": null
              }
            ],
            "receivers": [
              {
                "addr": "address1qnhak3ph0yqpxar3rrkzuasgnzlzfmq4pyn73m",
                "qos": "100",
                "qscs": null
              }
            ]
          }
        }
      ],
      "sigature": [
        {
          "pubkey": {
            "type": "tendermint/PubKeyEd25519",
            "value": "heAy23lzdDVvEDXHpkL8A+huCcslZDkLiFcK2Xk9J/E="
          },
          "signature": "McBrtjJskel+BB3U2sPKBCbNU2KoEVxrPf3vRz8uyBbC6YHriWIMTS3yh2rvuZ7EaDnqcCFRC6eJCtkLhQ3ACw==",
          "nonce": "2"
        }
      ],
      "chainid": "test-chain",
      "maxgas": "100000"
    }
  },
  "timestamp": "2019-08-06T03:38:09Z"
}
```