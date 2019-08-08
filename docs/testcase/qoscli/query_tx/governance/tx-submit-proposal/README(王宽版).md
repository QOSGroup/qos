# test case of qoscli tx submit-proposal

> `qoscli tx submit-proposal` 提交提议

---

qoscli tx submit-proposal 业务完整性的用例测试库, 包含以下部分:

* 情景1:
  
    提议类型为：Text，提交提议的质押（deposit）小于规定的MinDeposit的三分之一

* 情景2：

    提议类型为：Text，提交提议的质押（deposit）不小于规定的MinDeposit的三分之一

* 情景3：

    提议类型为：ParameterChange ，提交提议的proposer质押（deposit）超出本身账户拥有的qos数量。

* 情景4：

    提议类型为：ParameterChange ，提交提议的proposer质押（deposit）超未超出本身账户拥有的qos数量，且足够支付gas。

* 情景5：

    提议类型为：TaxUsage，提交提议的目标地址非guardian账号，社区费池提取比例介于0~1。

* 情景6：

    提议类型为：TaxUsage，提交提议的目标地址为guardian账号，社区费池提取比例低于0，或高于1。

* 情景7：

    提议类型为：TaxUsage，提交提议的目标地址为guardian账号，社区费池提取比例介于0~1。

* 情景8：

    当进行一笔交易（以提交提议为例）时候，质押数量过大，所消耗的gas大于系统默认的100000.

* 情景9：

    查询提议，查询单个提议，查询所有提议
