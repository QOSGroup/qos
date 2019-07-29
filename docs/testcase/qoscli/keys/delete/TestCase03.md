# Description
```
删除已存在的密钥
```
# Input
```
$ qoscli keys delete test02
```
# Output
若密码输入正确，删除成功：
```
$ qoscli keys delete test02
DANGER - enter password to permanently delete key:<输入密码>
Password deleted forever (uh oh!)
```
若密码输入错误，删除失败：
```
$ qoscli keys delete test02
DANGER - enter password to permanently delete key:<输入密码>
ERROR: Ciphertext decryption failed: Wrong Password
```