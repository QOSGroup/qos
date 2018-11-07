# QCP跨链交易
qos基于[qbase](https://www.github.com/QOSGroup/qbase)，提供跨链交易(QCP)支持

* inseq
Get max sequence received from inChain
```
qoscli qcp inseq --chain-id=xxx
```
* outseq
Get max sequence  to outChain
```
qoscli qcp outseq --chain-id=xxx
```
* outtx
Query qcp out tx
```
qoscli qcp outtx --chain-id=xxx --seq=x
```