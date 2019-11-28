# API调用说明

**API根据类型分为三类:**

- 查询API: 通过`GET`发送请求, 查询链上各种模块信息
- 创建交易: 通过`POST`发送请求, 创建不同类型的交易
- 广播交易: 通过`POST /txs` 发送请求, 将组装好的签名交易广播到链上

**通过API发送交易流程如下:**

1. 先根据具体要创建交易的API创建交易
2. 在本地使用私钥对返回的待签名串进行签名
3. 本地组装交易签名数据: 组装签名, 组装公钥信息
4. 调用`POST /txs` 将签名交易发送至链上


**首次创建交易与非首次创建交易区别:**

首次交易时, 链上没有用户的公钥信息, 无法对传来的交易进行签名校验. 
因此首次交易时,需要在交易体中携带用户的公钥信息. 交易完成后,链中会记录用户公钥信息以便后续交易校验. 

首次创建交易时, 交易体中返回的`pubkey`为空,需要手动设置公钥, 而在后续创建交易时,该字段会包含链中保存的用户公钥信息.

用户展示的公钥格式为以`qosaccpub1`开头的`bech32`格式的公钥, 比如: `qosaccpub1d9rez3ulq36p4x39zx94t0lnhd0f8359uwwzrlre99t6mr97nc8s6l8vql`, 
而在交易体中需要设置为`ed25519`格式公钥, 二者的转换关系可以通过`GET /accounts/pubkey/{bech32PubKey}`转换.

```
GET /accounts/pubkey/qosaccpub1d9rez3ulq36p4x39zx94t0lnhd0f8359uwwzrlre99t6mr97nc8s6l8vql

> 

< 

    {
    "type": "tendermint/PubKeyEd25519",
    "value": "aUeRR58EdBqaJRGLVb/zu16TxoXjnCH8eSlXrYy+ng8="
    }
```


**创建交易请求基础参数**

基础参数格式:

```json

  "base": {
    "from": "", 
    "chain_id": "",
    "nonce": "1",
    "max_gas": "",
    "height": "",
    "indent": false,
    "mode": "block"
  }

```


|参数名|参数说明|类型|是否必填|默认值|
| ------ | ---- |----- |----- | ---- |
|from|交易发起人|string|是| - |
|chain_id|链名称|string|否|若为空,则由RPC Server设置该字段|
|nonce|账户nonce值|int|否|若为空,则由RPC Server设置该字段|
|max_gas|最大消耗GAS值|int|否|400000|
|height|区块高度|int|否|该值仅在查询时生效, 创建交易时忽略该字段|
|indent|返回数据格式是否有缩进|boolean|否|false|
|mode|广播交易模式|string|否|参数仅在广播交易时生效, 可选值: sync, async, block. 默认值: block|


> 请求参数中,`int`类型的数据需要用双引号括起来, 否则RPC Server将无法正确解析数据
> 如: "max_gas": "200000"



下面以转账交易为例, 简要说明API使用流程:

## 转账交易

交易涉及到的API如下:

-  `POST /bank/accounts/{toAddress}/transfers`: 创建转账交易API
-  `POST /txs`: 广播交易API
-  `GET /accounts/pubkey/{bech32PubKey}`: 转换公钥格式API

### 首次创建转账交易

用户`qosacc1ywhfxyte0j9e3h2zq07cphkdnqseqk4mym0jtg`向用户`qosacc1urrplt425ka2jeypcttz55h97rt5yhl0vtphc2`发送`100000`QOS:

* 调用创建转账交易接口创建交易

