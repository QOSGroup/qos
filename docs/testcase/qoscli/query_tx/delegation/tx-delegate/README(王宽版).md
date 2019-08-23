# test case of qoscli tx delegate

> `qoscli tx delegate` 委托

---

qoscli tx delegate 业务完整性的用例测试库, 包含以下部分:

* 情景1:
  
    委托的账户没有代理验证节点。

* 情景2：

    发起委托的账户委托的tokens大于自身所持有的qos数量

* 情景3：

    委托的账户是验证节点，且绑定的tokens小于自身所持有的qos数量，且足够支付gas。
