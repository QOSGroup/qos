# Description
```
使用种子短语恢复原有密钥
```
# Input
原key name为test01.
使用种子短语恢复原有密钥(test01)到test01.
```
$ qoscli keys add test01 --recover
```
使用种子短语恢复原有密钥(test01)到test02.
```
$ qoscli keys add test02 --recover
```

# Output
使用种子短语恢复原有密钥(test01)到test01.
要求选择是否覆盖原有的test01.
若选择不覆盖，则会跳出创建过程。
```
$ qoscli keys add test01 --recover
override the existing name test01 [y/n]:n
```
若选择覆盖，则继续密钥恢复过程：首先输入新的密码， 然后输入种子短语。
```
$ qoscli keys add test01 --recover
override the existing name test01 [y/n]:y
Enter a passphrase for your key:<输入密码>
Repeat the passphrase:<再次输入密码>
> Enter your recovery seed phrase:
unique denial will mask night riot napkin meadow globe guide upgrade size differ run weekend pitch boring year indoor panic nominee refuse slow flame
NAME:   TYPE:   ADDRESS:                                                PUBKEY:
test01  local   address1qnhak3ph0yqpxar3rrkzuasgnzlzfmq4pyn73m  70UnpxP4b322BJYrf/ZcMBk+eifnNNkUc5kKSBJxM0U=
**Important** write this seed phrase in a safe place.
It is the only way to recover your account if you ever forget your password.
```
使用种子短语恢复原有密钥(test01)到test02.
```
$ qoscli keys add test02 --recover
Enter a passphrase for your key:<输入密码>
Repeat the passphrase:<再次输入密码>
> Enter your recovery seed phrase:
unique denial will mask night riot napkin meadow globe guide upgrade size differ run weekend pitch boring year indoor panic nominee refuse slow flame
NAME:   TYPE:   ADDRESS:                                                PUBKEY:
test02  local   address1qnhak3ph0yqpxar3rrkzuasgnzlzfmq4pyn73m  70UnpxP4b322BJYrf/ZcMBk+eifnNNkUc5kKSBJxM0U=
**Important** write this seed phrase in a safe place.
It is the only way to recover your account if you ever forget your password.
```