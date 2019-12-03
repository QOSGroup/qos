# ABCI

## EndBlocker

The end of each block will check the invariant-check flag, query the lock-info and release locked QOS.

### Invariant Check

Add an 'EventTypeInvariantCheck` event when there is a invariant check request:
```go
if NeedInvariantCheck(ctx) {
	ctx.EventManager().EmitEvent(btypes.NewEvent(qtypes.EventTypeInvariantCheck))
}
```
In
```go
func (app *QOSApp) EndBlocker(ctx context.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
    ...
}
```
checking whether there is a invariant check event to perform a network-wide data check.

### Lock-Release

If it exists [LockInfo](2_state.md#lockinfo), the release time is reached, a release operation will performed.

First, calculate release amount in this time
```go
if ReleaseTimes != 1 {
    releaseAmount = (TotalAmount - ReleasedAmount) / ReleaseTimes
} else {
    releaseAmount = TotalAmount - ReleasedAmount
}
```

then
- `LockedAccount` minus `releaseAmount`
- `Receiver` plus  `releaseAmount`
- `ReleasedAmount` plus `releaseAmount`
- `ReleaseTimes` = `ReleaseTimes - 1`
- `ReleaseTime` = `ReleaseTime + ReleaseInterval`

If the lock-release is all completed, the lock-release record is deleted.