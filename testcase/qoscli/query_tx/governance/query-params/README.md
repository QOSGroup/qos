# Test Cases

- 暂无

# Description
>     Query the parameters of the governance process. 

>     查询治理过程的参数。

# Example

`qoscli query params --module <module> --key <key_name>`

主要参数：

- `--module`       模块名称：`stake`、`gov`、`distribution`
- `--key`          参数名

查询所有参数：
```bash
$ qoscli query params --indent
[
  {
    "type": "stake",
    "value": {
      "max_validator_cnt": 10,
      "voting_status_len": 100,
      "voting_status_least": 50,
      "survival_secs": 600,
      "unbond_return_height": 10
    }
  },
  {
    "type": "distribution",
    "value": {
      "proposer_reward_rate": {
        "value": "0.040000000000000000"
      },
      "community_reward_rate": {
        "value": "0.010000000000000000"
      },
      "validator_commission_rate": {
        "value": "0.010000000000000000"
      },
      "delegator_income_period_height": "10",
      "gas_per_unit_cost": "10"
    }
  },
  {
    "type": "gov",
    "value": {
      "min_deposit": "10000000",
      "max_deposit_period": "172800000000000",
      "voting_period": "172800000000000",
      "quorum": "0.334000000000000000",
      "threshold": "0.500000000000000000",
      "veto": "0.334000000000000000",
      "penalty": "0.000000000000000000"
    }
  }
]
```

查询`gov`模块下参数：
```bash
$ qoscli query params --module gov --indent
{
  "type": "gov",
  "value": {
    "min_deposit": "10000000",
    "max_deposit_period": "172800000000000",
    "voting_period": "172800000000000",
    "quorum": "0.334000000000000000",
    "threshold": "0.500000000000000000",
    "veto": "0.334000000000000000",
    "penalty": "0.000000000000000000"
  }
}
```

查询`gov`模块下`min_deposit`参数值：
```bash
$ qoscli query params --module gov --key min_deposit
"10000000"
```

# Usage
```
  qoscli query params [flags]
```

# Available Commands

>无可用命令

# Flags

| ShortCut | Flag           | Required | Input Type | Default Input             | Input Range | Description                             |
|:---------|:---------------|:---------|:-----------|:--------------------------|:------------|:----------------------------------------|
| `-h`     | `--help`       | ✖        | -          | -                         | -           | 帮助文档                                    |
| -        | `--chain-id`   | ✖        | string     | -                         | -           | Tendermint节点的链ID                        |
| -        | `--height`     | ✖        | int        | -                         | -           | (可选)要查询的块高度，省略以获取最新的可证明块                |
| -        | `--indent`     | ✖        | -          | -                         | -           | 向JSON响应添加缩进                             |
| `-n`     | `--node`       | ✖        | string     | `"tcp://localhost:26657"` | -           | 为此链提供的Tendermint RPC接口: `<host>:<port>` |
| -        | `--trust-node` | ✖        | -          | -                         | -           | 是否信任连接的完整节点（不验证其响应证据）                   |
| -        | `--key`        | ✖        | string     | -                         | -           | (主要参数)参数名称                                    |
| -        | `--module`     | ✖        | string     | -                         | -           | (主要参数)模块名称                                    |


# Global Flags

| ShortCut | Flag         | Required | Input Type | Default Input | Input Range       | Description  |
|:---------|:-------------|:---------|:-----------|:--------------|:------------------|:-------------|
| `-e`     | `--encoding` | ✖        | string     | `hex`         | `hex`/`b64`/`btc` | 二进制编码        |
| -        | `--home`     | ✖        | string     | `/.qoscli`    | -                 | 配置和数据的目录     |
| `-o`     | `--output`   | ✖        | string     | `text`        | `text`/`json`     | 输出格式         |
| -        | `--trace`    | ✖        | -          | -             | -                 | 打印出错时的完整堆栈跟踪 |