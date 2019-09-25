# Description
```
参数`--coins`不合法
```
# Input
-- `--coins`没有单位
```
$ qoscli tx use-approve --coins 10000 --from test --to test01 --indent
```
-- `--from`没有`--coins`指定的QSC
```
$ qoscli tx use-approve --coins 10000Star --from test --to test01 --indent
```
# Output
-- `--coins`没有单位
```
$ qoscli tx use-approve --coins 10000 --from test --to test01 --indent
null
ERROR: coins str: 10000 parse faild
```
-- `--from`没有`--coins`指定的QSC
```
$ qoscli tx use-approve --coins 10000Star --from test --to test01 --indent
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
  "hash": "84E3E0E89C9293FD81A57BEB0CD6CB87A41A882E950983B1E5BBBC22F938B975",
  "height": "0"
}
ERROR: {"codespace":"sdk","code":1,"message":"TxStd's ITx ValidateData error:  ERROR:\nCodespace: approve\nCode: 102\nMessage: \"approve contains qsc that not exists\"\n"}
```