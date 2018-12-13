# Networks

## Single-node
* init
```
qosd init --chain-id=qos-test
{
  "chain_id": "qos-test",
  "node_id": "1c3100c28a44f1facf45aa83e9aa3d8ff8ac6b1f",
  "app_message": "null"
}
```
注意init 可添加--home flag指定配置文件地址，默认在$HOME/.qosd
`init`操作后,通过执行`qosd add-genesis-validator`添加validator

* add-genesis-accounts

使用`qosd add-genesis-accounts`初始化account账户到配置文件中.

> 使用`qoscli keys add `创建account账户

```

$ qoscli keys add qosInitAcc
Enter a passphrase for your key:
Repeat the passphrase:

$ qoscli keys list

NAME:   TYPE:   ADDRESS:                                                PUBKEY:
qosInitAcc      local   address1lly0audg7yem8jt77x2jc6wtrh7v96hgve8fh8  4MFA7MtUl1+Ak3WBtyKxGKvpcu4e5ky5TfAC26cN+mQ=

# 初始化账户
$ qosd add-genesis-accounts address1lly0audg7yem8jt77x2jc6wtrh7v96hgve8fh8,1000000qos

```

* config-root-ca

使用`qosd config-root-ca`初始化root CA公钥到配置文件.
```
Config root CA

Usage:
  qosd config-root-ca [root.pub] [flags]

Flags:
  -h, --help   help for config-root-ca

Global Flags:
      --home string        directory for config and data (default "/home/imuge/.qosd")
      --log_level string   Log level (default "main:info,state:info,*:error")
      --trace              print out full stack trace on errors
      
# 设置roort CA
$ qosd config-root-ca root.pub
```
查看genesis.json内容，确认配置成功。

* add-genesis-validator

使用`qosd add-genesis-validator`初始化validator到配置文件中.

```

$ qosd add-genesis-validator -h

Add genesis validator to genesis.json

Usage:
  qosd add-genesis-validator [flags]

Flags:
      --consPubkey string   validator's ed25519 consPubkey
  -h, --help                help for add-genesis-validator
      --operator string     operator address
      --power int           validator's voting power. default is 10 (default 10)

# 使用上面的初始化账户地址作为operator
$ qosd add-genesis-validator --operator address1lly0audg7yem8jt77x2jc6wtrh7v96hgve8fh8


```

* start
```
$ qosd start --with-tendermint
```
如果一切正常，会看到控制台输出打块信息

## Cluster

四个Validator节点为例

### 单台机器

* init
```
$ cd && mkdir node{1..4}
$ for i in {1..4}; do qosd init --chain-id=qos-test --name=node${i} --home=$HOME/node${i}; done
```

* 添加创世账户

```
$ qoscli keys add qosgenesisacc
Enter a passphrase for your key:
Repeat the passphrase:

$ qoscli keys  list
NAME:	TYPE:	ADDRESS:						PUBKEY:
qosgenesisacc	local	address1rh47fd6ykkj0kpkukkt9pskgppfl30lpv9n9pu	EnChknIClMgiwcqCKjIraYZdK4+wTaATAfp4GUNUIAo=


$ for i in {1..4}; do qosd add-genesis-account --addr $(qoscli keys  list | grep qosgenesisacc | awk '{print $3}')  --coins 100000qos,1000000qstar --home=$HOME/node${i}; done

//查看genesis.json配置
$ cat $HOME/node1/config/genesis.json

```


* 配置validators

查看每个节点的validator 配置文件:

```
//node1
$ cat $HOME/node1/config/priv_validator.json

//node2
$ cat $HOME/node2/config/priv_validator.json

//node3
$ cat $HOME/node3/config/priv_validator.json

//node4
$ cat $HOME/node4/config/priv_validator.json

```

使用`qosd add-genesis-validator`命令将每个节点的`validator`添加至`genesis.json`配置文件中

**添加`validator`需要指定操作者`operator`,可以通过`qoscli keys add`分别添加操作者`operator`**


