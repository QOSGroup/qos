# QOS测试( TxCreateQSC, TxIssue )

## 描述   
cli端：   qos\cmd\qosappcli.go  
qosd：    qos\cmd\qosd\main.go  
内置账户： qos\test\accountdefault.go中初始化creator账户，qos启动时加载  
CA证书：   qsc.crt, banker.crt为联盟链证书和banker证书.
## 编译  
命令行进入目录qos\cmd\qosd，编译得到qosd.exe  

    qod\cmd\qosd> go build  

命令行进入目录qos\cmd，编译得到cli端可执行文件qosappcli.exe(windows平台为例 qosappcli.exe)  

    qos\cmd > go build qosappcli.go  
## 步骤
### 1, qos初始化(qosd)  
    qosd.exe init --chain-id=qos  

### 2, qos启动(qosd)
    qosd.exe start --with-tendermint=true  
### 3, 发送TxCreateQSC(cli端)
    qosappcli.exe -m=txcreateqsc -pathqsc=d:\qsc.crt -pathbank=d:\banker.crt -qscchainid=qcptest -qoschainid=qos -maxgas=100 -nonce=1 -privkeycreator=rDwWppdGKFCv0wUxFqVID87GI/CFwLbL9p6EM6ug5brPbkXQoZMIH9+Rgi1/vFcNJUHp88fKZDNFdEif8dg73A==  

	 	3.1, pathqsc/pathbank: qsc/banker的CA文件路径  
 			 github.com/QOSGroup/kepler/examples/v1  (qsc.crt, banker.crt)  
		3.2, qscchainid:  联盟链chainid  
		3.3, qoschainid:  公链chainid  
		3.4, maxgas: 期望最大gas花费  
		3.5, privkeycreator: creator的 private key.  
### 4, 发送TxIssue(cli端)	
    qosappcli.exe -m=txissue -qscname=QSC -nonce=1 -qoschainid=qos -maxgas=100 -privkeybank=maD8NeYMqx6fHWHCiJdkV4/B+tDXFIpY4LX4vhrdmAYIKC67z/lpRje4NAN6FpaMBWuIjhWcYeI5HxMh2nTOQg==  

		4.1, qscname需和banker中的qscname相同，区分大小写  
		4.2, qoschainid:  公链chainid  
		4.3, privkeybank: banker的privatekey  

## 信息查询
查询账户信息(步骤2,3,4之后都可以执行查询账户信息，验证tx结果)  

### 查询创建的联盟链信息
	qosappcli.exe -m=qscquery -qscname=QSC  

		qscname: 要查询的联盟链名称  

### cli端查询banker  
    qosappcli.exe -m=accquery -addr=address1l7d3dc26adk9gwzp777s3a9p5tprn7m43p99cg  
### cli端查询acc1,acc2,acc3  
    //下面的账户在txcreatqsc时创建，并随机分发一定数量的联盟链币
    qosappcli.exe -m=accquery -addr=address1zsqzn6wdecyar6c6nzem3e8qss2ws95csr8d0r
    qosappcli.exe -m=accquery -addr=address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355  
    qosappcli.exe -m=accquery -addr=address1y9r4pjjnvkmpvw46de8tmwunw4nx4qnz2ax5ux  
