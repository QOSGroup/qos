# test case of qoscli tx modify-compound

> `qoscli tx modify-compound` 修改收益复投方式

---

qoscli tx modify-compound 业务完整性的用例测试库, 包含以下部分:

* 情景1:
  
    委托人在其没有委托的代理验证节点操作账户进行修改收益复投方式。

* 情景2：

    委托人正常修改在某一代理验证节点的收益复投方式。连续操作，由于还未超过收益期限，所以会导致失败。
