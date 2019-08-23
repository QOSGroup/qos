# Discription
>   Command group of validator.

>   命令组 - 验证人

# Commands

| Command                              | Alias                            | Has-Subcommand | Description               |
|:-------------------------------------|:---------------------------------|:---------------|:--------------------------|
| `qoscli query tendermint-validators` | `qoscli q tendermint-validators` | ✖              | 获取给定高度的tendermint验证人集合    |
| `qoscli query validator`             | `qoscli q validator`             | ✖              | 查询验证人的信息                  |
| `qoscli query validator-miss-vote`   | `qoscli q validator-miss-vote`   | ✖              | 查询最近投票窗口中的验证人漏块(错过投票)信息   |
| `qoscli query validator-period`      | `qoscli q validator-period`      | ✖              | 查询分发(distribution)验证人周期信息 |
| `qoscli query validators`            | `qoscli q validators`            | ✖              | 查询所有验证人的信息                |
| `qoscli tx active-validator`         | -                                | ✖              | 激活(Active)验证人             |
| `qoscli tx create-validator`         | -                                | ✖              | 创建用自委托初始化的新验证人            |
| `qoscli tx modify-validator`         | -                                | ✖              | 修改已存在的验证人账户               |
| `qoscli tx revoke-validator`         | -                                | ✖              | 撤销(Revoke)验证人             |
