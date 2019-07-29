# test case of qoscli tx redelegate

> `qoscli tx redelegate` 变更委托

---

qoscli tx redelegate 业务完整性的用例测试库, 包含以下部分:

* 情景1:
  
    变更委托的账户没有代理验证节点。

* 情景2：

    发起变更委托的账户在当前委托人验证节点绑定的tokens小于变更委托中指定的tokens数量。（这种情况下，是否可以使用发起变更委托账户的持有qos数量来填补不够的tokens？）

* 情景3：

    变更委托的账户是验证节点，且变更的的tokens小于在当前委托人验证节点绑定的tokens数量，且自身账户足够支付gas。
