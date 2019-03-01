# QOS Daemon server

`qosd`创建、初始化、启动QOS网络命令：

* `init`                  [初始化](#初始化) 
* `add-genesis-accounts`  [设置创世账户](#设置账户) 
* `add-genesis-validator` [设置验证节点](#设置验证节点) 
* `config-root-ca`        [设置CA](#设置ca) 
* `start`                 [启动](#启动) 
* `export        `        [状态导出](#状态导出) 
* `testnet`               [初始化测试网络](#初始化测试网络) 
* `unsafe-reset-all`      [重置](#重置) 
* `tendermint`            [Tendermint](#tendermint) 
* `version`               版本信息

全局参数：

| 参数 | 默认值 | 说明 |
| :--- | :---: | :--- |
|--home string        | "$HOME/.qosd" |directory for config and data (default "$HOMW/.qosd")|
|--log_level string   | "main:info,state:info,*:error" |Log level (default "main:info,state:info,*:error")|
|--trace              |  |print out full stack trace on errors|


## 初始化

`qosd init --moniker <your_custom_moniker> --chain-id <chain_id> --overwrite <overwrite>`

参数说明:

- `--moniker`      在P2P网络中的名称，与`config.toml`中`moniker`配置项对应，可后期修改
- `--chain-id`  链ID，链ID一致的节点才能组成同一个P2P网络
- `--overwrite` 是否覆盖已存在初始文件

初始化`genesis`、`priv-validator`、`p2p-node`文件

执行：
```bash
$ qosd init --moniker capricorn-1000
```
输出：
```bash
{
  "chain_id": "test-chain-9nlhQS",
  "node_id": "c427167c8d2838b00a46e33c4b325a7f05bd2c16",
  "app_message": null
}
```

会在`$HOME/.qosd/`下创建`data`和`config`两个目录。
`data`为空目录，用于存储网络启动后保存的数据，
`config`中会生成`config.toml`，`genesis.json`，`node_key.json`，`priv_validator.json`四个文件。

## 设置账户

`qosd add-genesis-accounts <account_coin_s>`

`<account_coin_s>`账户币种币值列表，eg:[address1],[coin1],[coin2];[address2],[coin1],[coin2]

添加创世账户至`genesis.json`文件：
```bash
$ qosd add-genesis-accounts address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy,10000QOS
```

会在`genesis.json`文件`app-state`中`accounts`部分添加地址为`address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy`，持有10000QOS的账户信息。

## 设置验证节点

`qosd add-genesis-validator --name <validator_name> --owner <account_address> --tokens <tokens> --description <description>`

参数说明参照[成为验证节点](qoscli.md#成为验证节点)

设置验证节点信息：
```bash
qoscli tx create-validator --name "Arya's node" --owner address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy --tokens 1000 --description "I am a validator."
```

会在`genesis.json`文件`app-state`中`validators`部分添加验证节点信息。

## 设置CA

`qosd config-root-ca --qcp <qcp_root.pub> --qsc <qsc_root.pub>`

`<qcp_root.pub>`、`<qsc_root.pub>`为根证书公钥文件路径

设置Root CA公钥信息，用于[联盟币](qoscli.md#联盟币（qsc）)和[联盟链](qoscli.md#联盟链（qcp）)涉及到证书操作的校验。

## 启动

`qosd start`

| 参数 | 默认值 | 说明 |
| :--- | :---: | :--- |
|--abci string                     | "socket" |Specify abci transport (socket | grpc) (default "socket")|
|--address string                  | "tcp://0.0.0.0:26658") |Listen address (default "tcp://0.0.0.0:26658")|
|--consensus.create_empty_blocks   | true |Set this to false to only produce blocks when there are txs or when the AppHash changes (default true)|
|--fast_sync                       | true |Fast blockchain syncing (default true)|
|--moniker string                  | <your_computer_name> |Node Name|
|--p2p.laddr string                | "tcp://0.0.0.0:26656" |Node listen address. (0.0.0.0:0 means any interface, any port) (default "tcp://0.0.0.0:26656")|
|--p2p.persistent_peers string     | "" |Comma-delimited ID@host:port persistent peers|
|--p2p.pex                         | true |Enable/disable Peer-Exchange (default true)|
|--p2p.private_peer_ids string     | "" |Comma-delimited private peer IDs|
|--p2p.seed_mode                   | false |Enable/disable seed mode|
|--p2p.seeds string                | "" |Comma-delimited ID@host:port seed nodes|
|--p2p.upnp                        | false |Enable/disable UPNP port forwarding|
|--priv_validator_laddr string     | "" |Socket address to listen on for connections from external priv_validator process|
|--proxy_app string                | "tcp://127.0.0.1:26658" |Proxy app address, or 'nilapp' or 'kvstore' for local testing. (default "tcp://127.0.0.1:26658")|
|--pruning string                  | "syncable" |Pruning strategy: syncable, nothing, everything (default "syncable")|
|--rpc.grpc_laddr string           | "" |GRPC listen address (BroadcastTx only). Port required|
|--rpc.laddr string                | "tcp://0.0.0.0:26657" |RPC listen address. Port required (default "tcp://0.0.0.0:26657")|
|--rpc.unsafe                      | false |Enabled unsafe rpc methods|
|--trace-store string              | false |Enable KVStore tracing to an output file|
|--with-tendermint                 | false |Run abci app embedded in-process with tendermint|

启动QOS网络

```bash
$ qosd start
```
启动QOS网络，并启动tendermint，如果正确[配置Validator](#设置验证节点)会看到打块信息。

## 状态导出

`qosd export --height <block_height> --for-zero-height <export_state_to_start_at_height_zero>`

主要参数：

- `--height`            指定导出区块高度
- `--for-zero-height`   是否导出状态从0高度重新启动网络

导出区块高度为4的状态数据：
```bash
qosd export --height 4
```

## 初始化测试网络

`qosd testnet`
| 参数 | 默认值 | 说明 |
| :--- | :---: | :--- |
|--chain-id string              | 随机字符串 |Chain ID|
|--genesis-accounts string      | "" |Add genesis accounts to genesis.json, eg: address16lwp3kykkjdc2gdknpjy6u9uhfpa9q4vj78ytd,1000000qos,1000000qstars. Multiple accounts separated by ';'|
|--hostname-prefix string       | "node" |Hostname prefix (node results in persistent peers list ID0@node0:26656, ID1@node1:26656, ...) (default "node")|
|--n int                        | 0 |Number of non-validators to initialize the testnet with|
|--name string                  | <your_computer_name> |Moniker|
|--node-dir-prefix string       | "node" |Prefix the directory name for each node with (node results in node0, node1, ...) (default "node")|
|--o string                     | "./mytestnet" |Directory to store initialization data for the testnet (default "./mytestnet")|
|--p2p-port int                 | 26656 |P2P Port (default 26656)|
|--populate-persistent-peers    | true |Update config of each node with the list of persistent peers build using either hostname-prefix or starting-ip-address (default true)|
|--root-ca string               | "" |Config pubKey of root CA|
|--starting-ip-address string   | "" |Starting IP address (192.168.0.1 results in persistent peers list ID0@192.168.0.1:26656, ID1@192.168.0.2:26656, ...)|
|--v int                        | 4 |Number of validators to initialize the testnet with (default 4)|

创建`v`+`n`个目录，每个目录创建必需的配置文件（`private validator`, `genesis`, `config`等等）。

根据实际服务器、网络配置，少量修改这`v`+`n`目录中的配置文件可以轻松搭建多个验证节点和非验证节点的QOS网络。

## 重置

`qosd unsafe-reset-all`

重置区块链数据库，删除地址簿文件，重置状态至创世状态。

## tendermint

tendermint子命令：

- `qosd tendermint show-address`    Shows this node's tendermint validator address
- `show-node-id`                    Show this node's ID
- `show-validator`                  Show this node's tendermint validator info
