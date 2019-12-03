# Description
```
查询区块链中已存在的交易
```
# Input
```
$ qoscli query txs --tags "message.gas.payer:address1hw43pwhtscealvu973r66vk83gus8myp40fy56" --indent
```
# Output
```
$ qoscli query txs --tags "message.gas.payer:address1hw43pwhtscealvu973r66vk83gus8myp40fy56" --indent
{
  "total_count": "2",
  "count": "2",
  "page_number": "1",
  "page_total": "1",
  "limit": "100",
  "txs": [
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
    },
    {
      "height": "3454",
      "txhash": "0CEF67AAED1ED02AB0BAD0FA4DBAB6B4806BD525BDEC0EC94C724AE32CCD7930",
      "gas_wanted": "100000",
      "gas_used": "16800",
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
              "value": "50000"
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
              "value": "50000"
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
                    "qos": "50000",
                    "qscs": null
                  }
                ],
                "receivers": [
                  {
                    "addr": "address1qnhak3ph0yqpxar3rrkzuasgnzlzfmq4pyn73m",
                    "qos": "50000",
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
              "signature": "OzqUU6vsXyLydTmhdl+A30xh0EIBuhBIjQ7icVNjgHq3riBkydZT8FaHVP3iH5rU8CdAnmBVnvC2ilvKCG3OAQ==",
              "nonce": "3"
            }
          ],
          "chainid": "test-chain",
          "maxgas": "100000"
        }
      },
      "timestamp": "2019-08-07T02:28:31Z"
    }
  ]
}
```
