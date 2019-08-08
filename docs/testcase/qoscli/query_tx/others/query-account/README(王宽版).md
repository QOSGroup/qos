# test case of qoscli query account

> `qoscli query account` 查询账户信息

---

qoscli query account 业务完整性的用例测试库, 包含以下部分:

* 查询账户信息情景1:
  
  在使用账户完成交易后，查看账户的状态信息。例如：查询地址，公钥，交易次数，账户qos及qscs余额等等

* 查询账户信息情景2：

  在节点本地的密钥库中删除了其中某一账号，但是qos网络中存在该账号，对该账号查询状态信息。
