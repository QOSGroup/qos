# test case of qoscli keys add

> `qoscli keys add` 新增密钥

---

## 情景说明

用户已经有账户，但是仍需要增加新账户地址。前提条件：用户已经有了账户：例如abc。用户新增密钥name与abc不一致，参考“新增密钥情景1”；否则满足当前测试情景。

## 测试命令

```bash
    qoscli keys add abc   //创建情景测试的前提条件
    qoscli keys add abc
```

## 测试结果

```bash
//
qoscli keys add abc
Enter a passphrase for your key:
Repeat the passphrase:
NAME:   TYPE:   ADDRESS:            PUBKEY:
abc local    address1n3g8c06yff0uys6n0gg9qd8r74zq48dhra3x79 UhYPH0rdbDJZHvB5n0J3A/rDs37r4JKZgVBa5WRa9WI=
**Important** write this seed phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

bring crouch loop elegant patrol share below piano round prison sister team soon raw lecture screen glare fat achieve wool town muffin tilt vote

qoscli keys add abc
override the existing name abc [y/n]:y
Enter a passphrase for your key:
Repeat the passphrase:
NAME:   TYPE:   ADDRESS:                        PUBKEY:
abc local   address1av27ppzk58njjczhy58wsc20a45m29hkhzze55  qD9xBdI6eKHnMjIsvSvCwQW8mwO/VOlyfAQO+IH+z6A=
**Important** write this seed phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

slice dutch fancy couple nothing dry vote segment dose give banner number odor shuffle staff use error inch list correct fiber flight sniff pattern

```

ps:
    从以上结果可以看出：会在QOS网络中产生两个账户address1n3g8c06yff0uys6n0gg9qd8r74zq48dhra3x79和address1av27ppzk58njjczhy58wsc20a45m29hkhzze55，但是当你以name：abc为关键字去查询时候，只能查到最后一次被覆盖后映射的地址address1av27ppzk58njjczhy58wsc20a45m29hkhzze55。之前产生的账户信息，只能够以地址为关键字去查询账户信息。
