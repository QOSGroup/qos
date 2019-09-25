# test case of qoscli tx unbond

> `qoscli tx unbond` 解除委托

---

## 情景说明

解除委托时选择的代理验证节点和解绑的tokens数量正常，或全部进行解除委托。前提条件：账户def在账户abc创建的验证节点进行过tokens大于10000的委托。

## 测试命令

```bash
    //验证测试结果命令
    qoscli query delegations def
    //解除委托10000
    qoscli tx unbond --owner abc --delegator def --tokens 10000
    //解除所有委托的qos数量
    qoscli tx unbond --owner abc --delegator def --tokens 10000 --all
    //验证测试结果命令
    qoscli query delegations def
```

## 测试结果

```bash
    qoscli query delegations def --indent
    [{"delegator_address":"address16xd8tzrm6f4jfrmtvy6sjafuy80lgj0gjwu8zt","owner_address":"address10xwx06gnrt3dlz7hfrx6a8wx3gyeghxm54rv7a","validator_pub_key":{"type":"tendermint/PubKeyEd25519","value":"aIm1GNnTm6ITZIpt7zfViP9Mc6jLrIF8TZtnSZ3TKB4="},"delegate_amount":"50000","is_compound":false}]  

    qoscli tx unbond --owner abc --delegator def --tokens 10000
    Password to sign with 'def':
    {"check_tx":{"gasWanted":"100000","gasUsed":"9174"},"deliver_tx":{"gasWanted":"100000","gasUsed":"49650","tags":[{"key":"YWN0aW9u","value":"dW5ib25kLWRlbGVnYXRpb24="},{"key":"dmFsaWRhdG9y","value":"YWRkcmVzczE1ZmNrdW1udGNhZ25oZzBocngzNXFmZDNwN2d6ZTZsNmhobHQ3dw=="},{"key":"ZGVsZWdhdG9y","value":"YWRkcmVzczE2eGQ4dHpybTZmNGpmcm10dnk2c2phZnV5ODBsZ2owZ2p3dTh6dA=="}]},"hash":"86340B3E6619072F252018C44CF34EFBCE901CF6768ECDFCFBD531C38303D55F","height":"33817"}

    qoscli tx unbond --owner abc --delegator def --tokens 10000 --all
    Password to sign with 'def':
    {"check_tx":{"gasWanted":"100000","gasUsed":"9174"},"deliver_tx":{"gasWanted":"100000","gasUsed":"49430","tags":[{"key":"YWN0aW9u","value":"dW5ib25kLWRlbGVnYXRpb24="},{"key":"dmFsaWRhdG9y","value":"YWRkcmVzczE1ZmNrdW1udGNhZ25oZzBocngzNXFmZDNwN2d6ZTZsNmhobHQ3dw=="},{"key":"ZGVsZWdhdG9y","value":"YWRkcmVzczE2eGQ4dHpybTZmNGpmcm10dnk2c2phZnV5ODBsZ2owZ2p3dTh6dA=="}]},"hash":"7BCD3D45F89169E00E979F8EDA2883404BE103877A48FA93DF3553B476B2B38D","height":"35962"}

    qoscli query delegations def
    [{"delegator_address":"address16xd8tzrm6f4jfrmtvy6sjafuy80lgj0gjwu8zt","owner_address":"address10xwx06gnrt3dlz7hfrx6a8wx3gyeghxm54rv7a","validator_pub_key":{"type":"tendermint/PubKeyEd25519","value":"aIm1GNnTm6ITZIpt7zfViP9Mc6jLrIF8TZtnSZ3TKB4="},"delegate_amount":"0","is_compound":false}]
```
