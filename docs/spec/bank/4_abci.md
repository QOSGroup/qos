# ABCI

## EndBlocker

每一块的结束会判断数据检查标志，查询锁定账户释放锁定QOS。

### 数据检查

存在数据检查请求时向Event中添加`EventTypeInvariantCheck`事件
```go
if NeedInvariantCheck(ctx) {
	ctx.EventManager().EmitEvent(btypes.NewEvent(qtypes.EventTypeInvariantCheck))
}
```
QOS网络会在
```go
func (app *QOSApp) EndBlocker(ctx context.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
    ...
}
```
中判断是否存在数据检查事件，从而执行全网数据检查。

### 锁定-释放

如果存在[锁定-释放信息](2_state.md#锁定-释放账户)，到达释放时间，且未释放QOS大于零，会执行释放操作。

当次释放QOS `releaseAmount`等于(待释放总量`TotalAmount`减去已释放总量`ReleasedAmount`)除以剩余释放次数`ReleaseTimes`，如果是最后一次则等于待释放总量`TotalAmount`减去已释放总量`ReleasedAmount`。

释放操作将从锁定账户`LockedAccount`扣除本次释放量`releaseAmount`，累加到接收账户`Receiver`，同时更新已释放总量`ReleasedAmount`，剩余释放次数`ReleaseTimes`减一，释放时间`ReleaseTime`增加一个释放周期`ReleaseInterval`。

如果锁定-释放全部完成，会删除锁定-释放记录。