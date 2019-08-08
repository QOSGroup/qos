# Description
```
可选参数ip格式不合法
```
# Input
```
$ qosd gentx --moniker TestNode --owner test01 --tokens 1000
```
# Output
```
$ qosd gentx --moniker TestNode --owner test01 --tokens 1000 --ip www.baidu.com
Password to sign with 'test01':

```
命令行无返回值, `$HOME/.qosd/config/gentx`目录下生成以`nodeID@IP`为文件名的已签名的交易数据文件.
这里生成的是`3752c6169c8466f2863ecf99eae53a3e6a0b77c7@www.baidu.com.json`, 其内容如下:
```
{
    "type": "qbase/txs/stdtx", 
    "value": {
        "itx": [
            {
                "type": "stake/txs/TxCreateValidator", 
                "value": {
                    "Owner": "address1n4u9hac9gv76xuxpdy9php6cenq8psv6h99cda", 
                    "PubKey": {
                        "type": "tendermint/PubKeyEd25519", 
                        "value": "quoukaDQIwv8/W16yVj6R334rrPSywAQUmJpW43mPkw="
                    }, 
                    "BondTokens": "1000", 
                    "IsCompound": false, 
                    "Description": {
                        "moniker": "TestNode", 
                        "logo": "", 
                        "website": "", 
                        "details": ""
                    }
                }
            }
        ], 
        "sigature": [
            {
                "pubkey": {
                    "type": "tendermint/PubKeyEd25519", 
                    "value": "72a3avx8mgth+knenrBSy2gh7b60oLKBQTx0cuGKMLA="
                }, 
                "signature": "MFcQdI9828O5XOF+TVOByRTiwfVlnHlv7VypyhaINUTxvbSdYE4fTUXqRZKg527BqgDfgLPBOPW4Pa5OfXimDA==", 
                "nonce": "1"
            }
        ], 
        "chainid": "test-chain", 
        "maxgas": "1000000"
    }
}
```