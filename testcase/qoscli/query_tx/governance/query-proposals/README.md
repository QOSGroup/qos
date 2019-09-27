# Test Cases

- 暂无

# Description
>     Query for a all proposals. 

>     使用可选的筛选器查询提案(proposal)。

您可以使用以下Flag标志筛选返回值:
```
$ qos query gov proposals --depositor wzj
```
```
$ qos query gov proposals --voter wzj
```
```
$ qos query gov proposals --status (DepositPeriod|VotingPeriod|Passed|Rejected).
```

# Example

`qoscli query proposals`

查询所有提议：
```bash
$ qoscli query proposals
```

查询结果：
```bash
[
  {
    "proposal_content": {
      "type": "gov/TextProposal",
      "value": {
        "title": "update qos",
        "description": "this is the description",
        "deposit": "100000000"
      }
    },
    "proposal_id": "1",
    "proposal_status": 2,
    "final_tally_result": {
      "yes": "0",
      "abstain": "0",
      "no": "0",
      "no_with_veto": "0"
    },
    "submit_time": "2019-04-03T08:20:34.99523986Z",
    "deposit_end_time": "2019-04-05T08:20:34.99523986Z",
    "total_deposit": "200000000",
    "voting_start_time": "2019-04-03T08:20:34.99523986Z",
    "voting_start_height": "700",
    "voting_end_time": "2019-04-05T08:20:34.99523986Z"
  }
]
```

# Usage
```
  qoscli query proposals [flags]
```

# Available Commands

>无可用命令

# Flags

| ShortCut | Flag           | Required | Input Type | Default Input             | Input Range                                             | Description                             |
|:---------|:---------------|:---------|:-----------|:--------------------------|:--------------------------------------------------------|:----------------------------------------|
| `-h`     | `--help`       | ✖        | -          | -                         | -                                                       | 帮助文档                                    |
| -        | `--chain-id`   | ✖        | string     | -                         | -                                                       | Tendermint节点的链ID                        |
| -        | `--height`     | ✖        | int        | -                         | -                                                       | (可选)要查询的块高度，省略以获取最新的可证明块                |
| -        | `--indent`     | ✖        | -          | -                         | -                                                       | 向JSON响应添加缩进                             |
| `-n`     | `--node`       | ✖        | string     | `"tcp://localhost:26657"` | -                                                       | 为此链提供的Tendermint RPC接口: `<host>:<port>` |
| -        | `--trust-node` | ✖        | -          | -                         | -                                                       | 是否信任连接的完整节点（不验证其响应证据）                   |
| -        | `--depositor`  | ✖        | string     | -                         | -                                                       | (主要参数)按抵押人抵押的提案筛选                       |
| -        | `--limit`      | ✖        | string     | -                         | -                                                       | (主要参数)限制为最新的[number]提案，默认为所有提案          |
| -        | `--status`     | ✖        | string     | -                         | `deposit_period`, `voting_period`, `passed`, `rejected` | (主要参数)按提案的状态筛选                          |
| -        | `--voter`      | ✖        | string     | -                         | -                                                       | (主要参数)按投票表决的提案筛选                        |


# Global Flags

| ShortCut | Flag         | Required | Input Type | Default Input | Input Range       | Description  |
|:---------|:-------------|:---------|:-----------|:--------------|:------------------|:-------------|
| `-e`     | `--encoding` | ✖        | string     | `hex`         | `hex`/`b64`/`btc` | 二进制编码        |
| -        | `--home`     | ✖        | string     | `/.qoscli`    | -                 | 配置和数据的目录     |
| `-o`     | `--output`   | ✖        | string     | `text`        | `text`/`json`     | 输出格式         |
| -        | `--trace`    | ✖        | -          | -             | -                 | 打印出错时的完整堆栈跟踪 |