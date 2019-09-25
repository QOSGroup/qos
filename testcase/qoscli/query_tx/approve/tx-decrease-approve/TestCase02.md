# Description
```
参数`--coins`不合法
```
# Input
-- `--coins`没有单位
```
$ qoscli tx decrease-approve --coins 10000 --from test --to test01 --indent
```
-- `--from`没有`--coins`指定的QSC
```
$ qoscli tx decrease-approve --coins 10000Star --from test --to test01 --indent
```
# Output
-- `--coins`没有单位
```
$ qoscli tx decrease-approve --coins 10000 --from test --to test01 --indent
null
ERROR: coins str: 10000 parse faild
```
-- `--from`没有`--coins`指定的QSC
```
$ qoscli tx decrease-approve --coins 10000Star --from test --to test01 --indent
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
  "hash": "30FBFC0513B516FADF19D98DEC04629DBDBCFEB07AACD1CF7427267E4278D67E",
  "height": "0"
}
ERROR: {"codespace":"sdk","code":1,"message":"TxStd's ITx ValidateData error:  ERROR:\nCodespace: approve\nCode: 102\nMessage: \"approve contains qsc that not exists\"\n"}
```