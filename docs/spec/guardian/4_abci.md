# ABCI

## BeginBlocker

每一块的开始会判断停网标记。

### 停止网络

存在停网标记时，直接`panic`，全网停止运行。