# Description
```
参数`--from`，`--to`不合法
```
# Input
- 参数`--from`不合法
```
$ qoscli query approve --from test05 --to test01 --indent
```
- 参数`--to`不合法
```
$ qoscli query approve --from test --to test05 --indent
```
# Output
- 参数`--from`不合法
```
$ qoscli query approve --from test05 --to test01 --indent
ERROR: Name: test05 not found
```
- 参数`--to`不合法
```
$ qoscli query approve --from test --to test05 --indent
ERROR: Name: test05 not found
```