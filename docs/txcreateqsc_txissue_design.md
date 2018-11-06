# TxCreateQSC设计
## 功能
创建联盟链  
公链中拥有一定数量QOS的账户，即可发起此Tx.
## 结构
    type TxCreateQSC struct {
        QscName     string         `json:"qscname"`     //从CA信息获取
        CreateAddr  btypes.Address `json:"createaddr"`  //QSC创建账户
        QscPubkey   crypto.PubKey  `json:"qscpubkey"`   //从CA信息获取
        Banker      btypes.Address `json:"banker"`      //从CA信息获取
        Extrate     string         `json:"extrate"`     //qcs:qos汇率(amino不支持binary形式的浮点数序列化，精度同qos erc20 [.0000])
        CAqsc       []byte         `json:"caqsc"`       //qsc的CA信息
        CAbanker    []byte         `json:"cabanker"`    //banker的CA信息
        Description string         `json:"description"` //描述信息
        AccInit     []AddrCoin     `json:"accinit"`     //初始化时接受qsc的账户
    }
## 执行
创建联盟链的核心工作：在公链上，保存联盟链信息（名称，publickey）及banker信息。
为保证创建TxCreateQSC数据的正确性，关键信息(如：联盟链名称，publickey, banker信息等)使用CA证书（CAqsc, CAbanker）传递。

TxCreateQSC实现Qbase::ITx接口，执行过程中，首先会校验CA证书的正确性，而后从CA证书中提取信息，校验Tx中字段是否正确。
完成校验，即可保存联盟链及banker信息；
QOS基于Qbase，故存储过程使用Qbase提供的mapper机制。
Tx的关键信息，存于mapper["base"]，storekey:"qsc/[联盟链名称]"
## 测试
参考：[txcreateqsc_txissue_test.md](https://github.com/QOSGroup/qos/tree/master/docs/txcreateqsc_txissue_test.md)  

# TxIssue设计
## 功能
向联盟链的Banker发币  
发币人其实是联盟链的Banker，即联盟链的Banker向自己发币。
## 结构
    type TxIssueQsc struct {
        QscName string         `json:"qscName"` //发币账户名
        Amount  btypes.BigInt  `json:"amount"`  //金额
        Banker  btypes.Address `json:"banker"`  //banker地址
    }
## 执行
TxIssue的核心工作：联盟链的banker向自己发币.
发币数量的限制机制和Gas机制，此版本暂未给出，待完善。

TxIssue实现Qbase::ITx接口，执行过程中，首先会校验Banker信息的正确性，若无问题，即会向banker账户存入一定数量的联盟链代币。
联盟链的Banker信息，可在DB中查询（mapper["base"]，storekey:"qsc/[联盟链名称]"）
## 测试
参考：[txcreateqsc_txissue_test.md](https://github.com/QOSGroup/qos/tree/master/docs/txcreateqsc_txissue_test.md)