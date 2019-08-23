# Description
```
缺失必须参数`--creator`，`--qsc.crt`
```
# Input
- 缺失必须参数`--creator`
```
$ qoscli tx create-qsc --qsc.crt ./qsc-star.crt
```
- 缺失必须参数`--qsc.crt`
```
$ qoscli tx create-qsc --creator starCreator
```
# Output
- 缺失必须参数`--creator`
```
$ qoscli tx create-qsc --qsc.crt ./qsc-star.crt
ERROR: required flag(s) "creator" not set
```
- 缺失必须参数`--qsc.crt`
```
$ qoscli tx create-qsc --creator starCreator
ERROR: required flag(s) "qsc.crt" not set
```