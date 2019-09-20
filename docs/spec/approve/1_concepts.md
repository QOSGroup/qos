# 概念

预授权是允许QOS网络上已存在账户授权其他账户使用自己所持有QOS或QSC代币的行为。
一个用户A，对另一个用户B进行预授权，该预授权只能创建（`create-approve`）一次，但可以增加（`increase-approve`）或减小（`decrease-approve`）授权数量。
被授权账户可以分批次也可以一次性提取所授权的QOS和QSCs，提取时扣除授权账户相应币种的指定数量。

一个用户A，对另一个用户B进行预授权，该预授权只能创建（`create-approve`）一次，但可以增加（`increase-approve`）或减小（`decrease-approve`）授权数量。

QOS鼓励用户使用预授权功能，所以QOS将预授权相关交易的交易`Gas`设置为0，仅收取数据库操作所产生的极少量`Gas`。