# Description
```
缺失必选参数`--receivers`和`--senders`
```
# Input
- 缺失必选参数`--receivers`和`--senders`
```
$ qoscli tx transfer
```
- 缺失必选参数`--receivers`
```
$ qoscli tx transfer --senders test,1000QOS
```
- 缺失必选参数`--senders`
```
$ qoscli tx transfer --receivers test01,1000QOS
```
# Output
- 缺失必选参数`--receivers`和`--senders`
```
$ qoscli tx transfer
ERROR: required flag(s) "receivers", "senders" not set
```
- 缺失必选参数`--receivers`
```
$ qoscli tx transfer --senders test,1000QOS
ERROR: required flag(s) "receivers" not set
```
- 缺失必选参数`--senders`
```
$ qoscli tx transfer --receivers test01,1000QOS
ERROR: required flag(s) "senders" not set
```