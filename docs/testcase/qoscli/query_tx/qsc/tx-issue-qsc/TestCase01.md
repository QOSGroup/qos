# Description
```
缺失必须参数`--amount`，`--banker`，`--qsc-name`
```
# Input
- 缺失必须参数`--amount`
```
$ qoscli tx issue-qsc --banker starBanker --qsc-name STAR
```
- 缺失必须参数`--banker`
```
$ qoscli tx issue-qsc --qsc-name STAR --amount 10000 
```
- 缺失必须参数`--qsc-name`
```
$ qoscli tx issue-qsc --banker starBanker --amount 10000 
```
# Output
- 缺失必须参数`--amount`
```
$ qoscli tx issue-qsc --banker starBanker --qsc-name STAR
ERROR: required flag(s) "amount" not set
```
- 缺失必须参数`--banker`
```
$ qoscli tx issue-qsc --qsc-name STAR --amount 10000 
ERROR: required flag(s) "banker" not set
```
- 缺失必须参数`--qsc-name`
```
$ qoscli tx issue-qsc --banker starBanker --amount 10000 
ERROR: required flag(s) "qsc-name" not set
```