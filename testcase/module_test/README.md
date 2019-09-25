# qoscli 模块测试

## 测试环境

- 操作系统: CentOS 7 18.10 (虚拟机环境)
- QOS版本 
```
{
 "version": "0.0.5",
 "commit": "16bc3df796f1492ea02a31de2cd6f834406c44a1",
 "go": "go version go1.12.5 linux/amd64"
}
```

## 待测模块

- [x] [Keys](keys.md) 本地密钥库
- [x] [Approve](approve.md) 预授权
- [ ] [Delegation](delegation.md) 委托
- [ ] [Governance](governance.md) 提案治理
- [ ] [Guardian](./guardian.md) 特权账户
- [ ] [QCP-QSC-Transfer](./qcp_qsc_transfer.md) QCP联盟链，QSC联盟币，以及QOS/QSCs转账
- [ ] [Validator](validator.md) 验证人节点

## 测试过程中涉及的密钥(keys)

- `node` - `Genesis Account: 100,000QOS`, `Genesis Guardian`, `Validator`
```bash
$ ./qoscli keys export node
Password to sign with 'node':
**Important** Don't leak your private key information to others.
Please keep your private key safely, otherwise your account will be attacked.

{"Name":"node","address":"address1qgwgmpsrd6anj3qjvjsqztj3xt9v24c4mh77x3","pubkey":{"type":"tendermint/PubKeyEd25519","value":"6ScYtz+2p/y1GZDe0XuibF0w1NPE6yY2o9GNgjUYFqo="},"privkey":{"type":"tendermint/PrivKeyEd25519","value":"tKFI6TwnYSV2ONVfW5xTappKiwnf0kvMClxPOmPKwFjpJxi3P7an/LUZkN7Re6JsXTDU08TrJjaj0Y2CNRgWqg=="}}
```

- `alice` - `Genesis Account: 200,000QOS`
```bash
$ ./qoscli keys export alice
Password to sign with 'alice':
**Important** Don't leak your private key information to others.
Please keep your private key safely, otherwise your account will be attacked.

{"Name":"alice","address":"address1eqqnaps04l6ht9xahtrfgg59ksllcq9qr8630q","pubkey":{"type":"tendermint/PubKeyEd25519","value":"EbIqepd8q2+8XnTjvlqjnWb1aptxLEkjiSuvX05nuBg="},"privkey":{"type":"tendermint/PrivKeyEd25519","value":"3UcHjus5TjQGpdJgXhNtViLbwQIKtcbVKDHbVa9w+iMRsip6l3yrb7xedOO+WqOdZvVqm3EsSSOJK69fTme4GA=="}}
```

- `bob` - `Genesis Account: 100,000QOS`
```bash
$ ./qoscli keys export bob
Password to sign with 'bob':
**Important** Don't leak your private key information to others.
Please keep your private key safely, otherwise your account will be attacked.

{"Name":"bob","address":"address15fc26swvguzy9wksha9506smj2gne5r3k7na3r","pubkey":{"type":"tendermint/PubKeyEd25519","value":"NtQ0jyi769kAoHy9kjhvzvvOQV5vcOuupg+/r+kDYVA="},"privkey":{"type":"tendermint/PrivKeyEd25519","value":"AlIiQkdzRu63iPE7OgPjEpR1SO2OOlWukVGUovTqBEQ21DSPKLvr2QCgfL2SOG/O+85BXm9w666mD7+v6QNhUA=="}}
```

- `charles` - `Genesis Account: 500,000QOS`
```bash
$ ./qoscli keys export charles
Password to sign with 'charles':
**Important** Don't leak your private key information to others.
Please keep your private key safely, otherwise your account will be attacked.

{"Name":"charles","address":"address1f37jvnehrfkpwzqtqtu5q9jx57034gvv0lz78s","pubkey":{"type":"tendermint/PubKeyEd25519","value":"UinjHf4xnFzTi0tdq+8ooehKg70WCpdXu48j37kIc8M="},"privkey":{"type":"tendermint/PrivKeyEd25519","value":"lfDVAXO1j+5cWk8uJQ5Dy9k1hw9GvnkZi7jOa9aDM/9SKeMd/jGcXNOLS12r7yih6EqDvRYKl1e7jyPfuQhzww=="}}
```

## 本地测试网搭建过程

### 本地测试网结构

```
 Node
  |--Alice
  |--Bob
  |--Charles
```

### 1. 虚拟机Node

1. 准备环境参数: 交易扫描路径, 监听地址(IP:Port), 关闭防火墙
2. 清除原有的`.qosd`和`.qoscli`目录
3. 初始化qos配置文件
4. 导入密钥`node`,`alice`,`bob`,`charles`，密码`12345678`
5. 添加账户`node`,`alice`,`bob`,`charles`到`genesis.json`文件，各账户余额如下:
    1. `node` - `1,000,000QOS`
    2. `alice` - `200,000QOS`
    3. `bob` - `100,000QOS`
    4. `charles` - `500,000QOS`
