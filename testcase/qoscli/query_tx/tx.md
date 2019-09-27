# Description

>     tx subcommands.

>     交易子命令。

# Usage
```
  qoscli tx [command]
```

# Available Commands

## Group 01 - QSC
| Command                | Alias | Has-Subcommand | Description |
|:-----------------------|:------|:---------------|:------------|
| `qoscli tx create-qsc` | -     | ✖              | 创建QSC       |
| `qoscli tx issue-qsc`  | -     | ✖              | 发布QSC       |

## Group 02 - QCP
| Command              | Alias | Has-Subcommand | Description |
|:---------------------|:------|:---------------|:------------|
| `qoscli tx init-qcp` | -     | ✖              | 初始化QCP      |

## Group 03 - Transfer
| Command              | Alias | Has-Subcommand | Description |
|:---------------------|:------|:---------------|:------------|
| `qoscli tx transfer` | -     | ✖              | 转账QOS和QSC   |

## Group 04 - Approve
| Command                      | Alias | Has-Subcommand | Description |
|:-----------------------------|:------|:---------------|:------------|
| `qoscli tx create-approve`   | -     | ✖              | 创建预授权       |
| `qoscli tx increase-approve` | -     | ✖              | 增加预授权       |
| `qoscli tx decrease-approve` | -     | ✖              | 减少预授权       |
| `qoscli tx use-approve`      | -     | ✖              | 使用预授权       |
| `qoscli tx cancel-approve`   | -     | ✖              | 取消预授权       |

## Group 05 - Validator
| Command                      | Alias | Has-Subcommand | Description    |
|:-----------------------------|:------|:---------------|:---------------|
| `qoscli tx create-validator` | -     | ✖              | 创建用自委托初始化的新验证人 |
| `qoscli tx modify-validator` | -     | ✖              | 修改已存在的验证人账户    |
| `qoscli tx revoke-validator` | -     | ✖              | 撤销(Revoke)验证人  |
| `qoscli tx active-validator` | -     | ✖              | 激活(Active)验证人  |

## Group 06 - Delegation
| Command                     | Alias | Has-Subcommand | Description           |
|:----------------------------|:------|:---------------|:----------------------|
| `qoscli tx delegate`        | -     | ✖              | 向验证人委托QOS             |
| `qoscli tx modify-compound` | -     | ✖              | 修改一个委托的复投信息           |
| `qoscli tx unbond`          | -     | ✖              | 从验证人解除QOS委托           |
| `qoscli tx redelegate`      | -     | ✖              | 将QOS从一个验证人重新委托到另一个验证人 |

## Group 07 - Proposal
| Command                     | Alias | Has-Subcommand | Description |
|:----------------------------|:------|:---------------|:------------|
| `qoscli tx submit-proposal` | -     | ✖              | 发起提案        |
| `qoscli tx deposit`         | -     | ✖              | 向提案抵押存款     |
| `qoscli tx vote`            | -     | ✖              | 向提案投票       |

## Group 08 - Guardian
| Command                     | Alias | Has-Subcommand | Description      |
|:----------------------------|:------|:---------------|:-----------------|
| `qoscli tx add-guardian`    | -     | ✖              | 添加特权用户(guardian) |
| `qoscli tx delete-guardian` | -     | ✖              | 删除特权用户(guardian) |

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
