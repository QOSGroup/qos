# Discription
>   Command group of delegation.

>   命令组 - 委托

# Commands

| Command                         | Alias                       | Has-Subcommand | Description                    |
|:--------------------------------|:----------------------------|:---------------|:-------------------------------|
| `qoscli query delegation`       | `qoscli q delegation`       | ✖              | 查询委托信息                         |
| `qoscli query delegations`      | `qoscli q delegations`      | ✖              | 查询一个委托人的所有委托                   |
| `qoscli query delegations-to`   | `qoscli q delegations-to`   | ✖              | 查询一个验证人所接收的所有委托                |
| `qoscli query delegator-income` | `qoscli q delegator-income` | ✖              | 查询分发(distribution)委托人收入信息      |
| `qoscli query redelegations`    | `qoscli q redelegations`    | ✖              | 查询一个委托人的所有重委托请求(redelegations) |
| `qoscli query unbondings`       | `qoscli q unbondings`       | ✖              | 查询一个委托人的所有正在解除绑定(unbonding)的委托 |
| `qoscli tx delegate`            | -                           | ✖              | 向验证人委托QOS                      |
| `qoscli tx modify-compound`     | -                           | ✖              | 修改一个委托的复投信息                    |
| `qoscli tx redelegate`          | -                           | ✖              | 将QOS从一个验证人重新委托到另一个验证人          |
| `qoscli tx unbond`              | -                           | ✖              | 从验证人解除QOS委托                    |
