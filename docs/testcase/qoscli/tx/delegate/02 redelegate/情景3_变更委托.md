# test case of qoscli tx delegate

> `qoscli tx delegate` 委托

---

## 情景说明

变更委托的账户是验证节点，且变更的的tokens小于在当前委托人验证节点绑定的tokens数量，且自身账户足够支付gas。

## 测试命令

```bash
    qoscli tx redelegate --from-owner address1nnvdqefva89xwppzs46vuskckr7klvzk8r5uaa --to-owner address1f66wr25emjtp5urfcpd02epwg5ply3xzcv2u20 --delegator jlgy02 --tokens 100 --max-gas 200000

    qoscli tx redelegate --from-owner address1nnvdqefva89xwppzs46vuskckr7klvzk8r5uaa --to-owner address1f66wr25emjtp5urfcpd02epwg5ply3xzcv2u20 --delegator jlgy02 --tokens 100 --max-gas 200000 --all
```

## 测试结果

```bash
    qoscli tx redelegate --from-owner address1nnvdqefva89xwppzs46vuskckr7klvzk8r5uaa --to-owner address1f66wr25emjtp5urfcpd02epwg5ply3xzcv2u20 --delegator jlgy02 --tokens 100 --max-gas 200000
    Password to sign with 'jlgy02':
    {"check_tx":{"gasWanted":"200000","gasUsed":"11852"},"deliver_tx":{"gasWanted":"200000","gasUsed":"91590","tags":[{"key":"YWN0aW9u","value":"Y3JlYXRlLXJlZGVsZWdhdGlvbg=="},{"key":"dmFsaWRhdG9y","value":"YWRkcmVzczFkZWNuNjhldWVjNWRzZ3hyanB2N3Q1eWR5OHR5ZDc1dzhncnlhZg=="},{"key":"bmV3LXZhbGlkYXRvcg==","value":"YWRkcmVzczFoeHl1dDJkeXZydnh1bGZ1OGZsYXl0MDl3eWhxN3IwNG05OGx2Ng=="},{"key":"ZGVsZWdhdG9y","value":"YWRkcmVzczFsMHduNjZnaDQ1bmZ0YTJyNHZxOHo1NHd1OWhnYXJzczI5OGU5Zw=="}]},"hash":"F9A859D542D2F2F28E6A6828750A1FA699BDABB3B457040994AD5DD66E92B7FB","height":"577817"}

    qoscli tx redelegate --from-owner address1nnvdqefva89xwppzs46vuskckr7klvzk8r5uaa --to-owner address1f66wr25emjtp5urfcpd02epwg5ply3xzcv2u20 --delegator jlgy02 --tokens 100 --max-gas 200000 --all
    Password to sign with 'jlgy02':
    {"check_tx":{"gasWanted":"200000","gasUsed":"11852"},"deliver_tx":{"gasWanted":"200000","gasUsed":"81910","tags":[{"key":"YWN0aW9u","value":"Y3JlYXRlLXJlZGVsZWdhdGlvbg=="},{"key":"dmFsaWRhdG9y","value":"YWRkcmVzczFkZWNuNjhldWVjNWRzZ3hyanB2N3Q1eWR5OHR5ZDc1dzhncnlhZg=="},{"key":"bmV3LXZhbGlkYXRvcg==","value":"YWRkcmVzczFoeHl1dDJkeXZydnh1bGZ1OGZsYXl0MDl3eWhxN3IwNG05OGx2Ng=="},{"key":"ZGVsZWdhdG9y","value":"YWRkcmVzczFsMHduNjZnaDQ1bmZ0YTJyNHZxOHo1NHd1OWhnYXJzczI5OGU5Zw=="}]},"hash":"BD74EB343A1118AC2D3D8F90176BD2B3E707B6CED1FC81AD1B885B4888EE9C09","height":"577861"}

```
