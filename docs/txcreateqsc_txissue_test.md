# QOS测试( TxCreateQSC, TxIssue )

## 描述   
cli端：   qos\cmd\qosappcli.go  
qosd：    qos\cmd\qosd\main.go  
内置账户： qos\test\accountdefault.go中初始化5个账户，qos启动时加载  
CA证书：   qsc.crt, banker.crt为联盟链证书和banker证书.
## 编译  
命令行进入目录qos\cmd\qosd，编译得到qosd  

    qod\cmd\qosd> go build  

命令行进入目录qos\cmd，编译得到cli端可执行文件qosappcli.exe(windows平台为例 qosappcli.exe)  

    qos\cmd > go build qosappcli.go  
## 步骤
### 1, qos初始化(qosd)  
    参数: qosd init --chain-id=qos  

### 2, qos启动(qosd)
    参数： qosd start --with-tendermint=true  
### 3, 发送TxCreateQSC(cli端)
    参数：	qosappcli.exe -m=txcreateqsc -pathqsc=d:\qsc.crt -pathbank=d:\banker.crt -chainid=qos -maxgas=100 -nonce=1  
说明：  
    构建此tx需要CA证书 ( qsc.crt, banker.crt )  
    3.1, pathqsc & pathbank 分别为qsc和banker的CA文件路径  
    3.2, example: d:\banker.crt  
    3.3, 测试获取路径: github.com/QOSGroup/kepler/examples/v1  ( qsc.crt, banker.crt )  
### 4, 发送TxIssue(cli端)	
    参数： qosappcli.exe -m=txissue -qscname=QSC -nonce=1 -chainid=qos -maxgas=100 -amount=20000  
说明：  
4.1, qscname需和banker.crt中的CSR.Subj.CN相同，区分大小写  

## 信息查询
查询账户信息(步骤2,3,4之后都可以执行查询账户信息，验证tx结果)  

### cli端查询banker  
    qosappcli.exe -m=accquery -addr=address1l7d3dc26adk9gwzp777s3a9p5tprn7m43p99cg  
### cli端查询acc1
    qosappcli.exe -m=accquery -addr=address1zsqzn6wdecyar6c6nzem3e8qss2ws95csr8d0r