6. 在`genesis.json`文件中，将账户`node`标记为特权用户(Guardian)
7. 生成创世交易：使账户`node`成为验证人(Validator), 绑定`500，000QOS`
8. 收集创世交易到`genesis.json`文件
9. 修改genesis.json文件中的governance参数`$min_deposit`,`$max_deposit_period`,`$voting_period`
9. 复制`genesis.json`文件, 覆盖到`~/genesis.json`中
10. 提取Node-ID, 生成Seeds
11. 修改`~/qos-init-slave.sh`脚本文件中的Seeds

以下shell脚本实现了上述步骤1~11:
```shell script
#!/bin/bash
# params
echo "1. 准备环境参数"
path=~/.qosd/config/gentx
echo "1.1 交易扫描路径: "$path
ip=$(ifconfig -a |grep inet |grep -v inet6 |grep -v 127 |head -n 1 |awk '{print $2}')
port=26656
echo "1.2 监听地址: "$ip":"$port
echo "1.3 关闭防火墙: "
pwd="admin"
expect -c "
set timeout 1
spawn sudo iptables -F
send \""$pwd"\r\"
interact
"
# main code
echo "2. 清除原有的.qosd和.qoscli目录"
rm -rf ./.qosd
rm -rf ./.qoscli
echo "3. 初始化qos配置文件"
./qosd init --moniker "local-test" --chain-id "test-chain"
echo "4.1 导入密钥node"
expect -c "
set timeout 1
spawn ~/qoscli keys import node
send \"tKFI6TwnYSV2ONVfW5xTappKiwnf0kvMClxPOmPKwFjpJxi3P7an/LUZkN7Re6JsXTDU08TrJjaj0Y2CNRgWqg==\r\"
expect \"*key:\" {send \"12345678\r\"}
expect \"*passphrase:\" {send \"12345678\r\"}
interact
"
echo "4.2 导入密钥alice"
expect -c "
set timeout 1
spawn ~/qoscli keys import alice
send \"3UcHjus5TjQGpdJgXhNtViLbwQIKtcbVKDHbVa9w+iMRsip6l3yrb7xedOO+WqOdZvVqm3EsSSOJK69fTme4GA==\r\"
expect \"*key:\" {send \"12345678\r\"}
expect \"*passphrase:\" {send \"12345678\r\"}
interact
"
echo "4.3 导入密钥bob"
expect -c "
set timeout 1
spawn ~/qoscli keys import bob
send \"AlIiQkdzRu63iPE7OgPjEpR1SO2OOlWukVGUovTqBEQ21DSPKLvr2QCgfL2SOG/O+85BXm9w666mD7+v6QNhUA==\r\"
expect \"*key:\" {send \"12345678\r\"}
expect \"*passphrase:\" {send \"12345678\r\"}
interact
"
echo "4.3 导入密钥charles"
expect -c "
set timeout 1
spawn ~/qoscli keys import charles
send \"lfDVAXO1j+5cWk8uJQ5Dy9k1hw9GvnkZi7jOa9aDM/9SKeMd/jGcXNOLS12r7yih6EqDvRYKl1e7jyPfuQhzww==\r\"
expect \"*key:\" {send \"12345678\r\"}
expect \"*passphrase:\" {send \"12345678\r\"}
interact
"
echo "5. 添加账户node,alice,bob,charles到genesis.json文件"
./qosd add-genesis-accounts node,1000000QOS
./qosd add-genesis-accounts alice,200000QOS
./qosd add-genesis-accounts bob,100000QOS
./qosd add-genesis-accounts charles,500000QOS
echo "6. 在genesis.json文件中，将账户node标记为特权用户(Guardian)"
./qosd add-guardian --address address1qgwgmpsrd6anj3qjvjsqztj3xt9v24c4mh77x3 --description "Genesis Guardian node"
echo "7. 生成创世交易：使账户node成为验证人(Validator), 绑定500，000QOS"
expect -c "
set timeout 1
spawn ~/qosd gentx --moniker central --owner node --tokens 500000 --compound
expect \"Password to sign with 'node':\" {send \"12345678\r\"}
interact
"
echo "8. 收集创世交易到genesis.json文件"
./qosd collect-gentxs
echo "9. 修改genesis.json文件中的governance参数"
sed -ri 's/"min_deposit": "10"/"min_deposit": "100000"/g' ~/.qosd/config/genesis.json
sed -ri 's/"max_deposit_period": "172800000000000"/"max_deposit_period": "300000000000"/g' ~/.qosd/config/genesis.json
sed -ri 's/"voting_period": "172800000000000"/"voting_period": "300000000000"/g' ~/.qosd/config/genesis.json
echo "10. 复制genesis.json文件, 覆盖到~/genesis.json中"
cp -f ~/.qosd/config/genesis.json ~/genesis.json
echo "11. 提取Node-ID, 生成Seeds"
files=$(ls $path)
items=${files//@127.0.0.1.json/@$ip:$port}
for item in $items
do
 line=$line$item\;
done
seeds=${line:0:-1}
echo $seeds
echo "12. 修改~/qos-init-slave脚本文件中的seeds"
sed -ri "s/seeds='(([0-9a-z]*@[0-9.]*:[0-9]*);*)*'/seeds='"$seeds"'/g" ~/qos-init-slave.sh
```

