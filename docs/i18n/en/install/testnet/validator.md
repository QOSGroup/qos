# Running a Validator Node

Validators are responsible for committing new blocks to the blockchain according to consensus protocal.
Before setting up your validator node, make sure you've already gone through the Full Node Setup guide and your node has finish the block synchronizing.

```bash
$ qoscli tendermint status
```
Make sure the value of `catching_up` is `false`, otherwise please wait.

## Get QOS Token

Becoming a validator requires the operator to hold a certain amount of QOS testnet token, you can get some by using the [faucet](http://explorer.qoschain.info/freecoin/get).

1. Create account

Use `qoscli keys add <name_of_key>` to create a new account, for example:
```bash
# create account named 'Peter'
$ qoscli keys add Peter
Enter a passphrase for your key: 
Repeat the passphrase: 
```
The output would look like:
```bash
NAME:   TYPE:   ADDRESS:                                                PUBKEY:
Peter local   address1epvxmtxx99gy5xv7k7sl55994pehxgqt03va2s  D+pHqEJVjQMiRzl5PbL8FraVZqWqxrxcTF7akcCIDfo=
**Important** write this seed phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

broom resource trash summer crop embrace stadium fish brief dolphin run decrease brief heart upgrade icon toe lift dawn regret dumb indoor drop glide
```
Remember `address1epvxmtxx99gy5xv7k7sl55994pehxgqt03va2s` is the `Address`，
`D+pHqEJVjQMiRzl5PbL8FraVZqWqxrxcTF7akcCIDfo=` is the `PubKey`，
`broom resource trash summer crop embrace stadium fish brief dolphin run decrease brief heart upgrade icon toe lift dawn regret dumb indoor drop glide` is the `seed phrase`,
you can use these 24 phrases to recover this account.

Run `qoscli keys --help` for more keybase tools information.

2. Claim tokens

You can get some QOS from [Faucet](http://explorer.qoschain.info/freecoin/get).

::: warning Note 
QOS got from faucet can only be used in the testnet.
:::

Then run the query command to confirm your account:
```bash
$ qoscli query account Peter --indent
{
  "type": "qbase/account/QOSAccount",
  "value": {
    "base_account": {
      "account_address": "address1epvxmtxx99gy5xv7k7sl55994pehxgqt03va2s",
      "public_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "D+pHqEJVjQMiRzl5PbL8FraVZqWqxrxcTF7akcCIDfo="
      },
      "nonce": "0"
    },
    "qos": "10000000",
    "qscs": null
  }
}
```
The `qos` indicates the amount of tokens in Peter's account. 

## Create Validator

1. Exec create command

```
qoscli tx create-validator --owner Peter --name "Peter's node" --tokens 20000000 --description "hi, my eth address: xxxxxx"
```
- `--owner` keybase name or address, eg, `Peter` or `address1epvxmtxx99gy5xv7k7sl55994pehxgqt03va2s`
- `--name`  validator name, anything you like
- `--tokens` amount of tokens deposited to the validator, less than the `qos` in the owner's account
- `--description` description info

The output looks like:
```bash
{"check_tx":{},"deliver_tx":{},"hash":"34A76D6D07D93FBE395DDC55E0596E4D312A02A9","height":"200"}
```

2. View validator info

Exec query command to show your validator info:

`qoscli query validator --owner <owner_address_of_validator>`
- `owner_address_of_validator` keybase name or address

```bash
$ qoscli query validator --owner address1epvxmtxx99gy5xv7k7sl55994pehxgqt03va2s
{
  "name": "Peter's node",
  "owner": "address1epvxmtxx99gy5xv7k7sl55994pehxgqt03va2s",
  "validatorPubkey": {
    "type": "tendermint/PubKeyEd25519",
    "value": "PJ58L4OuZp20opx2YhnMhkcTzdEWI+UayicuckdKaTo="
  },
  "bondTokens": "20000000",
  "description": "",
  "status": 0,
  "inactiveCode": 0,
  "inactiveTime": "0001-01-01T00:00:00Z",
  "inactiveHeight": "0",
  "bondHeight": "200"
}
```
The value `0` of `status` indicates that your node is a validator node.

### QOS Explorer

You can also see your validator on the [QOS explorer](http://explorer.qoschain.info/validator/list).