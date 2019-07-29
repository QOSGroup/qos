# Test Cases

- 暂无

# Description
>     Query details of a signal proposal.

>     查询单个提案(proposal)的详细信息。

# Example

`qoscli query proposal <proposal-id>`

查询`ProposalID`为1的提议：
```bash
$ qoscli query proposal 1 --indent
```

查询结果：
```bash
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
```

# Usage
```
  qoscli query proposal [id] [flags]
```

`[id]`为提案的`ProposalID`

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

# Global Flags

| ShortCut | Flag         | Required | Input Type | Default Input | Input Range       | Description  |
|:---------|:-------------|:---------|:-----------|:--------------|:------------------|:-------------|
| `-e`     | `--encoding` | ✖        | string     | `hex`         | `hex`/`b64`/`btc` | 二进制编码        |
| -        | `--home`     | ✖        | string     | `/.qoscli`    | -                 | 配置和数据的目录     |
| `-o`     | `--output`   | ✖        | string     | `text`        | `text`/`json`     | 输出格式         |
| -        | `--trace`    | ✖        | -          | -             | -                 | 打印出错时的完整堆栈跟踪 |