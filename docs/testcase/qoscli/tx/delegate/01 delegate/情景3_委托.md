# test case of qoscli tx delegate

> `qoscli tx delegate` 委托

---

## 情景说明

委托的账户是验证节点，且绑定的tokens小于自身所持有的qos数量，且足够支付gas。前提条件：在qos网络中abc账户为代理验证节点账户，账户def中所持有的qos数量大于50000

## 测试命令

```bash
    qoscli tx delegate --owner abc --delegator def --tokens 50000
```

## 测试结果

```bash
    qoscli tx delegate --owner abc --delegator def --tokens 50000
    Password to sign with 'def':
    {"check_tx":{"gasWanted":"100000","gasUsed":"9063"},"deliver_tx":{"gasWanted":"100000","gasUsed":"56410","tags":[{"key":"YWN0aW9u","value":"Y3JlYXRlLWRlbGVnYXRpb24="},{"key":"dmFsaWRhdG9y","value":"YWRkcmVzczE1ZmNrdW1udGNhZ25oZzBocngzNXFmZDNwN2d6ZTZsNmhobHQ3dw=="},{"key":"ZGVsZWdhdG9y","value":"YWRkcmVzczE2eGQ4dHpybTZmNGpmcm10dnk2c2phZnV5ODBsZ2owZ2p3dTh6dA=="}]},"hash":"97798C3E799E83D6C3579B6CE054D7AD5FD3ACB372C8B5833055E6D7378E4374","height":"5271"}
```
