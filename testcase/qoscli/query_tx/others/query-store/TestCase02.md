# Description
```
提供必要参数`--path`与`--data`
```
# Input
- 提供必要参数`--path`与`--data`：
```
$ qoscli query store --path /store/acc/subspace --data account
```
- 格式化输出：
```
$ qoscli query store --path /store/acc/subspace --data account --indent
```
- 提供错误的`--path`：
```
$ qoscli query store --path /store/name/subspace --data account --indent
```
- 提供错误的`--data`：
```
$ qoscli query store --path /store/acc/subspace --data name --indent
```
# Output
- 提供必要参数`--path`与`--data`：
```
$ qoscli query store --path /store/acc/subspace --data account
[{"key":"account:\ufffd\ufffd\u0010\ufffd\ufffd\ufffd3߳\ufffd\ufffdG\ufffd2Ǌ9\u0003\ufffd\ufffd","value":{"type":"qos/types/QOSAccount","value":{"base_account":{"account_address":"address1hw43pwhtscealvu973r66vk83gus8myp40fy56","public_key":{"type":"tendermint/PubKeyEd25519","value":"heAy23lzdDVvEDXHpkL8A+huCcslZDkLiFcK2Xk9J/E="},"nonce":"1"},"qos":"26303874372","qscs":null}}}]
```
- 格式化输出：
```
$ qoscli query store --path /store/acc/subspace --data account --indent
[
  {
    "key": "account:\ufffd\ufffd\u0010\ufffd\ufffd\ufffd3߳\ufffd\ufffdG\ufffd2Ǌ9\u0003\ufffd\ufffd",
    "value": {
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
        "qos": "26338279674",
        "qscs": null
      }
    }
  }
]
```
- 提供错误的`--path`：
```
$ qoscli query store --path /store/name/subspace --data account --indent
ERROR: response empty value
```
- 提供错误的`--data`：
```
$ qoscli query store --path /store/acc/subspace --data name --indent
null
```