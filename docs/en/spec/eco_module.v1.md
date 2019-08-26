# QOS Public Chain Economy Module v1

QOS Public chain is a double layered blockchain infrastructure constructed upon [Delegated Proof-of-Stake](https://multicoin.capital/wp-content/uploads/2018/03/DPoS_-Features-and-Tradeoffs.pdf) and [Byzantine fault tolerance](https://en.wikipedia.org/wiki/Byzantine_fault_tolerance)。

![Overview of QOS Economy Module](https://github.com/QOSGroup/static/blob/master/eco_overview.png?raw=true)

Contents:

* [Design Principals](#Design-Principals)
* [Roles](#Roles)
	* [Lightwallets](#Litewallets)

	* [Fullnodes](#Fullnodes)

	* [Validators](#Validators)

	* [Delegators](#Delegators)

* [Modules](#Modules)
	* [Inflation](#QOS-Inflation)

	* [Conmunity Fund](#Conmunity-Fund)

	* [Validate/Delegate](#Validate/Delegate)
		* [Validators](#Validators)

		* [Delegators](#Delegators)

	* [Slashing](#Slashing)
		* [Double signing](#Double-signing)

		* [Not voting in governance](#Not-voting-in-governance)

	* [gas](#gas)

## Design-Principals

* Tokens represent the same right of earning for both validators and delegators.
* Tokens represent the same right of voting for both validators and delegators.
* In case of slashing, not only validators being slashed, but also delegators bonded. Delegations should be commit after investigations.
* A validator could benefit from delegations, which makes its voting power increase. QOS does not set a limit to how much delegation a validator could receive. QOS believes by means of all information and statics of a validator being public, delegators could choose, hence defend the network.

## Roles

### Litewallets

QOS litewallets is decentralized light clients，it can execute transactions qoscli supports, by verifying small set of data in the block headers, a full copy of ledger is not necessary. They are suitable for low system-requirements situations, such as mobile devices.

The official QOS litewallet [Easyzone](https://github.com/QOSGroup/litewallet), also supports other mainstream public chains like Ethereum、Cosmos.

### Fullnodes

Same as other blockchain, QOS fullnodes are nodes holding full copy of ledger of the network.

### Validators

There is a set of validator nodes in QOS public chain, they act as a concrete implementation of the BFT consensus algorithm - each block in the network needs to collect at least 2/3 of the validators signature. Each block in the QOS public chain contains zero or more transactions. Validators verify transactions in the block, and every block passes the verification is signed with validator's private key and gets broadcasted to the network.

A validator must be the full node of the QOS public chain, but a full node needs to send [create certifier transaction] (validators/all_about_validators.md#create-validator) and meet [certain conditions] (validators/all_about_validators.md #How to become a QOS certifier) to become a validator.

A validator is bonded a certain amount of QOSs and undertakes the implementation of the DPOS algorithm. By the number of QOSs it bonds, it gains the benefits of QOS network mining.

To learn more about the validators or to become a validator in QOS, see [Verifiers Node Details] (validators/all_about_validators.md)

### Delegators

For QOS holders who do not have the ability or willingness to run a validator node themselves, but wish to gain the mining revenue, they can choose a validator and contribute their own QOS to the bonding of the validator by delegation, Increasing the validator's voting power and receive the corresponding mining income as a return.
For the calculation of the mining income, see [Validate/Delegate] (#Validate/Delegate)

The delegation could be done in litewallets, fullnodes are not a necessity.

Delegators shares the proceeds of the validator's mining profit, at the same time that they also share the validator's responsibilities and obligations. Once the validator is slashed for be disgust, the client will be punished accordingly.

In community autonomy (the function to be implemented), the principal and the verifier have the same voting rights.

Therefore, even if the full node is not running, the largest number of clients in the network still play an active and important role, that is, they must choose trusted and stable certifiers to increase the voting weight of these certifiers, and pay attention to the certifier's movements. To maintain the security and stability of the network.

## Modules

### QOS-Inflation

According to [QOS whitepaper](https://github.com/QOSGroup/whitepaper)，Inflation of QOS public chain is fixed every 4 years，during the first year of the mainnet launch, the amount of qos created in every block will be roughly the same.

Inflation plan in the mainnet：

Time|The 1st 4 years|The 2nd 4 years|The 3rd 4 years|The 4th 4 years|The 5th 4 years|The 6th 4 years|The 7th 4 years
:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:
Inflation(Million)|2550|1275|637.5|318.75|159.375|79.6875|79.6875

QOS defineds these 4-years as inflation_phrases, comprised of endtime and total_amount, the amount minted is identified by applied_amount. When the endtime is reached, the current inflation_phrases is end, the next inflation_phrases takes over.

A sample of inflation_phrases could be seen in the `mint`-`params`-`inflation_phrases` section of [genesis.json in QOS testnet](https://github.com/QOSGroup/qos-testnets), such as:

```
        "inflation_phrases": [
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

QOS supports modifying inflation_phrases throught governance votings.

Inflations of every blocks：

![Inflations of every blocks](https://github.com/QOSGroup/static/blob/master/rewardPerBlock.png?raw=true)

### Community-Fund

In every inflation of QOS, `$community_reward_rate` of QOSs will be attributed to community funds. Community funds will be used for community operation and construction, rewarding developers, valuable ecological promotion activities (such as a community-recognized QSC consortium chain).

The community funds' account is public and transparent. Any end user of QOS is entitled to initiate a government proposal of type `TaxUsage`, apply for some community funds to a certain QOS address for a claimed usage, and each QOS holder participating in the QOS economy module has the right to vote the proposal.

### Validate/Delegate

#### Validators

There is a set of validators in the QOS public chain, they acts as a concrete implementation of the BFT consensus algorithm - each block in the network needs to collect at least 2/3 of the voting power by validators to get recorded on chain. Each block contains zero or more transactions. Validators verify those transactions, and signed with its own private key if the verification passes and broadcasted to the network.

##### Validators' right

Every block is composed by one of the validators, who is thus a proposer. Proposer gains a $proposer_reward_rate of the total inflation of the block's inflation as an extra reward:

![Proposer's reward](https://github.com/QOSGroup/static/blob/master/proposerReward.png?raw=true)

The chances of proposing a block are proportional to the number of QOS a validator bonds, so the extra reward of the proposer does not change the voting power of each validator in the network.

Validators gain reward from validating, by means of bonding QOSs and undertakes the implementation of the DPOS algorithm.

![Reward for a validator(and its delegators) in one block](https://github.com/QOSGroup/static/blob/master/validatorReward.png?raw=true)

To know more about validators or wish to become one，please see[All about validators](validators/all_about_validators.md).

##### Slashing

In QOS, there is a observation window for every validator, each one of them must sign $ValidatorVotingStatusLeast blocks out of last $ValidatorVotingStatusLen blocks. Or else the validator and its delegator's bonded token will be slashed by a $SlashFractionDowntime portion of QOS.

#### Delegators

The QOS bound on a validator consists of two parts: self-bond tokens from the validator itself and the delegation-bond ones.

**A Validator's total binding (Voting Power) = Self-bond QOSs + ∑ Delegation-bonds**

For a delegator, QOSs it delegated can gain a corresponding proportion of the profit from the validator's total mining reward. Since the validator invested human and material resources, it can extract a certain percentage of commission from the total mining reward.

![Validator reward from one block](https://github.com/QOSGroup/static/blob/master/validatorSelfReward.png?raw=true)

![Deletgator reward from on block](https://github.com/QOSGroup/static/blob/master/delegatorReward.png?raw=true)

* Distribution cycle

After the delegation is created, a *distribution cycle*, whose length is defined by `$delegator_income_period_height` parameter is start. And the benefit/processing request will be only processed at the first block of a new cycle.

The operation of a delegator increasing/decreasing bonded tokens or modifying the configuration parameters, will not affect the current distribution cycle, will take effect at the begining of the next cycle.

The principal modifies the same configuration item multiple times in one cycle (for example, whether to re-invest), and applies to the next period based on the last modification in the period.

A delegator modifies one same configuration multiple times during one same cycle, only the last modification counts.

* Unbonding

Calculation of delegators unbonding by sending `TxUnbondDelegate` transactions

```

     |                x                             y          |
     |  --------------------------------|----------------------|
     |                                 unbond -----------------|-------------------------|

Last Distribution                                      Nest Distribution              Unbonded

```

The validator will add another checking point after an unbond, the amount unbonded will be returned to delegator's account after a period of `$unbond_return_height`.

Unbonding takes effect immediately，but the reward will be calculated first, distributed at the next cycle.

1. Next distribution: reward x + y；
1. when all QOSs bonded is unbonding: y = 0;
1. When only a part of QOSs bonded is unbonding， x > y > 0

* Re-investment

The delegator can specify and later modify whether to re-invest (`is_compound`). The re-investment means that the income generated in the previous cycle is automatically bond and participates in the next cycle. Otherwise, the profit will be automatically transferred to the delegator's account.

Re-investment can continuously and automatically expand the investment scale of entrusted mining, which is a good choice, but it should be noted that the bonded token redemption needs to pass a *freeze period* defined by the parameter `$unbond_return_height` in order to return to the client's account, blindly expanding the size of the delegation bond is not conducive to liquidity.

### Slashing

The reason for the validators being punished is due to its intentional/unintentional malicious behavior or failure to fulfill the obligations of a validator.

On the other hand, QOS network maintenance requires not only verifiers but also delegators. Delegators are not just a passive, only pursuing income role, but is able to actively choose and filter validators, and make his own voice heard in community governances.

Under the guidance of these principals, in QOS, the purpose of punishment is to reduce its bonding QOS, and the QOSs bond from delegators are also subject to proportional punishment.

#### Double-signing

Double-signing means that the same validator signs more than one time at the same height of the network and broadcasts different information to the network.

In the BFT network, double-signing validator is considered a Byzantine node. When the Byzantine node exceeds 2/3, the network will fork, so we regard the double-signing as a serious mistake and impose a higher penalty - destroy the verifier and its clients' bonded token in a hight proportion.

In practice, double-signing are often due to unintentional mistakes, including:
1. Stolen private key
1. Configuration error causes the same validator node starting twice or more

This indicates the security of private key and ability of set-ups and maintainance is critical to the validator.

#### Not-voting-in-governance

In governance voting，if a validator and its delegators are all absent to vote，both will be slashed by a `$penalty` proportion of bonded QOSs。

### gas

Other than hardware and networking costs, QOS defined different gas strategy according to transactions' nature.