```

POST /bank/accounts/qosacc1urrplt425ka2jeypcttz55h97rt5yhl0vtphc2/transfers
>
    {
      "base": {
        "from": "qosacc1ywhfxyte0j9e3h2zq07cphkdnqseqk4mym0jtg"
      },
      "qos":  "100000"
    }


<

Status Code: 200 OK

{
    "code": "200",
    "tx": "{\"type\":\"qbase/txs/stdtx\",\"value\":{\"itx\":[{\"type\":\"bank/txs/TxTransfer\",\"value\":{\"senders\":[{\"addr\":\"qosacc1ywhfxyte0j9e3h2zq07cphkdnqseqk4mym0jtg\",\"qos\":\"100000\",\"qscs\":null}],\"receivers\":[{\"addr\":\"qosacc1urrplt425ka2jeypcttz55h97rt5yhl0vtphc2\",\"qos\":\"100000\",\"qscs\":null}]}}],\"sigature\":[{\"pubkey\":null,\"signature\":null,\"nonce\":\"1\"}],\"chainid\":\"qos-test\",\"maxgas\":\"400000\"}}",
    "signer": "qosacc1ywhfxyte0j9e3h2zq07cphkdnqseqk4mym0jtg",
    "pubkey": null,
    "nonce": "1",
    "sign_bytes": "nEooSwoeChQjrpMReXyLmN1CA/2A3s2YIZBauxIGMTAwMDAwEh4KFODGH66qpbqpZIHC1ipS5fDXQl/vEgYxMDAwMDBxb3MtdGVzdAAAAAAABhqAAAAAAAAAAAFxb3MtdGVzdA=="
}

```

返回数据说明:

|参数|参数说明|
|---|---|
|tx|待组装数据的交易体|
|sign_bytes|base64编码的待签名字符串|
|signer|签名账户地址|
|pubkey|用户ed25519公钥信息. 若为空,说明链上无用户公钥数据|
|nonce|创建交易使用的nonce值|    


* 本地使用私钥信息对`sign_bytes`进行签名:

```
//js
let buffer: Uint8Array = qosKeys.SignBase64("nEooSwoeChQjrpMReXyLmN1CA/2A3s2YIZBauxIGMTAwMDAwEh4KFODGH66qpbqpZIHC1ipS5fDXQl/vEgYxMDAwMDBxb3MtdGVzdAAAAAAABhqAAAAAAAAAAAFxb3MtdGVzdA==");
let signature: string = qosKeys.EncodeBase64(buffer);

//java
SignDataResponse signDataResponse = liteWallet.SignBase64(address, password, "nEooSwoeChQjrpMReXyLmN1CA/2A3s2YIZBauxIGMTAwMDAwEh4KFODGH66qpbqpZIHC1ipS5fDXQl/vEgYxMDAwMDBxb3MtdGVzdAAAAAAABhqAAAAAAAAAAAFxb3MtdGVzdA==");
String signature = signDataResponse.getData();

//signature: NDcsNjcsMjUsMTIsMzEsMjM5LDI1LDE5Miw3NCwyOCwyNDIsOTcsMjQzLDIyLDIxNSwyMjQsMTM5LDEwOCw3MiwyNyw2OCw3MCw0MywxNTEsNDYsMTYsMTg4LDcxLDU1LDExNyw2OSwxODIsNTksMTQ2LDYzLDYxLDIwMyw2NCw3MiwxNjIsMTc1LDI0MiwyMjMsMTQsMzcsMTc4LDE2NiwyMTcsMTE2LDIyOCw3LDE2NCwxNzQsMTgwLDIwOSwxNTYsMjI3LDI0NSw2MCwxNDEsNTYsNTcsMjM1LDM=
```

* 调用`GET /accounts/pubkey/{bech32PubKey}`公钥转换接口转换公钥:

```
GET /accounts/pubkey/qosaccpub1d9rez3ulq36p4x39zx94t0lnhd0f8359uwwzrlre99t6mr97nc8s6l8vql

>

<

Status Code: 200 OK
    {
        "type": "tendermint/PubKeyEd25519",
        "value": "aUeRR58EdBqaJRGLVb/zu16TxoXjnCH8eSlXrYy+ng8="
    }

```

* 组装交易体

将待组装交易体`tx`中的字符串转换为JSON对象如下:

```json

{"type":"qbase/txs/stdtx","value":{"itx":[{"type":"bank/txs/TxTransfer","value":{"senders":[{"addr":"qosacc1ywhfxyte0j9e3h2zq07cphkdnqseqk4mym0jtg","qos":"100000","qscs":null}],"receivers":[{"addr":"qosacc1urrplt425ka2jeypcttz55h97rt5yhl0vtphc2","qos":"100000","qscs":null}]}}],"sigature":[{"pubkey":null,"signature":null,"nonce":"1"}],"chainid":"qos-test","maxgas":"400000"}}

```

其中需要组装的只有`sigature.pubkey` 和 `sigature.signature`两处, 将上两步获取的数据填充至对应字段中, 得到数据如下:

