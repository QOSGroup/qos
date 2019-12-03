# Description
```
参数`--coins`不合法
```
# Input
-- `--coins`没有单位
```
$ qoscli tx create-approve --coins 10000 --from test --to test01
```
-- `--from`没有`--coins`指定的QSC
```
$ qoscli tx create-approve --coins 10000Star --from test --to test01
```
# Output
-- `--coins`没有单位
```
$ qoscli tx create-approve --coins 10000 --from test --to test01
null
ERROR: coins str: 10000 parse faild
```
-- `--from`没有`--coins`指定的QSC
```
$ qoscli tx create-approve --coins 10000Star --from test --to test01 --indent
Password to sign with 'test':
{
  "check_tx": {
    "code": 1,
    "log": "{\"codespace\":\"sdk\",\"code\":1,\"message\":\"TxStd's ITx ValidateData error:  ERROR:\\nCodespace: approve\\nCode: 102\\nMessage: \\\"approve contains qsc that not exists\\\"\\n\"}",
    "gasWanted": "100000",
    "gasUsed": "1000",
    "events": []
  },
  "deliver_tx": {},
  "hash": "89117BAA1D2DF16D8032760D17EC6E6D657559B99451725C35DD4F7BF08A10FB",
  "height": "0"
}
ERROR: {"codespace":"sdk","code":1,"message":"TxStd's ITx ValidateData error:  ERROR:\nCodespace: approve\nCode: 102\nMessage: \"approve contains qsc that not exists\"\n"}
```