最终生成的`genesis.json`文件被另存为`~/genesis.json`, 其内容如下:
```json
{
  "genesis_time": "2019-08-20T02:35:22.040817499Z",
  "chain_id": "test-chain",
  "consensus_params": {
    "block": {
      "max_bytes": "1048576",
      "max_gas": "-1",
      "time_iota_ms": "1000"
    },
    "evidence": {
      "max_age": "100000"
    },
    "validator": {
      "pub_key_types": [
        "ed25519"
      ]
    }
  },
  "app_hash": "",
  "app_state": {
    "gen_txs": [
      {
        "itx": [
          {
            "type": "stake/txs/TxCreateValidator",
            "value": {
              "Owner": "address1qgwgmpsrd6anj3qjvjsqztj3xt9v24c4mh77x3",
              "PubKey": {
                "type": "tendermint/PubKeyEd25519",
                "value": "6k8y++Pgbgz0rppLFYsV4tW6vLgKvhmn93FJFu7R7GU="
              },
              "BondTokens": "500000",
              "IsCompound": true,
              "Description": {
                "moniker": "central",
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
              "value": "6ScYtz+2p/y1GZDe0XuibF0w1NPE6yY2o9GNgjUYFqo="
            },
            "signature": "5434giRHVQURWuf5WrE87PZhQehZnlJQui2mq1z8i3ZcV2zC67/794MRr6gUQVLPqK/araunpULqsxnP9QuGCw==",
            "nonce": "1"
          }
        ],
        "chainid": "test-chain",
        "maxgas": "1000000"
      }
    ],
    "accounts": [
      {
        "base_account": {
          "account_address": "address1qgwgmpsrd6anj3qjvjsqztj3xt9v24c4mh77x3",
          "public_key": null,
          "nonce": "0"
        },
        "qos": "1000000",
        "qscs": null
      },
      {
        "base_account": {
          "account_address": "address1eqqnaps04l6ht9xahtrfgg59ksllcq9qr8630q",
          "public_key": null,
          "nonce": "0"
        },
        "qos": "200000",
        "qscs": null
      },
      {
        "base_account": {
          "account_address": "address15fc26swvguzy9wksha9506smj2gne5r3k7na3r",
          "public_key": null,
          "nonce": "0"
        },
        "qos": "100000",
        "qscs": null
      },
      {
        "base_account": {
          "account_address": "address1f37jvnehrfkpwzqtqtu5q9jx57034gvv0lz78s",
          "public_key": null,
          "nonce": "0"
        },
        "qos": "500000",
        "qscs": null
      }
    ],
    "mint": {
      "params": {
        "inflation_phrases": [
          {
            "endtime": "2023-01-01T00:00:00Z",
            "total_amount": "2500000000000",
            "applied_amount": "0"
          },
          {
            "endtime": "2027-01-01T00:00:00Z",
            "total_amount": "12750000000000",
            "applied_amount": "0"
          },
          {
            "endtime": "2031-01-01T00:00:00Z",
            "total_amount": "6375000000000",
            "applied_amount": "0"
          },
          {
            "endtime": "2035-01-01T00:00:00Z",
            "total_amount": "3185000000000",
            "applied_amount": "0"
          }
        ]
      },
      "first_block_time": "0",
      "applied_qos_amount": "1800000"
    },
    "stake": {
      "params": {
        "max_validator_cnt": 10,
        "voting_status_len": 100,
        "voting_status_least": 50,
        "survival_secs": 600,
        "unbond_return_height": 10
      },
      "validators": null,
      "val_votes_info": null,
      "val_votes_in_window": null,
      "delegators_info": null,
      "delegator_unbond_info": null,
      "current_validators": null
    },
    "qcp": {
      "ca_root_pub_key": null,
      "qcps": null
    },
    "qsc": {
      "ca_root_pub_key": null,
      "qscs": null
    },
    "approve": {
      "approves": null
    },
    "distribution": {
      "community_fee_pool": "0",
      "last_block_proposer": "address1ah9uz0",
      "pre_distribute_amount": "0",
      "validators_history_period": null,
      "validators_current_period": null,
      "delegators_earning_info": null,
      "delegators_income_height": null,
      "validator_eco_fee_pools": null,
      "params": {
        "proposer_reward_rate": {
          "value": "0.040000000000000000"
        },
        "community_reward_rate": {
          "value": "0.010000000000000000"
        },
        "validator_commission_rate": {
          "value": "0.010000000000000000"
        },
        "delegator_income_period_height": "10",
        "gas_per_unit_cost": "10"
      }
    },
    "governance": {
      "starting_proposal_id": "1",
      "params": {
        "min_deposit": "100000",
        "min_proposer_deposit_rate": "0.334000000000000000",
        "max_deposit_period": "300000000000",
        "voting_period": "300000000000",
        "quorum": "0.334000000000000000",
        "threshold": "0.500000000000000000",
        "veto": "0.334000000000000000",
        "penalty": "0.000000000000000000",
        "burn_rate": "0.500000000000000000"
      },
      "proposals": null
    },
    "guardian": {
      "guardians": [
        {
          "description": "Genesis Guardian node",
          "guardian_type": 1,
          "address": "address1qgwgmpsrd6anj3qjvjsqztj3xt9v24c4mh77x3",
          "creator": "address1ah9uz0"
        }
      ]
    }
  }
}
```

