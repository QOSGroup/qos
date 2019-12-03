# test case of qoscli keys add

> `qoscli keys add` 新增密钥

---

## 情景说明

用户还没有钱包地址, 无需任何前提条件.

## 测试命令

【1】name 参数不使用双引号：“”（或单引号：‘’）进行包括

```bash

//正常字符：jack
1.qoscli keys add Jack
2.qoscli keys add Jack_hello

//包含特殊字符：& ! ( )
3.qoscli keys add 23&2
4.qoscli keys add 12!e
5.qoscli keys add 23(234)
...

//包含特殊字符：- + ，* / ? ~ # .
6.qoscli keys add 12,asd

//包含特殊字符：；
7.qoscli keys add 12;ads
```

### 测试结果

* 1.新增成功
* 2.新增成功
* 3.新增失败，且会导致本地keys库崩溃
* 4.新增失败，但不会导致本地keys库崩溃
* 5.新增失败，但不会导致本地keys库崩溃
* 6.新增成功，账户name为12,asd
* 7.新增成功，账户name为12

【2】name 参数使用双引号：“”（或单引号：‘’）进行包括

```bash
    1.qoscli keys add "12q&@#^$*()+-?:><~q;we"
    2.qoscli keys add "12!dsfs"
    3.qoscli keys add "12\!dsfs"
```

### 测试结果

* 1.新增成功
* 2.新增失败，但不会导致本地keys库崩溃
* 3.新增成功，账户name为12\!dsfs
  