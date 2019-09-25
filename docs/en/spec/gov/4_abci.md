# ABCI

## BeginBlocker

Checking software upgrade flag in the begining of each block. The upgrade process as follow:

1. Hard fork: panic, developers need to download `genesis.json` file and software with the right version, execute`qosd unsafe-reset-all`, then restart the network.
2. General upgrade, panic, developers need to download software with the right version the restart the network.

## EndBlocker

### Proposals in depositing status

Deleting all proposals that in depositing status and reached the `deposit_end_time`, sending all the deposit QOS to community fees.

### Proposals in voting status

Tallying all proposals that in depositing status and reached the `voting_end_time`. When the proposal passed:

- ProposalTypeText, updating the psoposal status only.
- ProposalTypeParameterChange, saving the new parameters.
- ProposalTypeTaxUsage, sending QOS from community fees to the specific address.
- ProposalTypeModifyInflation, saving new inflation phrases.   
- ProposalTypeSoftwareUpgrade, setting software upgrade flag.