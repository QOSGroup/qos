# 成为验证节点

在成为验证节点前，请确保已[启动完整节点](fullnode.md)，同步到最新高度：
```bash
$ qoscli tendermint status
```
其中`latest_block_height`为已同步高度，`catching_up`如果为`false`表示已同步到最新高度，否则请等待。可通过[区块链浏览器](http://explorer.qoschain.info/block/list)查看最新高度，及时了解同步进度。

成为验证节点前可查阅[QOS验证人详解](../../spec/validators/all_about_validators.md)和[QOS经济模型](../../spec/validators/eco_module.md)了解验证人的相关运行机制，然后执行**获取Token**和**成为验证节点**相关步骤成为验证节点。

## 获取Token

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

## 成为验证节点

1. 执行创建命令

创建验证节点需要**节点公钥**、**操作者账户地址**等信息

使用下面的命令执行创建验证节点：
```
qoscli tx create-validator --owner Peter --name "Peter's node" --tokens 20000000 --description "hi, my eth address: xxxxxx"
```
其中：
- `--owner` 操作者密钥库名字或账户地址，如使用之前创建的`Peter`或`Peter`对应的地址`address1epvxmtxx99gy5xv7k7sl55994pehxgqt03va2s`
- `--name`  给验证节点起个名字，如`Peter's node`
- `--nodeHome` 节点配置文件和数据所在目录，默认：`$HOME/.qosd`
- `--tokens` 将绑定到验证节点上的Token量，应小于等于操作者持有的QOS量
- `--description` 备注信息，可放置**以太坊地址**，方便接收QOS奖励

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