```json

{
    "type": "qbase/txs/stdtx",
    "value": {
        "itx": [
            {
                "type": "bank/txs/TxTransfer",
                "value": {
                    "senders": [
                        {
                            "addr": "qosacc1ywhfxyte0j9e3h2zq07cphkdnqseqk4mym0jtg",
                            "qos": "100000",
                            "qscs": null
                        }
                    ],
                    "receivers": [
                        {
                            "addr": "qosacc1urrplt425ka2jeypcttz55h97rt5yhl0vtphc2",
                            "qos": "100000",
                            "qscs": null
                        }
                    ]
                }
            }
        ],
        "sigature": [
            {
                "pubkey":     
                    {
                      "type": "tendermint/PubKeyEd25519",
                      "value": "aUeRR58EdBqaJRGLVb/zu16TxoXjnCH8eSlXrYy+ng8="
                    },
                "signature": "NDcsNjcsMjUsMTIsMzEsMjM5LDI1LDE5Miw3NCwyOCwyNDIsOTcsMjQzLDIyLDIxNSwyMjQsMTM5LDEwOCw3MiwyNyw2OCw3MCw0MywxNTEsNDYsMTYsMTg4LDcxLDU1LDExNyw2OSwxODIsNTksMTQ2LDYzLDYxLDIwMyw2NCw3MiwxNjIsMTc1LDI0MiwyMjMsMTQsMzcsMTc4LDE2NiwyMTcsMTE2LDIyOCw3LDE2NCwxNzQsMTgwLDIwOSwxNTYsMjI3LDI0NSw2MCwxNDEsNTYsNTcsMjM1LDM=",
                "nonce": "1"
            }
        ],
        "chainid": "qos-test",
        "maxgas": "400000"
    }
}
```

* 广播交易

```
POST /txs

>

{
    "tx": {
        "type": "qbase/txs/stdtx",
        "value": {
            "itx": [
                {
                    "type": "bank/txs/TxTransfer",
                    "value": {
                        "senders": [
                            {
                                "addr": "qosacc1ywhfxyte0j9e3h2zq07cphkdnqseqk4mym0jtg",
                                "qos": "100000",
                                "qscs": null
                            }
                        ],
                        "receivers": [
                            {
                                "addr": "qosacc1urrplt425ka2jeypcttz55h97rt5yhl0vtphc2",
                                "qos": "100000",
                                "qscs": null
                            }
                        ]
                    }
                }
            ],
            "sigature": [
                {
                    "pubkey": {
                        "type": "tendermint/PubKeyEd25519",
                        "value": "aUeRR58EdBqaJRGLVb/zu16TxoXjnCH8eSlXrYy+ng8="
                    },
                    "signature": "NDcsNjcsMjUsMTIsMzEsMjM5LDI1LDE5Miw3NCwyOCwyNDIsOTcsMjQzLDIyLDIxNSwyMjQsMTM5LDEwOCw3MiwyNyw2OCw3MCw0MywxNTEsNDYsMTYsMTg4LDcxLDU1LDExNyw2OSwxODIsNTksMTQ2LDYzLDYxLDIwMyw2NCw3MiwxNjIsMTc1LDI0MiwyMjMsMTQsMzcsMTc4LDE2NiwyMTcsMTE2LDIyOCw3LDE2NCwxNzQsMTgwLDIwOSwxNTYsMjI3LDI0NSw2MCwxNDEsNTYsNTcsMjM1LDM=",
                    "nonce": "1"
                }
            ],
            "chainid": "qos-test",
            "maxgas": "400000"
        }
    },
    "mode": "block"
}


<

 响应

```


### 首次创建转账交易

用户`qosacc1ywhfxyte0j9e3h2zq07cphkdnqseqk4mym0jtg`向用户`qosacc1urrplt425ka2jeypcttz55h97rt5yhl0vtphc2`发送`100000`QOS:

* 调用创建转账交易接口创建交易

