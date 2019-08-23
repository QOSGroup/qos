# Description
```
缺失必须参数[chainID]或`--seq`
```
# Input
- 缺失必须参数`[chainID]`
```
$ qoscli query qcp tx
```
- 缺失必须参数`--seq`
```
$ qoscli query qcp tx qcp-star
```
# Output
- 缺失必须参数`[chainID]`
```
$ qoscli query qcp tx
ERROR: accepts 1 arg(s), received 0
```
- 缺失必须参数`--seq`
```
$ qoscli query qcp tx qcp-star
ERROR: required flag(s) "seq" not set
```