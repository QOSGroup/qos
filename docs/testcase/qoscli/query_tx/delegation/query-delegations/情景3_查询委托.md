# test case of qoscli query delegation*

> `qoscli query delegation*` 查询委托

---

## 情景说明

查询某一委托人在所有验证节点的委托情况。前提条件：账号def在多个验证节点进行过委托。

## 测试命令

```bash
    qoscli query delegations def --indent
```

## 测试结果

```bash
    qoscli query delegations def --indent
    [
    {
        "delegator_address": "address1l0wn66gh45nfta2r4vq8z54wu9hgarss298e9g",
        "owner_address": "address1nnvdqefva89xwppzs46vuskckr7klvzk8r5uaa",
        "validator_pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "exGS/yWJthwY8za4dlrPRid2I9KE4G15nlJwO/+Off8="
        },
        "delegate_amount": "1900743434",
        "is_compound": false
    },
    {
        "delegator_address": "address1l0wn66gh45nfta2r4vq8z54wu9hgarss298e9g",
        "owner_address": "address1f66wr25emjtp5urfcpd02epwg5ply3xzcv2u20",
        "validator_pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "+GMpvsx/3zZJLCPf+EC3yqx6xaQ/tp3pKCan3TlwxWc="
        },
        "delegate_amount": "100000",
        "is_compound": false
    }
    ]

```
