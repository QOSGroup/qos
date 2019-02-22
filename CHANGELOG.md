# Changelog
## v0.0.4
2019.02.22

**DOWNLOAD**

[下载链接](https://github.com/QOSGroup/qos/blob/master/DOWNLOAD.md)

**FEATURES**
* [economic] 增加[公链经济模型](https://github.com/QOSGroup/qos/blob/master/docs/spec/validators/eco_module.md)
* [export] 增加公链数据导入导出
* [gas] 公链增加交易GAS的支持

**IMPROVEMENTS**
* [module] 重构模块代码
* [docs] 完善经济模型及验证人相关文档
* [client] 完善相关client command命令

**BREAKING CHANGES**
* [economic] 挖矿分配机制修改
* [kepler] 修改kepler版本依赖
* [validator] 验证人机制修改,详见[公链经济模型](https://github.com/QOSGroup/qos/blob/master/docs/spec/validators/eco_module.md)


## v0.0.3
2018.12.24

**DOWNLOAD**

[下载链接](https://github.com/QOSGroup/qos/blob/master/DOWNLOAD.md)

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
