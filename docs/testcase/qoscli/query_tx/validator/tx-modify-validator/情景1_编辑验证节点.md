# test case of qoscli tx modify-validator

> `qoscli tx modify-validator` 编辑验证节点

---

## 情景说明

验证节点的操作者，对验证节点的信息进行修改。增加或修改logo，网址，描述信息等。

## 测试命令

```bash
    //执行modify语句前 查询验证节点信息
    qoscli query validator jlgy01 --indent

    qoscli tx modify-validator --moniker jlgy666 --owner jlgy01 --logo "http://pic32.nipic.com/20130813/3347542_160503703000_2.jpg" --website "https://github.com/test" --details "jlgy23333333333"

    //执行modify语句后 查询验证节点信息
    qoscli query validator jlgy01 --indent
```

## 测试结果

```bash
    qoscli query validator jlgy01 --indent
    {
    "owner": "address1nnvdqefva89xwppzs46vuskckr7klvzk8r5uaa",
    "validatorAddress": "6E713D1F3CCE28D820C39059E5D08D21D646FA8E",
    "validatorPubkey": {
        "type": "tendermint/PubKeyEd25519",
        "value": "exGS/yWJthwY8za4dlrPRid2I9KE4G15nlJwO/+Off8="
    },
    "bondTokens": "2000000000",
    "description": {
        "moniker": "jlgy",
        "logo": "",
        "website": "",
        "details": ""
    },
    "status": "active",
    "InactiveDesc": "",
    "inactiveTime": "0001-01-01T00:00:00Z",
    "inactiveHeight": "0",
    "bondHeight": "617422"
    }

    qoscli tx modify-validator --moniker jlgy666 --owner jlgy01 --logo "http://pic32.nipic.com/20130813/3347542_160503703000_2.jpg" --website "https://github.com/wangkuanzzu" --details "jlgy23333333333"
    Password to sign with 'jlgy01':
    {"check_tx":{"gasWanted":"100000","gasUsed":"6703"},"deliver_tx":{"gasWanted":"100000","gasUsed":"17160","tags":[{"key":"YWN0aW9u","value":"bW9kaWZ5LXZhbGlkYXRvcg=="},{"key":"b3duZXI=","value":"YWRkcmVzczFubnZkcWVmdmE4OXh3cHB6czQ2dnVza2NrcjdrbHZ6azhyNXVhYQ=="},{"key":"ZGVsZWdhdG9y","value":"YWRkcmVzczFubnZkcWVmdmE4OXh3cHB6czQ2dnVza2NrcjdrbHZ6azhyNXVhYQ=="}]},"hash":"241AC66206A955AA44AE7A2555EAA9D17320241A700D4749AF74055EEC064C57","height":"617704"}

    qoscli query validator jlgy01 --indent
    {
    "owner": "address1nnvdqefva89xwppzs46vuskckr7klvzk8r5uaa",
    "validatorAddress": "6E713D1F3CCE28D820C39059E5D08D21D646FA8E",
    "validatorPubkey": {
        "type": "tendermint/PubKeyEd25519",
        "value": "exGS/yWJthwY8za4dlrPRid2I9KE4G15nlJwO/+Off8="
    },
    "bondTokens": "2000000000",
    "description": {
        "moniker": "jlgy666",
        "logo": "http://pic32.nipic.com/20130813/3347542_160503703000_2.jpg",
        "website": "https://github.com/wangkuanzzu",
        "details": "jlgy23333333333"
    },
    "status": "active",
    "InactiveDesc": "",
    "inactiveTime": "0001-01-01T00:00:00Z",
    "inactiveHeight": "0",
    "bondHeight": "617422"
    }

```
