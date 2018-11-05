# <font color=#0099ff>概要</font>
公链提供几种基础的 tx，每种 tx 有自己的结构和执行逻辑：  

    1，创建 QSC，   结构：TxCreateQSC  
    2，发币，       结构：TxIssueQsc  
    3，普通交易，   结构：TxTransform  
    4，Qcp执行结果，结构：QcpTxReasult
![Transaction_internal](https://github.com/QOSGroup/static/blob/master/transactionstruct_concreate.jpg?raw=true)

# <font color=#0099ff>TxCreateQSC</font>
功能：创建 QSC;  

发起用户（CreateAddr，需有一定数量的 qos）向CA申请，CA返回的信息含发起用户需要的QscName,QscPubKey,Banker publicKey信息，发起用户即可发起tx（tx中可向一组账户（AccInit）初始化一定数量的 qsc，若AccInit (类型[]AddrCoin) 中账户不存在则创建）

	[注]
	CreateAddr： QSC创建账户（即发起用户）
	QscName,QscPubKey,Banker publicKey：从CA中获取
	AccInit：    AddrCoin结构的数组，账户若不存在，创建；
	QscName：    长度为3~10个字符（数字+字母+下划线，区分大小写）  
```
type AddrCoin struct{
	Address		Address     //账户地址
	Amount		int64		//金额
}
```

# <font color=#0099ff>TxIssueQsc</font>
功能：发币；  

	QscName：Qsc名字
	Amount： 发币的数量

通过QscName，从TxCreateQSC执行后的state中查询到Banker，Banker向自己发行数量为Amount的qsc（签名者为Banker）。
若要向除Banker的其他人发币，则需再执行 TxTransform（Sender:Banker, Receive:接受者）
# <font color=#0099ff>TxTransform</font>
功能：转账交易；
	Senders：  转账人，AddrTrans结构
	Receivers：收款人，AddrTrans结构

Senders, Receivers为多对多（M <-> N）的关系，
转账处理逻辑要验证交易过程中，qsc的总数不变，
即： Senders的转出qsc总数 = Receivers的接受qsc总数  
签名过程：每个 Sender 都需对TxTransform整个结构的[]byte做签名，签名后按序存于TxStd.Signature（[]byte）中。