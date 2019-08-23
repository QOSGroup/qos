# Description
```
指定的`--receivers`和`--senders`在密钥库中不存在
```
# Input
- 指定的`--receivers`在密钥库中不存在
```
$ qoscli tx transfer --senders test,1000QOS  --receivers test04,1000QOS
```
- 指定的`--senders`在密钥库中不存在
```
$ qoscli tx transfer --senders test04,1000QOS  --receivers test01,1000QOS
```
# Output
- 指定的`--receivers`在密钥库中不存在
```
$ qoscli tx transfer --senders test,1000QOS  --receivers test04,1000QOS
null
ERROR: Name: test04 not found
```
- 指定的`--senders`在密钥库中不存在
```
$ qoscli tx transfer --senders test04,1000QOS  --receivers test01,1000QOS
null
ERROR: Name: test04 not found
```