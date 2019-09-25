# Description
```
查询已存在的[name or address]
```
# Input
```
$ qoscli query account test
```
# Output
原始输出：
```
$ qoscli query account test
{"type":"qos/types/QOSAccount","value":{"base_account":{"account_address":"address1hw43pwhtscealvu973r66vk83gus8myp40fy56","public_key":{"type":"tendermint/PubKeyEd25519","value":"heAy23lzdDVvEDXHpkL8A+huCcslZDkLiFcK2Xk9J/E="},"nonce":"1"},"qos":"25130669842","qscs":null}}
```
格式化后的输出：
```
$ qoscli query account test --indent
{
  "type": "qos/types/QOSAccount",
  "value": {
    "base_account": {
      "account_address": "address1hw43pwhtscealvu973r66vk83gus8myp40fy56",
      "public_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "heAy23lzdDVvEDXHpkL8A+huCcslZDkLiFcK2Xk9J/E="
      },
      "nonce": "1"
    },
    "qos": "25169604588",
    "qscs": null
  }
}
```