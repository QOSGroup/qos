# Networks

## Single-node
```
qosd init
qosd start --with-tendermint
```
默认init配置文件在$HOME/.qosd下，--home可指定存放目录，其他参考--help

如果一切正常，会看到控制台输出打块信息

```
...
I[11-06|08:07:57.717] Block{
  Header{
    ChainID:        test-chain-kwLu0n
    Height:         1
    Time:           2018-11-06 08:07:57.631158647 +0000 UTC
    NumTxs:         0
    TotalTxs:       0
    LastBlockID:    :0:000000000000
    LastCommit:
    Data:
    Validators:     EE657D91F4BDB46BC6E1252357EA03CE68343F8C
    App:
    Consensus:       D6B74BB35BDFFD8392340F2A379173548AE188FE
    Results:
    Evidence:
  }#9DB651ACE94EBEAAE5F0AF08F7CB4933DB4039E2
  Data{

  }#
  EvidenceData{

  }#
  Commit{
    BlockID:    :0:000000000000
    Precommits:
  }#
}#9DB651ACE94EBEAAE5F0AF08F7CB4933DB4039E2 module=consensus
...
```

默认创世[账户](../spec/account.md):
* address: address1k0m8ucnqug974maa6g36zw7g2wvfd4sug6uxay
* prikey: 0xa328891040ae9b773bcd30005235f99a8d62df03a89e4f690f9fa03abb1bf22715fc9ca05613f2d8061492e9f8149510b5b67d340d199ff24f34c85dbbbd7e0df780e9a6cc

启动完成，可进行[交易](../client/txs.md)操作
