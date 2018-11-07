# Genesis

## genesis.json
执行qosd init之后，默认会在$HOME/.qosd/config/目录下生成genesis.json配置文件

```
cat genesis.json
{
    ...
    "app_state":{
        // 添加自定义内容
    }
}
```

app_state部分可自定义，在qos app 创建时编写相应代码逻辑可将该部分信息保存到创世块中
```
// 设置 InitChainer
...
app.SetInitChainer(app.initChainer)
...
// 初始配置
func (app *QOSApp) initChainer(ctx context.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	// 上下文中获取mapper
	mainMapper := ctx.Mapper(mapper.GetMainStoreKey()).(*mapper.MainMapper)
	accountMapper := ctx.Mapper(account.AccountMapperName).(*account.AccountMapper)

	// 反序列化app_state
	stateJSON := req.AppStateBytes
	genesisState := &GenesisState{}
	err := accountMapper.GetCodec().UnmarshalJSON(stateJSON, genesisState)
	if err != nil {
		panic(err)
	}

	// 保存CA
	mainMapper.SetRootCA(genesisState.CAPubKey)

	// 保存初始账户
	for _, acc := range genesisState.Accounts {
		accountMapper.SetAccount(acc)
	}

	return abci.ResponseInitChain{}
}
...
```

## ca_pub_key
CA PubKey
```
"ca_pub_key": {
      "type": "tendermint/PubKeyEd25519",
      "value": "0SDDvhiMsqX9XLuscqovU8l24txbV7Mg4ecs+R6Swzk="
    }
```

## accounts
创世账户
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

## qcps
联盟链配置
```
"qcps": [
    {
        "name": "qstars-test",
        "chain_id": "qstars-test",
        "pub_key": {
            "type": "tendermint/PubKeyEd25519",
            "value": "O0XnpXPYkn4jWXwBPgG1wp1aCx0tuug9Ylc0WHBnl5Q="
        }
    }
]
```
联盟链的初始化可在qbase baseapp InitChain方法下执行，将解析的qcp信息保存到store中
```
func (app *BaseApp) InitChain(req abci.RequestInitChain) (res abci.ResponseInitChain) {
    ...
	// 保存初始QCP配置
	initQCP(app.deliverState.ctx, app.GetCdc(), req.AppStateBytes)
	...
}

func initQCP(ctx ctx.Context, cdc *go_amino.Codec, appState []byte) {
	if appState == nil {
		return
	}
	gs := types.GenesisState{}
	err := cdc.UnmarshalJSON(appState, &gs)
	if err != nil {
		panic(err)
	}
	if len(gs.QCPs) > 0 {
		qcpMapper := GetQcpMapper(ctx)
		for _, qcp := range gs.QCPs {
			qcpMapper.SetChainInTrustPubKey(qcp.ChainId, qcp.PubKey)
		}
	}
}
```
无需在qos app initChainer()中添加qcps处理逻辑