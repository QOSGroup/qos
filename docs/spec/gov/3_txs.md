# 交易

治理模块包含交易类型如下：

## 提议

QOS提议类型：

- ProposalTypeText           // 普通文本提议
- ProposalTypeParameterChange // 参数修改提议
- ProposalTypeTaxUsage       // 社区费池提取提议
- ProposalTypeModifyInflation // 修改通胀提议
- ProposalTypeSoftwareUpgrade // 软件升级提议

[发送提议](../../command/qoscli.md#提交提议)

### 文本提议

普通文本提议，提议可以是QOS网络建设性意见，新增或完善功能等等。

#### 结构

```go
type TxProposal struct {
	Title          string             `json:"title"`           //  标题
	Description    string             `json:"description"`     //  描述
	ProposalType   types.ProposalType `json:"proposal_type"`   //  类型
	Proposer       btypes.AccAddress  `json:"proposer"`        //  提议账户地址
	InitialDeposit btypes.BigInt      `json:"initial_deposit"` //  初始质押
}
```

#### 验证

必须通过以下校验交易才会执行：
- 标题不能为空且不能超过最大长度`MaxTitleLen`（默认200）
- 描述信息不能为空且不能超过最大长度`MaxDescriptionLen`(默认1000)
- 提议类型校验为`ProposalTypeText`

#### 签名

`proposer`

#### 交易费

0

### 参数修改提议

凡通过[参数查询](../../command/qoscli.md#参数查询)能查的参数均可通过此提议进行投票修改。

#### 结构

```go
type TxParameterChange struct {
	TxProposal                           // 基础数据，其中`ProposalType`为`ProposalTypeParameterChange`
	Params []types.Param `json:"params"` // 参数变更
}

type Param struct {
	Module string `json:"module"`   // 模块名
	Key    string `json:"key"` // 参数名
	Value  string `json:"value"` // 参数数值
}
```

#### 验证

必须通过以下校验交易才会执行：
- 标题不能为空且不能超过最大长度`MaxTitleLen`（默认200）
- 描述信息不能为空且不能超过最大长度`MaxDescriptionLen`(默认1000)
- 提议类型校验为`ProposalTypeParameterChange`
- 不存在未完成参数修改提议
- 参数列表不能为空
- 参数类型和数值正确


#### 签名

`proposer`

#### 交易费

0

### 社区费池提取提议

[系统账户](../guardian)可提交此提议从社区费池中提取QOS到指定账户。

#### 结构

```go
type TxTaxUsage struct {
	TxProposal                                          // 基础提议信息
	DestAddress btypes.AccAddress `json:"dest_address"` // 接收账户
	Percent     qtypes.Dec        `json:"percent"`      // 提取比例
}
```

#### 验证

必须通过以下校验交易才会执行：
- 标题不能为空且不能超过最大长度`MaxTitleLen`（默认200）
- 描述信息不能为空且不能超过最大长度`MaxDescriptionLen`(默认1000)
- 提议类型校验为`ProposalTypeTaxUsage`
- 接收账户不能为空，且接收账户为系统账户
- 提取比例范围(0, 1]

#### 签名

`proposer`

#### 交易费

0

### 修改通胀提议

QOS网络支持修改还未开始的通胀规则

#### 结构

```go
type TxModifyInflation struct {
	TxProposal              // 基础提议信息
	TotalAmount      btypes.BigInt         `json:"total_amount"`      // 总发行量
	InflationPhrases mint.InflationPhrases `json:"inflation_phrases"` // 通胀阶段
}
```

#### 验证

必须通过以下校验交易才会执行：
- 标题不能为空且不能超过最大长度`MaxTitleLen`（默认200）
- 描述信息不能为空且不能超过最大长度`MaxDescriptionLen`(默认1000)
- 提议类型校验为`ProposalTypeModifyInflation`
- 当前通胀和已完成通胀规则不可修改
- 修改后通胀规则总量与已流通量对应正确

#### 签名

`proposer`

#### 交易费

0

### 软件升级提议

通过软件升级提议完善现有功能逻辑，增加新特性。

#### 结构

```go
type TxSoftwareUpgrade struct {
	TxProposal
	Version       string `json:"version"`         // QOS版本
	DataHeight    int64  `json:"data_height"`     // 数据高度
	GenesisFile   string `json:"genesis_file"`    // `genesis.json`文件地址
	GenesisMD5    string `json:"genesis_md5"`     // `genesis.json`文件MD5
	ForZeroHeight bool   `json:"for_zero_height"` // 是否清除数据从0高度运行新网络
}
```

#### 验证

必须通过以下校验交易才会执行：
- 标题不能为空且不能超过最大长度`MaxTitleLen`（默认200）
- 描述信息不能为空且不能超过最大长度`MaxDescriptionLen`(默认1000)
- 提议类型校验为`ProposalTypeSoftwareUpgrade`
- 版本信息不能为空
- 如果是清除数据从0高度运行新网络，数据高度大于0，`genesis.json`及其MD5值均不能为空

#### 签名

`proposer`

#### 交易费

0

## 质押

可对未到达投票阶段的提议进行质押。

### 结构

```go
type TxDeposit struct {
	ProposalID int64             `json:"proposal_id"` // 提议ID
	Depositor  btypes.AccAddress `json:"depositor"`   // 质押账户
	Amount     btypes.BigInt     `json:"amount"`      // 质押QOS
}
```

### 验证

必须通过以下校验交易才会执行：
- `proposal_id`为正，提议存在且处于质押期
- `amount`为正
- 质押账户有足够QOS可质押

### 签名

`depositor`

### 交易费

0

## 投票

可对投票阶段的提议进行投票。

### 结构

```go
type TxVote struct {
	ProposalID int64             `json:"proposal_id"` // 提议ID
	Voter      btypes.AccAddress `json:"voter"`       // 投票账户地址
	Option     types.VoteOption  `json:"option"`      // 投票：Yes/Abstain/No/Nowithveto
}
```

### 验证

必须通过以下校验交易才会执行：
- 投票账户存在
- 投票类型有效
- 提议存在，且处于投票期

### 签名

`voter`

### 交易费

0