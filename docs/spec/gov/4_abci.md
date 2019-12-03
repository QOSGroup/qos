# ABCI

## BeginBlocker

每一块开始会检查软件升级标志，升级流程如下：

1. 硬分叉，停网，需节点开发者手动编译或下载最提议指定QOS版本，下载指定`genesis.json`文件替换本地文件，执行`ｑｏｓｄ unsafe-reset-all`,然后重新新启动。
2. 一般性升级，停网，开发者编译或下载提议指定QOS版本，然后重新启动。

## EndBlocker

### 质押期提议

遍历质押结束时间在当前区块时间之前的所有质押期提议，删除提议信息，扣除所有质押至社区费池。

### 投票期提议

遍历投票期结束时间在当前区块时间之前的所有投票期提议，统计投票数据。投票通过则针对不同提议类型执行不同生效逻辑：

- ProposalTypeText 设置提议状态
- ProposalTypeParameterChange 保存提议参数修改
- ProposalTypeTaxUsage 从社区费池提取相应比例QOS至提议接收账户
- ProposalTypeModifyInflation 保存最新通胀规则
- ProposalTypeSoftwareUpgrade 设置软件升级标志