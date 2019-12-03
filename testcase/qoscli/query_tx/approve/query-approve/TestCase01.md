# Description
```
缺失必须参数`--from`，`--to`
```
# Input
- 缺失必须参数`--from`
```
$ qoscli query approve --to test01
```
- 缺失必须参数`--to`
```
$ qoscli query approve --from test
```
# Output
- 缺失必须参数`--from`
```
$ qoscli query approve --to test01
ERROR: required flag(s) "from" not set
```
- 缺失必须参数`--to`
```
$ qoscli query approve --from test
ERROR: required flag(s) "to" not set
```