```
$ qoscli keys add node1oper
$ qoscli keys add node2oper
$ qoscli keys add node3oper
$ qoscli keys add node4oper
$ qoscli keys list
```

以下示例中使用创世账户`qosgenesisacc`作为4个节点`validator`的`operator`,也可以为节点指定不同的`operator`:

```
$ export node1pk=$(sed -n '/PubKey/,/value/p' $HOME/node1/config/priv_validator.json|sed 1d|awk -F\" '{print $4}')

$ export node2pk=$(sed -n '/PubKey/,/value/p' $HOME/node2/config/priv_validator.json|sed 1d|awk -F\" '{print $4}')

$ export node3pk=$(sed -n '/PubKey/,/value/p' $HOME/node3/config/priv_validator.json|sed 1d|awk -F\" '{print $4}')

$ export node4pk=$(sed -n '/PubKey/,/value/p' $HOME/node4/config/priv_validator.json|sed 1d|awk -F\" '{print $4}')

$ export operator=$(qoscli keys  list | grep qosgenesisacc | awk '{print $3}')

$ export node1id=$(qosd tendermint show-node-id --home=$HOME/node1)

#node1执行

$ for i in {1..4}; do qosd add-genesis-validator --home=$HOME/node1 --consPubkey $(eval echo '$'"node"${i}pk) --operator $operator --power 10;done

#将node1/config/genesis.json分别拷贝至其他节点config目录下
$ cp $HOME/node1/config/genesis.json $HOME/node2/config
$ cp $HOME/node1/config/genesis.json $HOME/node3/config
$ cp $HOME/node1/config/genesis.json $HOME/node4/config
```


* 查看node1 node id
```
$ qosd tendermint show-node-id --home=$HOME/node1
51bc172321ff3ea3e82f133d5116f0c11ac905d8
```

* 修改节点端口配置文件

`config/config.toml`定义了节点启动时不同的监听端口. 如果节点在不同的机器上运行,则只需修改`persistent_peers`选项即可:

主要修改端口配置选项如下

* moniker = ""
* proxy_app = "tcp://127.0.0.1:26658"
* prof_laddr = "localhost:6060"
* [rpc] laddr = "tcp://0.0.0.0:26657"
* [p2p] laddr = "tcp://0.0.0.0:26656"
* [p2p] persistent_peers = ""
* prometheus_listen_addr = ":26660"


node1配置如下:
* moniker = "node1"
* proxy_app = "tcp://127.0.0.1:16658"
* prof_laddr = "localhost:6061"
* [rpc] laddr = "tcp://0.0.0.0:16657"
* [p2p] laddr = "tcp://0.0.0.0:16656"
* [p2p] persistent_peers = ""
* prometheus_listen_addr = ":16660"

```
$ sed -i -e s#' *moniker.*'#'moniker = "node1"'#g   \
         -e s#' *proxy_app.*'#'proxy_app = "tcp://127.0.0.1:16658"'#g \
         -e s#' *prof_laddr.*'#'prof_laddr = "localhost:6061"'#g  \
         -e '/rpc/,/p2p/s#^ *laddr = .*$#laddr = "tcp://0.0.0.0:16657"#g'  \
         -e '/p2p/,/laddr/s#^ *laddr = .*$#laddr = "tcp://0.0.0.0:16656"#g' \
         -e s#' *prometheus_listen_addr.*'#'prometheus_listen_addr = ":16660"'#g \
         $HOME/node1/config/config.toml
```

node2配置如下:

* moniker = "node2"
* proxy_app = "tcp://127.0.0.1:26658"
* prof_laddr = "localhost:6062"
* [rpc] laddr = "tcp://0.0.0.0:26657"
* [p2p] laddr = "tcp://0.0.0.0:26656"
* [p2p] persistent_peers = ""
* prometheus_listen_addr = ":26660"

