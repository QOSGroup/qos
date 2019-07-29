# test case of qoscli query validator*

> `qoscli query validator*` 查询验证节点

---

## 情景说明

查询指定高度验证节点的集合。前提条件：当前的高度小于800000

## 测试命令

```bash
    qoscli query tendermint-validators 287730 --indent

    qoscli query tendermint-validators 0 --indent

    qoscli query tendermint-validators 800000 --indent
```

## 测试结果

```bash
    qoscli query tendermint-validators 287730 --indent
    current query height: 287730
    [
    {
        "Address": "address1rwqgpldfjs8l00lcwcpd2s2q99ect57fvluzp5",
        "VotingPower": "5000000000",
        "PubKey": {
        "type": "tendermint/PubKeyEd25519",
        "value": "MlwN7rHWA0w68dBJGbZgn2IQjmjI4Kyb4hvp3CDY/Ps="
        }
    },
    {
        "Address": "address122zarynhujlfydx59sw507k85vk4jp3uxezt0t",
        "VotingPower": "52500000000",
        "PubKey": {
        "type": "tendermint/PubKeyEd25519",
        "value": "XJEifEjxC6ik+UTOMka5V+HJVsjlhKE69CbNf6Yspas="
        }
    },
    {
        "Address": "address1tv6hlx9y2f3jmekmaa5qr7tves0dt3guhnzakt",
        "VotingPower": "5000000000",
        "PubKey": {
        "type": "tendermint/PubKeyEd25519",
        "value": "smATFP0NXD2uqbNazkiBDyMUj+GRgR3TerMahQyfiRo="
        }
    }
    ]

    qoscli query tendermint-validators 0 --indent
    ERROR: Validators: Response error: RPC error -32603 - Internal error: Height must be greater than 0

    qoscli query tendermint-validators 800000 --indent
    ERROR: Validators: Response error: RPC error -32603 - Internal error: Height must be less than or equal to the current blockchain height

```
