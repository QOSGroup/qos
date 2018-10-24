 QOS ABCI中包含以下存储集合
* base 保存区块链基本信息、QSC基本信息
* account 保存账户信息

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

