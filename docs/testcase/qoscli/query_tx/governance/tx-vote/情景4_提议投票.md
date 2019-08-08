# test case of qoscli tx vote

> `qoscli tx vote` 提议投票

---

## 情景说明

投票的提议编号传入正确，且投票选项在指定范围中。前提条件：在QOS网络中存在提议编号为5的提议，且该提议处于voting的状态。

## 测试命令

```bash
    //首次投票
    qoscli tx vote --proposal-id 5 --voter abc --option yes

    //再次投票
    qoscli tx vote --proposal-id 5 --voter abc --option no
```

## 测试结果

```bash
    qoscli tx vote --proposal-id 5 --voter abc --option yes
    Password to sign with 'abc':
    {"check_tx":{"gasWanted":"100000","gasUsed":"6883"},"deliver_tx":{"gasWanted":"100000","gasUsed":"9660","tags":[{"key":"YWN0aW9u","value":"dm90ZS1wcm9wb3NhbA=="},{"key":"cHJvcG9zYWwtaWQ=","value":"NQ=="},{"key":"dm90ZXI=","value":"YWRkcmVzczEweHd4MDZnbnJ0M2RsejdoZnJ4NmE4d3gzZ3llZ2h4bTU0cnY3YQ=="}]},"hash":"5A9C8C030F763CB638919A10BE52B6F6873E58067E980C2C649AFB7B01574126","height":"532117"}

    qoscli tx vote --proposal-id 5 --voter abc --option no
    Password to sign with 'abc':
    {"check_tx":{"gasWanted":"100000","gasUsed":"6883"},"deliver_tx":{"gasWanted":"100000","gasUsed":"9660","tags":[{"key":"YWN0aW9u","value":"dm90ZS1wcm9wb3NhbA=="},{"key":"cHJvcG9zYWwtaWQ=","value":"NQ=="},{"key":"dm90ZXI=","value":"YWRkcmVzczEweHd4MDZnbnJ0M2RsejdoZnJ4NmE4d3gzZ3llZ2h4bTU0cnY3YQ=="}]},"hash":"5A9C8C030F763CB638919A10BE52B6F6873E58067E980C2C649AFB7B01574126","height":"532190"}
```

ps：
    1. 使用同一个账号对同一个提议进行多次投票，系统记录的是最后一次投票，但是每次投票的gas费用都需要支付。
