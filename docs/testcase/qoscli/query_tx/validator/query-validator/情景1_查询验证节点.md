# test case of qoscli query validator*

> `qoscli query validator*` 查询验证节点

---

## 情景说明

根据验证节点操作者查找与其绑定的所有验证节点信息。前提条件：使用的账户地址是某一个验证节点的操作者。
ps：同一个账户只能是一个验证节点的操作者，不能同时是多个验证节点的操作者。

## 测试命令

```bash
    qoscli query validator address10xwx06gnrt3dlz7hfrx6a8wx3gyeghxm54rv7a
```

## 测试结果

```bash
    qoscli query validator address10xwx06gnrt3dlz7hfrx6a8wx3gyeghxm54rv7a
    {"owner":"address10xwx06gnrt3dlz7hfrx6a8wx3gyeghxm54rv7a","validatorAddress":"A2716E6E6BC7513BA1F719A34025B10F902CEBFA","validatorPubkey":{"type":"tendermint/PubKeyEd25519","value":"aIm1GNnTm6ITZIpt7zfViP9Mc6jLrIF8TZtnSZ3TKB4="},"bondTokens":"100000000000000","description":{"moniker":"validator201","logo":"","website":"","details":""},"status":"active","InactiveDesc":"","inactiveTime":"0001-01-01T00:00:00Z","inactiveHeight":"0","bondHeight":"0"}

```
