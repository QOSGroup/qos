# Concepts

## Guardian

Guardian is an account managed by the QOS Foundation to perform special functions on the normal operation of the QOS network.

There are two types of guardian:
- `Genesis` added in `genesis.json`
- `Ordinary` created by Genesis guardian

Both the `Genesis` and `Ordinary` guardians can execute `TxHaltNetwork`, submit `ProposalTypeTaxUsage` proposal,
and the `Genesis` guardian can also add/delete the `Oridinary` guardian.