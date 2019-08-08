# test case of qoscli tx delete-guardian

> `qoscli tx delete-guardian` 删除特权账户

---

## 情景说明

删除特权账户，同样执行该操作的账户必须是特权账户，且必须在genesis.json中配置。前提条件：在genesis.json文件中配置账户abc为特权账户，另外def和hij账户是通过tx add-guardian的方式添加成为特权账户。

## 测试命令

```bash
    //账户def是通过账户abc添加的特权账户
    qoscli tx delete-guardian --address hij --deleted-by def

    //账户abc是配置在genesis.json文件中的特权账户
    qoscli tx delete-guardian --address def --deleted-by abc
```

## 测试结果

```bash
    qoscli tx delete-guardian --address hij --deleted-by def
    Password to sign with 'def':
    {"check_tx":{"code":1,"log":"{\"codespace\":\"sdk\",\"code\":1,\"message\":\"TxStd's ITx ValidateData error:  ERROR:\\nCodespace: guardian\\nCode: 602\\nMessage: \\\"DeletedBy not exists or not init from genesis\\\"\\n\"}","gasWanted":"100000","gasUsed":"2432"},"deliver_tx":{},"hash":"FF29F789FB43248DC696835C9CCA6713BE2E1F5372AD46457470720C195DFFBD","height":"0"}
    ERROR: {"codespace":"sdk","code":1,"message":"TxStd's ITx ValidateData error:  ERROR:\nCodespace: guardian\nCode: 602\nMessage: \"DeletedBy not exists or not init from genesis\"\n"}

    qoscli tx delete-guardian --address def --deleted-by abc
    Password to sign with 'abc':
    {"check_tx":{"gasWanted":"100000","gasUsed":"8069"},"deliver_tx":{"gasWanted":"100000","gasUsed":"9069","tags":[{"key":"YWN0aW9u","value":"ZGVsZXRlLWd1YXJkaWFu"},{"key":"ZGVsZXRlLWJ5","value":"YWRkcmVzczEweHd4MDZnbnJ0M2RsejdoZnJ4NmE4d3gzZ3llZ2h4bTU0cnY3YQ=="},{"key":"Z3VhcmRpYW4=","value":"YWRkcmVzczFqajQ5NGE0dWd0NDhzeTgwbjNhbWc2ZHZoejB5M3lwOTRhM3B4dA=="}]},"hash":"0D5FE4776B02D8B7D6479FD38FDA954DCDEA8CB70C0E38CFD38D01B98C17EB15","height":"1856"}
```
