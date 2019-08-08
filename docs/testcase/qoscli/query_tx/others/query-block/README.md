# Test Cases

- [缺失参数[height]](./TestCase01.md)
- [指定参数[height]](./TestCase02.md)

# Description
>     Get block info at given height.

>     获取给定高度的区块信息。

# Example

`qoscli query block <height>`

其中`<height>`为区块高度

查询高度10区块信息：
```bash
$ qoscli query block 10 --indent
```

输出示例：
```bash
{
  "block_meta": {
    "block_id": {
      "hash": "A473CE3866A74277BC7F7B7AF70E55B40736B8A3CA3B8A55406AC8CF6E04ED50",
      "parts": {
        "total": "1",
        "hash": "B9C5DEF42EAA9D445E52B1F8DD34ECC96C02E537F43D1F7C8D829C84F8663127"
      }
    },
    "header": {
      "version": {
        "block": "10",
        "app": "0"
      },
      "chain_id": "Arya",
      "height": "20",
      "time": "2019-04-25T06:19:28.353298129Z",
      "num_txs": "0",
      "total_txs": "0",
      "last_block_id": {
        "hash": "BC153175007D7E5D5C6A27D22E3F7227224E43C537988DDCBF6C2F14A95DD432",
        "parts": {
          "total": "1",
          "hash": "EEFE6F3A761D9D28DBCA81424F9E50A8C716D0F4898FA7B3893CBB0AC7B55F4D"
        }
      },
      "last_commit_hash": "863F17ACB4909A5E043782DB06F3FE18C6DCF4988EE9B7C0CDA7D8337504FCFB",
      "data_hash": "",
      "validators_hash": "5CA1D1B7D703F2D2A9C270D1CD5819E7E0D439BA1C55645BCD8DB7B079389CA8",
      "next_validators_hash": "5CA1D1B7D703F2D2A9C270D1CD5819E7E0D439BA1C55645BCD8DB7B079389CA8",
      "consensus_hash": "294D8FBD0B94B767A7EBA9840F299A3586DA7FE6B5DEAD3B7EECBA193C400F93",
      "app_hash": "C31662F65DEE545FEDF15D98517CBF07034DC1821EF06DD87D2F956C315A0A9B",
      "last_results_hash": "",
      "evidence_hash": "",
      "proposer_address": "0E447E66089C9D97EFC2F4C172403F35740DD507"
    }
  },
  "block": {
    "header": {
      "version": {
        "block": "10",
        "app": "0"
      },
      "chain_id": "Arya",
      "height": "20",
      "time": "2019-04-25T06:19:28.353298129Z",
      "num_txs": "0",
      "total_txs": "0",
      "last_block_id": {
        "hash": "BC153175007D7E5D5C6A27D22E3F7227224E43C537988DDCBF6C2F14A95DD432",
        "parts": {
          "total": "1",
          "hash": "EEFE6F3A761D9D28DBCA81424F9E50A8C716D0F4898FA7B3893CBB0AC7B55F4D"
        }
      },
      "last_commit_hash": "863F17ACB4909A5E043782DB06F3FE18C6DCF4988EE9B7C0CDA7D8337504FCFB",
      "data_hash": "",
      "validators_hash": "5CA1D1B7D703F2D2A9C270D1CD5819E7E0D439BA1C55645BCD8DB7B079389CA8",
      "next_validators_hash": "5CA1D1B7D703F2D2A9C270D1CD5819E7E0D439BA1C55645BCD8DB7B079389CA8",
      "consensus_hash": "294D8FBD0B94B767A7EBA9840F299A3586DA7FE6B5DEAD3B7EECBA193C400F93",
      "app_hash": "C31662F65DEE545FEDF15D98517CBF07034DC1821EF06DD87D2F956C315A0A9B",
      "last_results_hash": "",
      "evidence_hash": "",
      "proposer_address": "0E447E66089C9D97EFC2F4C172403F35740DD507"
    },
    "data": {
      "txs": null
    },
    "evidence": {
      "evidence": null
    },
    "last_commit": {
      "block_id": {
        "hash": "BC153175007D7E5D5C6A27D22E3F7227224E43C537988DDCBF6C2F14A95DD432",
        "parts": {
          "total": "1",
          "hash": "EEFE6F3A761D9D28DBCA81424F9E50A8C716D0F4898FA7B3893CBB0AC7B55F4D"
        }
      },
      "precommits": [
        {
          "type": 2,
          "height": "19",
          "round": "0",
          "block_id": {
            "hash": "BC153175007D7E5D5C6A27D22E3F7227224E43C537988DDCBF6C2F14A95DD432",
            "parts": {
              "total": "1",
              "hash": "EEFE6F3A761D9D28DBCA81424F9E50A8C716D0F4898FA7B3893CBB0AC7B55F4D"
            }
          },
          "timestamp": "2019-04-25T06:19:28.353298129Z",
          "validator_address": "0E447E66089C9D97EFC2F4C172403F35740DD507",
          "validator_index": "0",
          "signature": "bfhVFCZMS/6hEmkFAaLfNwumKEUQNtRkGvnrMTTvezjpCbv/X0wSQQKq6g4crd5mI3WjZYp4vM+EA4SY55ucCw=="
        },
        {
          "type": 2,
          "height": "19",
          "round": "0",
          "block_id": {
            "hash": "BC153175007D7E5D5C6A27D22E3F7227224E43C537988DDCBF6C2F14A95DD432",
            "parts": {
              "total": "1",
              "hash": "EEFE6F3A761D9D28DBCA81424F9E50A8C716D0F4898FA7B3893CBB0AC7B55F4D"
            }
          },
          "timestamp": "2019-04-25T06:19:28.312339528Z",
          "validator_address": "E9816412631B42AE3921769FFD9DE121AA745422",
          "validator_index": "1",
          "signature": "vePZhdo+dRTEghf3aHhqWXJQgyXeSoB2q4o1WiIncxI1raXU5YTGKNEdD8Tq8TbmI2uDH5J6CAOGy9ru1DzODQ=="
        }
      ]
    }
  }
}
```

# Usage
```
  qoscli query block [height] [flags]
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