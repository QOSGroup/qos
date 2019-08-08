# test case of qoscli tx modify-compound

> `qoscli tx modify-compound` 修改收益复投方式

---

## 情景说明

委托人正常修改在某一代理验证节点的收益复投方式。ps:当前周期对绑定QOS的增减，对配置参数的修改，到下一周期开始时生效。周期由参数delegator_income_period_height决定。委托人在一个周期内多次修改同一配置项（例如是否复投），以该周期内最后一次修改为准，应用到下一周期。在每个周期交替时为委托人分配收益/处理请求。

## 测试命令

```bash
    //收益复投方式设置true
    qoscli tx modify-compound --owner address1nnvdqefva89xwppzs46vuskckr7klvzk8r5uaa --delegator jlgy02 --compound

    //收益复投方式设置为false
    qoscli tx modify-compound --owner address1nnvdqefva89xwppzs46vuskckr7klvzk8r5uaa --delegator jlgy02
```

## 测试结果

```bash
    qoscli tx modify-compound --owner address1nnvdqefva89xwppzs46vuskckr7klvzk8r5uaa --delegator jlgy02 --compound
    Password to sign with 'jlgy02':
    {"check_tx":{"gasWanted":"100000","gasUsed":"9165"},"deliver_tx":{"gasWanted":"100000","gasUsed":"18280","tags":[{"key":"YWN0aW9u","value":"bW9kaWZ5LWNvbXBvdW5k"},{"key":"dmFsaWRhdG9y","value":"YWRkcmVzczFkZWNuNjhldWVjNWRzZ3hyanB2N3Q1eWR5OHR5ZDc1dzhncnlhZg=="},{"key":"ZGVsZWdhdG9y","value":"YWRkcmVzczFsMHduNjZnaDQ1bmZ0YTJyNHZxOHo1NHd1OWhnYXJzczI5OGU5Zw=="}]},"hash":"3AB2AD7310182E9B854005219168CFA6199DFCC1BD92206195FD32809C15956F","height":"587259"}

    qoscli tx modify-compound --owner address1nnvdqefva89xwppzs46vuskckr7klvzk8r5uaa --delegator jlgy02
    Password to sign with 'jlgy02':
    {"check_tx":{"gasWanted":"100000","gasUsed":"9171"},"deliver_tx":{"gasWanted":"100000","gasUsed":"18230","tags":[{"key":"YWN0aW9u","value":"bW9kaWZ5LWNvbXBvdW5k"},{"key":"dmFsaWRhdG9y","value":"YWRkcmVzczFkZWNuNjhldWVjNWRzZ3hyanB2N3Q1eWR5OHR5ZDc1dzhncnlhZg=="},{"key":"ZGVsZWdhdG9y","value":"YWRkcmVzczFsMHduNjZnaDQ1bmZ0YTJyNHZxOHo1NHd1OWhnYXJzczI5OGU5Zw=="}]},"hash":"4C8A745A65D282E55292B0BB68AB7638C4D456692F52C2ABCFBC6D4C551AF4CE","height":"587658"}


```
