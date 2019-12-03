# Description
```
导出已存在的密钥
```
# Input
```
$ qoscli keys export test01
```
# Output
若密码输入正确，导出成功：
```
$ qoscli keys export test01
Password to sign with 'test01':<输入密码>
**Important** Don't leak your private key information to others.
Please keep your private key safely, otherwise your account will be attacked.

{"Name":"test01","address":"address1qnhak3ph0yqpxar3rrkzuasgnzlzfmq4pyn73m","pubkey":{"type":"tendermint/PubKeyEd25519","value":"70UnpxP4b322BJYrf/ZcMBk+eifnNNkUc5kKSBJxM0U="},"privkey":{"type":"tendermint/PrivKeyEd25519","value":"dzQ2ii+7KxVLzduw3PusyszjCtF/hgovYo+x4+ugfT7vRSenE/hvfbYElit/9lwwGT56J+c02RRzmQpIEnEzRQ=="}}
```
若密码输入错误，导出失败：
```
$ qoscli keys export test01
Password to sign with 'test01':<输入密码>
ERROR: Ciphertext decryption failed: Wrong Password
```