```


$ sed -i -e s#' *moniker.*'#'moniker = "node2"'#g   \
         -e s#' *proxy_app.*'#'proxy_app = "tcp://127.0.0.1:26658"'#g \
         -e s#' *prof_laddr.*'#'prof_laddr = "localhost:6062"'#g  \
         -e '/rpc/,/p2p/s#^ *laddr = .*$#laddr = "tcp://0.0.0.0:26657"#g'  \
         -e '/p2p/,/laddr/s#^ *laddr = .*$#laddr = "tcp://0.0.0.0:26656"#g' \
         -e /p2p/,/upnp/s#' *persistent_peers = .*'#"persistent_peers = \"${node1id}@127.0.0.1:16656\""#g \
         -e s#' *prometheus_listen_addr.*'#'prometheus_listen_addr = ":26660"'#g \
         $HOME/node2/config/config.toml
```



node3配置如下:

* moniker = "node3"
* proxy_app = "tcp://127.0.0.1:36658"
* prof_laddr = "localhost:6063"
* [rpc] laddr = "tcp://0.0.0.0:36657"
* [p2p] laddr = "tcp://0.0.0.0:36656"
* [p2p] persistent_peers = ""
* prometheus_listen_addr = ":36660"

```
$ sed -i -e s#' *moniker.*'#'moniker = "node3"'#g   \
         -e s#' *proxy_app.*'#'proxy_app = "tcp://127.0.0.1:36658"'#g \
         -e s#' *prof_laddr.*'#'prof_laddr = "localhost:6063"'#g  \
         -e '/rpc/,/p2p/s#^ *laddr = .*$#laddr = "tcp://0.0.0.0:36657"#g'  \
         -e '/p2p/,/laddr/s#^ *laddr = .*$#laddr = "tcp://0.0.0.0:36656"#g' \
         -e /p2p/,/upnp/s#' *persistent_peers = .*'#"persistent_peers = \"${node1id}@127.0.0.1:16656\""#g \
         -e s#' *prometheus_listen_addr.*'#'prometheus_listen_addr = ":36660"'#g \
         $HOME/node3/config/config.toml
```


node4配置如下:

* moniker = "node4"
* proxy_app = "tcp://127.0.0.1:46658"
* prof_laddr = "localhost:6064"
* [rpc] laddr = "tcp://0.0.0.0:46657"
* [p2p] laddr = "tcp://0.0.0.0:46656"
* [p2p] persistent_peers = ""
* prometheus_listen_addr = ":46660"

```
$ sed -i -e s#' *moniker.*'#'moniker = "node4"'#g   \
         -e s#' *proxy_app.*'#'proxy_app = "tcp://127.0.0.1:46658"'#g \
         -e s#' *prof_laddr.*'#'prof_laddr = "localhost:6064"'#g  \
         -e '/rpc/,/p2p/s#^ *laddr = .*$#laddr = "tcp://0.0.0.0:46657"#g'  \
         -e '/p2p/,/laddr/s#^ *laddr = .*$#laddr = "tcp://0.0.0.0:46656"#g' \
         -e /p2p/,/upnp/s#' *persistent_peers = .*'#"persistent_peers = \"${node1id}@127.0.0.1:16656\""#g \
         -e s#' *prometheus_listen_addr.*'#'prometheus_listen_addr = ":46660"'#g \
         $HOME/node4/config/config.toml
```



* start
```
$ qosd start --with-tendermint --home=$HOME/node1
$ qosd start --with-tendermint --home=$HOME/node2
$ qosd start --with-tendermint --home=$HOME/node3
$ qosd start --with-tendermint --home=$HOME/node4
```

### 四台机器

第一台node1 IP为ip1

* init

四台机器分别执行 init 命令,chain-id保持一致, name不同
```
$ qosd init --chain-id=qos-test --name=node1
```

在node1上执行:

```
$ qoscli keys add genAcc

$ qoscli keys list
NAME:	TYPE:	ADDRESS:						PUBKEY:
genAcc	local	address17k688l8afk4t42dr4z5ay0cpke39we7kxm9pzk	c5UuUZ/REvHExIY/eDcQvdjxiapE+aVSd37DulaxwBU=

```

