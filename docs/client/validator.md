# Validator命令行工具

参照[QOS验证人节点详解](validators/all_about_validators.md)了解validator相关实现机制。

含以下命令:

* `qoscli tx create-validator`: 成为validator
* `qoscli tx revoke-validator`: 撤消validator
* `qoscli tx active-validator`: 激活validator
* `qoscli query validator`: 查找validator
* `qoscli query validators`: validator列表

## create

```bash
$ qoscli tx create-validator --help

pubkey is a tendermint validator pubkey. the public key of the validator used in
Tendermint consensus.

owner is a keystore name or account address.

ex: pubkey: {"type":"tendermint/PubKeyEd25519","value":"VOn2rPx+t7Njdgi+eLb+jBuF175T1b7LAcHElsmIuXA="}

example:

         qoscli create-validator --name validatorName --owner ownerName --pubkey "VOn2rPx+t7Njdgi+eLb+jBuF175T1b7LAcHElsmIuXA=" --tokens 100

Usage:
  qoscli tx create-validator [flags]

Flags:
      --async                 broadcast transactions asynchronously
      --chain-id string       Chain ID of tendermint node
      --description string    description
  -h, --help                  help for create-validator
      --indent                add indent to json response
      --max-gas int           gas limit to set per tx
      --name string           name for validator
      --node string           <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")
      --nonce int             account nonce to sign the tx
      --nonce-node string     tcp://<host>:<port> to tendermint rpc interface for some chain to query account nonce
      --owner string          keystore name or account address
      --pubkey string         tendermint consensus validator public key
      --qcp                   enable qcp mode. send qcp tx
      --qcp-blockheight int   qcp mode flag. original tx blockheight, blockheight must greater than 0
      --qcp-extends string    qcp mode flag. qcp tx extends info
      --qcp-from string       qcp mode flag. qcp tx source chainID
      --qcp-seq int           qcp mode flag.  qcp in sequence
      --qcp-signer string     qcp mode flag. qcp tx signer key name
      --qcp-txindex int       qcp mode flag. original tx index
      --tokens int            bond tokens amount
      --trust-node            Trust connected full node (don't verify proofs for responses)

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "C:\\Users\\imuge/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```
主要参数：

- owner         操作者keystore name或 QOS account地址
- name          validator名称
- pubkey        priv_validator.json 中 `pub_key`内`value`部分
- tokens        绑定tokens，不能大于owner持有QOS数量
- description   备注

> 可以通过`qoscli keys import`导入*owner*账户，并保证owner持有QOS

## revoke
撤销validator，撤销之后validator处于Inactive状态

```bash
$ qoscli tx revoke-validator --help
Revoke validator

Usage:
  qoscli tx revoke-validator [flags]

Flags:
      --async                 broadcast transactions asynchronously
      --chain-id string       Chain ID of tendermint node
  -h, --help                  help for revoke-validator
      --indent                add indent to json response
      --max-gas int           gas limit to set per tx
      --node string           <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")
      --nonce int             account nonce to sign the tx
      --nonce-node string     tcp://<host>:<port> to tendermint rpc interface for some chain to query account nonce
      --owner string          owner keystore name or address
      --qcp                   enable qcp mode. send qcp tx
      --qcp-blockheight int   qcp mode flag. original tx blockheight, blockheight must greater than 0
      --qcp-extends string    qcp mode flag. qcp tx extends info
      --qcp-from string       qcp mode flag. qcp tx source chainID
      --qcp-seq int           qcp mode flag.  qcp in sequence
      --qcp-signer string     qcp mode flag. qcp tx signer key name
      --qcp-txindex int       qcp mode flag. original tx index
      --trust-node            Trust connected full node (don't verify proofs for responses)

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "C:\\Users\\imuge/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```
主要参数：

- owner         操作者keystore name或 QOS account地址

```bash
$ qoscli tx revoke-validator --owner <keystore_name_of_owner or address_of_owner_account>
```

## active

重新激活validator，可重新参与打块和投票

```bash
$ qoscli tx active-validator --help
Active validator

Usage:
  qoscli tx active-validator [flags]

Flags:
      --async                 broadcast transactions asynchronously
      --chain-id string       Chain ID of tendermint node
  -h, --help                  help for active-validator
      --indent                add indent to json response
      --max-gas int           gas limit to set per tx
      --node string           <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")
      --nonce int             account nonce to sign the tx
      --nonce-node string     tcp://<host>:<port> to tendermint rpc interface for some chain to query account nonce
      --owner string          owner keystore or address
      --qcp                   enable qcp mode. send qcp tx
      --qcp-blockheight int   qcp mode flag. original tx blockheight, blockheight must greater than 0
      --qcp-extends string    qcp mode flag. qcp tx extends info
      --qcp-from string       qcp mode flag. qcp tx source chainID
      --qcp-seq int           qcp mode flag.  qcp in sequence
      --qcp-signer string     qcp mode flag. qcp tx signer key name
      --qcp-txindex int       qcp mode flag. original tx index
      --trust-node            Trust connected full node (don't verify proofs for responses)

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "C:\\Users\\imuge/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```
主要参数：

- owner         操作者keystore name或 QOS account地址

```bash
$ qoscli tx active-validator --owner <keystore_name_of_owner or address_of_owner_account>
```

## query

查询validator信息

### validator

根据owner查找
```bash
$ qoscli query validator --help
Query validator's info

Usage:
  qoscli query validator [flags]

Flags:
      --chain-id string   Chain ID of tendermint node
      --height int        block height to query, omit to get most recent provable block
  -h, --help              help for validator
      --indent            add indent to json response
      --node string       <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")
      --owner string      validator's owner address
      --trust-node        Trust connected full node (don't verify proofs for responses)

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "C:\\Users\\imuge/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```
主要参数：
- owner     验证节点owner地址

查询实例：
```bash
$ qoscli query validator --owner address1gmllq4fgtlfe574dffaj90t3tkvy232phcukhq --indent
{                                                               
  "name": "imuge",                                              
  "owner": "address1gmllq4fgtlfe574dffaj90t3tkvy232phcukhq",    
  "validatorPubkey": {                                          
    "type": "tendermint/PubKeyEd25519",                         
    "value": "Z4gGHPR6l0sFJiC3zn9tNH5jivwpfBSWp/yPzvTkmIE="     
  },                                                            
  "bondTokens": "100",                                          
  "description": "",                                            
  "status": 0,                                                  
  "inactiveCode": 0,                                            
  "inactiveTime": "0001-01-01T00:00:00Z",                       
  "inactiveHeight": "0",                                        
  "bondHeight": "1"                                             
}                                                               
```

> 其中`status`表示当前validator状态，0 active，1 pending

### validators

validators列表，打印出当前网络中所有验证节点信息。
```bash
$ qoscli query validators --help
Query all validators info

Usage:
  qoscli query validators [flags]

Flags:
      --chain-id string   Chain ID of tendermint node
      --height int        block height to query, omit to get most recent provable block
  -h, --help              help for validators
      --indent            add indent to json response
      --node string       <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")
      --trust-node        Trust connected full node (don't verify proofs for responses)

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "C:\\Users\\imuge/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors

```