```

POST /bank/accounts/qosacc1urrplt425ka2jeypcttz55h97rt5yhl0vtphc2/transfers
>
    {
      "base": {
        "from": "qosacc1ywhfxyte0j9e3h2zq07cphkdnqseqk4mym0jtg"
      },
      "qos":  "100000"
    }


<

Status Code: 200 OK

{
    "code": "200",
    "tx": "{\"type\":\"qbase/txs/stdtx\",\"value\":{\"itx\":[{\"type\":\"bank/txs/TxTransfer\",\"value\":{\"senders\":[{\"addr\":\"qosacc1ywhfxyte0j9e3h2zq07cphkdnqseqk4mym0jtg\",\"qos\":\"100000\",\"qscs\":null}],\"receivers\":[{\"addr\":\"qosacc1urrplt425ka2jeypcttz55h97rt5yhl0vtphc2\",\"qos\":\"100000\",\"qscs\":null}]}}],\"sigature\":[{\"pubkey\":null,\"signature\":null,\"nonce\":\"1\"}],\"chainid\":\"qos-test\",\"maxgas\":\"400000\"}}",
    "signer": "qosacc1ywhfxyte0j9e3h2zq07cphkdnqseqk4mym0jtg",
    "pubkey": null,
    "nonce": "1",
    "sign_bytes": "nEooSwoeChQjrpMReXyLmN1CA/2A3s2YIZBauxIGMTAwMDAwEh4KFODGH66qpbqpZIHC1ipS5fDXQl/vEgYxMDAwMDBxb3MtdGVzdAAAAAAABhqAAAAAAAAAAAFxb3MtdGVzdA=="
}

```

返回数据说明:

|参数|参数说明|
|---|---|
|tx|待组装数据的交易体|
|sign_bytes|base64编码的待签名字符串|
|signer|签名账户地址|
|pubkey|用户ed25519公钥信息. 若为空,说明链上无用户公钥数据|
|nonce|创建交易使用的nonce值|    


* 本地使用私钥信息对`sign_bytes`进行签名:

```
//js
let buffer: Uint8Array = qosKeys.SignBase64("nEooSwoeChQjrpMReXyLmN1CA/2A3s2YIZBauxIGMTAwMDAwEh4KFODGH66qpbqpZIHC1ipS5fDXQl/vEgYxMDAwMDBxb3MtdGVzdAAAAAAABhqAAAAAAAAAAAFxb3MtdGVzdA==");
let signature: string = qosKeys.EncodeBase64(buffer);

//java
SignDataResponse signDataResponse = liteWallet.SignBase64(address, password, "nEooSwoeChQjrpMReXyLmN1CA/2A3s2YIZBauxIGMTAwMDAwEh4KFODGH66qpbqpZIHC1ipS5fDXQl/vEgYxMDAwMDBxb3MtdGVzdAAAAAAABhqAAAAAAAAAAAFxb3MtdGVzdA==");
String signature = signDataResponse.getData();

//signature: NDcsNjcsMjUsMTIsMzEsMjM5LDI1LDE5Miw3NCwyOCwyNDIsOTcsMjQzLDIyLDIxNSwyMjQsMTM5LDEwOCw3MiwyNyw2OCw3MCw0MywxNTEsNDYsMTYsMTg4LDcxLDU1LDExNyw2OSwxODIsNTksMTQ2LDYzLDYxLDIwMyw2NCw3MiwxNjIsMTc1LDI0MiwyMjMsMTQsMzcsMTc4LDE2NiwyMTcsMTE2LDIyOCw3LDE2NCwxNzQsMTgwLDIwOSwxNTYsMjI3LDI0NSw2MCwxNDEsNTYsNTcsMjM1LDM=
```

* 调用`GET /accounts/pubkey/{bech32PubKey}`公钥转换接口转换公钥:

```
GET /accounts/pubkey/qosaccpub1d9rez3ulq36p4x39zx94t0lnhd0f8359uwwzrlre99t6mr97nc8s6l8vql

>

<

Status Code: 200 OK
    {
        "type": "tendermint/PubKeyEd25519",
        "value": "aUeRR58EdBqaJRGLVb/zu16TxoXjnCH8eSlXrYy+ng8="
    }

```

* 组装交易体

将待组装交易体`tx`中的字符串转换为JSON对象如下:

```json

{"type":"qbase/txs/stdtx","value":{"itx":[{"type":"bank/txs/TxTransfer","value":{"senders":[{"addr":"qosacc1ywhfxyte0j9e3h2zq07cphkdnqseqk4mym0jtg","qos":"100000","qscs":null}],"receivers":[{"addr":"qosacc1urrplt425ka2jeypcttz55h97rt5yhl0vtphc2","qos":"100000","qscs":null}]}}],"sigature":[{"pubkey":null,"signature":null,"nonce":"1"}],"chainid":"qos-test","maxgas":"400000"}}

```

