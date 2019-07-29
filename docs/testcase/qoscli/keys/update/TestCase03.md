# Description
```
更新密钥库中已存在的密钥
首先输入当前密码，
然后输入并确认新密码。
```
# Input
```
$ qoscli keys update test03
```
# Output
- 当前密码输入错误，更新失败：
```
$ qoscli keys update test03
Enter the current passphrase:<输入密码>
ERROR: Ciphertext decryption failed: Wrong Password
```
- 新密码不足8位，更新失败： 
``` 
$ qoscli keys update test03
Enter the current passphrase:<输入密码>
Enter the new passphrase:<输入密码>
ERROR: password must be at least 8 characters
```
- 新密码两次输入不匹配，更新失败： 
```
$ qoscli keys update test03
Enter the current passphrase:<输入密码>
Enter the new passphrase:<输入密码>
Repeat the new passphrase:<输入密码>
ERROR: passphrases don't match
```
- 当前密码输入正确，且新密码输入正确，则更新成功
```
$ qoscli keys update test03
Enter the current passphrase:<输入密码>
Enter the new passphrase:<输入密码>
Repeat the new passphrase:<输入密码>
Password successfully updated!
```