# Stake

`MapperName` is `qcp` which has defined in [qbase](https://github.com/QOSGroup/qbase/tree/master/qcp).

## Root CA

Before the QOS network starts, we configure the public key information of the CA center for issuing QCP certificates in `genesis.json`. 
After the network starts, it will be saved as follows:

- rootca: `rootca -> amino(pubKey)` 

## QCP

see [qcp in qbase](https://github.com/QOSGroup/qbase/blob/master/docs/spec/qcp.md)
