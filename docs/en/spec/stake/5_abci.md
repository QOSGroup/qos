# ABCI

## BeginBlocker

### Double Sign

When there is an 'ABCIEvidenceTypeDuplicateVote' type of evidence in the evidence pool, verify the evidence based on the parameter 'max_evidence_age' and verify the height of the crime.
For the double-signed node that needs to be punished, calculate the penalty tokens according to the crime height node binding `tokens` and the parameter `slash_fraction_double_sign`, 
and preferentially deduct the relevant `unbondings` and `redelegations` from the validator node after the crime height.
The left portion is deducted from the current delegations of the validator node.

### Vote Info

Record validators' voting information each block.

According to the parameters `voting status len` and `voting status least`,
 the validator is required to participate in at least the `voting status least` block vote at the block height of each `voting status len`.
 Otherwise, deduct `tokens * slash_fraction_downtime` from the validator's current binding tokens.

## EndBlocker

### Unbondings

For the unbounding information that in completing time, return the binding `tokens` to the `delegator` account and delete the unbinding information.

### Redelegations

For the redelegation information that in completing time, create new delegation from `delegator` to `to-validator` and delete the redelegation information.