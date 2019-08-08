# Description
```
参数`--creator`，`--qcp.crt`不合法
```
# Input
- 参数`--creator`不存在
```
$ qoscli tx init-qcp --creator starBanker2 --qcp.crt ./qcp-star.crt --indent
```
- 参数`--qcp.crt`不存在
```
$ qoscli tx init-qcp --creator starBanker --qcp.crt ./qcp-star2.crt --indent
```
# Output
- 参数`--creator`不存在
```
$ qoscli tx init-qcp --creator starBanker2 --qcp.crt ./qcp-star.crt --indent
Password to sign with 'starBanker2':
{
  "check_tx": {
    "code": 1,
    "log": "{\"codespace\":\"sdk\",\"code\":1,\"message\":\"TxStd's ITx ValidateData error:  ERROR:\\nCodespace: qcp\\nCode: 404\\nMessage: \\\"creator not exists\\\"\\n\"}",
    "gasWanted": "100000",
    "gasUsed": "1000"
  },
  "deliver_tx": {},
  "hash": "3343F9DD0797ABD7FA6654363AB58DFC42A2E9C0008E428326C5EE311882772D",
  "height": "0"
}
ERROR: {"codespace":"sdk","code":1,"message":"TxStd's ITx ValidateData error:  ERROR:\nCodespace: qcp\nCode: 404\nMessage: \"creator not exists\"\n"}
```
- 参数`--qcp.crt`不存在
```
$ qoscli tx init-qcp --creator starBanker --qcp.crt ./qcp-star2.crt --indent
MustReadFile failed: open ./qcp-star2.crt: no such file or directory
```