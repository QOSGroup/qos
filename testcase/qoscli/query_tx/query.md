# Description

>     query(alias `q`) subcommands. 

>     查询子命令。

# Usage
```
  qoscli query [command]
```

# Alias
```
  qoscli q [command]
```

# Available Commands

## Group 01 - Others
| Command                  | Alias                | Has-Subcommand | Description                    |
|:-------------------------|:---------------------|:---------------|:-------------------------------|
| `qoscli query account`   | `qoscli q account`   | ✖              | 按地址(address)或名称(name)查询帐户信息    |
| `qoscli query store`     | `qoscli q store`     | ✖              | 按低级(low level)查询存储数据           |
| `qoscli query consensus` | `qoscli q consensus` | ✖              | 查询共识参数                         |
| `qoscli query qcp`       | `qoscli q qcp`       | ✔              | QCP子命令                         |
| `qoscli query qsc`       | `qoscli q qsc`       | ✖              | 按名称(name)查询QSC信息               |
| `qoscli query approve`   | `qoscli q approve`   | ✖              | 按来源(from)和目标(to)查询预授权(approve) |

## Group 02 - Validator & Delegation
| Command                            | Alias                          | Has-Subcommand | Description                    |
|:-----------------------------------|:-------------------------------|:---------------|:-------------------------------|
| `qoscli query validators`          | `qoscli q validators`          | ✖              | 查询所有验证人的信息                     |
| `qoscli query validator`           | `qoscli q validator`           | ✖              | 查询验证人的信息                       |
| `qoscli query validator-miss-vote` | `qoscli q validator-miss-vote` | ✖              | 查询最近投票窗口中的验证人错过投票的信息           |
| `qoscli query delegation`          | `qoscli q delegation`          | ✖              | 查询委托信息                         |
| `qoscli query delegations`         | `qoscli q delegations`         | ✖              | 查询一个委托人的所有委托                   |
| `qoscli query delegations-to`      | `qoscli q delegations-to`      | ✖              | 查询一个验证人所接收的所有委托                |
| `qoscli query unbondings`          | `qoscli q unbondings`          | ✖              | 查询一个委托人的所有正在解除绑定(unbonding)的委托 |
| `qoscli query redelegations`       | `qoscli q redelegations`       | ✖              | 查询一个委托人的所有重委托请求(redelegations) |

## Group 03 - Distribution
| Command                           | Alias                         | Has-Subcommand | Description               |
|:----------------------------------|:------------------------------|:---------------|:--------------------------|
| `qoscli query validator-period`   | `qoscli q validator-period`   | ✖              | 查询分发(distribution)验证人周期信息 |
| `qoscli query delegator-income`   | `qoscli q delegator-income`   | ✖              | 查询分发(distribution)委托人收入信息 |
| `qoscli query community-fee-pool` | `qoscli q community-fee-pool` | ✖              | 查询社区费用池                   |

## Group 04 - Proposal
| Command                  | Alias                | Has-Subcommand | Description                     |
|:-------------------------|:---------------------|:---------------|:--------------------------------|
| `qoscli query proposal`  | `qoscli q proposal`  | ✖              | 查询单个提案(proposal)的详细信息           |
| `qoscli query proposals` | `qoscli q proposals` | ✖              | 使用可选的筛选器查询提案(proposal)          |
| `qoscli query vote`      | `qoscli q vote`      | ✖              | 查询单个投票(vote)的详细信息               |
| `qoscli query votes`     | `qoscli q votes`     | ✖              | 查询对指定提案(proposal)的投票(vote)      |
| `qoscli query deposit`   | `qoscli q deposit`   | ✖              | 查询单个抵押存款(deposit)的详细信息          |
| `qoscli query deposits`  | `qoscli q deposits`  | ✖              | 查询对指定提案(proposal)的抵押存款(deposit) |
| `qoscli query tally`     | `qoscli q tally`     | ✖              | 获得提案投票的计票结果                     |
| `qoscli query params`    | `qoscli q params`    | ✖              | 查询治理过程的参数                       |

## Group 05 - Guardian
| Command                  | Alias                | Has-Subcommand | Description        |
|:-------------------------|:---------------------|:---------------|:-------------------|
| `qoscli query guardian`  | `qoscli q guardian`  | ✖              | 查询特权用户(guardian)   |
| `qoscli query guardians` | `qoscli q guardians` | ✖              | 查询特权用户(guardian)列表 |

## Group 06 - Status
| Command               | Alias             | Has-Subcommand | Description |
|:----------------------|:------------------|:---------------|:------------|
| `qoscli query status` | `qoscli q status` | ✖              | 查询远程节点的状态   |

## Group 07 - Basic
| Command                              | Alias                            | Has-Subcommand | Description            |
|:-------------------------------------|:---------------------------------|:---------------|:-----------------------|
| `qoscli query tendermint-validators` | `qoscli q tendermint-validators` | ✖              | 获取给定高度的tendermint验证人集合 |
| `qoscli query block`                 | `qoscli q block`                 | ✖              | 获取给定高度的区块信息            |

## Group 08 - Tx
| Command            | Alias          | Has-Subcommand | Description     |
|:-------------------|:---------------|:---------------|:----------------|
| `qoscli query txs` | `qoscli q txs` | ✖              | 分页查询与一组tag匹配的交易 |
| `qoscli query tx`  | `qoscli q tx`  | ✖              | 在提交的块中按哈希查询交易   |

# Flags

| ShortCut | Flag      | Input Type | Default Input | Input Range | Description            |
|:---------|:----------|:-----------|:--------------|:------------|:-----------------------|
| `-h`     | `--help`  | -          | -             | -           | (可选)帮助文档                   |

# Global Flags

| ShortCut | Flag         | Input Type | Default Input | Input Range       | Description  |
|:---------|:-------------|:-----------|:--------------|:------------------|:-------------|
| `-e`     | `--encoding` | string     | `hex`         | `hex`/`b64`/`btc` | 二进制编码        |
| -        | `--home`     | string     | `/.qoscli`    | -                 | 配置和数据的目录     |
| `-o`     | `--output`   | string     | `text`        | `text`/`json`     | 输出格式         |
| -        | `--trace`    | -          | -             | -                 | 打印出错时的完整堆栈跟踪 |
