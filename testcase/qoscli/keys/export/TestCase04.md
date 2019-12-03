# Description
```
导出已存在的密钥（只导出公钥）
```
# Input
```
$ qoscli keys export test01 --pubkey
```
# Output
```
$ qoscli keys export test01 --pubkey
**Important** Don't leak your private key information to others.
Please keep your private key safely, otherwise your account will be attacked.

{"Name":"test01","address":"address1qnhak3ph0yqpxar3rrkzuasgnzlzfmq4pyn73m","pubkey":{"type":"tendermint/PubKeyEd25519","value":"70UnpxP4b322BJYrf/ZcMBk+eifnNNkUc5kKSBJxM0U="},"privkey":null}
```