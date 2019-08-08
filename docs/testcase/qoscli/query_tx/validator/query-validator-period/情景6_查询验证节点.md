# test case of qoscli query validator*

> `qoscli query validator*` 查询验证节点

---

## 情景说明

查询验证节点窗口信息。查询的是实时的验证节点窗口信息。

## 测试命令

```bash
    qoscli query validator-period --owner abc --indent
```

## 测试结果

```bash
    qoscli query validator-period --owner abc --indent
    {
    "owner_address": "address1nnvdqefva89xwppzs46vuskckr7klvzk8r5uaa",
    "validator_pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "exGS/yWJthwY8za4dlrPRid2I9KE4G15nlJwO/+Off8="
    },
    "fees": "337792",
    "current_tokens": "3800743434",
    "current_period": "2814",
    "last_period": "2813",
    "last_period_fraction": {
        "value": "0.137991352306285518"
    }
}
```
