# Description
```
导入密钥库中不存在的密钥（从CA PRI文件导入）
```
# Input
```
$ qoscli keys import test03 --file /test03.pri
```
# Output
`--flie` 后面需要Flag参数，用于指定CA PRI文件路径.
- 若无参数, 导入失败：
```
$ qoscli keys import test03 --file
ERROR: flag needs an argument: --file
```
- 若指定的文件路径不存在，导入失败：
``` 
$ qoscli keys import test03 --file /test03.pri
ERROR: open /test03.pri: The system cannot find the file specified.
```
- 若指定的文件路径存在, 导入成功：
```
$ qoscli keys import test03 --file /.qoscli/test03.pri
> Enter a passphrase for your key:<输入密码>
> Repeat the passphrase:<输入密码>

```
-----
注意， 私钥文件`/.qoscli/test03.pri`的内容是export结果中的`privkey`部分。
- `qoscli keys export` 结果： 
``` 
{
    "Name": "test02", 
    "address": "address1qnhak3ph0yqpxar3rrkzuasgnzlzfmq4pyn73m", 
    "pubkey": {
        "type": "tendermint/PubKeyEd25519", 
        "value": "70UnpxP4b322BJYrf/ZcMBk+eifnNNkUc5kKSBJxM0U="
    }, 
    "privkey": {
        "type": "tendermint/PrivKeyEd25519", 
        "value": "dzQ2ii+7KxVLzduw3PusyszjCtF/hgovYo+x4+ugfT7vRSenE/hvfbYElit/9lwwGT56J+c02RRzmQpIEnEzRQ=="
    }
}
```
- 私钥文件`/.qoscli/test03.pri`内容：
``` 
{
    "type": "tendermint/PrivKeyEd25519", 
    "value": "dzQ2ii+7KxVLzduw3PusyszjCtF/hgovYo+x4+ugfT7vRSenE/hvfbYElit/9lwwGT56J+c02RRzmQpIEnEzRQ=="
}
```