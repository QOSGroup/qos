# Description
```
填写必选参数address
```
# Input
```
$ qosd add-guardian --address address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy
```
# Output
第一次调用:
```
$ qosd add-guardian --address address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy

```
命令行无返回值, `genesis.json`文件中`app-state`中`guardian`部分新增:
```
        {
          "description": "",
          "guardian_type": 1,
          "address": "address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy",
          "creator": "address1ah9uz0"
        }
```
第二次调用:
```
$ qosd add-guardian --address address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy
ERROR: guardian: address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy has already exists
```