# QOS Daemon server

`qosd`创建、初始化、启动QOS网络命令：

* `init`                  [初始化](#初始化)
* `add-genesis-accounts`  [设置创世账户](#设置账户)
* `add-guardian`          [添加系统账户](#添加系统账户)
* `gentx`                 [生成创世交易](#生成创世交易)
* `collect-gentxs`        [收集创世交易](#收集创世交易)
* `config-root-ca`        [设置CA](#设置ca)
* `start`                 [启动](#启动)
* `export        `        [状态导出](#状态导出)
* `testnet`               [初始化测试网络](#初始化测试网络)
* `unsafe-reset-all`      [重置](#重置)
* `tendermint`            [Tendermint](#tendermint)
* `version`               [版本信息](#版本)

全局参数：

| 参数 | 默认值 | 说明 |
| :--- | :---: | :--- |
|--home string        | "$HOME/.qosd" |directory for config and data (default "$HOME/.qosd")|
|--log_level string   | "main:info,state:info,*:error" |Log level (default "main:info,state:info,*:error")|
|--trace              |  |print out full stack trace on errors|


## 初始化

`qosd init --moniker <your_custom_moniker> --chain-id <chain_id> --overwrite <overwrite>`

参数说明:

- `--moniker`   在P2P网络中的名称，与`config.toml`中`moniker`配置项对应，可后期修改
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
$ qosd add-genesis-accounts qosacc1c7nh7qquvjm3p28xpsnfn420437ztvzy2hwdtk,10000QOS
```

会在`genesis.json`文件`app-state`中`accounts`部分添加地址为`qosacc1c7nh7qquvjm3p28xpsnfn420437ztvzy2hwdtk`，持有10000QOS的账户信息。

## 添加系统账户

`qosd add-guardian --address <address> --description <description>`

参数说明:

- `--address`     系统账户地址，可接收`TaxUsageProposal`提议从社区费池提取的QOS代币
- `--description` 描述

添加系统账户至`genesis.json`文件：
```bash
$ qosd add-guardian --address qosacc1c7nh7qquvjm3p28xpsnfn420437ztvzy2hwdtk --description "this is the description"
```

会在`genesis.json`文件`app-state`中`guardian`部分添加地址为`qosacc1c7nh7qquvjm3p28xpsnfn420437ztvzy2hwdtk`的系统账户。

## 生成创世交易

生成创建验证节点[TxCreateValidator](../spec/staking.md#TxCreateValidator)交易

`qosd gentx --moniker <validator_name> --owner <account_address> --tokens <tokens>`

参数说明参照[成为验证节点](qoscli.md#成为验证节点)

生成验证节点交易：
```bash
$ qosd gentx --moniker "Arya's node" --owner qosacc1c7nh7qquvjm3p28xpsnfn420437ztvzy2hwdtk --tokens 1000
```

默认会在`$HOME/.qosd/config/gentx`目录下生成以`nodeID@IP`为文件名的已签名的交易数据文件。

## 收集创世交易

`qosd collect-gentxs`

收集`gentx`目录下交易数据，填充到`genesis.json`中`app_state`下`gen_txs`中。

## 设置CA

`qosd config-root-ca --qcp <qcp_root.pub> --qsc <qsc_root.pub>`

`<qcp_root.pub>`、`<qsc_root.pub>`为根证书公钥文件路径

设置Root CA公钥信息，用于[代币](qoscli.md#代币)和[联盟链](qoscli.md#联盟链)涉及到证书操作的校验。

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
|--with-tendermint                 | true |Run abci app embedded in-process with tendermint|

启动QOS网络

```bash
$ qosd start
```

## 状态导出

`qosd export --height <block_height> --for-zero-height <export_state_to_start_at_height_zero> -o <directory for exported json file>`

主要参数：

- `--height`            指定导出区块高度
- `--for-zero-height`   是否导出状态从0高度重新启动网络
- `--o`                 导出文件位置

导出区块高度为4的状态数据：
```bash
qosd export --height 4
```

到处完默认成会在`$HOME/.qosd`下生成以`genesis-<height>-<timestamp>.json`命名文件。

## 初始化测试网络

`qosd testnet`命令行工具，可批量生成集群配置文件，相关命令参考：
```bash
testnet will create "v" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).

Note, strict routability for addresses is turned off in the config file.

Example:

	qosd testnet --chain-id=qostest --v=4 --o=./output --starting-ip-address=192.168.1.2 --genesis-accounts=qosacc1c7nh7qquvjm3p28xpsnfn420437ztvzy2hwdtk,1000000qos

Usage:
  qosd testnet [flags]

Flags:
      --chain-id string              Chain ID
      --compound                     whether the validator's income is calculated as compound interest, default: true (default true)
      --genesis-accounts string      Add genesis accounts to genesis.json, eg: qosacc1c7nh7qquvjm3p28xpsnfn420437ztvzy2hwdtk,1000000qos,1000000qstars. Multiple accounts separated by ';'
      --guardians string             addresses for guardian. Multiple addresses separated by ','
  -h, --help                         help for testnet
      --home-client string           directory for keybase (default "$HOME/.qoscli")
      --hostname-prefix string       Hostname prefix (node results in persistent peers list ID0@node0:26656, ID1@node1:26656, ...) (default "node")
      --node-dir-prefix string       Prefix the directory name for each node with (node results in node0, node1, ...) (default "node")
      --o string                     Directory to store initialization data for the testnet (default "./mytestnet")
      --qcp-root-ca string           Config pubKey of root CA for QSC
      --qsc-root-ca string           Config pubKey of root CA for QCP
      --starting-ip-address string   Starting IP address (192.168.0.1 results in persistent peers list ID0@192.168.0.1:26656, ID1@192.168.0.2:26656, ...)
      --v int                        Number of validators to initialize the testnet with (default 4)

Global Flags:
      --home string        directory for config and data (default "$HOME/.qosd")
      --log_level string   Log level (default "main:info,state:info,*:error")
      --trace              print out full stack trace on errors


```

主要参数说明：
- chain-id            链ID
- genesis-accounts    初始账户
- hostname-prefix     hostName前缀
- moniker             moniker
- qcp-root-ca         pubKey of root CA for QCP
- qsc-root-ca         pubKey of root CA for QSC
- compound            收益复投方式，默认true，即收益参与复投
- starting-ip-address 起始IP地址

## 重置

`qosd unsafe-reset-all`

重置区块链数据库，删除地址簿文件，重置状态至初始状态。

## Tendermint

tendermint子命令：

- `qosd tendermint show-address`    Show this node's tendermint validator address
- `qosd tendermint show-node-id`    Show this node's ID
- `qosd tendermint show-validator`  Show this node's tendermint validator info

## 版本
与[qoscli version](qoscli.md#版本)相同
