# test case of qoscli query validator*

> `qoscli query validator*` 查询验证节点

---

## 情景说明

查询全网中所有的验证节点信息。

## 测试命令

```bash
    qoscli query validators
```

## 测试结果

```bash
    qoscli query validators
    [{"owner":"address16cvparc8ek643ghues5xsd9yl0cjvtk63r4ppl","validatorAddress":"03BF0677FFEC7F62C616F61A5226451C57AE7C6B","validatorPubkey":{"type":"tendermint/PubKeyEd25519","value":"xaj42R1aRYCyS6+J7xmUxEA4K6EnTjVxsk/7t39XTDU="},"bondTokens":"5000000000","description":{"moniker":"瑞格钱包","logo":"http://easyzone.tokenxy.cn/logo/rgqb.jpeg","website":"","details":""},"status":"active","InactiveDesc":"","inactiveTime":"0001-01-01T00:00:00Z","inactiveHeight":"0","bondHeight":"287770"},{"owner":"address1mwntal4r4parp9eh36etgh58vktc5u6frxepex","validatorAddress":"1B8080FDA9940FF7BFF87602D54140297385D3C9","validatorPubkey":{"type":"tendermint/PubKeyEd25519","value":"MlwN7rHWA0w68dBJGbZgn2IQjmjI4Kyb4hvp3CDY/Ps="},"bondTokens":"5000000000","description":{"moniker":"缔联科技","logo":"http://easyzone.tokenxy.cn/logo/dlkj.png","website":"https://deallinker.com","details":""},"status":"active","InactiveDesc":"","inactiveTime":"0001-01-01T00:00:00Z","inactiveHeight":"0","bondHeight":"287717"}]

```
