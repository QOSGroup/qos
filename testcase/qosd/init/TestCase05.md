# Description
```
填写全部参数
```
# Input
```
$ qosd init --moniker testnet --chain-id test-chain-01 --overwrite
```
# Output
```
$ qosd init --moniker testnet --chain-id test-chain-01 --overwrite
{
 "moniker": "testnet",
 "chain_id": "test-chain-01",
 "node_id": "3752c6169c8466f2863ecf99eae53a3e6a0b77c7",
 "gentxs_dir": "",
 "app_message": {
  "gen_txs": null,
  "accounts": null,
  "mint": {
   "params": {
    "inflation_phrases": [
     {
      "end_time": "2023-01-01T00:00:00Z",
      "total_amount": "2500000000000",
      "applied_amount": "0"
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
   "first_block_time": "0",
   "applied_qos_amount": "0"
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
   "validators": null,
   "val_votes_info": null,
   "val_votes_in_window": null,
   "delegators_info": null,
   "delegator_unbond_info": null,
   "redelegations_info": null,
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
   "delegator_earning_info": null,
   "delegator_income_height": null,
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