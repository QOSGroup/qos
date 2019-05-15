# Setup A Full-node

## Upgrade

Stop the qosd thread, remove the outdated data.
```bash
rm -rf $HOME/.qosd
```

Then following [Init](#Init), [Config](#Config) and [Start](#Start) guides.

## Init

First, initialize the node and create necessary config files:

```bash
$ qosd init --moniker <your_custom_moniker>
```

`$HOME/.qosd/config` will be created to hold the config files

## Configure

Different QOS testnet has different config files, visit [testnets repo](https://github.com/QOSGroup/qos-testnets) to find yours.

For example, we are about to join the latest testnet:

### Replace `genesis.json`

Fetch the testnet's [genesis.json](https://raw.githubusercontent.com/QOSGroup/qos-testnets/master/latest/genesis.json) file into qosd's config directory:
```bash
$ curl https://raw.githubusercontent.com/QOSGroup/qos-testnets/master/latest/genesis.json > $HOME/.qosd/config/genesis.json
```

### Edit `config.toml`

Your node needs to know how to find the peers. You'll need to add healthy seed nodes to `$HOME/.qosd/config/config.toml`. Find `seeds` option in it then set the value:

```toml
# Comma separated list of seed nodes to connect to
seeds = "f1dbd6d0b931fe7f918a81e8248c21e2109caa97@47.103.79.28:26656"
```

## Start


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

The value of `catching_up` indicates whether your node has synchronized the latest block. When the value comes `false`, following [be a validator](validator.md) guide to become a validator. 