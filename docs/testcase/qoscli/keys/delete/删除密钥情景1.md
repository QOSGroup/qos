# test case of qoscli keys delete

> `qoscli keys delete` 删除密钥

---

## 情景说明

对于在当前节点本地存储的密钥信息，需要对某一密钥进行删除操作，使用此命令。前提条件：需要有账户，并知晓当前正确密码。

## 测试命令

```bash
    qoscli keys add abc   //创建情景测试的前提条件
    qoscli keys delete abc // 此处的参数只接受name，地址无效
```

## 测试结果

```bash
    [vagrant@vagrant-192-168-1-201 ~]$ qoscli keys delete abc
    DANGER - enter password to permanently delete key:
    Password deleted forever (uh oh!)
```

ps:

1. 此处是不是给的提示有误：Password deleted forever (uh oh!)
2. 在删除key之前，向该账户转入一笔钱，然后进行删除操作。删除后使用查询账户的命令，依据地址查询，可以查出地址信息。该账户的钱在丢失私钥后，账户就无法恢复了，在其中的qos同样无法转出到其他账户。
