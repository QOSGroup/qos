# QOS Client

Commands in `qoscli`:
* `keys`        [Keys](#keys)
* `query`       [Query](#query)
* `tx`          [Transactions](#transactions)
* `version`     [Version](#version)

for help using `--help`.

common flags:

| parameter | default | description |
| :--- | :---: | :--- |
|--nonce | 0 | account nonce to sign the tx |
|--max-gas| 0 | gas limit to set per tx |
|--chain-id| "" | Chain ID of tendermint node |
|--node| tcp://localhost:26657 | tcp://\<host\>:\<port\> to tendermint rpc interface for this chain |
|--height| 0 | block height to query, omit to get most recent provable block |
|--async| false | broadcast transactions asynchronously |
|--trust-node| false | Trust connected full node |
|--qcp| false | enable qcp mode. send qcp tx |
|--qcp-signer| "" | qcp mode flag. qcp tx signer key name |
|--qcp-seq| 0 | qcp mode flag.  qcp in sequence |
|--qcp-from| "" | qcp mode flag. qcp tx source |
|--qcp-blockheight| 0 | qcp mode flag. original tx blockheight |
|--qcp-txindex| 0 | qcp mode flag. original tx index |
|--qcp-extends| "" | qcp mode flag. qcp tx extends info |
|--indent| false | add indent to json response |
|--nonce-node| "" | tcp://\<host\>:\<port\> to tendermint rpc interface for some chain to query account nonce |

View [qbase-common-parameters](https://github.com/QOSGroup/qbase/blob/master/docs/client/command.md#客户端命令) for more information.

## keys

Commands for local keybase:
* `qoscli keys add`     [Add key](#add)
* `qoscli keys list`    [Keys list](#list)
* `qoscli keys update`  [Update password](#update)
* `qoscli keys delete`  [Delete key](#delete)
* `qoscli keys import`  [Import key](#import)
* `qoscli keys export`  [Export key](#export)

> Dafault keybase path: $HOME/.qoscli/keys/, deleting this path will clear all private keys stored locally.
Operating the key through the `keys` related command does not affect the account status in the QOS network.
Please keep your account private key information in a safe place.

#### add

`qoscli keys add <key_name>`

<key_name> is the name for the key.

Add a key named `Arya` to the local keybase:
```bash
$ qoscli keys add Arya
Enter a passphrase for your key:<Enter a password of no less than 8 digits>
Repeat the passphrase:<Repeat the password entered above>
NAME:	TYPE:	ADDRESS:						PUBKEY:
Arya	local	qosacc10327kf8v45a7uhev92llmuqwzkfgecvwckxt5m	qosaccpub1zcjduepqhn2n540cn0ts0qg7zd8xyrrwjg54lvaka228c3vs8gf5ph3eh27sy7nlzh
**Important** write this seed phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

thought frame must space few omit muffin fix merge mail ivory clump unveil dirt gadget load glove hub inner final crime churn crop stone
```
`qosacc10327kf8v45a7uhev92llmuqwzkfgecvwckxt5m` is the QOS account address, `qosaccpub1zcjduepqhn2n540cn0ts0qg7zd8xyrrwjg54lvaka228c3vs8gf5ph3eh27sy7nlzh` is the public key, `thought frame must space few omit muffin fix merge mail ivory clump unveil dirt gadget load glove hub inner final crime churn crop stone` is the mnemonic, it can be used to retrieve the private key of the account. Please keep the mnemonic in a safe place.

#### list

`qoscli keys list`
```bash
$ qoscli keys list
NAME:	TYPE:	ADDRESS:						PUBKEY:
Arya	local	qosacc10327kf8v45a7uhev92llmuqwzkfgecvwckxt5m	qosaccpub1zcjduepqhn2n540cn0ts0qg7zd8xyrrwjg54lvaka228c3vs8gf5ph3eh27sy7nlzh
```

#### update

`qoscli keys update <key_name>`

update `Arya`'s pass:
```bash
$ qoscli keys update Arya
Enter the current passphrase:<Enter current password>
Enter the new passphrase:<Enter a new password>
Repeat the new passphrase:<Repeat new password>
Password successfully updated!
```

#### export

`qoscli keys export <key_name>`

export `Arya`'s key:
```bash
qoscli keys export Arya
Password to sign with 'Arya':<Enter password>
**Important** Don't leak your private key information to others.
Please keep your private key safely, otherwise your account will be attacked.

{"Name":"Arya","address":"qosacc10327kf8v45a7uhev92llmuqwzkfgecvwckxt5m","pubkey":"qosaccpub1zcjduepqhn2n540cn0ts0qg7zd8xyrrwjg54lvaka228c3vs8gf5ph3eh27sy7nlzh","privkey":{"type":"tendermint/PrivKeyEd25519","value":"n/eCiOFjYFf22NCsacMfTXhxI6dV3DfI8cuxlZ48M0S81TpV+JvXB4EeE05iDG6SKV+ztuqUfEWQOhNA3jm6vQ=="}}
```
The exported key is the key information serialized by JSON. You can save the contents of the `privkey` part of the JSON string as a file and save it properly, it can be used for key import.

#### delete

`qoscli keys delete <key_name>`

delete `Arya`'s key:
```bash
$ qoscli keys delete Arya
DANGER - enter password to permanently delete key:<Enter password>
key deleted forever (uh oh!)
```

#### import

`qoscli keys import Arya --file <private key file path>`

import `Arya`'s key:
```bash
qoscli keys import Arya --file Arya.pri
> Enter a passphrase for your key:<Enter a password of no less than 8 digits>
> Repeat the passphrase:<Repeat the password entered above>

The contents of the Arya.pri file are:
{"type":"tendermint/PrivKeyEd25519","value":"n/eCiOFjYFf22NCsacMfTXhxI6dV3DfI8cuxlZ48M0S81TpV+JvXB4EeE05iDG6SKV+ztuqUfEWQOhNA3jm6vQ=="}

```

## Version
`qoscli version`

result:
```bash
{
 "version": "0.0.4-46-g5ec63bd", // qos version
 "commit": "5ec63bd74c2c92924c25ffd5be1ff0f232bfcda4", // commit ID
 "go": "go version go1.11.5 linux/amd64" //go version
}

```

## Query

* `qoscli query account`                [Account](#account)
* `qoscli query store`                  [Store](#store)
* `qoscli query consensus`              consensus parameters
* `qoscli query approve`                [Approve](#query-approve)
* `qoscli query qcp`                    [QCP](#query-qcp)
* `qoscli query qsc`                    [QSC](#query-qsc)
* `qoscli query qscs`                   [QSCs](#query-qscs)
* `qoscli query validators`             [Validators](#query-validators)
* `qoscli query validator`              [Validator](#query-validator)
* `qoscli query validator-miss-vote`    [Validator miss vote](#query-validator-miss-vote)
* `qoscli query validator-period`       [Validator period](#query-validator-period)
* `qoscli query community-fee-pool`     [Community fee pool](#query-community-fee-pool)
* `qoscli query delegation`             [Delegation](#query-delegation)
* `qoscli query delegations-to`         [Validator delegations](#query-validator-delegations)
* `qoscli query delegations`            [Delegator delegations](#query-delegator-delegations)
* `qoscli query delegator-income`       [Delegator income](#query-delegator-income)
* `qoscli query unbondings`             [Unbondings](#query-unbondings)
* `qoscli query redelegations`          [Redelegations](#query-redelegations)
* `qoscli query proposal`               [Proposal](#query-proposal)
* `qoscli query proposals`              [Proposals](#query-proposals)
* `qoscli query vote`                   [Vote](#query-vote)
* `qoscli query votes`                  [Votes](#query-votes)
* `qoscli query deposit`                [Deposit](#query-deposit)
* `qoscli query deposits`               [Deposits](#query-deposits)
* `qoscli query tally`                  [Tally](#tally)
* `qoscli query params`                 [Params](#params)
* `qoscli query inflation-phrases`      [Inflation rules](#query-inflation-rules)
* `qoscli query total-inflation  `      [Total inflation](#query-total-inflation)
* `qoscli query total-applied`          [Total applied](#query-total-applied)
* `qoscli query guardian`               [Guardian](#query-guardian)
* `qoscli query guardians`              [Guardians](#query-guardians)
* `qoscli query status`                 [Node status](#status)
* `qoscli query tendermint-validators`  [Tendermint-validators](#query-tendermint-validators)
* `qoscli query block`                  [Block](#block)
* `qoscli query txs`                    [Transactions by tags](#query-transactions-by-tags)
* `qoscli query tx`                     [Transaction by hash](#query-transaction-by-hash)

### Status
`qoscli query status --indent`

result:
```bash
{
  "node_info": {
    "protocol_version": {
      "p2p": "7",
      "block": "10",
      "app": "0"
    },
    "id": "4537e18828364c6e3529000e30bcf9f25b0fc50c",
    "listen_addr": "tcp://0.0.0.0:26656",
    "network": "imuge",
    "version": "0.30.1",
    "channels": "4020212223303800",
    "moniker": "node1",
    "other": {
      "tx_index": "on",
      "rpc_address": "tcp://0.0.0.0:26657"
    }
  },
  "sync_info": {
    "latest_block_hash": "4D935B625A5C2D63FD251C8448C9765916B289E435A0388F64401767DFA22BD5",
    "latest_app_hash": "29E08C36CE8CEA35EF4DE04B002C852505361B303950F3E07EBFC031F8DAB854",
    "latest_block_height": "396",
    "latest_block_time": "2019-04-25T06:53:11.777203643Z",
    "catching_up": false
  },
  "validator_info": {
    "address": "qoscons1glmfxaskefe75gat4se48fw0rdftn8h5dveeyw",
    "pub_key": "qosconspub1zcjduepquqyj75pn6auj725dcjmdyrd7ndlx6mnzhncelwju9lj9j7k0u6vq00rs3r",
    "voting_power": "10000000"
  }

}
```

When `catching up` is `false` means that the node has been synchronized to the latest height.

### Block
`qoscli query block <height>`

`<height>` is the block height

query block information in the height of 10:
```bash
$ qoscli query block 10 --indent
```

result:
```bash
{
  "block_meta": {
    "block_id": {
      "hash": "A473CE3866A74277BC7F7B7AF70E55B40736B8A3CA3B8A55406AC8CF6E04ED50",
      "parts": {
        "total": "1",
        "hash": "B9C5DEF42EAA9D445E52B1F8DD34ECC96C02E537F43D1F7C8D829C84F8663127"
      }
    },
    "header": {
      "version": {
        "block": "10",
        "app": "0"
      },
      "chain_id": "Arya",
      "height": "20",
      "time": "2019-04-25T06:19:28.353298129Z",
      "num_txs": "0",
      "total_txs": "0",
      "last_block_id": {
        "hash": "BC153175007D7E5D5C6A27D22E3F7227224E43C537988DDCBF6C2F14A95DD432",
        "parts": {
          "total": "1",
          "hash": "EEFE6F3A761D9D28DBCA81424F9E50A8C716D0F4898FA7B3893CBB0AC7B55F4D"
        }
      },
      "last_commit_hash": "863F17ACB4909A5E043782DB06F3FE18C6DCF4988EE9B7C0CDA7D8337504FCFB",
      "data_hash": "",
      "validators_hash": "5CA1D1B7D703F2D2A9C270D1CD5819E7E0D439BA1C55645BCD8DB7B079389CA8",
      "next_validators_hash": "5CA1D1B7D703F2D2A9C270D1CD5819E7E0D439BA1C55645BCD8DB7B079389CA8",
      "consensus_hash": "294D8FBD0B94B767A7EBA9840F299A3586DA7FE6B5DEAD3B7EECBA193C400F93",
      "app_hash": "C31662F65DEE545FEDF15D98517CBF07034DC1821EF06DD87D2F956C315A0A9B",
      "last_results_hash": "",
      "evidence_hash": "",
      "proposer_address": "0E447E66089C9D97EFC2F4C172403F35740DD507"
    }
  },
  "block": {
    "header": {
      "version": {
        "block": "10",
        "app": "0"
      },
      "chain_id": "Arya",
      "height": "20",
      "time": "2019-04-25T06:19:28.353298129Z",
      "num_txs": "0",
      "total_txs": "0",
      "last_block_id": {
        "hash": "BC153175007D7E5D5C6A27D22E3F7227224E43C537988DDCBF6C2F14A95DD432",
        "parts": {
          "total": "1",
          "hash": "EEFE6F3A761D9D28DBCA81424F9E50A8C716D0F4898FA7B3893CBB0AC7B55F4D"
        }
      },
      "last_commit_hash": "863F17ACB4909A5E043782DB06F3FE18C6DCF4988EE9B7C0CDA7D8337504FCFB",
      "data_hash": "",
      "validators_hash": "5CA1D1B7D703F2D2A9C270D1CD5819E7E0D439BA1C55645BCD8DB7B079389CA8",
      "next_validators_hash": "5CA1D1B7D703F2D2A9C270D1CD5819E7E0D439BA1C55645BCD8DB7B079389CA8",
      "consensus_hash": "294D8FBD0B94B767A7EBA9840F299A3586DA7FE6B5DEAD3B7EECBA193C400F93",
      "app_hash": "C31662F65DEE545FEDF15D98517CBF07034DC1821EF06DD87D2F956C315A0A9B",
      "last_results_hash": "",
      "evidence_hash": "",
      "proposer_address": "0E447E66089C9D97EFC2F4C172403F35740DD507"
    },
    "data": {
      "txs": null
    },
    "evidence": {
      "evidence": null
    },
    "last_commit": {
      "block_id": {
        "hash": "BC153175007D7E5D5C6A27D22E3F7227224E43C537988DDCBF6C2F14A95DD432",
        "parts": {
          "total": "1",
          "hash": "EEFE6F3A761D9D28DBCA81424F9E50A8C716D0F4898FA7B3893CBB0AC7B55F4D"
        }
      },
      "precommits": [
        {
          "type": 2,
          "height": "19",
          "round": "0",
          "block_id": {
            "hash": "BC153175007D7E5D5C6A27D22E3F7227224E43C537988DDCBF6C2F14A95DD432",
            "parts": {
              "total": "1",
              "hash": "EEFE6F3A761D9D28DBCA81424F9E50A8C716D0F4898FA7B3893CBB0AC7B55F4D"
            }
          },
          "timestamp": "2019-04-25T06:19:28.353298129Z",
          "validator_address": "0E447E66089C9D97EFC2F4C172403F35740DD507",
          "validator_index": "0",
          "signature": "bfhVFCZMS/6hEmkFAaLfNwumKEUQNtRkGvnrMTTvezjpCbv/X0wSQQKq6g4crd5mI3WjZYp4vM+EA4SY55ucCw=="
        },
        {
          "type": 2,
          "height": "19",
          "round": "0",
          "block_id": {
            "hash": "BC153175007D7E5D5C6A27D22E3F7227224E43C537988DDCBF6C2F14A95DD432",
            "parts": {
              "total": "1",
              "hash": "EEFE6F3A761D9D28DBCA81424F9E50A8C716D0F4898FA7B3893CBB0AC7B55F4D"
            }
          },
          "timestamp": "2019-04-25T06:19:28.312339528Z",
          "validator_address": "E9816412631B42AE3921769FFD9DE121AA745422",
          "validator_index": "1",
          "signature": "vePZhdo+dRTEghf3aHhqWXJQgyXeSoB2q4o1WiIncxI1raXU5YTGKNEdD8Tq8TbmI2uDH5J6CAOGy9ru1DzODQ=="
        }
      ]
    }
  }
}
```

### Account

Query account state:
`qoscli query account <key_name_or_account_address>`

<key_name_or_account_address> is the name of the key stored for the local keybase or the address of the corresponding account。

Assume in the local keybase `Arya`'s address is `qosacc1x5lcfaqxxq7g7dy4lj5vq0u6xamp78lsnza98y`:
```bash
qoscli query account Arya --indent
```
or
```bash
qoscli query account qosacc1x5lcfaqxxq7g7dy4lj5vq0u6xamp78lsnza98y --indent
```
Will display the following information:
```bash
{
  "type": "qos/types/QOSAccount",
  "value": {
    "account_address": "qosacc1smrus8jlc9z02gz5rm36u0q3fdctjxm4nrc639",
    "public_key": "qosaccpub1zcjduepqhewtk4rn8050sgyws8unnvr9l9rlshvfyargrds7qr9m40crphkscksjqd",
    "nonce": 1,
    "qos": "10000",
    "qscs": [
        {
            "coin_name": "AOE",
            "amount": "10000"
        }
    ]
  }
}
```
We can find out that `Arya` has 10000QOS and 10000sAOE.

### Store

`qoscli query store --path /store/<store_key>/subspace --data <query_data>`

main parameters:

- `--path`  location
- `--data`  query data

query root ca public key:
```bash
$ qoscli query store --path /store/acc/subspace --data account --indent
```

result:

```bash
[
  {
    "key": "account:\ufffd\ufffd\ufffd\u001e_\ufffdD\ufffd T\u001e\ufffd\ufffd\u003c\u0011Kp\ufffd\u001bu",
    "value": {
      "type": "qos/types/QOSAccount",
      "value": {
        "account_address": "qosacc1smrus8jlc9z02gz5rm36u0q3fdctjxm4nrc639",
        "public_key": "qosaccpub1zcjduepqhewtk4rn8050sgyws8unnvr9l9rlshvfyargrds7qr9m40crphkscksjqd",
        "nonce": 1,
        "qos": "10000",
        "qscs": [
            {
                "coin_name": "AOE",
                "amount": "10000"
            }
        ]
      }
    }
  }
]

```

### Query transactions
commands:

* `qoscli query tx`            [Query transaction by hash](#query-transaction-by-hash)
* `qoscli query txs`           [Query transactions by tags](#query-transactions-by-tags)

#### Query transaction by hash
After the transaction is executed, the transaction hash will be returned, and the transaction hash can be used to query the transaction details.

query transaction by `f5fc2c228cba754d5b95e49b02e81ff818f7b9140f1859d3797b09fb4aa12385`:

```bash
$ qoscli query tx f5fc2c228cba754d5b95e49b02e81ff818f7b9140f1859d3797b09fb4aa12385 --indent
```
result:

```bash
{
  "height": "153",
  "txhash": "8A317C18448BE7F2B76E05A930FCE0BA45B9C688A630CD50D886A05A7FB1A673",
  "gas_wanted": "9223372036854775807",
  "gas_used": "19500",
  "events": [
    {
      "type": "message",
      "attributes": [
        {
          "key": "module",
          "value": "transfer"
        },
        {
          "key": "gas.payer",
          "value": "qosacc1x5lcfaqxxq7g7dy4lj5vq0u6xamp78lsnza98y"
        }
      ]
    },
    {
      "type": "receive",
      "attributes": [
        {
          "key": "address",
          "value": "qosacc1smrus8jlc9z02gz5rm36u0q3fdctjxm4nrc639"
        },
        {
          "key": "qos",
          "value": "100"
        },
        {
          "key": "qscs"
        }
      ]
    },
    {
      "type": "send",
      "attributes": [
        {
          "key": "address",
          "value": "qosacc1x5lcfaqxxq7g7dy4lj5vq0u6xamp78lsnza98y"
        },
        {
          "key": "qos",
          "value": "100"
        },
        {
          "key": "qscs"
        }
      ]
    }
  ],
  "tx": {
    "type": "qbase/txs/stdtx",
    "value": {
      "itx": [
        {
          "type": "transfer/txs/TxTransfer",
          "value": {
            "senders": [
              {
                "addr": "qosacc1x5lcfaqxxq7g7dy4lj5vq0u6xamp78lsnza98y",
                "qos": "100",
                "qscs": null
              }
            ],
            "receivers": [
              {
                "addr": "qosacc1smrus8jlc9z02gz5rm36u0q3fdctjxm4nrc639",
                "qos": "100",
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
            "value": "2WRkEzr8Yd32reYXAmZURHwFm8CAWUTFAQc/IccnATg="
          },
          "signature": "tkyobbxwY5sjCTb1ibMEuMpOXHbDhKKHkZq0INZUvnyHODQloh4msRMmc62LLdswL1aY3vdld3OK+1pg2J86Dw==",
          "nonce": "1"
        }
      ],
      "chainid": "qos-test",
      "maxgas": "9223372036854775807"
    }
  },
  "timestamp": "2019-08-28T07:15:52Z"
}

```

#### Query transactions by tags
After the transaction is executed, transaction tags will be returned, and the transaction tags can be used to query the transaction information.

using `approve-from`=`qosacc1x5lcfaqxxq7g7dy4lj5vq0u6xamp78lsnza98y` to query transactions:

```bash
$ qoscli query txs --tags "create-approve.approve-from='qosacc1x5lcfaqxxq7g7dy4lj5vq0u6xamp78lsnza98y'" --indent
```
result:

```bash
[
  {
    "hash": "f5fc2c228cba754d5b95e49b02e81ff818f7b9140f1859d3797b09fb4aa12385",
    "height": "246",
    "tx": {
      "type": "qbase/txs/stdtx",
      "value": {
        "itx": [
          {
            "type": "approve/txs/TxCreateApprove",
            "value": {
              "Approve": {
                "from": "address1s348wvf49dfy64e6wafc90lcavp4lrd6xzhzhk",
                "to": "address1yqekgyy66v2cxzww6lqg6sdrsugjguxqws6mkf",
                "qos": "100",
                "qscs": null
              }
            }
          }
        ],
        "sigature": [
          {
            "pubkey": {
              "type": "tendermint/PubKeyEd25519",
              "value": "B/iatjhcJ4yFyHfGYKw2IneYGu2zG+ZOR8XmRUaji0A="
            },
            "signature": "VrsOsULJx86y8ch529zvl3Sh19TwGm/AldPlQhVWqhtg+calZmBrk25sD9HxCYijAt+ZUWMiLtPg3QZzCCqHAg==",
            "nonce": "1"
          }
        ],
        "chainid": "QOS",
        "maxgas": "100000"
      }
    },
    "result": {
      "gasWanted": "100000",
      "gasUsed": "15220",
      "tags": [
        {
          "key": "YWN0aW9u",
          "value": "Y3JlYXRlLWFwcHJvdmU="
        },
        {
          "key": "YXBwcm92ZS1mcm9t",
          "value": "YWRkcmVzczFzMzQ4d3ZmNDlkZnk2NGU2d2FmYzkwbGNhdnA0bHJkNnh6aHpoaw=="
        },
        {
          "key": "YXBwcm92ZS10bw==",
          "value": "YWRkcmVzczF5cWVrZ3l5NjZ2MmN4end3NmxxZzZzZHJzdWdqZ3V4cXdzNm1rZg=="
        }
      ]
    }
  }
]

```

visit [index](../spec/indexing.md) for more tags.

## Transactions

QOS supports the following transaction types:

* `qoscli tx transfer`         [Transfer](#transfer)
* `qoscli tx invariant-check`  [Invariant check](#invariant-check)
* `qoscli tx create-approve`   [Create approve](#create-approve)
* `qoscli tx increase-approve` [Increase approve](#increase-approve)
* `qoscli tx decrease-approve` [Decrease approve](#decrease-approve)
* `qoscli tx use-approve`      [Use approve](#use-approve)
* `qoscli tx cancel-approve`   [Cancel approve](#cancel-approve)
* `qoscli tx create-qsc`       [Create QSC](#create-qsc)
* `qoscli tx issue-qsc`        [Issue QSC](#issue-qsc)
* `qoscli tx init-qcp`         [Init QCP](#init-qcp)
* `qoscli tx create-validator` [Create validator](#create-validator)
* `qoscli tx create-validator` [Modify validator](#modify-validator)
* `qoscli tx revoke-validator` [Revoke validator](#revoke-validator)
* `qoscli tx active-validator` [Active validator](#active-validator)
* `qoscli tx delegate`         [Delegate](#delegate)
* `qoscli tx modify-compound`  [Modify compound](#modify-compound)
* `qoscli tx unbond`           [Unbond](#unbond)
* `qoscli tx redelegate`       [Redelegate](#redelegate)
* `qoscli tx submit-proposal`  [Submit proposal](#submit-proposal)
* `qoscli tx deposit`          [Deposit proposal](#deposit)
* `qoscli tx vote`             [Vote proposal](#vote)
* `qoscli tx add-guardian`     [Add guardian](#add-guardian)
* `qoscli tx delete-guardian`  [Delete guardian](#delete-guardian)
* `qoscli tx halt-network`     [Halt network](#halt-network)

Divided into [Bank](#bank),[Approve](#approve),[QSC](#qsc),[QCP](#qcp),[Validator](#validator),[Delegation](#delegation),[Governance](#governance),[Guardian](#guardian) these categories.

### Bank

See [Bank Spec](../spec/bank) to learn about bank module design.

* `qoscli tx transfer`          [Transfer](#transfer)
* `qoscli invariant-check`      [Invariant check](#invariant-check)

#### Transfer

`qoscli tx transfer --senders <senders_and_coins> --receivers <receivers_and_coins>`

main params:
- `--senders`   send set
- `--receivers` receive set

`Arya` send `qosacc1smrus8jlc9z02gz5rm36u0q3fdctjxm4nrc639` 1QOS and 1AOE:
```bash
$ qoscli tx transfer --senders Arya,1QOS,1AOE --receivers qosacc1smrus8jlc9z02gz5rm36u0q3fdctjxm4nrc639,1QOS,1AOE
Password to sign with 'Arya':<imput keybase password>
{"check_tx":{},"deliver_tx":{},"hash":"21ECB72C8F51B3BD8E3CB9D59765003B9D78BE75","height":"300"}
```

Execute [Query account](#account) to see the latest account state.

#### Invariant check

QOS has designed a [Invariant checking mechanism](../spec/bank), user can perform invariant checking operation by the following instructions:

`qoscli tx invariant-check --sender <sender's keybase name or address>`

main params:
- `--sender` sender's keybase name or address

::: warning Note 
This transaction sets a very large transaction fee. It is only submitted when the holder of the currency account finds that the QOS network data is abnormal. The data verification abnormality will stop the entire QOS network to protect the rights of the holder.
:::

`Arya` found a value overflow in the QOS network:
```bash
$ qoscli tx invariant-check --sender Arya
Password to sign with 'Arya':<input keystor password>
```
If there is no abnormality in the data, the normal transaction execution result will be returned, otherwise the whole network will stop running.

### Approve

[QOS Approve](../spec/approve/) includes:

* `qoscli tx create-approve`    [Create approve](#create-approve)
* `qoscli query approve`        [Query approve](#query-approve)
* `qoscli tx increase-approve`  [Increase approve](#increase-approve)
* `qoscli tx decrease-approve`  [Decrease approve](#decrease-approve)
* `qoscli tx use-approve`       [Use approve](#use-approve)
* `qoscli tx cancel-approve`    [Cancel approve](#cancel-approve)

> We use `qosacc1smrus8jlc9z02gz5rm36u0q3fdctjxm4nrc639` as `Sansa`'s address.

#### Create approve

`qoscli tx create-approve --from <key_name_or_account_address> --to <account_address> --coins <qos_and_qscs>`

main parameters:

- `--from`  approver's keybase name or address
- `--to`    approvee's keybase name or address
- `--coins` QOS and QSCs, [amount1][coin1],[amount2][coin2],..., comma separated.

`Arya` approve `Sansa` 100QOS,100AOE:
```
$ qoscli tx create-approve --from Arya --to qosacc1smrus8jlc9z02gz5rm36u0q3fdctjxm4nrc639 --coins 100QOS,100AOE
Password to sign with 'Arya':<input Arya's keybase password>
```
result:
```bash
{"check_tx":{},"deliver_tx":{},"hash":"9917953D8CDE80F457CD072DBCE73A36449B7A7C","height":"333"}
```

#### Query approve

`qoscli query approve --from <key_name_or_account_address> --to <account_address>`

query approve from `Arya` to `Sansa`:
```bash
qoscli query approve --from Arya --to qosacc1smrus8jlc9z02gz5rm36u0q3fdctjxm4nrc639
```
result:
```bash
{
  "from": "qosacc1x5lcfaqxxq7g7dy4lj5vq0u6xamp78lsnza98y",
  "to": "qosacc1smrus8jlc9z02gz5rm36u0q3fdctjxm4nrc639",
  "qos": "100",
  "qscs": [
    {
      "coin_name": "AOE",
      "amount": "100"
    }
  ]
}
```

#### Increase approve

`qoscli tx increase-approve --from <key_name_or_account_address> --to <account_address> --coins <qos_and_qscs>`

`Arya` add 100QOS,100AOE to the approve for `Sansa`:
```bash
$ qoscli tx increase-approve --from Arya --to qosacc1smrus8jlc9z02gz5rm36u0q3fdctjxm4nrc639 --coins 100QOS,100AOE
Password to sign with 'Arya':<input Arya's keybase password>
```

result：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"3C06676C53A5439D39CB4D0FBA3213C44DC1BA8E","height":"406"}
```

#### Decrease approve

`qoscli tx decrease-approve --from <key_name_or_account_address> --to <account_address> --coins <qos_and_qscs>`

`Arya` reduce 10QOS,10AOE of the approve to `Sansa`:
```bash
$ qoscli tx decrease-approve --from Arya --to qosacc1smrus8jlc9z02gz5rm36u0q3fdctjxm4nrc639 --coins 10QOS,10AOE
Password to sign with 'Arya':<input Arya's keybase  password>
```
result:
```bash
{"check_tx":{},"deliver_tx":{},"hash":"3C06676C53A5439D39CB4D0FBA3213C44DC1BA8E","height":"410"}
```

#### Use approve

`qoscli tx use-approve --from <account_address> --to <key_name_or_account_address> --coins <qos_and_qscs>`

`Sansa` use 10QOS,10AOE of the approve from `Arya`:
```bash
$ qoscli tx use-approve --from Arya --to Sansa --coins 10QOS,10AOE
Password to sign with 'Sansa':<input Sansa's keybase password>
```
result：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"3C06676C53A5439D39CB4D0FBA3213C44DC1BA8E","height":"430"}
```

We can use [query account](#account) to see the latest state of `Arya` and `Sansa`.

#### Cancel approve

`qoscli tx cancel-approve --from <account_address> --to <key_name_or_account_address>`

`Arya` cancel the approve to `Sansa`:
```bash
$ qoscli tx cancel-approve --from Arya --to qosacc1smrus8jlc9z02gz5rm36u0q3fdctjxm4nrc639
Password to sign with 'Arya':<input Arya's keybase password>
```
result:
```bash
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"484"}
```

### QSC

> Before creating QSC you need to apply for a [QSC certification](../spec/ca.md), visit [QSC spec](../spec/qsc) for more information.

commands:
* `qoscli tx create-qsc`    [Create QSC](#create-qsc)
* `qoscli query qsc`        [Query QSC](#query-qsc)
* `qoscli query qsc`        [Query QSCs](#query-qscs)
* `qoscli tx issue-qsc`     [Issue QSC](#issue-qsc)

#### Create QSC

`qoscli tx create-qsc --creator <key_name_or_account_address> --qsc.crt <qsc.crt_file_path> --accounts <account_qsc_s>`

main parameters:

- `--creator`       creator account
- `--qsc.crt`       crt file path
- `--accounts`      initial accounts, [addr1],[amount];[addr2],[amount2],..., optional

`Arya` create QSC token with name of `QOE`:
```bash
$ qoscli tx create-qsc --creator Arya --qsc.crt aoe.crt
Password to sign with 'Arya':<input Arya's keybase password>
```
> Assume `Arya` has `aoe.crt`, the `aoe.crt` contains `banker`'s public key and `banker`'s address is `qosacc1djtcyex03vluga35r8lattddkqt76s7f306xuq` which has been import into keybase as name `ATM`.

result:
```bash
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"200"}
```

#### Query QSC

`qoscli query qsc <qsc_name>`

`qsc_name` is the token name

query `AOE` information:
```bash
$ qoscli query qsc QOE --indent
```
result:
```bash
{
  "name": "AOE",
  "chain_id": "capricorn-1000",
  "extrate": "1:280.0000",
  "description": "",
  "banker": "qosacc1djtcyex03vluga35r8lattddkqt76s7f306xuq"
}
```

#### Query QSCs

`qoscli query qscs`

query all the QSCs:
```bash
$ qoscli query qscs --indent
```
result:
```bash
[
    {
      "name": "AOE",
      "chain_id": "capricorn-1000",
      "extrate": "1:280.0000",
      "description": "",
      "banker": "qosacc1djtcyex03vluga35r8lattddkqt76s7f306xuq"
    }
]
```

#### Issue QSC

After creating the token, you can issue tokens to the address of Banker:

`qoscli tx issue-qsc --qsc-name <qsc_name> --banker <key_name_or_account_address> --amount <qsc_amount>`

main parameters:

- `--qsc-name`  QSC token name
- `--banker`    Banker address or keybase name
- `--amount`    token amount

issue 10000 AOE:

```bash
$ qoscli tx issue-qsc --qsc-name AOE --banker ATM --amount 10000
Password to sign with 'ATM':<input ATM's keybase password>
```

result:
```bash
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"223"}
```

[query account](#account) to view the amount of AOE held in the `ATM` account.

### qcp

QCP is a QOS cross-chain protocol, supporting cross-chain transactions.

You need to apply before creating a federation chain
> Before initialing a QCP chain, you need to apply the [CA](../spec/ca.md). View [QCP spec](../spec/qcp) for more information。

* `qoscli tx init-qcp`: [Init QCP](#init-qcp)
* `qoscli query qcp`:   [Query QCP](#query-qcp)

#### Init QCP

`qoscli tx init-qcp --creator <key_name_or_account_address> --qcp.crt <qcp.crt_file_path>`

main parameters:

- `--creator`       creator assdress
- `--qcp.crt`       crt file path

> Suppose `Arya` has applied for the `qcp.crt` certificate at CA Center, and the chain ID is `aoe-1000` in `qcp.crt`.

`Arya` initializes the chain information in the QOS network:
```bash
$ qoscli tx init-qcp --creator Arya --qcp.crt qcp.crt
Password to sign with 'Arya':<input Arya's keybase password>
```

result:
```bash
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"243"}
```

#### Query QCP

- `qoscli query qcp list`
- `qoscli query qcp out`
- `qoscli query qcp in`
- `qoscli query qcp tx`

See [qbase-Qcp](https://github.com/QOSGroup/qbase/blob/master/docs/client/command.md#Qcp).

### Validator

Visit [stake spec](../spec/stake) learn about validator design.

* `qoscli tx create-validator`          [Create validator](#create-validator)
* `qoscli query validator`              [Query validator](#query-validator)
* `qoscli query validators`             [Query validators](#query-validators)
* `qoscli query validator-miss-vote`    [Query validator miss vote](#query-validator-miss-vote)
* `qoscli query validator-period`       [Validator period](#query-validator-period)
* `qoscli tx modify-validator`          [Modify validator](#modify-validator)
* `qoscli tx revoke-validator`          [Revoke validator](#revoke-validator)
* `qoscli tx active-validator`          [Active validator](#active-validator)

#### Create validator

`qoscli tx create-validator --moniker <validator_name> --owner <key_name_or_account_address> --tokens <tokens>`

main parameters:

- `--owner`         owner keybase name or address
- `--moniker`       name of validator, `len(moniker) <= 300`
- `--nodeHome`      default `$HOME/.qosd`
- `--tokens`        tokens
- `--compound`      whether the income is reinvested, default false
- `--logo`          logo, optional, `len(logo) <= 255`
- `--website`       website, optional,`len(website) <= 255`
- `--details`       description, optional,`len(details) <= 1000`

This transaction need information in `$HOME/.qosd/config/priv_validator.json`, if you have changed the path of this file, please use `--home` to specific it.

`Arya` has started a [full node](../install/testnet.md#setup-a-full-node), she can become a validator by executing the following command:
```bash
$ qoscli tx create-validator --moniker "Arya's node" --owner Arya --tokens 1000
```

result:
```bash
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"258"}
```

1000 QOS will be deducted from `Arya`'s account, bound to the validator node.

#### Modify validator

`qoscli tx modify-validator --owner <key_name_or_account_address> --validator <validator_address> --moniker <validator_name> --logo <logo_url> --website <website_url> --details <description info>`

main parameters:

- `--owner`         owner keybase name or address
- `--validator`     validator address
- `--moniker`       name of validator, `len(moniker) <= 300`
- `--logo`          logo, optional, `len(logo) <= 255`
- `--website`       website, optional,`len(website) <= 255`
- `--details`       description, optional,`len(details) <= 1000`

`Arya` executes `modify-validator` to modify her validator information:
```bash
$ qoscli tx modify-validator --moniker "Arya's node" --owner Arya --validaotor qosval1fzpaxwrmhqml7d90zuzvhmjfxsdqgvzrjpyvsl --logo "https://..." --website "https://..." --description "Long live Arya."
```

result:
```bash
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"265"}
```

#### Query validator

`qoscli query validator [validator-address]`

`validator-address`

query validator by validator address:

```bash
$ qoscli query validator qosval12kjmpgyg23l7axhzzne33jmd0r9y083wzt07hu --indent
```

result:
```bash
{
  "validator": "qosval12kjmpgyg23l7axhzzne33jmd0r9y083wzt07hu",
  "owner": "qosacc12kjmpgyg23l7axhzzne33jmd0r9y083w6mpa33",
  "consensusPubKey": "qosconspub1zcjduepq200x4crydcd6va90l2v2evz0ddjsv88vraqv7jqgds05e4n2xxfqqr6t45",
  "bondTokens": "10000000",
  "description": {
    "moniker": "Arya's node",
    "logo": "https://...",
    "website": "https://...",
    "details": "Long live Arya."
  },
  "commission": {
    "commission_rates": {
      "rate": "0.100000000000000000",
      "max_rate": "0.200000000000000000",
      "max_change_rate": "0.010000000000000000"
    },
    "update_time": "0001-01-01T00:00:00Z"
  },
  "status": "active",
  "InactiveDesc": "",
  "inactiveTime": "0001-01-01T00:00:00Z",
  "inactiveHeight": "0",
  "minPeriod": "0",
  "bondHeight": "0"
}

```

#### Query validators

`qoscli query validators`

query all validators:
```bash
$ qoscli query validators --indent
```

result:
```bash
validators:
[
  {
  "validator": "qosval12kjmpgyg23l7axhzzne33jmd0r9y083wzt07hu",
  "owner": "qosacc12kjmpgyg23l7axhzzne33jmd0r9y083w6mpa33",
  "consensusPubKey": "qosconspub1zcjduepq200x4crydcd6va90l2v2evz0ddjsv88vraqv7jqgds05e4n2xxfqqr6t45",
  "bondTokens": "10000000",
  "description": {
    "moniker": "Arya's node",
    "logo": "https://...",
    "website": "https://...",
    "details": "Long live Arya."
  },
  "commission": {
    "commission_rates": {
      "rate": "0.100000000000000000",
      "max_rate": "0.200000000000000000",
      "max_change_rate": "0.010000000000000000"
    },
    "update_time": "0001-01-01T00:00:00Z"
  },
  "status": "active",
  "InactiveDesc": "",
  "inactiveTime": "0001-01-01T00:00:00Z",
  "inactiveHeight": "0",
  "minPeriod": "0",
  "bondHeight": "0"
}
]
```

#### Query tendermint validators

`qoscli query tendermint-validators <height>`

query validators in the latest height:
```bash
$ qoscli query tendermint-validators --indent
```

result:
```bash
current query height: 100
[
  {
    "Address": "qoscons1stp35jzgkecv9qa2u0vh33gqclnrwu5wekgggz",
    "VotingPower": "10000000",
    "PubKey": "qosconspub1zcjduepq200x4crydcd6va90l2v2evz0ddjsv88vraqv7jqgds05e4n2xxfqqr6t45"
  }
]
```

#### Query validator miss vote

`qoscli query validator-miss-vote [validator-address]`

`validator-address` is the address of validator

query miss vote of `qosval1zlv7vhdcyqyvy9ljdxhcf2766nulwgvys3f6y2`:
```bash
$ qoscli query validator-miss-vote qosval1zlv7vhdcyqyvy9ljdxhcf2766nulwgvys3f6y2
```

result:
```bash
{"startHeight":"258","endHeight":"387","missCount":0,"voteDetail":[]}
```

#### Query validator period
`qoscli query validator-period  <validator-address>`

`validator-address` is the address of validator:

```bash
$ qoscli query validator-period qosval1zlv7vhdcyqyvy9ljdxhcf2766nulwgvys3f6y2
```

result:
```bash
{
  "validator_address": "qosval1zlv7vhdcyqyvy9ljdxhcf2766nulwgvys3f6y2",
  "consensus_pubkey": "qosconspub1zcjduepqma0j9e6ky7h06vgkmghrj3g9u3dwlkhxc8y87x9rn956ue8dulls9g3ax0",
  "fees": "0",
  "current_tokens": "10222379",
  "current_period": "105",
  "last_period": "104",
  "last_period_fraction": {
    "value": "0.000000000000000000"
  }
}
```

#### Query community fee pool
`qoscli query community-fee-pool`

query community fee pool:
```bash
$ qoscli query community-fee-pool
```

result:
```bash
123456
```

#### Revoke validator

`qoscli tx revoke-validator --owner <key_name_or_account_address> --validator <validator_address>`

`key_name_or_account_address` owner address or keybase name
`validator_address` address of validator

`Arya` revokes her validator node:
```bash
$ qoscli tx revoke-validator --owner Arya --validator qosval1zlv7vhdcyqyvy9ljdxhcf2766nulwgvys3f6y2
```

result:
```bash
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"268"}
```

#### Active validator

`qoscli tx active-validator --owner <key_name_or_account_address> --validator <validator_address>`

`key_name_or_account_address` owner address or keybase name
`validator_address` address of validator

`Arya` actives her revoked validator node:
```bash
$ qoscli tx active-validator --owner Arya --validator qosval1zlv7vhdcyqyvy9ljdxhcf2766nulwgvys3f6y2
```

result:
```bash
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"275"}
```



### Delegation

* `qoscli tx delegate`              [Delegate](#delete)
* `qoscli query delegation`         [Query delegation](#query-delegation)
* `qoscli query delegations-to`     [Query validator delegations](#query-validator-delegations)
* `qoscli query delegations`        [Query delegator delegations](#query-delegator-delegations)
* `qoscli query delegator-income`   [Query delegator income](#query-delegator-income)
* `qoscli tx modify-compound`       [Modify compound](#modify-compound)
* `qoscli tx unbond`                [Unbond](#unbond)
* `qoscli query unbondings`         [Query unbondings](#query-unbondings)
* `qoscli tx redelegate`            [Redelegate](#redelegate)
* `qoscli query redelegations`      [Query redelegations](#query-redelegations)

#### Delegate

`qoscli tx delegate --validator <validator_address> --delegator <delegator_key_name_or_account_address> --tokens <tokens> --compound <compound_or_not>`

main parameters:

- `--validator`     address of validator
- `--delegator`     delegator address or keybase name
- `--tokens`        tokens
- `--compound`      whether the income is reinvested, default false

`Sansa` delegates 100 QOS to `Arya`'s validator:
```bash
$ qoscli tx delegate --validator qosval1zlv7vhdcyqyvy9ljdxhcf2766nulwgvys3f6y2 --delegator Sansa --tokens 100
```

#### Query delegation

`qoscli query delegation --validator <validator_address> --delegator <delegator_key_name_or_account_address>`

main parameters:

- `--validator`     address of validator
- `--delegator`     delegator address or keybase name

query the delegation information of `Sansa` on `Arya`'s validator node:
```bash
$ qoscli query delegation --validator qosval1zlv7vhdcyqyvy9ljdxhcf2766nulwgvys3f6y2 --delegator Sansa
```

result:
```bash
{
  "delegator_address": "qosacc12tr0v5uv9xpns79w8q34plakz8gh66855arlrd",
  "validator_address": "qosval12tr0v5uv9xpns79w8q34plakz8gh6685vddu9q",
  "validator_cons_pub_key": "qosconspub1zcjduepqgvpv0ky5nzkt238aus9gh90ct96wm3nddmz7pt2qwyn3lwme8dpskrfu76",
  "delegate_amount": "10000000",
  "is_compound": false
}
```

#### Query validator delegations

`qoscli query delegations-to [validator-address]`

main parameters:

- `validator-address`     address of validator

query all delegations on `Arya`'s validator node:
```bash
$ qoscli query delegations-to qosval12tr0v5uv9xpns79w8q34plakz8gh6685vddu9q
```

result:
```bash
[
  {
    "delegator_address": "qosacc12tr0v5uv9xpns79w8q34plakz8gh66855arlrd",
    "validator_address": "qosval12tr0v5uv9xpns79w8q34plakz8gh6685vddu9q",
    "validator_cons_pub_key": "qosconspub1zcjduepqgvpv0ky5nzkt238aus9gh90ct96wm3nddmz7pt2qwyn3lwme8dpskrfu76",
    "delegate_amount": "10000000",
    "is_compound": false
  }
...
]

```

#### Query delegator delegations

`qoscli query delegations [delegator]`

main parameters:

- `delegator`     delegator address or keybase name

query `Sansa`'s delegations:
```bash
$ qoscli query delegations Sansa
```

result:
```bash
[
  {
    "delegator_address": "qosacc12tr0v5uv9xpns79w8q34plakz8gh66855arlrd",
    "validator_address": "qosval12tr0v5uv9xpns79w8q34plakz8gh6685vddu9q",
    "validator_cons_pub_key": "qosconspub1zcjduepqgvpv0ky5nzkt238aus9gh90ct96wm3nddmz7pt2qwyn3lwme8dpskrfu76",
    "delegate_amount": "10000000",
    "is_compound": false
  }
]

```

#### Query delegator income

`qoscli query delegator-income --validator <validator_address> --delegator <delegator_key_name_or_account_address`

main parameters:

- `--validator`  address of validator
- `--delegator`  delegator address or keybase name

query `Sansa`'s delegation income on `qosval12tr0v5uv9xpns79w8q34plakz8gh6685vddu9q`:
```bash
$ qoscli query delegator-income --delegator Sansa --validator qosval12tr0v5uv9xpns79w8q34plakz8gh6685vddu9q --indent

```

result:
```bash
{
  "validator_address": "qosval12tr0v5uv9xpns79w8q34plakz8gh6685vddu9q",
  "consensus_pubkey": "qosconspub1zcjduepqgvpv0ky5nzkt238aus9gh90ct96wm3nddmz7pt2qwyn3lwme8dpskrfu76",
  "previous_validator_period": "0",
  "bond_token": "10000000",
  "earns_starting_height": "0",
  "first_delegate_height": "0",
  "historical_rewards": "3221431",
  "last_income_calHeight": "0",
  "last_income_calFees": "0"
}

```

#### Modify compound

`qoscli tx modify-compound --validator <validator_address> --delegator <delegator_key_name_or_account_address> --compound <compound_or_not>`

main parameters:

- `--validator`  address of validator
- `--delegator`  delegator address or keybase name
- `--compound`   whether the income is reinvested, default false

change compound value of `Sansa`'s delegation on `qosval12tr0v5uv9xpns79w8q34plakz8gh6685vddu9q`:
```bash
$ qoscli tx modify-compound --delegator Sansa --validator qosval12tr0v5uv9xpns79w8q34plakz8gh6685vddu9q --compound
```

#### Unbond

`qoscli tx unbond --validator <validator_address> --delegator <delegator_key_name_or_account_address> --tokens <tokens> --all <unbond_all>`

main parameters:

- `--validator`  address of validator
- `--delegator`  delegator address or keybase name
- `--tokens`     tokens
- `--all`        whether unbond all, default false

`Sansa` unbond 50 QOS from delegation on `qosval12tr0v5uv9xpns79w8q34plakz8gh6685vddu9q`:
```bash
$ qoscli tx unbond --delegator Sansa --validator qosval12tr0v5uv9xpns79w8q34plakz8gh6685vddu9q --tokens 50
```

#### Query unbondings

`qoscli query unbondings <delegator_key_name_or_account_address>`

query `Sansa`'s unbonding QOS:
```bash
$ qoscli query unbondings Sansa
```

#### Redelegate

`qoscli tx redelegate --from-validator <validator_address> --to-validator <validator_address> --delegator <delegator_key_name_or_account_address> --tokens <tokens> --all <unbond_all>`

main parameters:

- `--from-validator`    origin validator address
- `--to-validator`      target validator address
- `--delegator`         delegator address or keybase name
- `--tokens`            tokens 
- `--compound`          whether the income is reinvested, default false
- `--all`               whether redelegate all tokens, default false

`Sansa` redelegate 10 QOS from `qosval12tr0v5uv9xpns79w8q34plakz8gh6685vddu9q` to `qosval12tr0v5uv9xpns79w8q34plakz8gh6685vddu10`：
```bash
$ qoscli tx redelegate --from-validator qosval12tr0v5uv9xpns79w8q34plakz8gh6685vddu9q --to-validator qosval67werwer98sr76asdf0sdfsd98 --delegator Sansa --tokens 10
```

#### Query redelegations

`qoscli query redelegations <delegator_key_name_or_account_address>`

query redelegations of `Sansa`:
```bash
$ qoscli query redelegations Sansa
```

### Governance

* `qoscli tx submit-proposal`  [Submit proposal](#submit-proposal)
* `qoscli query proposal`      [Query proposal](#query-proposal)
* `qoscli query proposals`     [Query proposals](#query-proposals)
* `qoscli tx deposit`          [Deposit](#deposit)
* `qoscli query deposit`       [Query deposit](#query-deposit)
* `qoscli query deposits`      [Query deposits](#query-deposits)
* `qoscli tx vote`             [Vote](#vote)
* `qoscli query vote`          [Query vote](#query-vote)
* `qoscli query votes`         [Query votes](#query-votes)
* `qoscli query tally`         [Tally](#tally)
* `qoscli query params`        [Query params](#query-params)

#### Submit proposal

`qoscli tx submit-proposal
    --title <proposal_title>
    --proposal-type <proposal_type>
    --proposer <proposer_key_name_or_account_address>
    --deposit <deposit_amount_of_qos>
    --description <description>`

main parameters:

required parameters:

- `--title`             title
- `--proposal-type`     proposal type: `Text`/`ParameterChange`/`TaxUsage`/`ModifyInflation`/`SoftwareUpgrade`
- `--proposer`          proposal address or keybase name
- `--deposit`           initial deposit, must gt `min_deposit * min_proposer_deposit_rate`
- `--description`       description

`TaxUsage` unique parameters:

- `--dest-address`      address for accepting QOS
- `--percent`           percent, ranges (0, 1]

`ParameterChange` unique parameters:

- `--params`            parameters list, `module:key_name:value,module:key_name:value`

`ModifyInflation` unique parameters:

- `inflation`           inflation rules, '[{"end_time":"2023-10-20T00:00:00Z","total_amount":"25500000000000","applied_amount":"0"},{"end_time":"2027-10-20T00:00:00Z","total_amount":"12750000000000","applied_amount":"0"},{"end_time":"2031-10-20T00:00:00Z","total_amount":"6375000000000","applied_amount":"0"},{"end_time":"2035-10-20T00:00:00Z","total_amount":"3187500000000","applied_amount":"0"},{"end_time":"2039-10-20T00:00:00Z","total_amount":"1593750000000","applied_amount":"0"},{"end_time":"2043-10-20T00:00:00Z","total_amount":"796875000000","applied_amount":"0"},{"end_time":"2047-10-20T00:00:00Z","total_amount":"796875000000","applied_amount":"0"}]'
- `total-amount`        total amount of QOS

`SoftwareUpgrade` unique parameters:

- `--version`           version
- `--data-height`       data height
- `--genesis-file`      genesis.json url
- `--genesis-md5`       genesis.json md5
- `--for-zero-height`   hard fork ?

`Arya` submit a simple text proposal:
```bash
$ qoscli tx submit-proposal --title 'update qos' --proposal-type Text --proposer Arya --deposit 10000000 --description 'this is the description'
```

`Arya` submit a `ParameterChange` proposal:
```bash
$ qoscli tx submit-proposal --title 'update parameters' --proposal-type ParameterChange --proposer Arya --deposit 10000000 --description 'this is the description' --params gov:min_deposit:1000
```

Assume `Arya is a guardian.`Arya` submit a `TaxUsage` proposal:
```bash
$ qoscli tx submit-proposal --title 'use tax' --proposal-type TaxUsage --proposer Arya --deposit 10000000 --description 'this is the description' --dest-address Sansa --percent 0.5
```

Arya` submit a `AddInflationPhrase` proposal:
```bash
$ qoscli tx submit-proposal --title 'add inflation phrase' --proposal-type ModifyInflation --proposer Arya --deposit 100000000 --description 'this is the description' --total-amount 10000000000000 --inflation-phrases '[{"end_time":"2023-10-20T00:00:00Z","total_amount":"25500000000000","applied_amount":"0"},{"end_time":"2027-10-20T00:00:00Z","total_amount":"12750000000000","applied_amount":"0"},{"end_time":"2031-10-20T00:00:00Z","total_amount":"6375000000000","applied_amount":"0"},{"end_time":"2035-10-20T00:00:00Z","total_amount":"3187500000000","applied_amount":"0"},{"end_time":"2039-10-20T00:00:00Z","total_amount":"1593750000000","applied_amount":"0"},{"end_time":"2043-10-20T00:00:00Z","total_amount":"796875000000","applied_amount":"0"},{"end_time":"2047-10-20T00:00:00Z","total_amount":"796875000000","applied_amount":"0"}]'
```

`Arya` submit a `SoftwareUpgrade` proposal:
```bash
$ qoscli tx submit-proposal --title 'update qos' --proposal-type SoftwareUpgrade --proposer Arya --deposit 10000000 --description 'upgrade qos to v0.0.6 with genesis file exporting in height 100' --genesis-file "https://.../genesis.json" --data-height 110 --version "0.0.6" --genesis-md5 88c4827158d194116b66b561691e83ef
```

#### Query proposal

`qoscli query proposal <proposal-id>`

query the first proposal:
```bash
$ qoscli query proposal 1 --indent
```

result:
```bash
{
  "proposal_content": {
    "type": "gov/TextProposal",
    "value": {
      "title": "update qos",
      "description": "this is the description",
      "deposit": "100000000"
    }
  },
  "proposal_id": "1",
  "proposal_status": 2,
  "final_tally_result": {
    "yes": "0",
    "abstain": "0",
    "no": "0",
    "no_with_veto": "0"
  },
  "submit_time": "2019-04-03T08:20:34.99523986Z",
  "deposit_end_time": "2019-04-05T08:20:34.99523986Z",
  "total_deposit": "200000000",
  "voting_start_time": "2019-04-03T08:20:34.99523986Z",
  "voting_start_height": "700",
  "voting_end_time": "2019-04-05T08:20:34.99523986Z"
}
```

#### Query proposals

`qoscli query proposals`

query all proposals:
```bash
$ qoscli query proposals
```

result:
```bash
[
  {
    "proposal_content": {
      "type": "gov/TextProposal",
      "value": {
        "title": "update qos",
        "description": "this is the description",
        "deposit": "100000000"
      }
    },
    "proposal_id": "1",
    "proposal_status": 2,
    "final_tally_result": {
      "yes": "0",
      "abstain": "0",
      "no": "0",
      "no_with_veto": "0"
    },
    "submit_time": "2019-04-03T08:20:34.99523986Z",
    "deposit_end_time": "2019-04-05T08:20:34.99523986Z",
    "total_deposit": "200000000",
    "voting_start_time": "2019-04-03T08:20:34.99523986Z",
    "voting_start_height": "700",
    "voting_end_time": "2019-04-05T08:20:34.99523986Z"
  }
]
```

#### Deposit

`qoscli tx deposit --proposal-id <proposal_id> --depositor <depositor_key_name_or_account_address> --amount <amount_of_qos>`

main parameters:

- `--proposal-id`       proposal ID
- `--depositor`         depositor address or keybase name
- `--amount`            amount of qos to deposit

`Arya` deposit 100000 QOS to the first proposal:
```bash
$ qoscli tx deposit --proposal-id 1 --depositor Arya --amount 100000
```

#### Query deposit

`qoscli query deposit <proposal-id> <depositer>`

main parameters:

- `--proposal-id`       proposal ID
- `--depositor`         depositor address or keybase name

query `Arya`'s deposit information to the first proposal:
```bash
$ qoscli query deposit 1 Arya --indent
```

result:
```bash
{
  "depositor": "address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy",
  "proposal_id": "1",
  "amount": "100000000"
}
```

#### Query deposits

`qoscli query deposits <proposal-id>`

main parameters:

- `--proposal-id`       proposal ID

query all deposits to the first proposal: 
```bash
$ qoscli query deposits 1 --indent
```

result:
```bash
[
  {
    "depositor": "address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy",
    "proposal_id": "1",
    "amount": "100000000"
  }
]
```

#### Vote

`qoscli tx vote --proposal-id <proposal_id> --voter <voter_key_name_or_account_address> --option <vote_option>`

main parameters:

- `--proposal-id`       proposal ID
- `--voter`             voter address or keybase name
- `--option`            vote option, `Yes`,`Abstain`,`No`,`NoWithVeto`

`Arya` vote `Yes` to the first proposal:
```bash
$ qoscli tx vote --proposal-id 1 --voter Arya --option Yes
```

#### Query vote

`qoscli query vote <proposal-id> <voter>`

main parameters:

- `--proposal-id`       proposal ID
- `--voter`             voter address or keybase name

query `Arya`'s vote to the first proposal:
```bash
$ qoscli query vote 1 Arya --indent
```

result:
```bash
{
  "voter": "address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy",
  "proposal_id": "1",
  "option": "Yes"
}
```

#### Query votes

`qoscli query votes <proposal-id>`

main parameters:

- `--proposal-id`       proposal ID

query votes:
```bash
$ qoscli query votes 1 --indent
```

result:
```bash
[
  {
    "voter": "address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy",
    "proposal_id": "1",
    "option": "Yes"
  }
]
```

#### Tally

`qoscli query tally <proposal-id>`

main parameters:

- `--proposal-id`       proposal ID

query the real-time statistics on the first proposal:
```bash
$ qoscli query tally 1 --indent
```

result:
```bash
{
  "yes": "100",
  "abstain": "0",
  "no": "0",
  "no_with_veto": "0"
}
```

#### Params

`qoscli query params --module <module> --key <key_name>`

main parameters:

- `--module`       module name: `stake`,`gov`,`distribution`
- `--key`          parameter key

query all parameters:
```bash
$ qoscli query params --indent
[
  {
    "type": "stake",
    "value": {
      "max_validator_cnt": 10,
      "voting_status_len": 100,
      "voting_status_least": 50,
      "survival_secs": 600,
      "unbond_return_height": 10
    }
  },
  {
    "type": "distribution",
    "value": {
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
  {
    "type": "gov",
    "value": {
      "min_deposit": "10000000",
      "max_deposit_period": "172800000000000",
      "voting_period": "172800000000000",
      "quorum": "0.334000000000000000",
      "threshold": "0.500000000000000000",
      "veto": "0.334000000000000000",
      "penalty": "0.000000000000000000"
    }
  }
]
```

query parameters in `gov` mosule:
```bash
$ qoscli query params --module gov --indent
{
  "type": "gov",
  "value": {
    "min_deposit": "10000000",
    "max_deposit_period": "172800000000000",
    "voting_period": "172800000000000",
    "quorum": "0.334000000000000000",
    "threshold": "0.500000000000000000",
    "veto": "0.334000000000000000",
    "penalty": "0.000000000000000000"
  }
}
```

query `min_deposit` in `gov` module:
```bash
$ qoscli query params --module gov --key min_deposit
"10000000"
```


### Inflation

* `qoscli query inflation-phrases`      [Inflation rules](#query-inflation-rules)
* `qoscli query total-inflation  `      [Total inflation](#query-total-inflation)
* `qoscli query total-applied`          [Total applied](#query-total-applied)

#### Query inflation rules

`qoscli query inflation-phrases`

result:
```bash
[
  {
    "end_time": "2023-10-20T00:00:00Z",
    "total_amount": "25500000000000",
    "applied_amount": "0"
  },
  {
    "end_time": "2027-10-20T00:00:00Z",
    "total_amount": "12750000000000",
    "applied_amount": "0"
  },
  {
    "end_time": "2031-10-20T00:00:00Z",
    "total_amount": "6375000000000",
    "applied_amount": "0"
  },
  {
    "end_time": "2035-10-20T00:00:00Z",
    "total_amount": "3187500000000",
    "applied_amount": "0"
  },
  {
    "end_time": "2039-10-20T00:00:00Z",
    "total_amount": "1593750000000",
    "applied_amount": "0"
  },
  {
    "end_time": "2043-10-20T00:00:00Z",
    "total_amount": "796875000000",
    "applied_amount": "0"
  },
  {
    "end_time": "2047-10-20T00:00:00Z",
    "total_amount": "796875000000",
    "applied_amount": "0"
  }
]
```

#### Query total inflation

`qoscli query total-inflation`

result:
```bash
"100000000000000"
```

#### Query total applied

`qoscli query total-applied`

result:
```bash
"49000130122714"
```

### Guardian

* `qoscli query guardian`      [Query guardian](#query-guardian)
* `qoscli query guardians`     [Query guardians](#query-guardians)
* `qoscli tx add-guardian`     [Add guardian](#add-guardian)
* `qoscli tx delete-guardian`  [Delete guardian](#delete-guardian)
* `qoscli tx halt-network`     [Halt network](#halt-network)


#### Query guardian

`qoscli query guardian <guardian_key_name_or_account_address>`

query `Arya`'s information as guardian:
```bash
$ qoscli query guardian Arya --indent
```

result:
```bash
{
  "description": "Arya",
  "guardian_type": 1,
  "address": "qosacc1ctmavdk57x0q7c9t98v7u79607222ars4qczcy",
  "creator": "qosacc1ah9uz0"
}
```

#### Query guardians

`qoscli query guardians`

query all guardians:
```bash
$ qoscli query guardians --indent
```

result:
```bash
[
  {
    "description": "Arya",
    "guardian_type": 1,
    "address": "qosacc1ctmavdk57x0q7c9t98v7u79607222ars4qczcy",
    "creator": "qosacc1ah9uz0"
  }
]
```

#### Add guardian

`Genesis` guardian can execute `add-guardian` to add `Oridinary` guardian.

`qoscli tx add-guardian --address <new_guardian_key_name_or_account_address> --creator <creator_key_name_or_account_address> --description <description>`

main params:

- `--address`         address of guardian to be added
- `--creator`         address of guardian who execute this transaction
- `--description`     description

`Arya` add `Sansa` as guardian:
```bash
$ qoscli tx add-guardian --address Sansa --creator Arya --description 'set Sansa to be a guardian'
```

#### Delete guardian

`Genesis` guardian can execute `delete-guardian` to delete `Oridinary` guardian.

`qoscli tx delete-guardian --address <new_guardian_key_name_or_account_address> --deleted-by <delete_operator_key_name_or_account_address>`

main params:

- `--address`         address of guardian to be removed
- `--deleted-by`      address of guardian who execute this transaction

`Arya` remove `Sansa` from guardians:
```bash
$ qoscli tx delete-guardian --address Sansa --deleted-by Arya
```

#### Halt network

In an emergency, the guardians can stop the network in time.

`qoscli tx halt-network --guardian <guardian_key_name_or_account_address> --reason <reason_for_halting_network>`

main params:

- `--guardian`   address of guardian
- `--reason`     reason for halting the network

A major bug occurred on the network, `Arya` stop the QOS network:
```bash
$ qoscli tx halt-network --guardian Arya --reason 'bug'
```
