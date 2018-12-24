# Changelog

## v0.0.3
2018.12.24

**FEATURES**
* [abci] 实现最小化经济模型:验证者加入、撤销及惩罚机制
* [abci] 实现验证者挖矿奖励
* [tx] 增加TxCreateValidator、TxRevokeValidator、TxActiveValidator等transaction
* [command] 增加create-validator、revoke-validator、active-validator等操作命令
* [command] 增加validator、validators等查询命令
* [command] 增加配置root ca公钥命令


**IMPROVEMENTS**
* [docs] 最小化经济模型文档完善
* [docs] 安装文档及加入测试网文档完善
* [config] 统一add-genesis-accounts、testnet中accounts添加格式
* [command] transfer command优化

## v0.0.2
2018.11.30

**BREAKING CHANGES**
* [client] client command基于qbase重构
* [server] qosd init流程优化
* [qbase]  qbase依赖版本升级至v0.0.7

**FEATURES**
* [tx] 增加CreateValidatorTX
* [abci] endblocker增加validator更新
* [abci] 增加挖矿奖励

**IMPROVEMENTS**
* [client] 重构Approve Command

**BUG FIXES**
* [bug] fix-65:qoscli tx transfer 负值 index out of range
* [bug] 解决transfer tx GetSignData方法qos nil导致签名不一致问题
