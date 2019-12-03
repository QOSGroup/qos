# All about QOS validators

> All parameters defined in genesis.json and mentioned in this document begins with $,

## QOS vaidators\' right

* Validating transactions
* Gain minting rewards by its bonded QOS
* Gain commission fee from delegations
* Gain Transaction fee

## QOS validators\' obligation

* Validating transactions
* Stability of validating
* Security of its own private key
* Paticipating governance

## How to become QOS validator

### Hardware requirement

One must become a full-node before becoming a validator. To ensure the stability of validating and broadcasting, necessary security meassure must be taken.
In testnets, we recommend validators to reach hardware requirements as below:

* Both cloud services and individual computers will do, provided that they can be continuously online
* Low-latency public network, bandwidth > 4M
* 1 core CPU, 2G memory, disk > 200G

### Installation and setups

Please see [Join the testnet](../../install/testnet.md)

### Number of validator

In QOS, validators are sorted by voting power from largest to smallest, the first $max_validator_cnt can validate.

### State of validators

![State of validators](https://github.com/QOSGroup/static/blob/master/validator_status.png?raw=true)

* **Active state**

The state of validator continuously validating and broadcasting transactions.

An active validator has all the [validators\' rights](#QOS vaidators\' right)

An ordinary full-node may become an active validator by sending [create-validator](#create-validator) transaction

But not any full-node can become a validator by the method above, since the network has limited the maximum number of validators.

At a certain point of time, validators need to stay active by validating at least `$voting_status_least` of the last `$voting_status_len` blocks. For those can\'t reach such requirement, the state will switch to `Inactive` compulsively

For a new validator or a recently re-actived validator, though the total number of blocks it has chance to validate has not reached `$voting_status_len`, it would be inactivated for already skipping `$voting_status_least` blocks.

* **Inactive state**

An active validator may turn to an inactive validator for not being able to validate at least `$voting_status_least` of the last `$voting_status_len` blocks, or sending [revoke-validator](#revoke-validator) transaction.

The inactive state is the intermediate state between an active state and an unbonded state.

A validator could stay in inactive state for at most `$survival_secs` seconds. if an inactive will unbond after `$survival_secs` without doing anything.

An inactive validator can not benefit validators\' right nor excute validator\'s duty, after `$unbond_return_height`, the QOS bonded will return to delegators\' account.

* **Unbond state**

When a validators unbonds, the QOS bonded on it will return to the delegators, including the self-bonded part.

An unbondded validator is regarded as an full-node.

### voting power of validators

As a DPOS network, validaors of QOS need to bond a certain amount of tokens to represent their voting power.

A validator must have self-bond tokens at the initialization. After the creation, its voting power may contributed by delegators.

For more information of QOS validator, please see [eco module](eco_module.v1.md)
