# 加入QOS测试网络

QOS测试网络以十二星座命名，我们在*2018年12月26日*启动第一个测试网络，以摩羯座命名：`capricorn-1000`。
最新的测试网信息请查看[qos-testnets](https://github.com/QOSGroup/qos-testnets)。

**视频教程**(v0.0.4参考视频，最新版有些许差异)：
- [油管](https://youtu.be/-eFdx0rIPb4)
- [优酷](http://v.youku.com/v_show/id_XNDA4NTA1MDM1Ng==.html)

按下面步骤，可加入我们的测试网络：

## 安装 QOS

参照[安装引导](http://docs.qoschain.info/qos/install/testnet/installation.html)
和[qos-testnets](https://github.com/QOSGroup/qos-testnets)安装对应版本的QOS。

## 启动全节点

### 升级
在运行过其他版本QOS测试网络的机器上参与新的测试网络，请务必**停掉之前网络并删除运行数据**。

默认`$HOME/.qosd`目录下文件：

```bash
rm -rf $HOME/.qosd
```

然后执行下面的[初始化](#初始化)和[启动](#启动)操作。

### 初始化

适用`qosd init`命令初始化节点、创建必要的配置文件。
默认的配置和数据存储目录为 `$HOME/.qosd`，可以添加`--home`修改存储位置。

```bash
$ qosd init --moniker <your_custom_moniker>
```

执行完`qosd init`会在`$HOME/.qosd/config`下生成`genesis.json`、`config.toml`等配置文件。

### 配置运行网络

不同QOS运行网络对应不同的配置，可访问[testnets repo](https://github.com/QOSGroup/qos-testnets)了解不同网络的运行配置。

下面操作以**最新测试网**为例。

#### 替换`genesis.json`

默认路径`$HOME/.qosd/config/genesis.json`

下载[`genesis.json`文件](https://raw.githubusercontent.com/QOSGroup/qos-testnets/master/latest/genesis.json)替换本地文件。
若没有更改默认存储位置，也可通过下面的命令执行替换：
```bash
$ curl https://raw.githubusercontent.com/QOSGroup/qos-testnets/master/latest/genesis.json > $HOME/.qosd/config/genesis.json
```

#### 编辑`config.toml`:

默认路径`$HOME/.qosd/config/config.toml`

修改`config.toml`，找到`seeds`配置项，添加seed节点：
```toml
# Comma separated list of seed nodes to connect to
seeds = "1233b1c5bad7561d7c5a28b4a2149760a8b673d2@47.103.79.28:26656"
```

### 启动

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
$ qoscli query status
```

如果看到`catching_up`为`false`，说明节点已经同步完成，否则还在同步区块。

## 成为验证节点

在成为验证节点前，请确保已[启动全节点](#启动全节点)，同步到最新高度：
```bash
$ qoscli query status
```
其中`latest_block_height`为已同步高度，`catching_up`如果为`false`表示已同步到最新高度，否则请等待。可通过[区块链浏览器](http://explorer.qoschain.info/block/list)查看最新高度，及时了解同步进度。

成为验证节点前可查阅[QOS验证人详解](../spec/validators/all_about_validators.md)和[QOS经济模型](../spec/validators/eco_module.md)了解验证人的相关运行机制，然后执行**获取Token**和**成为验证节点**相关步骤成为验证节点。

### 获取Token

成为验证节点需要账户作为操作者，与节点绑定。如果还没有用于操作的账户，可通过下面步骤创建。
同时成为验证人需要操作者持有一定量的Token，测试网络可通过[水龙头](http://explorer.qoschain.info/freecoin/get)免费获取。

1. 创建操作账户

可通过`qoscli keys add <name_of_key>`创建密钥保存在本地密钥库。`<name_of_key>`可自定义，仅作为本地密钥库存储名字。
下面以创建Peter账户为例：
```bash
$ qoscli keys add Peter
// 输入不少于8位的密码，请牢记密码信息
Enter a passphrase for your key: 
Repeat the passphrase: 
```
会输出如下信息：
```bash
NAME:   TYPE:   ADDRESS:                                                PUBKEY:
Peter local   address1epvxmtxx99gy5xv7k7sl55994pehxgqt03va2s  D+pHqEJVjQMiRzl5PbL8FraVZqWqxrxcTF7akcCIDfo=
**Important** write this seed phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

broom resource trash summer crop embrace stadium fish brief dolphin run decrease brief heart upgrade icon toe lift dawn regret dumb indoor drop glide
```
其中`address1epvxmtxx99gy5xv7k7sl55994pehxgqt03va2s`为QOS账户地址，
`D+pHqEJVjQMiRzl5PbL8FraVZqWqxrxcTF7akcCIDfo=`为公钥信息，
`broom resource trash summer crop embrace stadium fish brief dolphin run decrease brief heart upgrade icon toe lift dawn regret dumb indoor drop glide`为助记词，可用于账号恢复，***请牢记助记词***。

更多密钥库相关操作执行请执行`qoscli keys --help`查阅。

2. 获取QOS

测试网络的QOS可访问[水龙头](http://explorer.qoschain.info/freecoin/get)免费获取。

::: warning Note 
从水龙头获取的QOS仅可用于QOS测试网络使用
:::

可通过下面的指令查询Peter账户信息：
```bash
$ qoscli query account Peter --indent
```

会看到类似如下信息：
```bash
{
  "type": "qbase/account/QOSAccount",
  "value": {
    "base_account": {
      "account_address": "address1epvxmtxx99gy5xv7k7sl55994pehxgqt03va2s",
      "public_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "D+pHqEJVjQMiRzl5PbL8FraVZqWqxrxcTF7akcCIDfo="
      },
      "nonce": "0"
    },
    "qos": "10000000",
    "qscs": null
  }
}
```
其中`qos`值代表持有的QOS量，大于0就可以按照接下来的操作成为验证节点。

### 成为验证节点

1. 执行创建命令

创建验证节点需要**节点公钥**、**操作者账户地址**等信息

使用下面的命令执行创建验证节点：
```
qoscli tx create-validator --owner Peter --moniker "Peter's node" --tokens 20000000
```
其中：
- `--owner` 操作者密钥库名字或账户地址，如使用之前创建的`Peter`或`Peter`对应的地址`address1epvxmtxx99gy5xv7k7sl55994pehxgqt03va2s`
- `--moniker`给验证节点起个名字，如`Peter's node`
- `--logo`   logo
- `--website` 网址
- `--details`  详细描述信息
- `--nodeHome` 节点配置文件和数据所在目录，默认：`$HOME/.qosd`
- `--tokens` 将绑定到验证节点上的Token量，应小于等于操作者持有的QOS量
- `--commission-rate` 佣金比例，默认值`0.1`
- `--commission-max-rate` 最高佣金比例，默认值`0.2`
- `--commission-max-change-rate` 24小时内佣金最大变化范围，默认`0.01`

会输出出类似如下信息：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"34A76D6D07D93FBE395DDC55E0596E4D312A02A9","height":"200"}
```

2. 查看节点信息

成功执行创建操作，可通过`qoscli query validator --owner <owner_address_of_validator>`指令查询验证节点信息，其中`owner_address_of_validator`为操作者账户地址。
```bash
qoscli query validator --owner address1epvxmtxx99gy5xv7k7sl55994pehxgqt03va2s
```
会输出出类似如下信息：
```bash
{
  "name": "Peter's node",
  "owner": "address1epvxmtxx99gy5xv7k7sl55994pehxgqt03va2s",
  "validatorPubkey": {
    "type": "tendermint/PubKeyEd25519",
    "value": "PJ58L4OuZp20opx2YhnMhkcTzdEWI+UayicuckdKaTo="
  },
  "bondTokens": "20000000",
  "description": "",
  "status": 0,
  "inactiveCode": 0,
  "inactiveTime": "0001-01-01T00:00:00Z",
  "inactiveHeight": "0",
  "bondHeight": "200"
}
```
其中`status`为0说明已成功成为验证节点，将参与相关网络的打块和投票任务。


### QOS区块链浏览器

成为验证节点后，如果绑定的Token在对应网络中排在前100名，可通过[区块链浏览器](http://explorer.qoschain.info/validator/list)查看节点信息。