 QOS ABCI中包含以下存储集合
* base 保存区块链基本信息、QSC基本信息
* account 保存账户信息
* qcp 保存qcp交易sequence,tx列表

### base

* 区块链基本信息

```
base/chainId //qos或者qscname
base/rootca
```

* QSC初始信息

```
qsc/[name]:{name,pubkey,bankerAddress,createAddress,exrate,CA,description}
```

### account


```
[accountAddress]:{}
```

### qcp
```
[chainId]/out/sequence //需要输出到"chainId"的qcp tx最大序号
[chainId]/out/tx_[sequence] //需要输出到"chainId"的每个qcp tx
[chainId]/in/sequence //已经接受到来自"chainId"的qcp tx最大序号
[chainId]/in/pubkey //接受来自"chainId"的合法公钥
```
