# Description
```
缺失必须参数`--coins`，`--from`，`--to`
```
# Input
- 缺失必须参数`--coins`
```
$ qoscli tx decrease-approve --from test --to test01
```
- 缺失必须参数`--from`
```
$ qoscli tx decrease-approve --coins 10000QSC --to test01
```
- 缺失必须参数`--to`
```
$ qoscli tx decrease-approve --coins 10000QSC --from test
```
# Output
- 缺失必须参数`--coins`
```
$ qoscli tx decrease-approve --from test --to test01
ERROR: required flag(s) "coins" not set
```
- 缺失必须参数`--from`
```
$ qoscli tx decrease-approve --coins 10000QSC --to test01
ERROR: required flag(s) "from" not set
```
- 缺失必须参数`--to`
```
$ qoscli tx decrease-approve --coins 10000QSC --from test
ERROR: required flag(s) "to" not set
```