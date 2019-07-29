# Test Cases

- [缺失必选参数](./TestCase01.md)
- [单个账户, 仅QOS](./TestCase02.md)
- [单个账户, QOS与QSC](./TestCase03.md)
- [多个账户, 整齐格式](./TestCase04.md)
- [多个账户, 混杂格式](./TestCase05.md)

# Description
>     add-genesis-accounts [accounts] will add [accounts] into app_state.
>     Multiple accounts separated by ';'.

>     add-genesis-accounts [accounts] 将添加 [accounts] 到 app_state。
>     多个帐户之间以';'进行分隔。

添加创世账户至`genesis.json`文件, 例如：
```bash
$ qosd add-genesis-accounts address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy,10000QOS
```
会在`genesis.json`文件`app-state`中`accounts`部分添加地址为`address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy`，持有10000QOS的账户信息。

# Example
```
  qosd add-genesis-accounts "address1lly0audg7yem8jt77x2jc6wtrh7v96hgve8fh8,1000000qos;address1auhqphrnk74jx2c5n80m9pdgl0ln79tyz32xlc,100000qos"
```

# Usage
```
  qosd add-genesis-accounts [accounts] [flags]
```

其中, `[accounts]`为账户币种币值列表，形如:`[address1],[coin1],[coin2];[address2],[coin1],[coin2]`

# Available Commands

>无可用命令

# Flags

| ShortCut | Flag            | Required | Input Type | Default Input | Input Range | Description       |
|:---------|:----------------|:---------|:-----------|:--------------|:------------|:------------------|
| `-h`     | `--help`        | ✖        | -          | -             | -           | 帮助文档              |
| -        | `--home-client` | ✖        | string     | `/.qoscli`    | -           | (主要参数)keybase所在目录 |

# Global Flags

| ShortCut | Flag          | Required | Input Type | Default Input                    | Input Range | Description  |
|:---------|:--------------|:---------|:-----------|:---------------------------------|:------------|:-------------|
| -        | `--home`      | ✖        | string     | `/.qosd`                         | -           | 配置和数据的目录     |
| -        | `--log_level` | ✖        | string     | `"main:info,state:info,*:error"` | -           | 日志级别         |
| -        | `--trace`     | ✖        | -          | -                                | -           | 打印出错时的完整堆栈跟踪 |