在node1配置文件`genesis.json`中添加初始账户:

```
$ qosd add-genesis-account --addr $(qoscli keys  list | grep genAcc | awk '{print $3}')  --coins 100000qos,1000000qstar
```

在node1配置文件`genesis.json`中添加node1节点为validator:

```
$ qosd add-genesis-validator --operator $(qoscli keys  list | grep genAcc | awk '{print $3}')
```

在node1配置文件`genesis.json`中分别添加其他节点为validator:

```
$ qosd add-genesis-validator --consPubkey $NODES_VALIDATOR_PUBKEY --operator $(qoscli keys  list | grep genAcc | awk '{print $3}')
```


* 将node1中`genesis.json`文件拷贝至其他节点的$HOME/.qosd/config目录下

* 查看node1 node id</br>
在node1上运行：
```
$ qosd tendermint show-node-id
b70c6ce13a11e14ee14bc793cbef835aa1b4b6bb
```

* 修改node2配置
```
$ cd $HOME/.qosd/config
$ vi config.toml

# Comma separated list of nodes to keep persistent connections to
persistent_peers = "b70c6ce13a11e14ee14bc793cbef835aa1b4b6bb@ip1:26656"

```
* 修改node3配置
```
$ cd $HOME/.qosd/config
$ vi config.toml

# Comma separated list of nodes to keep persistent connections to
persistent_peers = "b70c6ce13a11e14ee14bc793cbef835aa1b4b6bb@ip1:26656"

```
* 修改node4配置
```
$ cd $HOME/.qosd/config
$ vi config.toml

# Comma separated list of nodes to keep persistent connections to
persistent_peers = "b70c6ce13a11e14ee14bc793cbef835aa1b4b6bb@ip1:26656"

```

* start</br>
四台机器上分别执行
```
$ qosd start --with-tendermint
```

### testnet cmd

qosd testnet命令行工具，可批量生成集群配置文件，相关命令参考：
```
$ qosd testnet --help
testnet will create "v" + "n" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).

Note, strict routability for addresses is turned off in the config file.

Optionally, it will fill in persistent_peers list in config file using either hostnames or IPs.

Example:

	qosd testnet --chain-id=qostest --v=4 --o=./output --starting-ip-address=192.168.1.2 --genesis-accounts=address16lwp3kykkjdc2gdknpjy6u9uhfpa9q4vj78ytd,1000000qos,1000000qstars

Usage:
  qosd testnet [flags]

Flags:
      --chain-id string              Chain ID
      --genesis-accounts string      Add genesis accounts to genesis.json, eg: address16lwp3kykkjdc2gdknpjy6u9uhfpa9q4vj78ytd,1000000qos,1000000qstars. Multiple accounts separated by ';'
  -h, --help                         help for testnet
      --hostname-prefix string       Hostname prefix (node results in persistent peers list ID0@node0:26656, ID1@node1:26656, ...) (default "node")
      --moniker string               Moniker
      --n int                        Number of non-validators to initialize the testnet with
      --node-dir-prefix string       Prefix the directory name for each node with (node results in node0, node1, ...) (default "node")
      --o string                     Directory to store initialization data for the testnet (default "./mytestnet")
      --p2p-port int                 P2P Port (default 26656)
      --populate-persistent-peers    Update config of each node with the list of persistent peers build using either hostname-prefix or starting-ip-address (default true)
      --root-ca string               Config root CA
      --starting-ip-address string   Starting IP address (192.168.0.1 results in persistent peers list ID0@192.168.0.1:26656, ID1@192.168.0.2:26656, ...)
      --v int                        Number of validators to initialize the testnet with (default 4)

Global Flags:
      --home string        directory for config and data (default "/home/imuge/.qosd")
      --log_level string   Log level (default "main:info,state:info,*:error")
      --trace              print out full stack trace on errors
```
app_state下validators列表中的operator与其cons_pubkey对应