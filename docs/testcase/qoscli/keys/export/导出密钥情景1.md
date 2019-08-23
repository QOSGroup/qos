# test case of qoscli keys export

> `qoscli keys export` 导出密钥

---

## 情景说明

对于在当前节点本地存储的密钥信息，需要对某一密钥进行导出操作，使用此命令。前提条件：需要有账户，并知晓正确密码。

## 测试命令

```bash
    qoscli keys add abc   //创建情景测试的前提条件
    qoscli keys export abc // 此处的参数只接受name，地址无效
```

## 测试结果

```bash
    qoscli keys export abc
    Password to sign with 'abc':
    **Important** Don't leak your private key information to others.
Please keep your private key safely, otherwise your account will be attacked.

{"Name":"abc","address":"address1l9dw4l67mcgpxfvccg8as54k96zz2spglrc6dn","pubkey":{"type":"tendermint/PubKeyEd25519","value":"KzZkv6avo8D4yoKrUl/lZ0v8BfIwDNfmKfjENLEzh1E="},"privkey":{"type":"tendermint/PrivKeyEd25519","value":"/+tLfiTH+FKGvnoqIuI/sEWnQDtmh7+z84Ni4aY942MrNmS/pq+jwPjKgqtSX+VnS/wF8jAM1+Yp+MQ0sTOHUQ=="}}
```

ps：
    此处实现的导出：是以json字符串格式将密钥的信息打印出来，并不会自动生成对应的文件进行存储。如果需要进行存储，需要进一步操作，以文件的方式将json字符串保存至本地。
