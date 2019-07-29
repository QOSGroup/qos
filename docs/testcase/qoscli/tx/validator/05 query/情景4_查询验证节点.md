# test case of qoscli query validator*

> `qoscli query validator*` 查询验证节点

---

## 情景说明

查询某一验证节点的漏块信息。前提条件：abc账户为某一验证节点的操作者。

## 测试命令

```bash
    qoscli query validator-miss-vote abc
```

## 测试结果

```bash
    qoscli query validator-miss-vote abc
    ERROR: response empty value. getStakeConfig is empty

```

ps：当前版本v0.0.5的qos，使用该命令存在报错。属于系统bug，有待修复。
