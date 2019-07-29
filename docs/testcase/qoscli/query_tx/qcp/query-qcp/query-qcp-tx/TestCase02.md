# Description
```
参数[chainID]或`--seq`不合法
```
# Input
- 参数[chainID]不合法
```
$ qoscli query qcp tx qcp-star2 --seq 1
```
- 参数`--seq`不合法
```
$ qoscli query qcp tx qcp-star --seq 0
```
# Output
- 参数[chainID]不合法
```
$ qoscli query qcp tx qcp-star2 --seq 1
ERROR: GetGetOutChainTx return empty. there is not exists qcp-star2/1 out tx
```
- 参数`--seq`不合法
```
$ qoscli query qcp tx qcp-star --seq 0
ERROR: GetGetOutChainTx return empty. there is not exists qcp-star/0 out tx
```