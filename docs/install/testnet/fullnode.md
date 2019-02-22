# 启动完整节点

在启动完整节点前，请确保已按照[安装引导](../installation.md)正确安装QOS。

## 初始化

适用`qosd init`命令初始化节点、创建必要的配置文件。
默认的配置和数据存储目录为 `$HOME/.qosd`，可以添加`--home`修改存储位置。

```bash
$ qosd init --moniker <your_custom_moniker>
```
::: warning Note
`name`仅支持ASCII字符，使用Unicode字符将使节点无法访问
:::

执行完`qosd init`会在`$HOME/.qosd/config`下生成`genesis.json`、`config.toml`等配置文件。

## 配置运行网络

不同QOS运行网络对应不同的配置，可访问[testnets repo](https://github.com/QOSGroup/qos-testnets)了解不同网络的运行配置。
下面操作以测试网`capricorn-1000`为例。

### 替换`genesis.json`

默认路径`$HOME/.qosd/config/genesis.json`

下载`capricorn-1000`对应[`genesis.json`文件](https://raw.githubusercontent.com/QOSGroup/qos-testnets/master/capricorn-1000/genesis.json)替换本地文件。
若没有更改默认存储位置，也可通过下面的命令执行替换操作：
```bash
$ curl https://raw.githubusercontent.com/QOSGroup/qos-testnets/master/capricorn-1000/genesis.json > $HOME/.qosd/config/genesis.json
```

### 编辑`config.toml`:

默认路径`$HOME/.qosd/config/config.toml`

修改`config.toml`，找到`seeds`配置项，添加seed节点：
```toml
# Comma separated list of seed nodes to connect to
seeds = "5d9fcba29ce9a066cdd6e4c45001567a4bd1dbf4@47.100.231.9:26656"
```

## 启动

运行启动命令：

```bash
$ qosd start --log_level debug
```

控制台开始打印启动日志，提示`This node is not a validator`说明节点不是验证节点，节点开始同步QOS网络区块信息。
You can see the node is running, your node is not a validator, and your node is synchronizing blocks from the QOS testnet.
```bash
Starting ABCI with Tendermint                module=main 
Starting multiAppConn                        module=proxy impl=multiAppConn
Starting localClient                         module=abci-client connection=query impl=localClient
Starting localClient                         module=abci-client connection=mempool impl=localClient
Starting localClient                         module=abci-client connection=consensus impl=localClient
ABCI Handshake                               module=consensus appHeight=0 appHash=
ABCI Replay Blocks                           module=consensus appHeight=0 storeHeight=0 stateHeight=0
update Validators                            module=main len=4
Completed ABCI Handshake - Tendermint and App are synced module=consensus appHeight=0 appHash=
This node is not a validator                 module=consensus addr=666A495A6B05C975B241880785665417B5CEA2A6 pubKey=PubKeyEd25519{36BA673E7CC36F09C353720441C439A96E81B54689BAC219F0D24C52C3D23E65}
Starting Node                                module=node impl=Node
Starting EventBus                            module=events impl=EventBus
Local listener                               module=p2p ip=0.0.0.0 port=26656
Starting DefaultListener                     module=p2p impl=Listener(@172.31.230.212:26656)
P2P Node ID                                  module=node ID=db49a8d5a902910e0f8aee19e1b4889d6a235a91 file=/root/.qosd/config/node_key.json
Add our address to book                      module=p2p book=/root/.qosd/config/addrbook.json addr=db49a8d5a902910e0f8aee19e1b4889d6a235a91@172.31.230.212:26656
Starting RPC HTTP server on tcp://0.0.0.0:26657 module=rpc-server 
Starting P2P Switch                          module=p2p impl="P2P Switch"
Starting EvidenceReactor                     module=evidence impl=EvidenceReactor
Starting PEXReactor                          module=p2p impl=PEXReactor
Starting AddrBook                            module=p2p book=/root/.qosd/config/addrbook.json impl=AddrBook
Starting MempoolReactor                      module=mempool impl=MempoolReactor
Starting BlockchainReactor                   module=blockchain impl=BlockchainReactor
Starting BlockPool                           module=blockchain impl=BlockPool
Starting ConsensusReactor                    module=consensus impl=ConsensusReactor
ConsensusReactor                             module=consensus fastSync=true
Saving AddrBook to file                      module=p2p book=/root/.qosd/config/addrbook.json size=1
Starting IndexerService                      module=txindex impl=IndexerService
Ensure peers                                 module=p2p numOutPeers=0 numInPeers=0 numDialing=0 numToDial=10
Will dial address                            module=p2p addr=5d9fcba29ce9a066cdd6e4c45001567a4bd1dbf4@47.100.231.9:26656
Dialing peer                                 module=p2p address=5d9fcba29ce9a066cdd6e4c45001567a4bd1dbf4@47.100.231.9:26656
Successful handshake with peer               module=p2p peer=47.100.231.9:26656 peerNodeInfo="NodeInfo{id: 5d9fcba29ce9a066cdd6e4c45001567a4bd1dbf4, moniker: qos0, network: capricorn-1000 [listen 172.19.222.64:26656], version: 0.23.1 ([amino_version=0.12.0 p2p_version=0.5.0 consensus_version=v1/0.2.2 rpc_version=0.7.0/3 tx_index=on rpc_addr=tcp://0.0.0.0:26657])}"
Starting Peer                                module=p2p peer=47.100.231.9:26656 impl="Peer{MConn{47.100.231.9:26656} 5d9fcba29ce9a066cdd6e4c45001567a4bd1dbf4 out}"
Starting MConnection                         module=p2p peer=47.100.231.9:26656 impl=MConn{47.100.231.9:26656}
Added peer                                   module=p2p peer="Peer{MConn{47.100.231.9:26656} 5d9fcba29ce9a066cdd6e4c45001567a4bd1dbf4 out}"
Dialing peer                                 module=p2p address=5d9fcba29ce9a066cdd6e4c45001567a4bd1dbf4@47.100.231.9:26656
update Validators                            module=main len=0
Executed block                               module=state height=1 validTxs=0 invalidTxs=0
Committed state                              module=state height=1 txs=0 appHash=E5C9EABCC5C3ACB7EA6D8ED4D17B997BFCDD6F4F
Recheck txs                                  module=mempool numtxs=0 height=1
Indexed block                                module=txindex height=1
mint reward                                  module=main predict=8085999 actual=8085999
validatorVoteInfo                            module=main height=2 address1nfsgxj0l4gtgje0ydmjg6harsfvmduxtq8fdwa="not vote"
update Validators                            module=main len=0
Executed block                               module=state height=2 validTxs=0 invalidTxs=0
Committed state                              module=state height=2 txs=0 appHash=F70CB6559B9DA8015A63547696DC011032B7161F
Recheck txs                                  module=mempool numtxs=0 height=2
Indexed block                                module=txindex height=2
mint reward                                  module=main predict=8085999 actual=8085999
validatorVoteInfo                            module=main height=3 address1nfsgxj0l4gtgje0ydmjg6harsfvmduxtq8fdwa="not vote"
...
```

可运行下面命令检查节点运行状态：

```bash
$ qoscli tendermint status
```

如果看到`catching_up`为`false`，说明节点已经同步完成，否则还在同步区块。
同步完成后，可参照[成为验证节点](validator.md)引导成为对应网络验证节点。