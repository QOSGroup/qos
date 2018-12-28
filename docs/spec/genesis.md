# Genesis

执行qosd init之后，默认会在$HOME/.qosd/config/目录下生成genesis.json配置文件

```bash
cat genesis.json
{
    ...
    "app_state":{
        // 添加GenesisState结构对应数据内容
    }
}
```

app_state对应数据结构
```go
type GenesisState struct {
	CAPubKey   crypto.PubKey         `json:"ca_pub_key"`
	Accounts   []*account.QOSAccount `json:"accounts"`
	Validators []types.Validator     `json:"validators"`

	SPOConfig   types.SPOConfig   `json:"spo_config"`
	StakeConfig types.StakeConfig `json:"stake_config"`
}
```

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

Root [CA](ca.md) PubKey，用于联盟链签名认证等

> 可通过`qosd config-root-ca <path of root.pub>`配置。

```
"ca_pub_key": {
      "type": "tendermint/PubKeyEd25519",
      "value": "0SDDvhiMsqX9XLuscqovU8l24txbV7Mg4ecs+R6Swzk="
    }
```

## accounts

初始账户

> 初始账户可以通过`qosd init add-genesis-account`命令添加

```
"accounts": [
    {
        "base_account": {
            "account_address": "address1k0m8ucnqug974maa6g36zw7g2wvfd4sug6uxay"
        },
        "qos": "100000000",
        "qscs": [
            {
                "coin_name": "qstar",
                "amount": "100000000"
            }
        ]
    }
]
```

## validators

QOS使用app_state中包含的validators作为验证节点，启动前务必正确配置validators。

> 验证节点可通过`qosd add-genesis-validator`命令添加

```
"validators": [                                                           
  {                                                                       
    "name": "imuge",                                                      
    "owner": "address1gmllq4fgtlfe574dffaj90t3tkvy232phcukhq",            
    "validatorPubkey": {                                                  
      "type": "tendermint/PubKeyEd25519",                                 
      "value": "GXs+JzWZEorSt4eN2/UzSZPRfcAPvT67A9Lkn5QufaA="             
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