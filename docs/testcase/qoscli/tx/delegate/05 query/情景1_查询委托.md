# test case of qoscli query delegation*

> `qoscli query delegation*` 查询委托

---

## 情景说明

查询某一委托人在某一个验证节点的委托信息，不会查询出委托人的历史委托情况，只有最终的一个委托状态。前提条件：账户def在账户abc创建的验证节点上进行过委托。

## 测试命令

```bash
    qoscli query delegation --owner abc --delegator def --indent
```

## 测试结果

```bash
    qoscli query delegation --owner abc --delegator def --indent
    {
    "delegator_address": "address1l0wn66gh45nfta2r4vq8z54wu9hgarss298e9g",
    "owner_address": "address1nnvdqefva89xwppzs46vuskckr7klvzk8r5uaa",
    "validator_pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "exGS/yWJthwY8za4dlrPRid2I9KE4G15nlJwO/+Off8="
    },
    "delegate_amount": "1900743434",
    "is_compound": false
    }

```
