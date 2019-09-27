# test case of qoscli tx *guardian

> `qoscli tx *guardian` 系统用户操作

---

qoscli tx *guardian 业务完整性的用例测试库, 包含以下部分:

* 情景1:
  
    特权账户查询，查询单个guardian，可以通过指定keys_name 或是 account_address。

* 情景2：

    特权账户查询，以列表形式查询所有guardians。

* 情景3：

    添加特权账户，特权账户的添加只能是由特权账户来完成，也就是在genesis.json中配置的特权账户。

* 情景4：

    删除特权账户，同样执行该操作的账户必须是特权账户，且必须在genesis.json中配置。
