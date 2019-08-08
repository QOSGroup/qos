# test case of qoscli tx modify-compound

> `qoscli tx modify-compound` 修改收益复投方式

---

## 情景说明

委托人在其没有委托的代理验证节点操作账户进行修改收益复投方式。前提条件：委托人delegator在验证人owner没有委托情况。

## 测试命令

```bash
    qoscli tx modify-compound --owner address1f66wr25emjtp5urfcpd02epwg5ply3xzcv2u20 --delegator jlgy02 --compound true
```

## 测试结果

```bash
    qoscli tx modify-compound --owner address1f66wr25emjtp5urfcpd02epwg5ply3xzcv2u20 --delegator jlgy02 --compound
    Password to sign with 'jlgy02':
    {"check_tx":{"code":1,"log":"{\"codespace\":\"sdk\",\"code\":1,\"message\":\"TxStd's ITx ValidateData error:  ERROR:\\nCodespace: stake\\nCode: 501\\nMessage: \\\"delegator not delegate the owner's validator\\\"\\n\"}","gasWanted":"100000","gasUsed":"3687"},"deliver_tx":{},"hash":"DE85BB5C30196720AF5DC3CBDF4A0913777F09619B94FF17EA5D817BF46BF408","height":"0"}
    ERROR: {"codespace":"sdk","code":1,"message":"TxStd's ITx ValidateData error:  ERROR:\nCodespace: stake\nCode: 501\nMessage: \"delegator not delegate the owner's validator\"\n"}

```
