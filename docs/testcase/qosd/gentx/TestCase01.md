# Description
```
缺失必选参数
```
# Input
```
$ qosd gentx
$ qosd gentx --moniker TestNode
$ qosd gentx --owner address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy
$ qosd gentx --tokens 1000
$ qosd gentx --moniker TestNode --owner address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy
$ qosd gentx --moniker TestNode --tokens 1000
$ qosd gentx --owner address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy --tokens 1000
```
# Output
```
$ qosd gentx
ERROR: required flag(s) "moniker", "owner", "tokens" not set
```
```
$ qosd gentx --moniker TestNode
ERROR: required flag(s) "owner", "tokens" not set
```
```
$ qosd gentx --owner address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy
ERROR: required flag(s) "moniker", "tokens" not set
```
```
$ qosd gentx --tokens 1000
ERROR: required flag(s) "moniker", "owner" not set
```
```
$ qosd gentx --moniker TestNode --owner address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy
ERROR: required flag(s) "tokens" not set
```
```
$ qosd gentx --moniker TestNode --tokens 1000
ERROR: required flag(s) "owner" not set
```
```
$ qosd gentx --owner address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy --tokens 1000
ERROR: required flag(s) "moniker" not set
```
