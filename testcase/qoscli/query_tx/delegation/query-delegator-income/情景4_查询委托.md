# test case of qoscli query delegation*

> `qoscli query delegation*` 查询委托

---

## 情景说明

查询某一委托人在某一验证节点的收益信息。前提条件：账户def在账户abc创建的验证节点上进行过委托。

## 测试命令

```bash
    qoscli query delegator-income --owner abc --delegator def --indent
```

## 测试结果

```bash
    qoscli query delegator-income --owner abc --delegator def --indent
    {
    "owner_address": "address1nnvdqefva89xwppzs46vuskckr7klvzk8r5uaa",
    "validator_pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "exGS/yWJthwY8za4dlrPRid2I9KE4G15nlJwO/+Off8="
    },
    "previous_validaotr_period": "2758",
    "bond_token": "1900000000",
    "earns_starting_height": "604856",
    "first_delegate_height": "467156",
    "historical_rewards": "11352",
    "last_income_calHeight": "604856",
    "last_income_calFees": "207481"
    }

```
