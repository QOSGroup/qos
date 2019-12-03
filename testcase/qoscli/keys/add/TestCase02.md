# Description
```
添加新密钥
```
# Input
```
$ qoscli keys add test01
```
# Output
两次输入密码一致，操作成功：
```
$ qoscli keys add test01
Enter a passphrase for your key:<输入密码>
Repeat the passphrase:<再次输入密码>
NAME:   TYPE:   ADDRESS:                                                PUBKEY:
test01  local   address1qjuh59naanercd6pnum4gqe34r7mmvumcxtyha  hXEWR+DeNN2LHgKgQ6iCVPxjD+SEq4VwVPWUsPSPlwI=
**Important** write this seed phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

enemy next odor omit crew canvas scrap fatigue worry expand artefact car table moment parrot ozone now explain six disagree buffalo large gadget tank
```
两次输入密码不一致，操作失败：
```
$ qoscli keys add test01
Enter a passphrase for your key:<输入密码>
Repeat the passphrase:<再次输入密码>
ERROR: passphrases don't match
```
输入密码小于8位，操作失败：
```
$ qoscli keys add test01
Enter a passphrase for your key:<输入密码>
Repeat the passphrase:<再次输入密码>
ERROR: password must be at least 8 characters
```
若name已存在，会提示是否覆盖已存在的key.
若选择不覆盖，则会跳出创建过程。
```
$ qoscli keys add test01
override the existing name test01 [y/n]:n
```
若选择覆盖，则继续创建过程，与正常路径相同：
```
$ qoscli keys add test01
override the existing name test01 [y/n]:y
Enter a passphrase for your key:<输入密码>
Repeat the passphrase:<再次输入密码>
NAME:   TYPE:   ADDRESS:                                                PUBKEY:
test01  local   address1qnhak3ph0yqpxar3rrkzuasgnzlzfmq4pyn73m  70UnpxP4b322BJYrf/ZcMBk+eifnNNkUc5kKSBJxM0U=
**Important** write this seed phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

unique denial will mask night riot napkin meadow globe guide upgrade size differ run weekend pitch boring year indoor panic nominee refuse slow flame

```