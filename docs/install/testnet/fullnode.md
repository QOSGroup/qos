# 启动完整节点

在启动完整节点前，请确保已按照[安装引导](../installation.md)正确安装QOS。

## 升级
在运行过其他版本QOS测试网络的机器上参与新的测试网络，请务必**停掉之前网络并删除运行数据**。

默认`$HOME/.qosd`目录下文件：

```bash
rm -rf $HOME/.qosd
```

然后执行下面的[初始化](#初始化)和[启动](#启动)操作。

## 初始化

适用`qosd init`命令初始化节点、创建必要的配置文件。
默认的配置和数据存储目录为 `$HOME/.qosd`，可以添加`--home`修改存储位置。

```bash
$ qosd init --moniker <your_custom_moniker>
```

执行完`qosd init`会在`$HOME/.qosd/config`下生成`genesis.json`、`config.toml`等配置文件。

## 配置运行网络

不同QOS运行网络对应不同的配置，可访问[testnets repo](https://github.com/QOSGroup/qos-testnets)了解不同网络的运行配置。

下面操作以**最新测试网**为例。

### 替换`genesis.json`

默认路径`$HOME/.qosd/config/genesis.json`

下载[`genesis.json`文件](https://raw.githubusercontent.com/QOSGroup/qos-testnets/master/latest/genesis.json)替换本地文件。
若没有更改默认存储位置，也可通过下面的命令执行替换：
```bash
$ curl https://raw.githubusercontent.com/QOSGroup/qos-testnets/master/latest/genesis.json > $HOME/.qosd/config/genesis.json
```

### 编辑`config.toml`:

默认路径`$HOME/.qosd/config/config.toml`

修改`config.toml`，找到`seeds`配置项，添加seed节点：
```toml
# Comma separated list of seed nodes to connect to
seeds = "f1dbd6d0b931fe7f918a81e8248c21e2109caa97@47.105.156.172:26656"
```

## 启动

运行启动命令：

```bash
$ qosd start
```

`--log_level debug`会打印很多debug日志，由于日志文件会很大，可以不添加`--log_level`参数，用默认的日志配置。

控制台开始打印启动日志，提示`This node is not a validator`说明节点不是验证节点，节点开始同步QOS网络区块信息。
```bash
I[26026-02-26|17:19:55.657] Starting ABCI with Tendermint                module=main 
I[26026-02-26|17:19:55.752] Starting multiAppConn                        module=proxy impl=multiAppConn
I[26026-02-26|17:19:55.752] Starting localClient                         module=abci-client connection=query impl=localClient
I[26026-02-26|17:19:55.752] Starting localClient                         module=abci-client connection=mempool impl=localClient
I[26026-02-26|17:19:55.752] Starting localClient                         module=abci-client connection=consensus impl=localClient
I[26026-02-26|17:19:55.752] ABCI Handshake App Info                      module=consensus height=0 hash= software-version= protocol-version=0
I[26026-02-26|17:19:55.759] ABCI Replay Blocks                           module=consensus appHeight=0 storeHeight=0 stateHeight=0
I[26026-02-26|17:19:55.761] update Validators                            module=main len=2
I[26026-02-26|17:19:55.768] Completed ABCI Handshake - Tendermint and App are synced module=consensus appHeight=0 appHash=
I[26026-02-26|17:19:55.768] Version info                                 module=node software=0.27.3 block=8 p2p=5
I[26026-02-26|17:19:55.768] This node is not a validator                 module=consensus addr=5616023310A4FE28C1138C5780F7F9CBBE997AE1 pubKey=PubKeyEd25519{54EE7F393278E40B4A22159890B6B3EA6076BBE284E3E59B7CBCC2795D65E56B}
I[26026-02-26|17:19:55.820] P2P Node ID                                  module=p2p ID=dd31bf449fdca95236bf54e9e3a216d27cdef7e0 file=/home/imuge/.qosd/config/node_key.json
I[26026-02-26|17:19:55.820] Starting Node                                module=node impl=Node
I[26026-02-26|17:19:55.820] Starting EventBus                            module=events impl=EventBus
I[26026-02-26|17:19:55.820] Starting PubSub                              module=pubsub impl=PubSub
I[26026-02-26|17:19:55.821] Starting P2P Switch                          module=p2p impl="P2P Switch"
I[26026-02-26|17:19:55.821] Starting MempoolReactor                      module=mempool impl=MempoolReactor
I[26026-02-26|17:19:55.821] Starting BlockchainReactor                   module=blockchain impl=BlockchainReactor
I[26026-02-26|17:19:55.821] Starting RPC HTTP server on [::]:26657       module=rpc-server 
I[26026-02-26|17:19:55.821] Starting BlockPool                           module=blockchain impl=BlockPool
I[26026-02-26|17:19:55.821] Starting ConsensusReactor                    module=consensus impl=ConsensusReactor
I[26026-02-26|17:19:55.821] ConsensusReactor                             module=consensus fastSync=true
I[26026-02-26|17:19:55.821] Starting EvidenceReactor                     module=evidence impl=EvidenceReactor
I[26026-02-26|17:19:55.821] Starting PEXReactor                          module=pex impl=PEXReactor
I[26026-02-26|17:19:55.821] Starting AddrBook                            module=p2p book=/home/imuge/.qosd/config/addrbook.json impl=AddrBook
I[26026-02-26|17:19:55.821] Starting IndexerService                      module=txindex impl=IndexerService
I[26026-02-26|17:19:55.821] Ensure peers                                 module=pex numOutPeers=0 numInPeers=0 numDialing=0 numToDial=10
I[26026-02-26|17:19:55.821] No addresses to dial nor connected peers. Falling back to seeds module=pex 
I[26026-02-26|17:19:55.821] Dialing peer                                 module=p2p address=f1dbd6d0b931fe7f918a81e8248c21e2109caa97@47.105.156.172:26656
I[26026-02-26|17:19:55.890] Starting Peer                                module=p2p peer=f1dbd6d0b931fe7f918a81e8248c21e2109caa97@0.0.0.0:26656 impl="Peer{MConn{47.105.156.172:26656} f1dbd6d0b931fe7f918a81e8248c21e2109caa97 out}"
I[26026-02-26|17:19:55.890] Starting MConnection                         module=p2p peer=f1dbd6d0b931fe7f918a81e8248c21e2109caa97@0.0.0.0:26656 impl=MConn{47.105.156.172:26656}
D[26026-02-26|17:19:55.890] Send                                         module=p2p peer=f1dbd6d0b931fe7f918a81e8248c21e2109caa97@0.0.0.0:26656 channel=64 conn=MConn{47.105.156.172:26656} msgBytes=5A433AB9
D[26026-02-26|17:19:55.891] Request addrs                                module=pex from="Peer{MConn{47.105.156.172:26656} f1dbd6d0b931fe7f918a81e8248c21e2109caa97 out}"
D[26026-02-26|17:19:55.891] Send                                         module=p2p peer=f1dbd6d0b931fe7f918a81e8248c21e2109caa97@0.0.0.0:26656 channel=0 conn=MConn{47.105.156.172:26656} msgBytes=723A31CD
I[26026-02-26|17:19:55.891] Added peer                                   module=p2p peer="Peer{MConn{47.105.156.172:26656} f1dbd6d0b931fe7f918a81e8248c21e2109caa97 out}"
D[26026-02-26|17:19:55.891] No votes to send, sleeping                   module=consensus peer="Peer{MConn{47.105.156.172:26656} f1dbd6d0b931fe7f918a81e8248c21e2109caa97 out}" rs.Height=1 prs.Height=0 localPV=BA{2:__} peerPV=nil-BitArray localPC=BA{2:__} peerPC=nil-BitArray
D[26026-02-26|17:19:55.991] Flush                                        module=p2p peer=f1dbd6d0b931fe7f918a81e8248c21e2109caa97@0.0.0.0:26656 conn=MConn{47.105.156.172:26656}
D[26026-02-26|17:19:55.992] Read PacketMsg                               module=p2p peer=f1dbd6d0b931fe7f918a81e8248c21e2109caa97@0.0.0.0:26656 conn=MConn{47.105.156.172:26656} packet="PacketMsg{40:5A433AB908D902 T:1}"
D[26026-02-26|17:19:55.992] Received bytes                               module=p2p peer=f1dbd6d0b931fe7f918a81e8248c21e2109caa97@0.0.0.0:26656 chID=64 msgBytes=5A433AB908D902
D[26026-02-26|17:19:55.992] Receive                                      module=blockchain src="Peer{MConn{47.105.156.172:26656} f1dbd6d0b931fe7f918a81e8248c21e2109caa97 out}" chID=64 msg="[bcStatusResponseMessage 345]"
D[26026-02-26|17:19:55.992] Read PacketMsg                               module=p2p peer=f1dbd6d0b931fe7f918a81e8248c21e2109caa97@0.0.0.0:26656 conn=MConn{47.105.156.172:26656} packet="PacketMsg{20:C96A6FA808DA021804 T:1}"
D[26026-02-26|17:19:55.992] Received bytes                               module=p2p peer=f1dbd6d0b931fe7f918a81e8248c21e2109caa97@0.0.0.0:26656 chID=32 msgBytes=C96A6FA808DA021804
D[26026-02-26|17:19:55.992] Receive                                      module=consensus src="Peer{MConn{47.105.156.172:26656} f1dbd6d0b931fe7f918a81e8248c21e2109caa97 out}" chId=32 msg="[NewRoundStep H:346 R:0 S:RoundStepPrevote LCR:0]"
D[26026-02-26|17:19:55.992] TrySend                                      module=p2p peer=f1dbd6d0b931fe7f918a81e8248c21e2109caa97@0.0.0.0:26656 channel=64 conn=MConn{47.105.156.172:26656} msgBytes=BB1DC4F20805
D[26026-02-26|17:19:55.992] Read PacketMsg                               module=p2p peer=f1dbd6d0b931fe7f918a81e8248c21e2109caa97@0.0.0.0:26656 conn=MConn{47.105.156.172:26656} packet="PacketMsg{20:1919B3D508DA0218012001 T:1}"
D[26026-02-26|17:19:55.992] Received bytes                               module=p2p peer=f1dbd6d0b931fe7f918a81e8248c21e2109caa97@0.0.0.0:26656 chID=32 msgBytes=1919B3D508DA0218012001
D[26026-02-26|17:19:55.992] Receive                                      module=consensus src="Peer{MConn{47.105.156.172:26656} f1dbd6d0b931fe7f918a81e8248c21e2109caa97 out}" chId=32 msg="[HasVote VI:1 V:{346/00/1}]"
D[26026-02-26|17:19:55.992] setHasVote                                   module=consensus peerH/R=346/0 H/R=346/0 type=1 index=1
D[26026-02-26|17:19:55.992] Read PacketMsg                               module=p2p peer=f1dbd6d0b931fe7f918a81e8248c21e2109caa97@0.0.0.0:26656 conn=MConn{47.105.156.172:26656} packet="PacketMsg{20:87E347CB08DA021A2408011220B3451F0BDBAE0C469B1BA5BA1963B9DA60ADA5216F30D65559BD4C4599B790B022050801120101 T:1}"
...
```

可运行下面命令检查节点运行状态：

```bash
$ qoscli tendermint status
```

如果看到`catching_up`为`false`，说明节点已经同步完成，否则还在同步区块。
同步完成后，可参照[成为验证节点](validator.md)引导成为对应网络验证节点。