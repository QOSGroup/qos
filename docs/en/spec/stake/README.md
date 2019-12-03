# Stake module

## Abstract

This module mainly manages the validator and delegation information.

## Contents

1. **[Concepts](1_concepts.md)**
2. **[State](2_state.md)**
3. **[Transactions](3_txs.md)**
    - [Create validator](3_txs.md#txcreatevalidator)
    - [Modify validator](3_txs.md#txmodifyvalidator)
    - [Revoke validator](3_txs.md#txrevokevalidator)
    - [Active validator](3_txs.md#txactivevalidator)
    - [Create delegate](3_txs.md#txcreatedelegate)
    - [Modify compound](3_txs.md#txmodifycompound)
    - [Unbond delegation](3_txs.md#txunbonddelegation)
    - [Create redelegation](3_txs.md#txcreateredelegation)
4. **[Events](4_events.md)**
    - [Transactions](4_events.md#transactions)
    - [Begin Blocker](4_events.md#beginblocker)
5. **[ABCI](5_abci.md)**
    - [Begin Blocker](5_abci.md#beginblocker)
    - [End Blocker](5_abci.md#endblocker)
6. **[Params](6_params.md)**
