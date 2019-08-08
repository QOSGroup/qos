# test case of qoscli tx add-guardian

> `qoscli tx add-guardian` 添加特权账户

---

## 情景说明

添加特权账户，特权账户的添加只能是由特权账户来完成，也就是在genesis.json中配置的特权账户。前提条件：知晓在genesis.json文件中的特权账户的密码或是私钥。测试命令执行的前提是账户abc配置在genesis.json文件中，本地密钥库也保存账户abc的信息。

## 测试命令

```bash
    //账户abc是配置在genesis.json文件中的特权账户
    qoscli tx add-guardian --address def --creator abc --description 'set def to be a guardian'

    //账户def是通过账户abc添加的特权账户
    qoscli tx add-guardian --address hij --creator def --description 'set hij to be a guardian'
```

## 测试结果

```bash
    qoscli tx add-guardian --address def --creator abc --description 'set def to be a guardian'
    Password to sign with 'abc':
    {"check_tx":{"gasWanted":"100000","gasUsed":"7856"},"deliver_tx":{"gasWanted":"100000","gasUsed":"12046","tags":[{"key":"YWN0aW9u","value":"YWRkLWd1YXJkaWFu"},{"key":"Y3JlYXRvcg==","value":"YWRkcmVzczEweHd4MDZnbnJ0M2RsejdoZnJ4NmE4d3gzZ3llZ2h4bTU0cnY3YQ=="},{"key":"Z3VhcmRpYW4=","value":"YWRkcmVzczFsNmp1YXF5OWZrMGRwczBmbjVkY2c0ZnB5MzZ6bXJ5cDhteTR1eA=="}]},"hash":"857BF0332E9FB1F0B89378833FCA1D06E1543464465E7DF4BF46AC417935CCEC","height":"203"}

    qoscli tx add-guardian --address hij --creator def --description 'set hij to be a guardian'
    Password to sign with 'def':
    {"check_tx":{"code":1,"log":"{\"codespace\":\"sdk\",\"code\":1,\"message\":\"TxStd's ITx ValidateData error:  ERROR:\\nCodespace: guardian\\nCode: 602\\nMessage: \\\"Creator not exists or not init from genesis\\\"\\n\"}","gasWanted":"100000","gasUsed":"2213"},"deliver_tx":{},"hash":"A935D93832DC2AEE23EC00813BA77D5C0280ECB0C8C7B50E34E6CAB75030360F","height":"0"}
    ERROR: {"codespace":"sdk","code":1,"message":"TxStd's ITx ValidateData error:  ERROR:\nCodespace: guardian\nCode: 602\nMessage: \"Creator not exists or not init from genesis\"\n"}

```
