# Networks

## Single-node
* init
```
qosd init --chain-id=qos-test
{
  "chain_id": "qos-test",
  "node_id": "c1c44c8ab99b894b559e6afd9f442d00e667dd9e",
  "app_message": {
    "name": "Arya",
    "pass": "12345678",
    "address": "address1cnfqru6rts4nz224mvrf58ne427uthmcut4kc3",
    "secret": "course mimic uncover all man staff economy robust endorse series boring order apology document same retreat pelican choose skate round buzz habit transfer spoon"
  }
}
```
注意init 可添加--home flag指定配置文件地址，默认在$HOME/.qosd
* start
```
qosd start --with-tendermint
```
如果一切正常，会看到控制台输出打块信息

可查看创世[账户](../spec/account.md)，进行[交易](../client/txs.md)操作

## Cluster

四个Validator节点为例

### 单台机器

* init
```
qosd init --chain-id=qos-test --name=node --home=$HOME/node1
qosd init --chain-id=qos-test --name=node --home=$HOME/node2
qosd init --chain-id=qos-test --name=node --home=$HOME/node3
qosd init --chain-id=qos-test --name=node --home=$HOME/node4
```

* 配置validators

查看genesis.json validators部分
```
# node1
cat $HOME/node1/config/genesis.json
"validators": [
    {
      "pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "uiKoUyDs5+5SrtTeZjPhO5Q87B2GTDcXKcUkn1k5J/g="
      },
      "power": "10",
      "name": ""
    }
]

# node2
cat $HOME/node2/config/genesis.json
"validators": [
    {
      "pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "nDWoLDHfEMLqRwxQJR6oK/XOGzwZPUk8X4a5J6UVoMY="
      },
      "power": "10",
      "name": ""
    }
]

# node3
cat $HOME/node3/config/genesis.json
"validators": [
    {
      "pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "gQrmTG5NurpaLlzEvbC61fFKSMAQK+2BeYquJSIq4aI="
      },
      "power": "10",
      "name": ""
    }
]

# node4
cat $HOME/node4/config/genesis.json
"validators": [
    {
      "pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "exEwlLZSv6mnv8hD1MtISnU2/dsKn+Gi9o2s2sRwEuE="
      },
      "power": "10",
      "name": ""
    }
]
```
统一四个节点的genesis.json validators部分
```
"validators": [
    {
      "pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "uiKoUyDs5+5SrtTeZjPhO5Q87B2GTDcXKcUkn1k5J/g="
      },
      "power": "10",
      "name": ""
    },
    {
      "pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "nDWoLDHfEMLqRwxQJR6oK/XOGzwZPUk8X4a5J6UVoMY="
      },
      "power": "10",
      "name": ""
    },
    {
      "pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "gQrmTG5NurpaLlzEvbC61fFKSMAQK+2BeYquJSIq4aI="
      },
      "power": "10",
      "name": ""
    },
    {
      "pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "exEwlLZSv6mnv8hD1MtISnU2/dsKn+Gi9o2s2sRwEuE="
      },
      "power": "10",
      "name": ""
    }
]
```

* app_state

统一genesis.json中 spp_state信息，参照[genesis](../spec/genesis.md)

这里我们可以统一用node1的genesis.json的app_state部分


* 查看node1 node id
```
qosd tendermint show-node-id --home=$HOME/node1
b70c6ce13a11e14ee14bc793cbef835aa1b4b6bb
```

* 修改node2配置

