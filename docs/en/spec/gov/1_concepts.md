# Concepts

## Proposal procedure

### Deposit

Every proposal have to raise up a certain amount of QOS, defined by `$min_deposit`, before enter voting phrase, in order to be resistant to DDos.
When a proposal having been deposited for more than `$max_deposit_period` but still have not reached `$min_deposit`, the proposal would be dropped and deleted.
For the proposer, the minimum deposit should be `$min_deposit * $min_proposer_deposit_rate`.
A proposal having entered voting phrase will still be open to deposit.

* `min_deposit`: Minimum deposit amount of QOS
* `min_proposer_deposit_rate`: Minimum deposit amount of QOS for proposer
* `max_deposit_period`: How many blocks the deposit period could last at maximum, if the amount does not reach `$min_deposit`.
* `voting_period`: How many blocks the voting could last at maximum, if the tally result does not reach an conclusion.
* `quorum`: The least proportion of voting power participated for a proposal to tally.
* `threshold`: The minimum proportion of `Yes` for a proposal to be `PASS`
* `veto`: The minimum proportion of `Veto` for a proposal to be `NoWithVeto`
* `penalty`: The proportion of voting power to be slashed if a validator does not vote
* `burn_rate`: The proportion of deposit burnt when the tally result of `PASS` or `REJECT`

### Voting

The voting options are：
* `Yes`
* `Abstain`
* `No`
* `NoWithVeto`

Only validators and delegators may vote，for repeated votes only the last one counts.

### Tally

Votes from validators/delegators weights by their bonded QOS。Token holders can only gain voting right by bonding QOS to validator(s).So that QOS bonded is regarded as voting power.

The tally result is calculated as `Voting power of the voter/All Voting power of the network`

Tally result would be：

* Invalid：`Voting power participated/Total Voting power` < `$quorum`

* PASS：Voting power voted `Yes`/Total Voting power > `$threshold`

* NoWithVeto：Voting power voted `NoWithVeto`/Total voting power > `$veto`

* REJECT：Other results other than above.

### burning

Ether the tally result of a proposal is `PASS` or `REJECT`, `Deposit * $burn_rate` of QOS would be burnt as proposal fee，the remaining `Deposit` will return to where it comes。

If the tally result is `NoWithVeto`，the all `Deposit` will go to community fund.

### Slashing

If a validator is bonded on both of a proposal\'s voting start and end heights, but not participated in the vote. Then all QOS it bonded will slashed by proportion of `$penalty`

The participation of a validator is not just the vote from the operator of validator, any vote from the delegators of the validator counts. Similarly, if a validator get slashed for not voting, every delegator bonded will will get slashed by the proportion of bonded QOS.



