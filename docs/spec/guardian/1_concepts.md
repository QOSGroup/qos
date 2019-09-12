# 概念

## 系统账户

系统账户是由QOS基金会管理，对QOS网络正常运行行使特殊职能的账户。

系统账户有两种：
- `Genesis` 在`genesis.json`中初始添加的类型
- `Ordinary` 由`Genesis`类型系统账户创建

`Genesis`和`Ordinary`系统账户均执行停止网络，发起从社区费池提取QOS提议操作，另外`Genesis`账户还可以添加/删除`Oridinary`类型系统账户。