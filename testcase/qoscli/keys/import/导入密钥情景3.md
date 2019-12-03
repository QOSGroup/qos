# test case of qoscli keys import

> `qoscli keys import` 导入密钥

---

## 情景说明

对于在当前节点本地存储的密钥信息，需要将外部的密钥（没有保存导出密钥命令信息，只记录了私钥信息）导入至节点本地存储库中，使用此命令。前提条件：需要有账户，某种情形下（例如：节点发生宕机导致本地密钥库奔溃）发生了账户被删除。

私钥信息：/+tLfiTH+FKGvnoqIuI/sEWnQDtmh7+z84Ni4aY942MrNmS/pq+jwPjKgqtSX+VnS/wF8jAM1+Yp+MQ0sTO`HUA`==

修改某些字符后的私钥：/+tLfiTH+FKGvnoqIuI/sEWnQDtmh7+z84Ni4aY942MrNmS/pq+jwPjKgqtSX+VnS/wF8jAM1+Yp+MQ0sTO`ABC`==

## 测试命令

```bash
    qoscli keys delete abc   //创建情景测试的前提条件
    qoscli keys import abc
```

## 测试结果

```bash
    qoscli keys import abc
    > Enter ed25519 private key:
    /+tLfiTH+FKGvnoqIuI/sEWnQDtmh7+z84Ni4aY942MrNmS/pq+jwPjKgqtSX+VnS/wF8jAM1+Yp+MQ0sTOABC==
    > Enter a passphrase for your key:
    > Repeat the passphrase:

    qoscli keys list  // 验证导入是否成功
```

ps：
    从结果可以看出，导入是成功的，但是地址和公钥都不是我们原来的，相当于创建了新的账户和地址，与原来的账户和地址毫无关联。
