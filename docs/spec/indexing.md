# Index

## Indexing

QOS基于tendermint开发，`indexing`底层说明参照[tendermint indexing](https://tendermint.com/docs/app-dev/indexing-transactions.html).

## Tags

QOS执行交易后会添加`tag`信息，可根据`tag`[查询交易详情](../command/qoscli.md#根据标签查找交易)

### Approve

| KEY | VALUE | 说明 |
| :--- | :---: | :--- |
|   action   |   create-approve,increase-approve,decrease-approve,use-approve,cancel-approve    |   操作   |
| approve-from | address | 预授权账户地址 |
| approve-to | address | 接收预授权账户地址 |

### Gov

| KEY | VALUE | 说明 |
| :--- | :---: | :--- |
|   action   |   submit-proposal,deposit-proposal,vote-proposal   |   操作   |
| proposer | address | 提议账户地址 |
| proposal-id | proposal-id | 提议ID |
| depositor | depositor | 质押账户地址 |
| voter | address | 投票账户地址 |
| proposal-result | proposal-dropped,proposal-passed,proposal-rejected | 提议结果 |
| proposal-type | Text,Parameter,TaxUsage | 提议类型 |

### Guardian

| KEY | VALUE | 说明 |
| :--- | :---: | :--- |
|   action   |   add-guardian,delete-guardian  |   操作   |
|   guardian   |   address  |   账户地址   |
|   creator   |   address  |   创建账户地址   |
|   delete-by   |   address  |   删除账户地址   |

### QCP

| KEY | VALUE | 说明 |
| :--- | :---: | :--- |
|   action   |   init-qcp  |   操作   |
|   qcp   |   chain-id  |   联盟链chain-id  |
|   creator   |   address  |   创建账户地址   |

### QSC

| KEY | VALUE | 说明 |
| :--- | :---: | :--- |
|   action   |   create-qsc,issue-qsc  |   操作   |
|   qsc   |   qsc-name  |   联盟币名称  |
|   creator   |   address  |   创建账户地址   |
|   banker   |   address  |   联盟币发放接收地址   |

### Stake

| KEY | VALUE | 说明 |
| :--- | :---: | :--- |
|   action   |   create-validator,revoke-validator,active-validator,create-delegation,modify-compound,unbond-delegation,create-redelegation  |   操作   |
|   validator   |   address  |   验证节点地址   |
|   new-validator   |   address  |   新验证节点地址   |
|   owner   |   address  |   在验证节点所有者地址   |
|   delegator   |   address  |   代理账户地址   |

### Transfer

| KEY | VALUE | 说明 |
| :--- | :---: | :--- |
|   sender   |   address  |   发送账户地址   |
|   receiver   |   address  |   接收账户地址   |