# Changelog

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
