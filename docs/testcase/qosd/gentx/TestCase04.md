# Description
```
必选参数owner指定的账户不存在
```
# Input
```
$ qosd gentx --moniker TestNode --owner test --tokens 1000
$ qosd gentx --moniker TestNode --owner address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy --tokens 1000
```
# Output
```
$ qosd gentx --moniker TestNode --owner test --tokens 1000
ERROR: Name: test not found
```
```
$ qosd gentx --moniker TestNode --owner address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy --tokens 1000
ERROR: key with address address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy not found
```