# 存储

操作QCP的`MapperName`为`qcp`，定义在了[qbase](https://github.com/QOSGroup/qbase/tree/master/qcp)

## 根证书

QOS网络启动前需要在`genesis.json`中配置好CA中心用于签发QCP证书的公钥信息。网络启动后会如下保存：

- rootca: `rootca -> amino(pubKey)` 

## QCP

see [qcp in qbase](https://github.com/QOSGroup/qbase/blob/master/docs/spec/qcp.md)
