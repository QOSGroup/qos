# Test Cases

- [直接调用](./TestCase01.md)

# Description
>     Query remote node for status.

>     查询远程节点的状态。

# Example

`qoscli query status --indent`

输出示例：
```bash
{
  "node_info": {
    "protocol_version": {
      "p2p": "7",
      "block": "10",
      "app": "0"
    },
    "id": "4537e18828364c6e3529000e30bcf9f25b0fc50c",
    "listen_addr": "tcp://0.0.0.0:26656",
    "network": "imuge",
    "version": "0.30.1",
    "channels": "4020212223303800",
    "moniker": "node1",
    "other": {
      "tx_index": "on",
      "rpc_address": "tcp://0.0.0.0:26657"
    }
  },
  "sync_info": {
    "latest_block_hash": "4D935B625A5C2D63FD251C8448C9765916B289E435A0388F64401767DFA22BD5",
    "latest_app_hash": "29E08C36CE8CEA35EF4DE04B002C852505361B303950F3E07EBFC031F8DAB854",
    "latest_block_height": "396",
    "latest_block_time": "2019-04-25T06:53:11.777203643Z",
    "catching_up": false
  },
  "validator_info": {
    "address": "0E447E66089C9D97EFC2F4C172403F35740DD507",
    "pub_key": {
      "type": "tendermint/PubKeyEd25519",
      "value": "FIGPykhLqi5X5HYrFMiI7hus7x2rNVg18pPevIBRLoU="
    },
    "voting_power": "26590937"
  }
}
```

其中`catching_up`为`false`表示节点已同步到最新高度。

# Usage
```
  qoscli query status [flags]
```

# Available Commands

>无可用命令

# Flags

| ShortCut | Flag           | Required | Input Type | Default Input             | Input Range | Description                             |
|:---------|:---------------|:---------|:-----------|:--------------------------|:------------|:----------------------------------------|
| `-h`     | `--help`       | ✖        | -          | -                         | -           | 帮助文档                                    |
| -        | `--indent`     | ✖        | -          | -                         | -           | 向JSON响应添加缩进                             |
| `-n`     | `--node`       | ✖        | string     | `"tcp://localhost:26657"` | -           | 为此链提供的Tendermint RPC接口: `<host>:<port>` |

# Global Flags

| ShortCut | Flag         | Required | Input Type | Default Input | Input Range       | Description  |
|:---------|:-------------|:---------|:-----------|:--------------|:------------------|:-------------|
| `-e`     | `--encoding` | ✖        | string     | `hex`         | `hex`/`b64`/`btc` | 二进制编码        |
| -        | `--home`     | ✖        | string     | `/.qoscli`    | -                 | 配置和数据的目录     |
| `-o`     | `--output`   | ✖        | string     | `text`        | `text`/`json`     | 输出格式         |
| -        | `--trace`    | ✖        | -          | -             | -                 | 打印出错时的完整堆栈跟踪 |