config.toml
```
cd $HOME/node2/config
vi config.toml

# TCP or UNIX socket address of the ABCI application,
# or the name of an ABCI application compiled in with the Tendermint binary
proxy_app = "tcp://127.0.0.1:26668"

# TCP or UNIX socket address for the profiling server to listen on
prof_laddr = "localhost:6061"

# TCP or UNIX socket address for the RPC server to listen on
laddr = "tcp://0.0.0.0:26667"

# Address to listen for incoming connections
laddr = "tcp://0.0.0.0:26666"

# Comma separated list of nodes to keep persistent connections to
persistent_peers = "b70c6ce13a11e14ee14bc793cbef835aa1b4b6bb@127.0.0.1:26656"

# Address to listen for Prometheus collector(s) connections
prometheus_listen_addr = ":26670"
```
* 修改node3配置
```
cd $HOME/node3/config
vi config.toml

# TCP or UNIX socket address of the ABCI application,
# or the name of an ABCI application compiled in with the Tendermint binary
proxy_app = "tcp://127.0.0.1:26678"

# TCP or UNIX socket address for the profiling server to listen on
prof_laddr = "localhost:6062"

# TCP or UNIX socket address for the RPC server to listen on
laddr = "tcp://0.0.0.0:26677"

# Address to listen for incoming connections
laddr = "tcp://0.0.0.0:26676"

# Comma separated list of nodes to keep persistent connections to
persistent_peers = "b70c6ce13a11e14ee14bc793cbef835aa1b4b6bb@127.0.0.1:26656"

# Address to listen for Prometheus collector(s) connections
prometheus_listen_addr = ":26680"
```
* 修改node4配置
```
cd $HOME/node4/config
vi config.toml

# TCP or UNIX socket address of the ABCI application,
# or the name of an ABCI application compiled in with the Tendermint binary
proxy_app = "tcp://127.0.0.1:26688"

# TCP or UNIX socket address for the profiling server to listen on
prof_laddr = "localhost:6063"

# TCP or UNIX socket address for the RPC server to listen on
laddr = "tcp://0.0.0.0:26687"

# Address to listen for incoming connections
laddr = "tcp://0.0.0.0:26686"

# Comma separated list of nodes to keep persistent connections to
persistent_peers = "b70c6ce13a11e14ee14bc793cbef835aa1b4b6bb@127.0.0.1:26656"

# Address to listen for Prometheus collector(s) connections
prometheus_listen_addr = ":26690"
```

* start
```
qosd start --with-tendermint --home=$HOME/node1
qosd start --with-tendermint --home=$HOME/node2
qosd start --with-tendermint --home=$HOME/node3
qosd start --with-tendermint --home=$HOME/node4
```

### 四台机器

第一台IP为ip1

* init

四台机器分别执行 init 命令
```
qosd init --chain-id=qos-test --name=node
```

* 统一app_state

统一genesis.json中 spp_state信息，参照[genesis](../spec/genesis.md)

* 查看node1 node id</br>
在第一台机器上运行：
```
qosd tendermint show-node-id
b70c6ce13a11e14ee14bc793cbef835aa1b4b6bb
```

* 修改node2配置
```
cd $HOME/.qosd/config
vi config.toml

# Comma separated list of nodes to keep persistent connections to
persistent_peers = "b70c6ce13a11e14ee14bc793cbef835aa1b4b6bb@ip1:26656"

```
* 修改node3配置
```
cd $HOME/.qosd/config
vi config.toml

# Comma separated list of nodes to keep persistent connections to
persistent_peers = "b70c6ce13a11e14ee14bc793cbef835aa1b4b6bb@ip1:26656"

```
* 修改node4配置
```
cd $HOME/.qosd/config
vi config.toml

# Comma separated list of nodes to keep persistent connections to
persistent_peers = "b70c6ce13a11e14ee14bc793cbef835aa1b4b6bb@ip1:26656"

```

* start</br>
四台机器上分别执行
```
qosd start --with-tendermint
```

### Tendermint testnet cmd

qosd已集成tendermint testnet命令行工具，可批量生成集群配置文件，相关命令参考：
```
qosd testnet --help
testnet will create "v" + "n" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).

Note, strict routability for addresses is turned off in the config file.

Optionally, it will fill in persistent_peers list in config file using either hostnames or IPs.

Example:

	tendermint testnet --v 4 --o ./output --populate-persistent-peers --starting-ip-address 192.168.10.2

Usage:
  qosd testnet [flags]

Flags:
  -h, --help                         help for testnet
      --hostname-prefix string       Hostname prefix (node results in persistent peers list ID0@node0:26656, ID1@node1:26656, ...) (default "node")
      --n int                        Number of non-validators to initialize the testnet with
      --node-dir-prefix string       Prefix the directory name for each node with (node results in node0, node1, ...) (default "node")
      --o string                     Directory to store initialization data for the testnet (default "./mytestnet")
      --p2p-port int                 P2P Port (default 26656)
      --populate-persistent-peers    Update config of each node with the list of persistent peers build using either hostname-prefix or starting-ip-address (default true)
      --starting-ip-address string   Starting IP address (192.168.0.1 results in persistent peers list ID0@192.168.0.1:26656, ID1@192.168.0.2:26656, ...)
      --v int                        Number of validators to initialize the testnet with (default 4)

Global Flags:
      --home string        directory for config and data (default "/home/imuge/.qosd")
      --log_level string   Log level (default "main:info,state:info,*:error")
      --trace              print out full stack trace on errors
```
