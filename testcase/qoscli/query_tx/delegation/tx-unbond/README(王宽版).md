# test case of qoscli tx unbond

> `qoscli tx unbond` 解除委托

---

qoscli tx unbond 业务完整性的用例测试库, 包含以下部分:

* 情景1:
  
    解除委托时选择的代理验证节点owner错误，委托人未曾向该验证节点进行过委托。

* 情景2：

    解除委托时解绑的tokens大于委托人在该验证节点委托的tokens数量。

* 情景3：

    解除委托时选择的代理验证节点和解绑的tokens数量正常，或全部进行解除委托。
