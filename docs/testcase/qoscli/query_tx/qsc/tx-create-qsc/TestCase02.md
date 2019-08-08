# Description
```
参数`--creator`，`--qsc.crt`不合法
```
# Input
- 参数`--creator`不存在
```
$ qoscli tx create-qsc --creator starCreator2 --qsc.crt ./qsc-star.crt --indent
```
- 参数`--qcp.crt`不存在
```
$ qoscli tx create-qsc --creator starCreator --qsc.crt ./qsc-star2.crt --indent
```
# Output
- 参数`--creator`不存在
```
$ qoscli tx create-qsc --creator starCreator --qcp.crt ./qcp-star.crt --indent
Password to sign with 'starCreator2':
{
  "check_tx": {
    "code": 1,
    "log": "{\"codespace\":\"sdk\",\"code\":1,\"message\":\"TxStd's ITx ValidateData error:  ERROR:\\nCodespace: qsc\\nCode: 305\\nMessage: \\\"creator not exists\\\"\\n\"}",
    "gasWanted": "100000",
    "gasUsed": "3111"
  },
  "deliver_tx": {},
  "hash": "8E99B4B957B7F090A8879BD56B867AE9351E8E3394D38594AC002A5CB5768E97",
  "height": "0"
}
ERROR: {"codespace":"sdk","code":1,"message":"TxStd's ITx ValidateData error:  ERROR:\nCodespace: qsc\nCode: 305\nMessage: \"creator not exists\"\n"}
```
- 参数`--qsc.crt`不存在
```
$ qoscli tx create-qsc --creator starCreator --qsc.crt ./qsc-star2.crt --indent
MustReadFile failed: open ./qsc-star2.crt: no such file or directory
```