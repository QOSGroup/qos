# Mint 模块

## 简介

参照[QOS白皮书](https://github.com/QOSGroup/whitepaper)，QOS发行总量100亿，已流通49亿，剩余51亿待主网上线分七个通胀阶段释放完成。
未开始的通胀阶段可通过`gov`模块提交`ModifyInflation`提议进行投票修改。

## 目录

1. **[概念](1_concepts.md)**
2. **[状态](2_state.md)**
3. **[ABCI](3_abci.md)**
    - [EndBlocker](3_abci.md#endblocker)
4. **[事件](4_events.md)**
    - [EndBlocker](4_events.md#endblocker)