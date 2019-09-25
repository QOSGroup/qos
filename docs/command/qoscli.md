# QOS Client

提供与QOS网络交互的命令行工具`qoscli`，主要提供以下命令行功能：
* `keys`        [本地密钥库](#密钥（keys）)
* `query`       [信息查询](#查询（query）)
* `tx`          [交易](#交易（tx）)
* `version`     [版本信息](#版本（version）)

所有命令均可通过添加`--help`获取命令说明

命令中涉及的通用参数：

| 参数 | 默认值 | 说明 |
| :--- | :---: | :--- |
|--nonce | 0 | account nonce to sign the tx |
|--max-gas| 0 | gas limit to set per tx |
|--chain-id| "" | Chain ID of tendermint node |
|--node| tcp://localhost:26657 | tcp://\<host\>:\<port\> to tendermint rpc interface for this chain |
|--height| 0 | block height to query, omit to get most recent provable block |
|--async| false | broadcast transactions asynchronously |
|--trust-node| false | Trust connected full node |
|--qcp| false | enable qcp mode. send qcp tx |
|--qcp-signer| "" | qcp mode flag. qcp tx signer key name |
|--qcp-seq| 0 | qcp mode flag.  qcp in sequence |
|--qcp-from| "" | qcp mode flag. qcp tx source |
|--qcp-blockheight| 0 | qcp mode flag. original tx blockheight |
|--qcp-txindex| 0 | qcp mode flag. original tx index |
|--qcp-extends| "" | qcp mode flag. qcp tx extends info |
|--indent| false | add indent to json response |
|--nonce-node| "" | tcp://\<host\>:\<port\> to tendermint rpc interface for some chain to query account nonce |

更多说明参照[qbase-通用参数](https://github.com/QOSGroup/qbase/blob/master/docs/client/command.md#客户端命令)

## 密钥（keys）

本地密钥库主要包含以下指令：
* `qoscli keys add`     [新增密钥](#新增（add）)
* `qoscli keys list`    [显示密钥列表](#列表（list）)
* `qoscli keys update`  [更新密钥保存密码](#更新（update）)
* `qoscli keys delete`  [从密钥库删除密钥](#删除（delete）)
* `qoscli keys import`  [导入密钥](#导入（import）)
* `qoscli keys export`  [导出密钥](#导出（export）)

> 密钥库为本地存储，默认存储位置为：$HOME/.qoscli/keys/，删除存储文件会清空本地存储所有私钥。通过`keys`相关指令操作密钥不影响QOS网络中账户状态，请妥善保管账户私钥信息。

### 新增（add）

`qoscli keys add <key_name>`

<key_name>可随意填写，仅作为本地密钥库密钥区分。

如下指令将生成一个名字为`Arya`的密钥到本地密钥库：
```bash
$ qoscli keys add Arya
Enter a passphrase for your key:<输入不少于8位的密码>
Repeat the passphrase:<重复上面输入的密码>
NAME:	TYPE:	ADDRESS:						PUBKEY:
Arya	local	qosacc10327kf8v45a7uhev92llmuqwzkfgecvwckxt5m	qosaccpub1zcjduepqhn2n540cn0ts0qg7zd8xyrrwjg54lvaka228c3vs8gf5ph3eh27sy7nlzh
**Important** write this seed phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

thought frame must space few omit muffin fix merge mail ivory clump unveil dirt gadget load glove hub inner final crime churn crop stone
```
其中`qosacc10327kf8v45a7uhev92llmuqwzkfgecvwckxt5m`为适用于QOS网络的账户地址，`qosaccpub1zcjduepqhn2n540cn0ts0qg7zd8xyrrwjg54lvaka228c3vs8gf5ph3eh27sy7nlzh`为账户公钥信息，`thought frame must space few omit muffin fix merge mail ivory clump unveil dirt gadget load glove hub inner final crime churn crop stone`为助记词，可用于账户私钥找回，请妥善保管助记词。

### 列表（list）

`qoscli keys list`
```bash
$ qoscli keys list
NAME:	TYPE:	ADDRESS:						PUBKEY:
Arya	local	qosacc10327kf8v45a7uhev92llmuqwzkfgecvwckxt5m	qosaccpub1zcjduepqhn2n540cn0ts0qg7zd8xyrrwjg54lvaka228c3vs8gf5ph3eh27sy7nlzh
```

### 更新（update）

`qoscli keys update <key_name>`

更新`Arya`存储密码：
```bash
$ qoscli keys update Arya
Enter the current passphrase:<输入当前密码>
Enter the new passphrase:<输入新密码>
Repeat the new passphrase:<重复新密码>
Password successfully updated!
```

### 导出（export）

`qoscli keys export <key_name>`

导出`Arya`密钥信息：
```bash
qoscli keys export Arya
Password to sign with 'Arya':<输入>
**Important** Don't leak your private key information to others.
Please keep your private key safely, otherwise your account will be attacked.

{"Name":"Arya","address":"qosacc10327kf8v45a7uhev92llmuqwzkfgecvwckxt5m","pubkey":"qosaccpub1zcjduepqhn2n540cn0ts0qg7zd8xyrrwjg54lvaka228c3vs8gf5ph3eh27sy7nlzh","privkey":{"type":"tendermint/PrivKeyEd25519","value":"n/eCiOFjYFf22NCsacMfTXhxI6dV3DfI8cuxlZ48M0S81TpV+JvXB4EeE05iDG6SKV+ztuqUfEWQOhNA3jm6vQ=="}}
```
导出的密钥是通过JSON序列化后的密钥信息，可以将JSON字符串中的`privkey`部分内容保存为文件并妥善保存，可用于密钥导入。

### 删除（delete）

`qoscli keys delete <key_name>`

删除`Arya`密钥信息：
```bash
$ qoscli keys delete Arya
DANGER - enter password to permanently delete key:<输入密码>
key deleted forever (uh oh!)
```

### 导入（import）

`qoscli keys import Arya --file <私钥文件路径>`

导入上面通过`export`导出的私钥文件：
```bash
qoscli keys import Arya --file Arya.pri
> Enter a passphrase for your key:<输入不少于8位的密码>
> Repeat the passphrase:<重复上面输入的密码>

其中Arya.pri文件内容为:
{"type":"tendermint/PrivKeyEd25519","value":"n/eCiOFjYFf22NCsacMfTXhxI6dV3DfI8cuxlZ48M0S81TpV+JvXB4EeE05iDG6SKV+ztuqUfEWQOhNA3jm6vQ=="}

```

## 版本（version）
`qoscli version`

输出示例：
```bash
{
 "version": "0.0.4-46-g5ec63bd", //QOS版本信息
 "commit": "5ec63bd74c2c92924c25ffd5be1ff0f232bfcda4", //QOS源码commit ID
 "go": "go version go1.11.5 linux/amd64" //go 版本信息
}

```

## 查询（query）

* `qoscli query account`                [账户查询](#账户)
* `qoscli query store`                  [存储查询](#存储（store）)
* `qoscli query consensus`              共识参数查询
* `qoscli query approve`                [预授权](#查询预授权)
* `qoscli query qcp`                    [跨链相关信息查询](#查询联盟链)
* `qoscli query qsc`                    [代币查询](#查询代币)
* `qoscli query qscs`                   [所有代币查询](#查询所有代币)
* `qoscli query validators`             [验证节点列表](#验证节点列表)
* `qoscli query validator`              [验证节点查询](#查询验证节点)
* `qoscli query validator-miss-vote`    [验证节点漏块信息](#查询验证节点漏块信息)
* `qoscli query validator-period`       [验证节点窗口信息](#验证节点窗口信息)
* `qoscli query community-fee-pool`     [社区收益池](#社区收益池)
* `qoscli query delegation`             [委托查询](#委托查询)
* `qoscli query delegations-to`         [验证节点委托列表](#验证节点委托列表)
* `qoscli query delegations`            [代理用户委托列表](#代理用户委托列表)
* `qoscli query delegator-income`       [委托收益查询](#委托收益查询)
* `qoscli query unbondings`             [待返还委托](#待返还委托)
* `qoscli query redelegations`          [待执行委托变更](#待执行委托变更)
* `qoscli query community-fee-pool`     [社区费池](#社区费池)
* `qoscli query proposal`               [提议查询](#提议查询)
* `qoscli query proposals`              [提议列表](#提议列表)
* `qoscli query vote`                   [投票查询](#投票查询)
* `qoscli query votes`                  [投票列表](#投票列表)
* `qoscli query deposit`                [抵押查询](#抵押查询)
* `qoscli query deposits`               [抵押列表](#抵押列表)
* `qoscli query tally`                  [投票统计](#投票统计)
* `qoscli query params`                 [参数查询](#参数查询)
* `qoscli query inflation-phrases`      [通胀规则查询](#通胀规则查询)
* `qoscli query total-inflation  `      [发行总量查询](#发行总量查询)
* `qoscli query total-applied`          [流通总量查询](#流通总量查询)
* `qoscli query guardian`               [系统账户查询](#系统账户查询)
* `qoscli query guardians`              [系统账户列表](#系统账户列表)
* `qoscli query status`                 [查询节点状态](#状态（status）)
* `qoscli query tendermint-validators`  [获取指定高度验证节点集合](#获取指定高度验证节点集合)
* `qoscli query block`                  [获取指定高度区块信息](#区块（block）)
* `qoscli query txs`                    [根据标签查找交易](#根据标签查找交易)
* `qoscli query tx`                     [根据交易hash查询交易信息](#根据交易hash查询交易信息)

查询的具体指令将在各自模块进行介绍。

### 状态（status）
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
    "address": "qoscons1glmfxaskefe75gat4se48fw0rdftn8h5dveeyw",
    "pub_key": "qosconspub1zcjduepquqyj75pn6auj725dcjmdyrd7ndlx6mnzhncelwju9lj9j7k0u6vq00rs3r",
    "voting_power": "10000000"
  }

}
```

其中`catching_up`为`false`表示节点已同步到最新高度。

### 区块（block）
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

### 账户

查询账户
`qoscli query account <key_name_or_account_address>`

<key_name_or_account_address>为本地密钥库存储的密钥名字或对应账户的地址。

假设本地密钥库中`Arya`地址为`qosacc1x5lcfaqxxq7g7dy4lj5vq0u6xamp78lsnza98y`，且QOS网络中已经创建了`qosacc1x5lcfaqxxq7g7dy4lj5vq0u6xamp78lsnza98y`对应账号，可执行：
```bash
qoscli query account Arya --indent
```
或
```bash
qoscli query account qosacc1x5lcfaqxxq7g7dy4lj5vq0u6xamp78lsnza98y --indent
```
输出类似如下信息：
```bash
{
  "type": "qos/types/QOSAccount",
  "value": {
    "account_address": "qosacc1smrus8jlc9z02gz5rm36u0q3fdctjxm4nrc639",
    "public_key": "qosaccpub1zcjduepqhewtk4rn8050sgyws8unnvr9l9rlshvfyargrds7qr9m40crphkscksjqd",
    "nonce": 1,
    "qos": "10000",
    "qscs": [
        {
            "coin_name": "AOE",
            "amount": "10000"
        }
    ]
  }
}
```
可以看到`Arya`持有10000个QOS、10000个AOE，更多账户说明请阅读[QOS账户设计](../spec/account.md)文档。

### 存储（store）

QOS网络的存储内容均可通过下面指令查找：

`qoscli query store --path /store/<store_key>/subspace --data <query_data>`

主要参数：

- `--path`  存储位置
- `--data`  查询内容，以<query_data>开头的数据会被查出来

查询QOS网络中存储的ROOT CA 信息：

```bash
$ qoscli query store --path /store/acc/subspace --data account --indent
```

执行结果：

```bash
[
  {
    "key": "account:\ufffd\ufffd\ufffd\u001e_\ufffdD\ufffd T\u001e\ufffd\ufffd\u003c\u0011Kp\ufffd\u001bu",
    "value": {
      "type": "qos/types/QOSAccount",
      "value": {
        "account_address": "qosacc1smrus8jlc9z02gz5rm36u0q3fdctjxm4nrc639",
        "public_key": "qosaccpub1zcjduepqhewtk4rn8050sgyws8unnvr9l9rlshvfyargrds7qr9m40crphkscksjqd",
        "nonce": 1,
        "qos": "10000",
        "qscs": [
            {
                "coin_name": "AOE",
                "amount": "10000"
            }
        ]
      }
    }
  }
]

```

### 交易（tx query）
支持的查询命令：

* `qoscli query tx`            [根据交易hash查询交易信息](#根据交易hash查询交易信息)
* `qoscli query txs`           [根据标签查找交易](#根据标签查找交易)

#### 根据交易hash查询交易信息
执行交易后会返回交易hash，通过交易hash可查询交易详细信息。

根据hash `f5fc2c228cba754d5b95e49b02e81ff818f7b9140f1859d3797b09fb4aa12385` 查询交易信息：

```bash
$ qoscli query tx f5fc2c228cba754d5b95e49b02e81ff818f7b9140f1859d3797b09fb4aa12385 --indent
```
输出示例：

```bash
{
  "height": "153",
  "txhash": "8A317C18448BE7F2B76E05A930FCE0BA45B9C688A630CD50D886A05A7FB1A673",
  "gas_wanted": "9223372036854775807",
  "gas_used": "19500",
  "events": [
    {
      "type": "message",
      "attributes": [
        {
          "key": "module",
          "value": "transfer"
        },
        {
          "key": "gas.payer",
          "value": "qosacc1x5lcfaqxxq7g7dy4lj5vq0u6xamp78lsnza98y"
        }
      ]
    },
    {
      "type": "receive",
      "attributes": [
        {
          "key": "address",
          "value": "qosacc1smrus8jlc9z02gz5rm36u0q3fdctjxm4nrc639"
        },
        {
          "key": "qos",
          "value": "100"
        },
        {
          "key": "qscs"
        }
      ]
    },
    {
      "type": "send",
      "attributes": [
        {
          "key": "address",
          "value": "qosacc1x5lcfaqxxq7g7dy4lj5vq0u6xamp78lsnza98y"
        },
        {
          "key": "qos",
          "value": "100"
        },
        {
          "key": "qscs"
        }
      ]
    }
  ],
  "tx": {
    "type": "qbase/txs/stdtx",
    "value": {
      "itx": [
        {
          "type": "transfer/txs/TxTransfer",
          "value": {
            "senders": [
              {
                "addr": "qosacc1x5lcfaqxxq7g7dy4lj5vq0u6xamp78lsnza98y",
                "qos": "100",
                "qscs": null
              }
            ],
            "receivers": [
              {
                "addr": "qosacc1smrus8jlc9z02gz5rm36u0q3fdctjxm4nrc639",
                "qos": "100",
                "qscs": null
              }
            ]
          }
        }
      ],
      "sigature": [
        {
          "pubkey": {
            "type": "tendermint/PubKeyEd25519",
            "value": "2WRkEzr8Yd32reYXAmZURHwFm8CAWUTFAQc/IccnATg="
          },
          "signature": "tkyobbxwY5sjCTb1ibMEuMpOXHbDhKKHkZq0INZUvnyHODQloh4msRMmc62LLdswL1aY3vdld3OK+1pg2J86Dw==",
          "nonce": "1"
        }
      ],
      "chainid": "qos-test",
      "maxgas": "9223372036854775807"
    }
  },
  "timestamp": "2019-08-28T07:15:52Z"
}

```

#### 根据标签查找交易
执行交易后会同时会返回QOS为交易所打tag，通过交易tag可查询交易信息。

根据`approve-from`=`qosacc1x5lcfaqxxq7g7dy4lj5vq0u6xamp78lsnza98y`查询预授权交易信息：

```bash
$ qoscli query txs --tags "create-approve.approve-from='qosacc1x5lcfaqxxq7g7dy4lj5vq0u6xamp78lsnza98y'" --indent
```
输出示例：

```bash
[
  {
    "hash": "f5fc2c228cba754d5b95e49b02e81ff818f7b9140f1859d3797b09fb4aa12385",
    "height": "246",
    "tx": {
      "type": "qbase/txs/stdtx",
      "value": {
        "itx": [
          {
            "type": "approve/txs/TxCreateApprove",
            "value": {
              "Approve": {
                "from": "address1s348wvf49dfy64e6wafc90lcavp4lrd6xzhzhk",
                "to": "address1yqekgyy66v2cxzww6lqg6sdrsugjguxqws6mkf",
                "qos": "100",
                "qscs": null
              }
            }
          }
        ],
        "sigature": [
          {
            "pubkey": {
              "type": "tendermint/PubKeyEd25519",
              "value": "B/iatjhcJ4yFyHfGYKw2IneYGu2zG+ZOR8XmRUaji0A="
            },
            "signature": "VrsOsULJx86y8ch529zvl3Sh19TwGm/AldPlQhVWqhtg+calZmBrk25sD9HxCYijAt+ZUWMiLtPg3QZzCCqHAg==",
            "nonce": "1"
          }
        ],
        "chainid": "QOS",
        "maxgas": "100000"
      }
    },
    "result": {
      "gasWanted": "100000",
      "gasUsed": "15220",
      "tags": [
        {
          "key": "YWN0aW9u",
          "value": "Y3JlYXRlLWFwcHJvdmU="
        },
        {
          "key": "YXBwcm92ZS1mcm9t",
          "value": "YWRkcmVzczFzMzQ4d3ZmNDlkZnk2NGU2d2FmYzkwbGNhdnA0bHJkNnh6aHpoaw=="
        },
        {
          "key": "YXBwcm92ZS10bw==",
          "value": "YWRkcmVzczF5cWVrZ3l5NjZ2MmN4end3NmxxZzZzZHJzdWdqZ3V4cXdzNm1rZg=="
        }
      ]
    }
  }
]

```

更多交易Tag请查阅[index](../spec/indexing.md)

## 交易（tx）

QOS支持以下几种交易类型：

* `qoscli tx transfer`         [转账](#转账)
* `qoscli tx invariant-check`  [数据检查](#数据检查)
* `qoscli tx create-approve`   [创建预授权](#创建预授权)
* `qoscli tx increase-approve` [增加预授权](#增加预授权)
* `qoscli tx decrease-approve` [减少预授权](#减少预授权)
* `qoscli tx use-approve`      [使用预授权](#使用预授权)
* `qoscli tx cancel-approve`   [取消预授权](#取消预授权)
* `qoscli tx create-qsc`       [创建代币](#创建代币)
* `qoscli tx issue-qsc`        [发放代币](#发放代币)
* `qoscli tx init-qcp`         [初始化联盟链](#初始化联盟链)
* `qoscli tx create-validator` [成为验证节点](#成为验证节点)
* `qoscli tx create-validator` [编辑验证节点](#编辑验证节点)
* `qoscli tx revoke-validator` [撤销验证节点](#撤销验证节点)
* `qoscli tx active-validator` [激活验证节点](#激活验证节点)
* `qoscli tx delegate`         [委托](#委托)
* `qoscli tx modify-compound`  [修改收益复投方式](#修改收益复投方式)
* `qoscli tx unbond`           [解除委托](#解除委托)
* `qoscli tx redelegate`       [变更委托验证节点](#变更委托验证节点)
* `qoscli tx submit-proposal`  [提交提议](#提交提议)
* `qoscli tx deposit`          [提议抵押](#提议抵押)
* `qoscli tx vote`             [提议投票](#提议投票)
* `qoscli tx add-guardian`     [添加系统账户](#添加系统账户)
* `qoscli tx delete-guardian`  [删除系统账户](#删除系统账户)
* `qoscli tx halt-network`     [停止网络](#停止网络)

主要分为[Bank](#bank)，[预授权](#预授权)，[代币](#代币)，[联盟链](#联盟链)，[验证节点](#验证节点（validator）)，[治理](#治理)，[系统账户](#系统账户)这几大类。

### Bank

See [Bank模块](../spec/bank) to learn about bank module design.

* `qoscli tx transfer`          [转账](#转账)
* `qoscli invariant-check`      [数据检查](#数据检查)

#### 转账

查阅[转账设计](../spec/bank)了解QOS转账交易设计。

`qoscli tx transfer --senders <senders_and_coins> --receivers <receivers_and_coins>`

支持一次转账中包含多币种，多账户

主要参数：
- `--senders`   发送集合，账户传keystore name 或 address，多个账户半角分号分隔
- `--receivers` 接收集合，账户传keystore name 或 address，多个账户半角分号分隔

`Arya`向地址`qosacc1smrus8jlc9z02gz5rm36u0q3fdctjxm4nrc639`转账1个QOS，1个AOE
```bash
$ qoscli tx transfer --senders Arya,1QOS,1AOE --receivers qosacc1smrus8jlc9z02gz5rm36u0q3fdctjxm4nrc639,1QOS,1AOE
Password to sign with 'Arya':<输入密码>
{"check_tx":{},"deliver_tx":{},"hash":"21ECB72C8F51B3BD8E3CB9D59765003B9D78BE75","height":"300"}
```

转账成功可通过[账户查询](#账户)查看最新账户状态，交易执行可能会有一定时间的延迟。

#### 数据检查

QOS设计了一套[数据检查机制](../spec/bank)，用户可以通过下面的指令执行数据检查操作：

`qoscli tx invariant-check --sender <sender's keybase name or address>`

主要参数：
- `--sender` 发送此交易的账户keystore name 或 address

::: warning Note 
此交易设置了特别大的交易费，仅限持币账户发现QOS网络数据异常时提交，数据验证异常会停止整个QOS网络，以保护持币账户权益。
:::

`Arya`发现QOS网络中某处数值溢出，发现数据检查：
```bash
$ qoscli tx invariant-check --sender Arya
Password to sign with 'Arya':<输入密码>
```
如果数据并无异常，将返回正常交易执行结果，否则全网停止运行。

### 预授权

[QOS预授权设计](../spec/approve.md)包含以下操作指令：

* `qoscli tx create-approve`    [创建预授权](#创建预授权)
* `qoscli query approve`        [查询预授权](#查询预授权)
* `qoscli tx increase-approve`  [增加预授权](#增加预授权)
* `qoscli tx decrease-approve`  [减少预授权](#减少预授权)
* `qoscli tx use-approve`       [使用预授权](#使用预授权)
* `qoscli tx cancel-approve`    [取消预授权](#取消预授权)

> 下面实例中假设`Sansa`地址为`qosacc1smrus8jlc9z02gz5rm36u0q3fdctjxm4nrc639`

#### 创建预授权

`qoscli tx create-approve --from <key_name_or_account_address> --to <account_address> --coins <qos_and_qscs>`

主要参数：

- `--from`  授权账户本地密钥库名字或账户地址
- `--to`    被授权账户地址
- `--coins` 授权币种、币值列表，[amount1][coin1],[amount2][coin2],...，以半角逗号相隔

`Arya`向`Sansa`授权100个QOS，100个AOE：
```
$ qoscli tx create-approve --from Arya --to qosacc1smrus8jlc9z02gz5rm36u0q3fdctjxm4nrc639 --coins 100QOS,100AOE
Password to sign with 'Arya':<输入Arya本地密钥库密码>
```
执行结果：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"9917953D8CDE80F457CD072DBCE73A36449B7A7C","height":"333"}
```

#### 查询预授权

`qoscli query approve --from <key_name_or_account_address> --to <account_address>`

查询`Arya`对`Sansa`的预授权：
```bash
qoscli query approve --from Arya --to qosacc1smrus8jlc9z02gz5rm36u0q3fdctjxm4nrc639
```
执行结果：
```bash
{
  "from": "qosacc1x5lcfaqxxq7g7dy4lj5vq0u6xamp78lsnza98y",
  "to": "qosacc1smrus8jlc9z02gz5rm36u0q3fdctjxm4nrc639",
  "qos": "100",
  "qscs": [
    {
      "coin_name": "AOE",
      "amount": "100"
    }
  ]
}
```

#### 增加预授权

`qoscli tx increase-approve --from <key_name_or_account_address> --to <account_address> --coins <qos_and_qscs>`

`Arya`向`Sansa`增加授权100个QOS，100个AOE：
```bash
$ qoscli tx increase-approve --from Arya --to qosacc1smrus8jlc9z02gz5rm36u0q3fdctjxm4nrc639 --coins 100QOS,100AOE
Password to sign with 'Arya':<输入Arya本地密钥库密码>
```

执行结果：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"3C06676C53A5439D39CB4D0FBA3213C44DC1BA8E","height":"406"}
```

#### 减少预授权

`qoscli tx decrease-approve --from <key_name_or_account_address> --to <account_address> --coins <qos_and_qscs>`

`Arya`向`Sansa`减少授权10个QOS，10个AOE：
```bash
$ qoscli tx decrease-approve --from Arya --to qosacc1smrus8jlc9z02gz5rm36u0q3fdctjxm4nrc639 --coins 10QOS,10AOE
Password to sign with 'Arya':<输入Arya本地密钥库密码>
```
执行结果：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"3C06676C53A5439D39CB4D0FBA3213C44DC1BA8E","height":"410"}
```

#### 使用预授权

`qoscli tx use-approve --from <account_address> --to <key_name_or_account_address> --coins <qos_and_qscs>`

`Sansa`使用`Arya`向自己预授权中的10个QOS，10个AOE：
```bash
$ qoscli tx use-approve --from Arya --to Sansa --coins 10QOS,10AOE
Password to sign with 'Sansa':<输入Sansa本地密钥库密码>
```
执行结果：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"3C06676C53A5439D39CB4D0FBA3213C44DC1BA8E","height":"430"}
```

可通过[账户查询](#账户（account）)查看`Arya`和`Sansa`最新账户状态

#### 取消预授权

`qoscli tx cancel-approve --from <account_address> --to <key_name_or_account_address>'

`Arya`取消对`Sansa`的授权：
```bash
$ qoscli tx cancel-approve --from Arya --to qosacc1smrus8jlc9z02gz5rm36u0q3fdctjxm4nrc639
Password to sign with 'Arya':<输入Arya本地密钥库密码>
```
执行结果：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"484"}
```

### 代币

> 创建联盟币前需要申请[CA](../spec/ca.md)，点击[QSC设计文档](../spec/qsc)了解更多。

联盟币相关指令：
* `qoscli tx create-qsc`    [创建代币](#创建代币)
* `qoscli query qsc`        [查询代币](#查询代币)
* `qoscli query qsc`        [查询所有代币](#查询所有代币)
* `qoscli tx issue-qsc`     [发放代币](#发放代币)

#### 创建代币

`qoscli tx create-qsc --creator <key_name_or_account_address> --qsc.crt <qsc.crt_file_path> --accounts <account_qsc_s>`

主要参数：

- `--creator`       创建账号
- `--qsc.crt`       证书位置
- `--accounts`      初始发放地址币值集合，[addr1],[amount];[addr2],[amount2],...，该参数可为空，即只创建联盟币

`Arya`在QOS网络中创建`QOE`，不含初始发放地址币值信息：
```bash
$ qoscli tx create-qsc --creator Arya --qsc.crt aoe.crt
Password to sign with 'Arya':<输入Arya本地密钥库密码>
```
> 假设`Arya`已在CA中心申请`aoe.crt`证书，`aoe.crt`中包含`banker`公钥，对应地址`qosacc1djtcyex03vluga35r8lattddkqt76s7f306xuq`，已经导入到本地私钥库中，名字为`ATM`，。

执行结果：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"200"}
```

#### 查询代币

`qoscli query qsc <qsc_name>`

`qsc_name`为联盟币名称

查询`AOE`信息：
```bash
$ qoscli query qsc QOE --indent
```
执行结果：
```bash
{
  "name": "AOE",
  "chain_id": "capricorn-1000",
  "extrate": "1:280.0000",
  "description": "",
  "banker": "qosacc1djtcyex03vluga35r8lattddkqt76s7f306xuq"
}
```

#### 查询所有代币

`qoscli query qscs`

查询所有代币信息：
```bash
$ qoscli query qscs --indent
```
执行结果：
```bash
[
    {
      "name": "AOE",
      "chain_id": "capricorn-1000",
      "extrate": "1:280.0000",
      "description": "",
      "banker": "qosacc1djtcyex03vluga35r8lattddkqt76s7f306xuq"
    }
]
```

#### 发放代币

针对使用包含`Banker`公钥创建的联盟币，可向`Banker`地址发放（增发）对应联盟币：

`qoscli tx issue-qsc --qsc-name <qsc_name> --banker <key_name_or_account_address> --amount <qsc_amount>`

主要参数：
- `--qsc-name`  联盟币名字
- `--banker`    Banker地址或私钥库中私钥名
- `--amount`    联盟币发放（增发）量

向联盟币AOE `Banker`中发放（增发）10000AOE：

```bash
$ qoscli tx issue-qsc --qsc-name AOE --banker ATM --amount 10000
Password to sign with 'ATM':<输入ATM本地密钥库密码>
```

执行结果：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"223"}
```

可通过[账户查询](#账户)查看`ATM`账户所持有AOE数量。

### 联盟链

QOS跨链协议QCP，支持跨链交易

> 创建联盟链前需要申请[CA](../spec/ca.md)，点击[联盟链设计文档](../spec/qcp.md)了解更多。

联盟链相关指令：
* `qoscli tx init-qcp`: [初始化联盟链](#初始化联盟链)
* `qoscli query qcp`:   [查询qcp信息](#查询联盟链)

#### 初始化联盟链

`qoscli tx init-qcp --creator <key_name_or_account_address> --qcp.crt <qcp.crt_file_path>`

主要参数：

- `--creator`       创建账号
- `--qcp.crt`       证书位置

> 假设`Arya`已在CA中心申请`qcp.crt`证书，`qcp.crt`中联盟链ID为`aoe-1000`

`Arya`在QOS网络中初始化联盟链信息：
```bash
$ qoscli tx init-qcp --creator Arya --qcp.crt qcp.crt
Password to sign with 'Arya':<输入Arya本地密钥库密码>
```

执行结果：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"243"}
```

#### 查询联盟链

跨链协议是[qbase](https://www.github.com/QOSGroup/qbase)提供支持，主要有以下四个查询指令：
- `qoscli query qcp list`
- `qoscli query qcp out`
- `qoscli query qcp in`
- `qoscli query qcp tx`

指令说明请参照[qbase-Qcp](https://github.com/QOSGroup/qbase/blob/master/docs/client/command.md#Qcp)。

### 验证节点（validator）

验证节点相关概念和机制请参阅[验证人详解](../spec/validators/all_about_validators.md)和[QOS经济模型](../spec/validators/eco_module.md)。验证节点包含以下子命令：

* `qoscli tx create-validator`          [成为验证节点](#成为验证节点)
* `qoscli query validator`              [查询验证节点](#查询验证节点)
* `qoscli query validators`             [验证节点列表](#验证节点列表)
* `qoscli query validator-miss-vote`    [验证节点漏块信息](#查询验证节点漏块信息)
* `qoscli query community-fee-pool`     [社区收益池](#社区收益池)
* `qoscli tx revoke-validator`          [撤消验证节点](#撤销验证节点)
* `qoscli tx active-validator`          [激活验证节点](#激活验证节点)

#### 成为验证节点

`qoscli tx create-validator --moniker <validator_name> --owner <key_name_or_account_address> --tokens <tokens>`

主要参数：

- `--owner`         操作者账户地址或密钥库中密钥名字
- `--moniker`       验证节点名字，`len(moniker) <= 300`
- `--nodeHome`      节点配置文件和数据所在目录，默认：`$HOME/.qosd`
- `--tokens`        绑定tokens，不能大于操作者持有QOS数量
- `--compound`      是否收益复投
- `--logo`          logo, 可选参数，`len(logo) <= 255`
- `--website`       网址, 可选参数，`len(website) <= 255`
- `--details`       详细描述信息, 可选参数，`len(details) <= 1000`

创建的validator基于本地的配置文件取`$HOME/.qosd/config/priv_validator.json`内信息，如果更改过默认位置，请使用`--home`指定`config`所在目录。

`Arya`初始化了一个[全节点](../install/testnet.md#启动全节点)，可通过下面指令成为验证节点：
```bash
$ qoscli tx create-validator --moniker "Arya's node" --owner Arya --tokens 1000
```

执行结果：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"258"}
```

执行成为验证节点命令后将从`Arya`账户扣除1000QOS，绑定到验证节点中，验证节点参与投票、打块所获得的挖矿收益将直接增加到`Arya`账户。

#### 编辑验证节点

`qoscli tx modify-validator --owner <key_name_or_account_address> --validator <validator_address> --moniker <validator_name> --logo <logo_url> --website <website_url> --details <description info>`

主要参数：

- `--owner`         操作者账户地址或密钥库中密钥名字
- `--validator`     待修改的验证人地址
- `--moniker`       验证节点名字，`len(moniker) <= 300`
- `--nodeHome`      节点配置文件和数据所在目录，默认：`$HOME/.qosd`
- `--compound`      是否收益复投
- `--logo`          logo, 可选参数，`len(logo) <= 255`
- `--website`       网址, 可选参数，`len(website) <= 255`
- `--details`       详细描述信息, 可选参数，`len(details) <= 1000`

`Arya`可通过`modify-validator`添加/修改节点信息：
```bash
$ qoscli tx modify-validator --moniker "Arya's node" --owner Arya --validaotor qosval1fzpaxwrmhqml7d90zuzvhmjfxsdqgvzrjpyvsl --logo "https://..." --website "https://..." --description "Long live Arya."
```

执行结果：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"265"}
```

执行成为验证节点命令后将从`Arya`账户扣除1000QOS，绑定到验证节点中，验证节点参与投票、打块所获得的挖矿收益将直接增加到`Arya`账户。

#### 查询验证节点

`qoscli query validator [validator-address]`

`validator-address`为验证人地址

可根据操作者查找与其绑定的验证节点信息。

```bash
$ qoscli query validator qosval12kjmpgyg23l7axhzzne33jmd0r9y083wzt07hu --indent
```

执行结果：
```bash
{
  "validator": "qosval12kjmpgyg23l7axhzzne33jmd0r9y083wzt07hu",
  "owner": "qosacc12kjmpgyg23l7axhzzne33jmd0r9y083w6mpa33",
  "consensusPubKey": "qosconspub1zcjduepq200x4crydcd6va90l2v2evz0ddjsv88vraqv7jqgds05e4n2xxfqqr6t45",
  "bondTokens": "10000000",
  "description": {
    "moniker": "Arya's node",
    "logo": "https://...",
    "website": "https://...",
    "details": "Long live Arya."
  },
  "commission": {
    "commission_rates": {
      "rate": "0.100000000000000000",
      "max_rate": "0.200000000000000000",
      "max_change_rate": "0.010000000000000000"
    },
    "update_time": "0001-01-01T00:00:00Z"
  },
  "status": "active",
  "InactiveDesc": "",
  "inactiveTime": "0001-01-01T00:00:00Z",
  "inactiveHeight": "0",
  "minPeriod": "0",
  "bondHeight": "0"
}

```

#### 验证节点列表

`qoscli query validators`

查询所有验证节点：
```bash
$ qoscli query validators --indent
```

执行结果：
```bash
validators:
[
  {
  "validator": "qosval12kjmpgyg23l7axhzzne33jmd0r9y083wzt07hu",
  "owner": "qosacc12kjmpgyg23l7axhzzne33jmd0r9y083w6mpa33",
  "consensusPubKey": "qosconspub1zcjduepq200x4crydcd6va90l2v2evz0ddjsv88vraqv7jqgds05e4n2xxfqqr6t45",
  "bondTokens": "10000000",
  "description": {
    "moniker": "Arya's node",
    "logo": "https://...",
    "website": "https://...",
    "details": "Long live Arya."
  },
  "commission": {
    "commission_rates": {
      "rate": "0.100000000000000000",
      "max_rate": "0.200000000000000000",
      "max_change_rate": "0.010000000000000000"
    },
    "update_time": "0001-01-01T00:00:00Z"
  },
  "status": "active",
  "InactiveDesc": "",
  "inactiveTime": "0001-01-01T00:00:00Z",
  "inactiveHeight": "0",
  "minPeriod": "0",
  "bondHeight": "0"
}
]
```

#### 获取指定高度验证节点集合

`qoscli query tendermint-validators <height>`

查询最新高度所有验证节点：
```bash
$ qoscli query tendermint-validators --indent
```

执行结果：
```bash
current query height: 100
[
  {
    "Address": "qoscons1stp35jzgkecv9qa2u0vh33gqclnrwu5wekgggz",
    "VotingPower": "10000000",
    "PubKey": "qosconspub1zcjduepq200x4crydcd6va90l2v2evz0ddjsv88vraqv7jqgds05e4n2xxfqqr6t45"
  }
]
```

#### 查询验证节点漏块信息

`qoscli query validator-miss-vote [validator-address]`

`validator-address`为验证人地址

查询`qosval1zlv7vhdcyqyvy9ljdxhcf2766nulwgvys3f6y2`的节点漏块信息：
```bash
$ qoscli query validator-miss-vote qosval1zlv7vhdcyqyvy9ljdxhcf2766nulwgvys3f6y2
```

执行结果：
```bash
{"startHeight":"258","endHeight":"387","missCount":0,"voteDetail":[]}
```

#### 验证节点窗口信息
`qoscli query validator-period  <validator-address>`

`validator-address`为验证人地址

查询`Arya`的节点漏块信息：
```bash
$ qoscli query validator-period qosval1zlv7vhdcyqyvy9ljdxhcf2766nulwgvys3f6y2
```

执行结果：
```bash
{
  "validator_address": "qosval1zlv7vhdcyqyvy9ljdxhcf2766nulwgvys3f6y2",
  "consensus_pubkey": "qosconspub1zcjduepqma0j9e6ky7h06vgkmghrj3g9u3dwlkhxc8y87x9rn956ue8dulls9g3ax0",
  "fees": "0",
  "current_tokens": "10222379",
  "current_period": "105",
  "last_period": "104",
  "last_period_fraction": {
    "value": "0.000000000000000000"
  }
}
```

#### 社区收益池
`qoscli query community-fee-pool`

查询社区收益：
```bash
$ qoscli query community-fee-pool
```

执行结果：
```bash
123456
```

#### 撤销验证节点

`qoscli tx revoke-validator --owner <key_name_or_account_address> --validator <validator_address>`

`key_name_or_account_address`为验证人所有者的账户地址或密钥库中密钥名字
`validator_address` 为验证人地址

`Arya`将自己的节点撤销为为验证节点：
```bash
$ qoscli tx revoke-validator --owner Arya --validator qosval1zlv7vhdcyqyvy9ljdxhcf2766nulwgvys3f6y2
```

执行结果：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"268"}
```

执行撤销命令后`Arya`的节点将处于pending状态，不再参与投票和打块。

#### 激活验证节点

`qoscli tx active-validator --owner <key_name_or_account_address> --validator <validator_address>`

`key_name_or_account_address`为操作者账户地址或密钥库中密钥名字
`validator_address` 为验证人地址

`Arya`将自己处于pending状态的节点重新激活为验证节点：
```bash
$ qoscli tx active-validator --owner Arya --validator qosval1zlv7vhdcyqyvy9ljdxhcf2766nulwgvys3f6y2
```

执行结果：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"275"}
```

执行成功，`Arya`的节点将继续参与投票、打块等共识职能，并获得挖矿奖励。



### 委托（delegate）

* `qoscli tx delegate`              [委托](#委托)
* `qoscli query delegation`         [委托查询](#委托查询)
* `qoscli query delegations-to`     [验证节点委托列表](#验证节点委托列表)
* `qoscli query delegations`        [代理用户委托列表](#代理用户委托列表)
* `qoscli query delegator-income`   [委托收益查询](#委托收益查询)
* `qoscli tx modify-compound`       [修改收益复投方式](#修改收益复投方式)
* `qoscli tx unbond`                [解除委托](#解除委托)
* `qoscli query unbondings`         [待返还委托](#待返还委托)
* `qoscli tx redelegate`            [变更委托验证节点](#变更委托验证节点)
* `qoscli query redelegations`      [待执行委托变更](#待执行委托变更)

#### 委托

`qoscli tx delegate --validator <validator_address> --delegator <delegator_key_name_or_account_address> --tokens <tokens> --compound <compound_or_not>`

主要参数：

- `--validator`     验证人地址
- `--delegator`     委托人账户地址或秘钥库中秘钥名字
- `--tokens`        绑定tokens，不能大于`delegator`持有QOS数量
- `--compound`      收益是否复投，默认`false`

`Sansa`将自己的100个QOS代理给`Arya`创建的验证节点：
```bash
$ qoscli tx delegate --owner Arya --delegator Sansa --tokens 100
```

#### 委托查询

`qoscli query delegation --validator <validator_address> --delegator <delegator_key_name_or_account_address>`

主要参数：

- `--validator`     验证人地址
- `--delegator`     委托人账户地址或秘钥库中秘钥名字

`Sansa`在`Arya`上的代理信息：
```bash
$ qoscli query delegation --owner Arya --delegator Sansa
```

查询结果：
```bash
{
  "delegator_address": "qosacc12tr0v5uv9xpns79w8q34plakz8gh66855arlrd",
  "validator_address": "qosval12tr0v5uv9xpns79w8q34plakz8gh6685vddu9q",
  "validator_cons_pub_key": "qosconspub1zcjduepqgvpv0ky5nzkt238aus9gh90ct96wm3nddmz7pt2qwyn3lwme8dpskrfu76",
  "delegate_amount": "10000000",
  "is_compound": false
}
```

#### 验证节点委托列表

`qoscli query delegations-to [validator-address]`

主要参数：

- `validator-address`     验证人地址

`Arya`验证节点上的所有代理信息：
```bash
$ qoscli query delegations-to qosval12tr0v5uv9xpns79w8q34plakz8gh6685vddu9q
```

查询结果示例：
```bash
[
  {
    "delegator_address": "qosacc12tr0v5uv9xpns79w8q34plakz8gh66855arlrd",
    "validator_address": "qosval12tr0v5uv9xpns79w8q34plakz8gh6685vddu9q",
    "validator_cons_pub_key": "qosconspub1zcjduepqgvpv0ky5nzkt238aus9gh90ct96wm3nddmz7pt2qwyn3lwme8dpskrfu76",
    "delegate_amount": "10000000",
    "is_compound": false
  }
...
]

```

#### 代理用户委托列表

`qoscli query delegations [delegator]`

主要参数：

- `delegator`     委托人账户地址或秘钥库中秘钥名字

`Sansa`的所有代理信息：
```bash
$ qoscli query delegations Sansa
```

查询结果：
```bash
[
  {
    "delegator_address": "qosacc12tr0v5uv9xpns79w8q34plakz8gh66855arlrd",
    "validator_address": "qosval12tr0v5uv9xpns79w8q34plakz8gh6685vddu9q",
    "validator_cons_pub_key": "qosconspub1zcjduepqgvpv0ky5nzkt238aus9gh90ct96wm3nddmz7pt2qwyn3lwme8dpskrfu76",
    "delegate_amount": "10000000",
    "is_compound": false
  }
]

```

#### 社区费池

`qoscli query community-fee-pool`

社区费池查询：
```bash
$ qoscli query community-fee-pool
```

查询结果：
```bash
"27211098"
```

#### 委托收益查询

`qoscli query delegator-income --validator <validator_address> --delegator <delegator_key_name_or_account_address`

主要参数：

- `--validator`  验证人地址
- `--delegator`  委托人账户地址或秘钥库中秘钥名字

`Sansa`查询代理给`qosval12tr0v5uv9xpns79w8q34plakz8gh6685vddu9q`的收益信息：
```bash
$ qoscli query delegator-income --delegator Sansa --validator qosval12tr0v5uv9xpns79w8q34plakz8gh6685vddu9q --indent

```

查询结果：
```bash
{
  "validator_address": "qosval12tr0v5uv9xpns79w8q34plakz8gh6685vddu9q",
  "consensus_pubkey": "qosconspub1zcjduepqgvpv0ky5nzkt238aus9gh90ct96wm3nddmz7pt2qwyn3lwme8dpskrfu76",
  "previous_validator_period": "0",
  "bond_token": "10000000",
  "earns_starting_height": "0",
  "first_delegate_height": "0",
  "historical_rewards": "3221431",
  "last_income_calHeight": "0",
  "last_income_calFees": "0"
}

```

#### 修改收益复投方式

`qoscli tx modify-compound --validator <validator_address> --delegator <delegator_key_name_or_account_address> --compound <compound_or_not>`

主要参数：

- `--validator`  验证人地址
- `--delegator`  委托人账户地址或秘钥库中秘钥名字
- `--compound`      收益是否复投，默认`false`

`Sansa`将收益设置为复投方式：
```bash
$ qoscli tx modify-compound --delegator Sansa --validator qosval12tr0v5uv9xpns79w8q34plakz8gh6685vddu9q --compound
```

#### 解除委托

`qoscli tx unbond --validator <validator_address> --delegator <delegator_key_name_or_account_address> --tokens <tokens> --all <unbond_all>`

主要参数：

- `--validator`  验证人地址
- `--delegator`  委托人账户地址或秘钥库中秘钥名字
- `--tokens`        解绑tokens，不能大于目前代理的QOS数量
- `--all`           是否取消全部QOS代理，默认false

`Sansa`解除代理给`qosval12tr0v5uv9xpns79w8q34plakz8gh6685vddu9q`的50个QOS：
```bash
$ qoscli tx unbond --delegator Sansa --validator qosval12tr0v5uv9xpns79w8q34plakz8gh6685vddu9q --tokens 50
```

#### 待返还委托

`qoscli query unbondings <delegator_key_name_or_account_address>`

根据质押用户地址查询该用户下所有待返还质押

查询未返还`Sansa`的质押数据：
```bash
$ qoscli query unbondings Sansa
```

#### 变更委托验证节点

`qoscli tx redelegate --from-validator <validator_address> --to-validator <validator_address> --delegator <delegator_key_name_or_account_address> --tokens <tokens> --all <unbond_all>`

主要参数：

- `--from-validator`    原始验证人地址
- `--to-validator`      新委托的验证人地址
- `--delegator`         委托人账户地址或秘钥库中秘钥名字
- `--tokens`        解绑并代理给新代理的tokens，不能大于目前代理的QOS数量
- `--compound`      新代理收益是否复投，默认`false`
- `--all`           是否从`from-owner`完全解绑，全部代理给`to-owner`，默认false

`Sansa`将代理给`qosval12tr0v5uv9xpns79w8q34plakz8gh6685vddu9q`的10个QOS转移到`qosval67werwer98sr76asdf0sdfsd98`的验证节点上：
```bash
$ qoscli tx redelegate --from-validator qosval12tr0v5uv9xpns79w8q34plakz8gh6685vddu9q --to-owner qosval67werwer98sr76asdf0sdfsd98 --delegator Sansa --tokens 10
```

#### 待执行委托变更

`qoscli query redelegations <delegator_key_name_or_account_address>`

根据质押用户地址查询该用户下所有待执行委托变更

查询未返还`Sansa`的待执行委托变更：
```bash
$ qoscli query redelegations Sansa
```

### 治理

* `qoscli tx submit-proposal`  [提交提议](#提交提议)
* `qoscli query proposal`      [提议查询](#提议查询)
* `qoscli query proposals`     [提议列表](#提议列表)
* `qoscli tx deposit`          [提议抵押](#提议抵押)
* `qoscli query deposit`       [抵押查询](#抵押查询)
* `qoscli query deposits`      [抵押列表](#抵押列表)
* `qoscli tx vote`             [提议投票](#提议投票)
* `qoscli query vote`          [投票查询](#投票查询)
* `qoscli query votes`         [投票列表](#投票列表)
* `qoscli query tally`         [投票统计](#投票统计)
* `qoscli query params`        [参数查询](#参数查询)

#### 提交提议

`qoscli tx submit-proposal
    --title <proposal_title>
    --proposal-type <proposal_type>
    --proposer <proposer_key_name_or_account_address>
    --deposit <deposit_amount_of_qos>
    --description <description>`

主要参数：

必填参数：

- `--title`             标题
- `--proposal-type`     提议类型：`Text`、`ParameterChange`、`TaxUsage`
- `--proposer`          提议账户，账户地址或密钥库中密钥名字
- `--deposit`           提议押金，不能小于`min_deposit`与`min_proposer_deposit_rate`乘积
- `--description`       描述信息

`TaxUsage`类型提议特有参数：

- `--dest-address`      目标地址，用于接收QOS
- `--percent`           社区费池提取比例，小数0~1

`ParameterChange`类型提议特有参数：

- `--params`            参数列表，格式：`module:key_name:value,module:key_name:value`，如：gov:min_deposit:10000

`ModifyInflation`类型提议特有参数：

- `inflation-phrases`   完整通胀规则，例如'[{"end_time":"2023-10-20T00:00:00Z","total_amount":"25500000000000","applied_amount":"0"},{"end_time":"2027-10-20T00:00:00Z","total_amount":"12750000000000","applied_amount":"0"},{"end_time":"2031-10-20T00:00:00Z","total_amount":"6375000000000","applied_amount":"0"},{"end_time":"2035-10-20T00:00:00Z","total_amount":"3187500000000","applied_amount":"0"},{"end_time":"2039-10-20T00:00:00Z","total_amount":"1593750000000","applied_amount":"0"},{"end_time":"2043-10-20T00:00:00Z","total_amount":"796875000000","applied_amount":"0"},{"end_time":"2047-10-20T00:00:00Z","total_amount":"796875000000","applied_amount":"0"}]'
- `total-amount`        QOS总发行量

`SoftwareUpgrade`类型提议特有参数：

- `--version`           QOS软件版本
- `--data-height`       数据版本
- `--genesis-file`      genesis.json文件url
- `--genesis-md5`       genesis.json文件md5
- `--for-zero-height`   清除本地数据，从第0高度重新开始

`Arya`提交一个文本提议：
```bash
$ qoscli tx submit-proposal --title 'update qos' --proposal-type Text --proposer Arya --deposit 100000000 --description 'this is the description'
```

`Arya`提交一个参数修改提议：
```bash
$ qoscli tx submit-proposal --title 'update parameters' --proposal-type ParameterChange --proposer Arya --deposit 100000000 --description 'this is the description' --params gov:min_deposit:1000
```

假设`Arya`在QOS初始化时已经通过[添加系统账户](qosd.md#添加系统账户) 添加到了`genesis.json`，`Arya`提交一个提取费池提议：
```bash
$ qoscli tx submit-proposal --title 'use tax' --proposal-type TaxUsage --proposer Arya --deposit 100000000 --description 'this is the description' --dest-address Sansa --percent 0.5
```

Arya`提交一个修改通胀规则提议：
```bash
$ qoscli tx submit-proposal --title 'add inflation phrase' --proposal-type ModifyInflation --proposer Arya --deposit 100000000 --description 'this is the description' --total-amount 10000000000000 --inflation-phrases '[{"end_time":"2023-10-20T00:00:00Z","total_amount":"25500000000000","applied_amount":"0"},{"end_time":"2027-10-20T00:00:00Z","total_amount":"12750000000000","applied_amount":"0"},{"end_time":"2031-10-20T00:00:00Z","total_amount":"6375000000000","applied_amount":"0"},{"end_time":"2035-10-20T00:00:00Z","total_amount":"3187500000000","applied_amount":"0"},{"end_time":"2039-10-20T00:00:00Z","total_amount":"1593750000000","applied_amount":"0"},{"end_time":"2043-10-20T00:00:00Z","total_amount":"796875000000","applied_amount":"0"},{"end_time":"2047-10-20T00:00:00Z","total_amount":"796875000000","applied_amount":"0"}]'
```

`Arya`提交一个软件升级提议：
```bash
$ qoscli tx submit-proposal --title 'update qos' --proposal-type SoftwareUpgrade --proposer Arya --deposit 100000000 --description 'upgrade qos to v0.0.6 with genesis file exporting in height 100' --genesis-file "https://.../genesis.json" --data-height 110 --version "0.0.6" --genesis-md5 88c4827158d194116b66b561691e83ef
```

#### 提议查询

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

#### 提议列表

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

#### 提议抵押

提议在抵押、投票阶段都可以执行下面的抵押交易：

`qoscli tx deposit --proposal-id <proposal_id> --depositor <depositor_key_name_or_account_address> --amount <amount_of_qos>`

主要参数：

- `--proposal-id`       提议ID
- `--depositor`         抵押账户，地址或密钥库名字
- `--amount`            抵押QOS数量

`Arya`抵押100000个QOS到3号提议：
```bash
$ qoscli tx deposit --proposal-id 1 --depositor Arya --amount 100000
```

#### 抵押查询

`qoscli query deposit <proposal-id> <depositer>`

主要参数：

- `--proposal-id`       提议ID
- `--depositor`         抵押账户，地址或密钥库名字

查询`Arya`在编号为1的提议上的抵押：
```bash
$ qoscli query deposit 1 Arya --indent
```

查询结果：
```bash
{
  "depositor": "address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy",
  "proposal_id": "1",
  "amount": "100000000"
}
```

#### 抵押列表

`qoscli query deposits <proposal-id>`

主要参数：

- `--proposal-id`       提议ID

查询编号为1的提议上的所有抵押：
```bash
$ qoscli query deposits 1 --indent
```

查询结果：
```bash
[
  {
    "depositor": "address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy",
    "proposal_id": "1",
    "amount": "100000000"
  }
]
```

#### 提议投票

进入投票阶段的提议可通过下面指令进行投票操作：

`qoscli tx vote --proposal-id <proposal_id> --voter <voter_key_name_or_account_address> --option <vote_option>`

主要参数：

- `--proposal-id`       提议ID
- `--voter`             投票账户，地址或密钥库名字
- `--option`            投票选项，可选值：`Yes`,`Abstain`,`No`,`NoWithVeto`

`Arya`给1号提议投票`Yes`：
```bash
$ qoscli tx vote --proposal-id 1 --voter Arya --option Yes
```

#### 投票查询

`qoscli query vote <proposal-id> <voter>`

主要参数：

- `--proposal-id`       提议ID
- `--voter`             投票账户，地址或密钥库名字

查询`Arya`在编号为1的提议上的投票信息：
```bash
$ qoscli query vote 1 Arya --indent
```

查询结果：
```bash
{
  "voter": "address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy",
  "proposal_id": "1",
  "option": "Yes"
}
```

#### 投票列表

`qoscli query votes <proposal-id>`

主要参数：

- `--proposal-id`       提议ID

查询编号为1的提议上的所有投票：
```bash
$ qoscli query votes 1 --indent
```

查询结果：
```bash
[
  {
    "voter": "address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy",
    "proposal_id": "1",
    "option": "Yes"
  }
]
```

#### 投票统计

`qoscli query tally <proposal-id>`

主要参数：

- `--proposal-id`       提议ID

查询编号为1的提议上实时统计结果：
```bash
$ qoscli query tally 1 --indent
```

查询结果：
```bash
{
  "yes": "100",
  "abstain": "0",
  "no": "0",
  "no_with_veto": "0"
}
```

#### 参数查询

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

### 通胀

* `qoscli query inflation-phrases`      [通胀规则查询](#通胀规则查询)
* `qoscli query total-inflation`        [发行总量查询](#发行总量查询)
* `qoscli query total-applied`          [流通总量查询](#流通总量查询)

#### 通胀规则查询

`qoscli query inflation-phrases`

查询结果：
```bash
[
  {
    "end_time": "2023-10-20T00:00:00Z",
    "total_amount": "25500000000000",
    "applied_amount": "0"
  },
  {
    "end_time": "2027-10-20T00:00:00Z",
    "total_amount": "12750000000000",
    "applied_amount": "0"
  },
  {
    "end_time": "2031-10-20T00:00:00Z",
    "total_amount": "6375000000000",
    "applied_amount": "0"
  },
  {
    "end_time": "2035-10-20T00:00:00Z",
    "total_amount": "3187500000000",
    "applied_amount": "0"
  },
  {
    "end_time": "2039-10-20T00:00:00Z",
    "total_amount": "1593750000000",
    "applied_amount": "0"
  },
  {
    "end_time": "2043-10-20T00:00:00Z",
    "total_amount": "796875000000",
    "applied_amount": "0"
  },
  {
    "end_time": "2047-10-20T00:00:00Z",
    "total_amount": "796875000000",
    "applied_amount": "0"
  }
]
```

#### 发行总量查询

`qoscli query total-inflation`

查询结果：
```bash
"100000000000000"
```

#### 流通总量查询

`qoscli query total-applied`

查询结果：
```bash
"49000130122714"
```

### 系统账户

* `qoscli query guardian`      [系统账户查询](#系统账户查询)
* `qoscli query guardians`     [系统账户列表](#系统账户列表)
* `qoscli tx add-guardian`     [添加系统账户](#添加系统账户)
* `qoscli tx delete-guardian`  [删除系统账户](#删除系统账户)
* `qoscli tx halt-network`     [停止网络](#停止网络)

#### 系统账户查询

`qoscli query guardian <guardian_key_name_or_account_address>`

查询`Arya`系统账户信息：
```bash
$ qoscli query guardian Arya --indent
```

查询结果：
```bash
{
  "description": "Arya",
  "guardian_type": 1,
  "address": "qosacc1ctmavdk57x0q7c9t98v7u79607222ars4qczcy",
  "creator": "qosacc1ah9uz0"
}
```

#### 系统账户列表

`qoscli query guardians`

查询所有系统账户：
```bash
$ qoscli query guardians --indent
```

查询结果：
```bash
[
  {
    "description": "Arya",
    "guardian_type": 1,
    "address": "qosacc1ctmavdk57x0q7c9t98v7u79607222ars4qczcy",
    "creator": "qosacc1ah9uz0"
  }
]
```

#### 添加系统账户

在`genesis.json`中配置的系统账户可通过下面的添加指令添加新的系统账户：

`qoscli tx add-guardian --address <new_guardian_key_name_or_account_address> --creator <creator_key_name_or_account_address> --description <description>`

主要参数：

- `--address`         新系统账户，账户地址或密钥库中密钥名字
- `--creator`         创建账户，账户地址或密钥库中密钥名字
- `--description`     描述

`Arya`添加`Sansa`为系统账户：
```bash
$ qoscli tx add-guardian --address Sansa --creator Arya --description 'set Sansa to be a guardian'
```

#### 删除系统账户

在`genesis.json`中配置的系统账户可通过下面的指令删除非`genesis.json`中配置的系统账户：

`qoscli tx delete-guardian --address <new_guardian_key_name_or_account_address> --deleted-by <delete_operator_key_name_or_account_address>`

主要参数：

- `--address`         系统账户，账户地址或密钥库中密钥名字
- `--deleted-by`      删除操作账户，账户地址或密钥库中密钥名字

`Arya`将`Sansa`从系统账户中删除：
```bash
$ qoscli tx delete-guardian --address Sansa --deleted-by Arya
```

#### 停止网络

紧急情况下，系统账户可及时停止网络运行。

`qoscli tx halt-network --guardian <guardian_key_name_or_account_address> --reason <reason_for_halting_network>`

主要参数：

- `--guardian`   系统账户，账户地址或密钥库中密钥名字
- `--reason`     停网原因

网络出现重大bug，`Arya`停止QOS网络：
```bash
$ qoscli tx halt-network --guardian Arya --reason 'bug'
```