其中需要组装的只有`sigature.pubkey` 和 `sigature.signature`两处, 将上两步获取的数据填充至对应字段中, 得到数据如下:

```json

{
    "type": "qbase/txs/stdtx",
    "value": {
        "itx": [
            {
                "type": "bank/txs/TxTransfer",
                "value": {
                    "senders": [
                        {
                            "addr": "qosacc1ywhfxyte0j9e3h2zq07cphkdnqseqk4mym0jtg",
                            "qos": "100000",
                            "qscs": null
                        }
                    ],
                    "receivers": [
                        {
                            "addr": "qosacc1urrplt425ka2jeypcttz55h97rt5yhl0vtphc2",
                            "qos": "100000",
                            "qscs": null
                        }
                    ]
                }
            }
        ],
        "sigature": [
            {
                "pubkey":     
                    {
                      "type": "tendermint/PubKeyEd25519",
                      "value": "aUeRR58EdBqaJRGLVb/zu16TxoXjnCH8eSlXrYy+ng8="
                    },
                "signature": "NDcsNjcsMjUsMTIsMzEsMjM5LDI1LDE5Miw3NCwyOCwyNDIsOTcsMjQzLDIyLDIxNSwyMjQsMTM5LDEwOCw3MiwyNyw2OCw3MCw0MywxNTEsNDYsMTYsMTg4LDcxLDU1LDExNyw2OSwxODIsNTksMTQ2LDYzLDYxLDIwMyw2NCw3MiwxNjIsMTc1LDI0MiwyMjMsMTQsMzcsMTc4LDE2NiwyMTcsMTE2LDIyOCw3LDE2NCwxNzQsMTgwLDIwOSwxNTYsMjI3LDI0NSw2MCwxNDEsNTYsNTcsMjM1LDM=",
                "nonce": "1"
            }
        ],
        "chainid": "qos-test",
        "maxgas": "400000"
    }
}
```

* 广播交易

```
POST /txs

>

{
    "tx": {
        "type": "qbase/txs/stdtx",
        "value": {
            "itx": [
                {
                    "type": "bank/txs/TxTransfer",
                    "value": {
                        "senders": [
                            {
                                "addr": "qosacc1ywhfxyte0j9e3h2zq07cphkdnqseqk4mym0jtg",
                                "qos": "100000",
                                "qscs": null
                            }
                        ],
                        "receivers": [
                            {
                                "addr": "qosacc1urrplt425ka2jeypcttz55h97rt5yhl0vtphc2",
                                "qos": "100000",
                                "qscs": null
                            }
                        ]
                    }
                }
            ],
            "sigature": [
                {
                    "pubkey": {
                        "type": "tendermint/PubKeyEd25519",
                        "value": "aUeRR58EdBqaJRGLVb/zu16TxoXjnCH8eSlXrYy+ng8="
                    },
                    "signature": "NDcsNjcsMjUsMTIsMzEsMjM5LDI1LDE5Miw3NCwyOCwyNDIsOTcsMjQzLDIyLDIxNSwyMjQsMTM5LDEwOCw3MiwyNyw2OCw3MCw0MywxNTEsNDYsMTYsMTg4LDcxLDU1LDExNyw2OSwxODIsNTksMTQ2LDYzLDYxLDIwMyw2NCw3MiwxNjIsMTc1LDI0MiwyMjMsMTQsMzcsMTc4LDE2NiwyMTcsMTE2LDIyOCw3LDE2NCwxNzQsMTgwLDIwOSwxNTYsMjI3LDI0NSw2MCwxNDEsNTYsNTcsMjM1LDM=",
                    "nonce": "1"
                }
            ],
            "chainid": "qos-test",
            "maxgas": "400000"
        }
    },
    "mode": "block"
}


<

 响应

```

### 非首次创建转账交易

非首次创建交易时, 返回的交易体`tx`中将带有用户公钥信息,  可以不需要调用`GET /accounts/pubkey/{bech32PubKey}`接口进行公钥转换. 组装交易时,只需把签名数据进行组装就可以了.





