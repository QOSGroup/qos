# Join QOS Testnet

QOS testnet is named after 12 constellations, we use Capricorn to start the first testnet named capricorn-1000 at *Dec 26, 2018*.
The latest information of QOS testnet please see [qos-testnets](https://github.com/QOSGroup/qos-testnets).

**Video**(v0.0.4)：
- [youtube](https://youtu.be/-eFdx0rIPb4)
- [youku](http://v.youku.com/v_show/id_XNDA4NTA1MDM1Ng==.html)

To join the QOS testnet, please follow the steps:

## Install QOS on Your Server

See [Installation guide](http://docs.qoschain.info/qos/en/install/installation.html) to install QOS locally.

## Setup A Full-node

### Upgrade

Stop the qosd thread, remove the outdated data.
```bash
rm -rf $HOME/.qosd
```

Then following [Init](#init), [Config](#configure) and [Start](#start) guides.

### Init

First, initialize the node and create necessary config files:

```bash
$ qosd init --moniker <your_custom_moniker>
```

`$HOME/.qosd/config` will be created to hold the config files

### Configure

Different QOS testnet has different config files, visit [testnets repo](https://github.com/QOSGroup/qos-testnets) to find yours.

For example, we are about to join the latest testnet:

#### Replace `genesis.json`

Fetch the testnet's [genesis.json](https://raw.githubusercontent.com/QOSGroup/qos-testnets/master/latest/genesis.json) file into qosd's config directory:
```bash
$ curl https://raw.githubusercontent.com/QOSGroup/qos-testnets/master/latest/genesis.json > $HOME/.qosd/config/genesis.json
```

#### Edit `config.toml`

Your node needs to know how to find the peers. You'll need to add healthy seed nodes to `$HOME/.qosd/config/config.toml`. Find `seeds` option in it then set the value:

```toml
# Comma separated list of seed nodes to connect to
seeds = "f1dbd6d0b931fe7f918a81e8248c21e2109caa97@47.103.79.28:26656"
```

### Start


Start the full node with this command:
```bash
$ qosd start
```

Check if everything is running smoothly:

```bash
$ qoscli tendermint status --indent
{
  "node_info": {
    "protocol_version": {
      "p2p": "5",
      "block": "8",
      "app": "0"
    },
    "id": "b558f76507ecc348d1c096e3e1e935c1c4364c8b",
    "listen_addr": "tcp://0.0.0.0:26656",
    "network": "capricorn-2000",
    "version": "0.27.3",
    "channels": "4020212223303800",
    "moniker": "qos",
    "other": {
      "tx_index": "on",
      "rpc_address": "tcp://0.0.0.0:26657"
    }
  },
  "sync_info": {
    "latest_block_hash": "016C36688A3AAEF9CD5506B05045FDD4F25CEFF67786BA2EA33B52C7DB62F934",
    "latest_app_hash": "E206E6701494315B7E11E13B47F280C65BF0F250FDC3808317BFE57401F57BE3",
    "latest_block_height": "1761",
    "latest_block_time": "2019-04-10T08:07:15.411666471Z",
    "catching_up": false
  },
  "validator_info": {
    "address": "91601AE127D5656210542EC56C2E718A17711933",
    "pub_key": {
      "type": "tendermint/PubKeyEd25519",
      "value": "4CswaJ+kece3c8rMgVBhUtsWvOV+IkAnQKy5A147PM8="
    },
    "voting_power": "1000"
  }
}

```

The value of `catching_up` indicates whether your node has synchronized the latest block. When the value comes `false`, following the next step to become a validator. 


## Running a Validator Node

Validators are responsible for committing new blocks to the blockchain according to consensus protocal.
Before setting up your validator node, make sure you've already gone through the Full Node Setup guide and your node has finish the block synchronizing.

```bash
$ qoscli tendermint status
```
Make sure the value of `catching_up` is `false`, otherwise please wait.

### Get QOS Token

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

### Create Validator

1. Exec create command

```
qoscli tx create-validator --owner Peter --moniker "Peter's node" --tokens 20000000
```
- `--owner` keybase name or address, eg, `Peter` or `address1epvxmtxx99gy5xv7k7sl55994pehxgqt03va2s`
- `--moniker`  validator name, anything you like
- `--tokens` amount of tokens deposited to the validator, less than the `qos` in the owner's account

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