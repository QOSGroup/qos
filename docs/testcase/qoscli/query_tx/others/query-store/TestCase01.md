# Description
```
缺失必要参数`--path`与`--data`
```
# Input
- 缺失必要参数`--path`与`--data`：
```
$ qoscli query store
```
- 缺失必要参数`--path`：
```
$ qoscli query store --data account
```
- 缺失必要参数`--data`：
```
$ qoscli query store --path /store/acc/subspace
```
# Output
- 缺失必要参数`--path`与`--data`：
```
$ qoscli query store
ERROR: required flag(s) "data", "path" not set
```
- 缺失必要参数`--path`：
```
$ qoscli query store --data account
ERROR: required flag(s) "path" not set
```
- 缺失必要参数`--data`：
```
$ qoscli query store --path /store/acc/subspace
ERROR: required flag(s) "data" not set
```