# Description
```
参数`--amount`，`--banker`，`--qsc-name`不合法
```
# Input
- 参数`--amount`为负数
```
$ qoscli tx issue-qsc --banker starBanker --qsc-name STAR --amount -10000 --indent
```
- 参数`--amount`为零
```
$ qoscli tx issue-qsc --banker starBanker --qsc-name STAR --amount 0 --indent
```
- 参数`--banker`不存在
```
$ qoscli tx issue-qsc --banker starBanker2 --qsc-name STAR --amount 100000 --indent
```
- 参数`--qsc-name`不存在
```
$ qoscli tx issue-qsc --banker starBanker --qsc-name STAR2 --amount 100000 --indent
```
# Output
- 参数`--amount`为负数
```
$ qoscli tx issue-qsc --banker starBanker --qsc-name STAR --amount -10000 --indent
Password to sign with 'starBanker':
{
  "check_tx": {
    "code": 1,
    "log": "{\"codespace\":\"sdk\",\"code\":1,\"message\":\"TxStd's ITx ValidateData error:  ERROR:\\nCodespace: qsc\\nCode: 301\\nMessage: \\\"invalid tx msg\\\"\\n\"}",
    "gasWanted": "100000"
  },
  "deliver_tx": {},
  "hash": "71721D80E226609A46CD5173147D0E37F34B63FA9F34FEF81D208C37AD55B77C",
  "height": "0"
}
ERROR: {"codespace":"sdk","code":1,"message":"TxStd's ITx ValidateData error:  ERROR:\nCodespace: qsc\nCode: 301\nMessage: \"invalid tx msg\"\n"}
```
- 参数`--amount`为零
```
$ qoscli tx issue-qsc --banker starBanker --qsc-name STAR --amount 0 --indent
Password to sign with 'starBanker':
{
  "check_tx": {
    "code": 1,
    "log": "{\"codespace\":\"sdk\",\"code\":1,\"message\":\"TxStd's ITx ValidateData error:  ERROR:\\nCodespace: qsc\\nCode: 301\\nMessage: \\\"invalid tx msg\\\"\\n\"}",
    "gasWanted": "100000"
  },
  "deliver_tx": {},
  "hash": "9D76691C05CE49FB243E42FFDA0588A6BBDC8D9608E6E7D942DFBF4953AC87C1",
  "height": "0"
}
ERROR: {"codespace":"sdk","code":1,"message":"TxStd's ITx ValidateData error:  ERROR:\nCodespace: qsc\nCode: 301\nMessage: \"invalid tx msg\"\n"}
```
- 参数`--banker`不存在
```
$ qoscli tx issue-qsc --banker starBanker2 --qsc-name STAR --amount 100000 --indent
Password to sign with 'starBanker2':
{
  "check_tx": {
    "code": 1,
    "log": "{\"codespace\":\"sdk\",\"code\":1,\"message\":\"TxStd's ITx ValidateData error:  ERROR:\\nCodespace: qsc\\nCode: 301\\nMessage: \\\"invalid tx msg\\\"\\n\"}",
    "gasWanted": "100000",
    "gasUsed": "1141"
  },
  "deliver_tx": {},
  "hash": "452C808F721DF731C228C58435FF7577F1A0B16EB205AF0935937395833697E2",
  "height": "0"
}
ERROR: {"codespace":"sdk","code":1,"message":"TxStd's ITx ValidateData error:  ERROR:\nCodespace: qsc\nCode: 301\nMessage: \"invalid tx msg\"\n"}
```
- 参数`--qsc-name`不存在
```
$ qoscli tx issue-qsc --banker starBanker --qsc-name STAR2 --amount 100000 --indent
Password to sign with 'starBanker':
{
  "check_tx": {
    "code": 1,
    "log": "{\"codespace\":\"sdk\",\"code\":1,\"message\":\"TxStd's ITx ValidateData error:  ERROR:\\nCodespace: qsc\\nCode: 307\\nMessage: \\\"qsc not exists\\\"\\n\"}",
    "gasWanted": "100000",
    "gasUsed": "1000"
  },
  "deliver_tx": {},
  "hash": "00F80DA4A8E656AC3FB1392894E787711585E4D4ED6C63F0289AF087B90B1DB4",
  "height": "0"
}
ERROR: {"codespace":"sdk","code":1,"message":"TxStd's ITx ValidateData error:  ERROR:\nCodespace: qsc\nCode: 307\nMessage: \"qsc not exists\"\n"}
```