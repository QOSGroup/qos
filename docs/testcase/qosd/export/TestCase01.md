# Description
```
不指定参数
```
# Input
```
$ qosd export
```
# Output
```
$ qosd export
I[2019-08-01|16:22:23.918] DATA CONFIRM                                 module=main height=106 accounts=24988167 delegations=100000 unbond=0 redelegation=0 feepool=251739 pre=0 valshared=825876 total=26165782 applied=26165782 diff=0
export success: \.qosd\genesis-106-1564647743.json
```
导出的`genesis-106-1564647743.json`文件内容如下:
```
{
 "genesis_time": "2019-08-01T08:10:41.0546339Z",
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
  "gen_txs": null,
  "accounts": [
   {
    "base_account": {
     "account_address": "address1hw43pwhtscealvu973r66vk83gus8myp40fy56",
     "public_key": {
      "type": "tendermint/PubKeyEd25519",
      "value": "heAy23lzdDVvEDXHpkL8A+huCcslZDkLiFcK2Xk9J/E="
     },
     "nonce": "1"
    },
    "qos": "24988167",
    "qscs": null
   }
  ],
  "mint": {
   "params": {
    "inflation_phrases": [
     {
      "end_time": "2023-01-01T00:00:00Z",
      "total_amount": "2500000000000",
      "applied_amount": "25165782"
     },
     {
      "end_time": "2027-01-01T00:00:00Z",
      "total_amount": "12750000000000",
      "applied_amount": "0"
     },
     {
      "end_time": "2031-01-01T00:00:00Z",
      "total_amount": "6375000000000",
      "applied_amount": "0"
     },
     {
      "end_time": "2035-01-01T00:00:00Z",
      "total_amount": "3185000000000",
      "applied_amount": "0"
     }
    ]
   },
   "first_block_time": "1564647041",
   "applied_qos_amount": "26165782"
  },
  "stake": {
   "params": {
    "max_validator_cnt": 10,
    "voting_status_len": 100,
    "voting_status_least": 50,
    "survival_secs": 600,
    "unbond_return_height": 10,
    "redelegation_height": 10
   },
   "validators": [
    {
     "owner": "address1hw43pwhtscealvu973r66vk83gus8myp40fy56",
     "pub_key": {
      "type": "tendermint/PubKeyEd25519",
      "value": "stxAw3cY2oTc5abe/8190af7FXlmxWUz/vIhkn/cgKw="
     },
     "bond_tokens": "100000",
     "description": {
      "moniker": "TestNode",
      "logo": "",
      "website": "",
      "details": ""
     },
     "status": 0,
     "inactive_code": 0,
     "inactive_time": "0001-01-01T00:00:00Z",
     "inactive_height": "0",
     "min_period": "0",
     "bond_height": "0"
    }
   ],
   "val_votes_info": [
    {
     "validator_pub_key": {
      "type": "tendermint/PubKeyEd25519",
      "value": "stxAw3cY2oTc5abe/8190af7FXlmxWUz/vIhkn/cgKw="
     },
     "vote_info": {
      "startHeight": "2",
      "indexOffset": "105",
      "missedBlocksCounter": "0"
     }
    }
   ],
   "val_votes_in_window": null,
   "delegators_info": [
    {
     "delegator_addr": "address1hw43pwhtscealvu973r66vk83gus8myp40fy56",
     "validator_pub_key": {
      "type": "tendermint/PubKeyEd25519",
      "value": "stxAw3cY2oTc5abe/8190af7FXlmxWUz/vIhkn/cgKw="
     },
     "delegate_amount": "100000",
     "is_compound": false
    }
   ],
   "delegator_unbond_info": null,
   "redelegations_info": null,
   "current_validators": [
    {
     "owner": "address1hw43pwhtscealvu973r66vk83gus8myp40fy56",
     "pub_key": {
      "type": "tendermint/PubKeyEd25519",
      "value": "stxAw3cY2oTc5abe/8190af7FXlmxWUz/vIhkn/cgKw="
     },
     "bond_tokens": "100000",
     "description": {
      "moniker": "TestNode",
      "logo": "",
      "website": "",
      "details": ""
     },
     "status": 0,
     "inactive_code": 0,
     "inactive_time": "0001-01-01T00:00:00Z",
     "inactive_height": "0",
     "min_period": "0",
     "bond_height": "0"
    }
   ]
  },
  "qcp": {
   "ca_root_pub_key": null,
   "qcps": []
  },
  "qsc": {
   "ca_root_pub_key": null,
   "qscs": []
  },
  "approve": {
   "approves": []
  },
  "distribution": {
   "community_fee_pool": "251739",
   "last_block_proposer": "address1qfsgxueg9h6vjfqqjdtdjn0k3gwcnu66lwf6lw",
   "pre_distribute_amount": "0",
   "validators_history_period": [
    {
     "validator_pubkey": {
      "type": "tendermint/PubKeyEd25519",
      "value": "stxAw3cY2oTc5abe/8190af7FXlmxWUz/vIhkn/cgKw="
     },
     "period": "8",
     "summary": {
      "value": "202.685640000000000000"
     }
    },
    {
     "validator_pubkey": {
      "type": "tendermint/PubKeyEd25519",
      "value": "stxAw3cY2oTc5abe/8190af7FXlmxWUz/vIhkn/cgKw="
     },
     "period": "9",
     "summary": {
      "value": "215.762140000000000000"
     }
    },
    {
     "validator_pubkey": {
      "type": "tendermint/PubKeyEd25519",
      "value": "stxAw3cY2oTc5abe/8190af7FXlmxWUz/vIhkn/cgKw="
     },
     "period": "10",
     "summary": {
      "value": "228.838640000000000000"
     }
    }
   ],
   "validators_current_period": [
    {
     "validator_pub_key": {
      "type": "tendermint/PubKeyEd25519",
      "value": "stxAw3cY2oTc5abe/8190af7FXlmxWUz/vIhkn/cgKw="
     },
     "current_period_summary": {
      "fees": "784590",
      "period": "11"
     }
    }
   ],
   "delegator_earning_info": [
    {
     "validator_pub_key": {
      "type": "tendermint/PubKeyEd25519",
      "value": "stxAw3cY2oTc5abe/8190af7FXlmxWUz/vIhkn/cgKw="
     },
     "delegator_address": "address1hw43pwhtscealvu973r66vk83gus8myp40fy56",
     "earning_start_info": {
      "previous_period": "10",
      "bond_token": "100000",
      "earns_starting_height": "100",
      "first_delegate_height": "0",
      "historical_rewards": "41286",
      "last_income_calHeight": "100",
      "last_income_calFees": "1376460"
     }
    }
   ],
   "delegator_income_height": [
    {
     "validator_pub_key": {
      "type": "tendermint/PubKeyEd25519",
      "value": "stxAw3cY2oTc5abe/8190af7FXlmxWUz/vIhkn/cgKw="
     },
     "delegator_address": "address1hw43pwhtscealvu973r66vk83gus8myp40fy56",
     "height": "110"
    }
   ],
   "validator_eco_fee_pools": [
    {
     "validator_address": "address1qfsgxueg9h6vjfqqjdtdjn0k3gwcnu66lwf6lw",
     "eco_fee_pool": {
      "proposerTotalRewardFee": "1006588",
      "commissionTotalRewardFee": "239001",
      "preDistributeTotalRewardFee": "23668454",
      "preDistributeRemainTotalFee": "825876"
     }
    }
   ],
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
    "min_deposit": "10",
    "min_proposer_deposit_rate": "0.334000000000000000",
    "max_deposit_period": "172800000000000",
    "voting_period": "172800000000000",
    "quorum": "0.334000000000000000",
    "threshold": "0.500000000000000000",
    "veto": "0.334000000000000000",
    "penalty": "0.000000000000000000",
    "burn_rate": "0.500000000000000000"
   },
   "proposals": null
  },
  "guardian": {
   "guardians": null
  }
 }
}
```