### 虚拟机Alice, Bob, Charles

1. 清除原有的`.qosd`和`.qoscli`目录
2. 初始化qos配置文件
3. 用之前生成的测试网配置`~/genesis.json`文件替换生成的`genesis.json`文件
4. 导入密钥`node`,`alice`,`bob`,`charles`，密码`12345678`
5. 在`~/.qosd/config/config.toml`文件中, 设置`seeds`, 格式: `Node-ID@IP:Port`

以下shell脚本实现了上述步骤1~5:
```shell script
#!/bin/bash
seeds='9215b628a74f457daa2a051cff736dbc96162036@192.168.61.128:26656'
echo "1. 清除原有的.qosd和.qoscli目录"
rm -rf ./.qosd
rm -rf ./.qoscli
echo "2. 初始化qos配置文件"
./qosd init --moniker "local-test" --chain-id "test-chain"
echo "3. 使用测试网genesis.json文件替换生成的genesis.json文件"
cp -f ~/genesis.json ~/.qosd/config/genesis.json
echo "4.1 导入密钥node"
expect -c "
set timeout 1
spawn ~/qoscli keys import node
send \"tKFI6TwnYSV2ONVfW5xTappKiwnf0kvMClxPOmPKwFjpJxi3P7an/LUZkN7Re6JsXTDU08TrJjaj0Y2CNRgWqg==\r\"
expect \"*key:\" {send \"12345678\r\"}
expect \"*passphrase:\" {send \"12345678\r\"}
interact
"
echo "4.2 导入密钥alice"
expect -c "
set timeout 1
spawn ~/qoscli keys import alice
send \"3UcHjus5TjQGpdJgXhNtViLbwQIKtcbVKDHbVa9w+iMRsip6l3yrb7xedOO+WqOdZvVqm3EsSSOJK69fTme4GA==\r\"
expect \"*key:\" {send \"12345678\r\"}
expect \"*passphrase:\" {send \"12345678\r\"}
interact
"
echo "4.3 导入密钥bob"
expect -c "
set timeout 1
spawn ~/qoscli keys import bob
send \"AlIiQkdzRu63iPE7OgPjEpR1SO2OOlWukVGUovTqBEQ21DSPKLvr2QCgfL2SOG/O+85BXm9w666mD7+v6QNhUA==\r\"
expect \"*key:\" {send \"12345678\r\"}
expect \"*passphrase:\" {send \"12345678\r\"}
interact
"
echo "4.3 导入密钥charles"
expect -c "
set timeout 1
spawn ~/qoscli keys import charles
send \"lfDVAXO1j+5cWk8uJQ5Dy9k1hw9GvnkZi7jOa9aDM/9SKeMd/jGcXNOLS12r7yih6EqDvRYKl1e7jyPfuQhzww==\r\"
expect \"*key:\" {send \"12345678\r\"}
expect \"*passphrase:\" {send \"12345678\r\"}
interact
"
echo "5. 配置到Node节点的连接: 在config.toml文件中, 设置seeds: "$seeds
sed -ri 's#seeds = \"(([0-9a-z]*@[0-9.]*:[0-9]*);*)*\"#seeds = \"'$seeds'\"#g' ~/.qosd/config/config.toml
```