# Description
```
缺失必须参数`--creator`，`--qcp.crt`
```
# Input
- 缺失必须参数`--creator`
```
$ qoscli tx init-qcp --qcp.crt ./qcp-star.crt
```
- 缺失必须参数`--qcp.crt`
```
$ qoscli tx init-qcp --creator starBanker
```
# Output
- 缺失必须参数`--creator`
```
$ qoscli tx init-qcp --qcp.crt ./qcp-star.crt
ERROR: required flag(s) "creator" not set
```
- 缺失必须参数`--qcp.crt`
```
$ qoscli tx init-qcp --creator starBanker
ERROR: required flag(s) "qcp.crt" not set
```