# Description
```
多个账户, 混杂格式
```
# Input
```
$ qosd add-genesis-accounts address1n4u9hac9gv76xuxpdy9php6cenq8psv6h99cda,10000QOS;address14d3jwnhzmfvnv5gqx7ze204hvh26y8nmwd0q3t,30000QOS,60000AOE
```
# Output
```
$ qosd add-genesis-accounts address1n4u9hac9gv76xuxpdy9php6cenq8psv6h99cda,10000QOS;address14d3jwnhzmfvnv5gqx7ze204hvh26y8nmwd0q3t,30000QOS,60000AOE

```
命令行无返回值, `genesis.json`文件中`app-state`中`accounts`部分新增:
```
      {
        "base_account": {
          "account_address": "address1n4u9hac9gv76xuxpdy9php6cenq8psv6h99cda",
          "public_key": null,
          "nonce": "0"
        },
        "qos": "10000",
        "qscs": []
      },
      {
        "base_account": {
          "account_address": "address14d3jwnhzmfvnv5gqx7ze204hvh26y8nmwd0q3t",
          "public_key": null,
          "nonce": "0"
        },
        "qos": "30000",
        "qscs": [
          {
            "coin_name": "AOE",
            "amount": "60000"
          }
        ]
      }
```