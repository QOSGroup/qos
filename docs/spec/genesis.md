# Genesis

执行`qosd init`之后，默认会在`$HOME/.qosd/config/`目录下生成`genesis.json`文件

```bash
cat genesis.json
{
    ...
    "app_state":{
        // 添加GenesisState结构对应数据内容
    }
}
```

`app_state`对应数据结构`GenesisState`：
```go
type GenesisState struct {
	CAPubKey   crypto.PubKey         `json:"ca_pub_key"`
	Accounts   []*account.QOSAccount `json:"accounts"`
	Validators []types.Validator     `json:"validators"`

	SPOConfig   types.SPOConfig   `json:"spo_config"`
	StakeConfig types.StakeConfig `json:"stake_config"`
}
```

处理逻辑：

在qos app 创建时编写相应代码逻辑可将该部分信息保存到创世块中
```go
// 设置 InitChainer
app.SetInitChainer(app.initChainer)

// 初始配置
func (app *QOSApp) initChainer(ctx context.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	
	//初始化app_state内容
	
	return abci.ResponseInitChain{}
}
...
```

## ca_pub_key

[QOS CA](ca.md)根证书公钥，用于联盟币、联盟链相关业务逻辑签名认证。

> 可通过`qosd config-root-ca <path of root.pub>`配置，`root.pub`为[go-amino](https://github.com/tendermint/go-amino)JSON 序列化ed25519编码的公钥信息。

执行
```bash
$ qosd config-root-ca root.pub
```
会在`genesis.json`中添加如下部分：
```bash
"ca_pub_key": {
    "type": "tendermint/PubKeyEd25519",
    "value": "0SDDvhiMsqX9XLuscqovU8l24txbV7Mg4ecs+R6Swzk="
}
```

## accounts

初始账户

> 初始账户可以通过`qosd add-genesis-accounts`命令添加

执行
```bash
$ qosd add-genesis-accounts address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy,100000000QOS,100000000AOE
```
会在`genesis.json`中`accounts`部分添加对应账户信息：
```bash
"accounts": [
    {
        "base_account": {
            "account_address": "address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy"
        },
        "qos": "100000000",
        "qscs": [
            {
                "coin_name": "AOE",
                "amount": "100000000"
            }
        ]
    }
]
```

## validators

QOS使用app_state中包含的validators作为验证节点，启动前务必正确配置validators。

> 验证节点可通过`qosd add-genesis-validator`命令添加。

执行：
```bash
$ qosd add-genesis-validator --name "Arya's node" --owner address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy --tokens 100
```
会在`genesis.json`中`validators`部分添加对应validator信息：
```bash
"validators": [                                                           
  {                                                                       
    "name": "Arya's node",                                                      
    "owner": "address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy",            
    "validatorPubkey": {                                                  
      "type": "tendermint/PubKeyEd25519",                                 
      "value": "agD9zt3RhmAq6YnF7UKn0Kw53wQvZsiPmRYdG+dyaDk=""             
    },                                                                    
    "bondTokens": "100",                                                  
    "description": "",                                                    
    "status": 0,                                                          
    "inActiveCode": 0,                                                    
    "inActiveTime": "0001-01-01T00:00:00Z",                               
    "inActiveHeight": "0",                                                
    "bondHeight": "1"                                                     
  }                                                                       
]                                                                        
```