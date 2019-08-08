# Description
```
指定参数[height]
```
# Input
- 指定的参数[height] < 0
```
$ qoscli query block -100
```
- 指定的参数[height] == 0
```
$ qoscli query block 0
```
- 指定的参数[height] > 0 且 <= 区块链当前高度
```
$ qoscli query block 100
```
- 指定的参数[height] > 区块链当前高度
```
$ qoscli query block -100
```
# Output
- 指定的参数[height] < 0
```
$ qoscli query block -100
ERROR: unknown shorthand flag: '1' in -100
```
- 指定的参数[height] == 0
```
$ qoscli query block 0
ERROR: Block: response error: RPC error -32603 - Internal error: Height must be greater than 0
```
- 指定的参数[height] > 0 且 <= 区块链当前高度 
  - 原始输出：
  ```
    $ qoscli query block 100
    {"block_meta":{"block_id":{"hash":"8753579A1FFA9BA96B3EF4AB346920045F4551C584DB4372705B5E98DBDEB364","parts":{"total":"1","hash":"5510BEDEA1EF7D3C0F6B47389CE3AA8EA4D622DEE7EE39D20CF8D7AF8809E4FD"}},"header":{"version":{"block":"10","app":"0"},"chain_id":"test-chain","height":"100","time":"2019-08-01T08:20:53.9023678Z","num_txs":"0","total_txs":"0","last_block_id":{"hash":"BC980B00D54BFE46C953D4616B96B2853555974CA092824726BE9601AB68F102","parts":{"total":"1","hash":"3D40CD62B503D650538328402151E04349135611844D99F53FBB9DD1C5B66391"}},"last_commit_hash":"349D7D0B1B0F413427A81E63C7511FBE83A6E441EA75DD92939D6AAF3D7B29FE","data_hash":"","validators_hash":"7A3AFBB37239FCE14B64AED973DA5D48C9292AFF80FD674A577214B06BEF2E15","next_validators_hash":"7A3AFBB37239FCE14B64AED973DA5D48C9292AFF80FD674A577214B06BEF2E15","consensus_hash":"294D8FBD0B94B767A7EBA9840F299A3586DA7FE6B5DEAD3B7EECBA193C400F93","app_hash":"72C87A44182FE9585174267DE8B808AD82E9D26931E3A846FAC33B87BA6E8BFE","last_results_hash":"","evidence_hash":"","proposer_address":"02608373282DF4C924009356D94DF68A1D89F35A"}},"block":{"header":{"version":{"block":"10","app":"0"},"chain_id":"test-chain","height":"100","time":"2019-08-01T08:20:53.9023678Z","num_txs":"0","total_txs":"0","last_block_id":{"hash":"BC980B00D54BFE46C953D4616B96B2853555974CA092824726BE9601AB68F102","parts":{"total":"1","hash":"3D40CD62B503D650538328402151E04349135611844D99F53FBB9DD1C5B66391"}},"last_commit_hash":"349D7D0B1B0F413427A81E63C7511FBE83A6E441EA75DD92939D6AAF3D7B29FE","data_hash":"","validators_hash":"7A3AFBB37239FCE14B64AED973DA5D48C9292AFF80FD674A577214B06BEF2E15","next_validators_hash":"7A3AFBB37239FCE14B64AED973DA5D48C9292AFF80FD674A577214B06BEF2E15","consensus_hash":"294D8FBD0B94B767A7EBA9840F299A3586DA7FE6B5DEAD3B7EECBA193C400F93","app_hash":"72C87A44182FE9585174267DE8B808AD82E9D26931E3A846FAC33B87BA6E8BFE","last_results_hash":"","evidence_hash":"","proposer_address":"02608373282DF4C924009356D94DF68A1D89F35A"},"data":{"txs":null},"evidence":{"evidence":null},"last_commit":{"block_id":{"hash":"BC980B00D54BFE46C953D4616B96B2853555974CA092824726BE9601AB68F102","parts":{"total":"1","hash":"3D40CD62B503D650538328402151E04349135611844D99F53FBB9DD1C5B66391"}},"precommits":[{"type":2,"height":"99","round":"0","block_id":{"hash":"BC980B00D54BFE46C953D4616B96B2853555974CA092824726BE9601AB68F102","parts":{"total":"1","hash":"3D40CD62B503D650538328402151E04349135611844D99F53FBB9DD1C5B66391"}},"timestamp":"2019-08-01T08:20:53.9023678Z","validator_address":"02608373282DF4C924009356D94DF68A1D89F35A","validator_index":"0","signature":"oNRTIo4cU2y6/ielWCoGY8/5ayTWSr0PvxMqllcAcZ8Hh92V4qaLzg2pezq+PUITC+mqMMMNkwyILsfSfCEKDA=="}]}}}
  ```
  - 格式化输出：
  ```
    $ qoscli query block 100 --indent
    {
      "block_meta": {
        "block_id": {
          "hash": "8753579A1FFA9BA96B3EF4AB346920045F4551C584DB4372705B5E98DBDEB364",
          "parts": {
            "total": "1",
            "hash": "5510BEDEA1EF7D3C0F6B47389CE3AA8EA4D622DEE7EE39D20CF8D7AF8809E4FD"
          }
        },
        "header": {
          "version": {
            "block": "10",
            "app": "0"
          },
          "chain_id": "test-chain",
          "height": "100",
          "time": "2019-08-01T08:20:53.9023678Z",
          "num_txs": "0",
          "total_txs": "0",
          "last_block_id": {
            "hash": "BC980B00D54BFE46C953D4616B96B2853555974CA092824726BE9601AB68F102",
            "parts": {
              "total": "1",
              "hash": "3D40CD62B503D650538328402151E04349135611844D99F53FBB9DD1C5B66391"
            }
          },
          "last_commit_hash": "349D7D0B1B0F413427A81E63C7511FBE83A6E441EA75DD92939D6AAF3D7B29FE",
          "data_hash": "",
          "validators_hash": "7A3AFBB37239FCE14B64AED973DA5D48C9292AFF80FD674A577214B06BEF2E15",
          "next_validators_hash": "7A3AFBB37239FCE14B64AED973DA5D48C9292AFF80FD674A577214B06BEF2E15",
          "consensus_hash": "294D8FBD0B94B767A7EBA9840F299A3586DA7FE6B5DEAD3B7EECBA193C400F93",
          "app_hash": "72C87A44182FE9585174267DE8B808AD82E9D26931E3A846FAC33B87BA6E8BFE",
          "last_results_hash": "",
          "evidence_hash": "",
          "proposer_address": "02608373282DF4C924009356D94DF68A1D89F35A"
        }
      },
      "block": {
        "header": {
          "version": {
            "block": "10",
            "app": "0"
          },
          "chain_id": "test-chain",
          "height": "100",
          "time": "2019-08-01T08:20:53.9023678Z",
          "num_txs": "0",
          "total_txs": "0",
          "last_block_id": {
            "hash": "BC980B00D54BFE46C953D4616B96B2853555974CA092824726BE9601AB68F102",
            "parts": {
              "total": "1",
              "hash": "3D40CD62B503D650538328402151E04349135611844D99F53FBB9DD1C5B66391"
            }
          },
          "last_commit_hash": "349D7D0B1B0F413427A81E63C7511FBE83A6E441EA75DD92939D6AAF3D7B29FE",
          "data_hash": "",
          "validators_hash": "7A3AFBB37239FCE14B64AED973DA5D48C9292AFF80FD674A577214B06BEF2E15",
          "next_validators_hash": "7A3AFBB37239FCE14B64AED973DA5D48C9292AFF80FD674A577214B06BEF2E15",
          "consensus_hash": "294D8FBD0B94B767A7EBA9840F299A3586DA7FE6B5DEAD3B7EECBA193C400F93",
          "app_hash": "72C87A44182FE9585174267DE8B808AD82E9D26931E3A846FAC33B87BA6E8BFE",
          "last_results_hash": "",
          "evidence_hash": "",
          "proposer_address": "02608373282DF4C924009356D94DF68A1D89F35A"
        },
        "data": {
          "txs": null
        },
        "evidence": {
          "evidence": null
        },
        "last_commit": {
          "block_id": {
            "hash": "BC980B00D54BFE46C953D4616B96B2853555974CA092824726BE9601AB68F102",
            "parts": {
              "total": "1",
              "hash": "3D40CD62B503D650538328402151E04349135611844D99F53FBB9DD1C5B66391"
            }
          },
          "precommits": [
            {
              "type": 2,
              "height": "99",
              "round": "0",
              "block_id": {
                "hash": "BC980B00D54BFE46C953D4616B96B2853555974CA092824726BE9601AB68F102",
                "parts": {
                  "total": "1",
                  "hash": "3D40CD62B503D650538328402151E04349135611844D99F53FBB9DD1C5B66391"
                }
              },
              "timestamp": "2019-08-01T08:20:53.9023678Z",
              "validator_address": "02608373282DF4C924009356D94DF68A1D89F35A",
              "validator_index": "0",
              "signature": "oNRTIo4cU2y6/ielWCoGY8/5ayTWSr0PvxMqllcAcZ8Hh92V4qaLzg2pezq+PUITC+mqMMMNkwyILsfSfCEKDA=="
            }
          ]
        }
      }
    }
  ```
- 指定的参数[height] > 区块链当前高度
```
$ qoscli query block 10000
ERROR: Block: response error: RPC error -32603 - Internal error: Height must be less than or equal to the current blockchain height
```