# Setup A Full-node

Before setting up your validator node, make sure you already had QOS installed by following this [guide](installation.md)

## Init Your Node

These instructions are for setting up a brand new full node from scratch.

First, initialize the node and create the necessary config files:

```bash
qosd init --name <your_custom_moniker> --chain-id <the_lastest_qos_testnet_chain_id>
```

::: warning Note
Only ASCII characters are supported for the `--name`. Using Unicode characters will render your node unreachable.
:::

You can edit `moniker` later, in the `~/.qosd/config/config.toml` file:

```toml
# A custom human readable name for this node
moniker = "<your_custom_moniker>"
```

And you can also edit `chain-id` later, int the `~/.qosd/config/genesis.json`
```toml
# The lastest QOS testnet chain-id
"chain_id": "<the_lastest_qos_testnet_chain_id>"
```

Your full node has been initialized!

## Get Configuration Files

Fetch the testnet's `genesis.json`, `config.toml` file into `qosd`'s config directory.

```bash
curl https://raw.githubusercontent.com/QOSGroup/qos-testnets/master/lastest/genesis.json > $HOME/.gaiad/config/genesis.json
curl https://raw.githubusercontent.com/QOSGroup/qos-testnets/master/lastest/config.toml > $HOME/.gaiad/config/config.toml
```

Note we use the `latest` directory in the [testnets repo](https://github.com/QOSGroup/qos-testnets)
which contains details for the latest testnet. If you are connecting to a different testnet, ensure you get the right files.


## Add Seed Nodes

// TODO seed nodes

## Run a Full Node

Start the full node with this command:

```bash
qosd start --with-tendermint
```

Check that everything is running smoothly:

```bash
qoscli status
```

If you see the catching_up is false, it means your node is fully synced with the network, otherwise your node is still downloading blocks. Once fully synced, you could upgrade your node to a validator node. The instructions is in [here](validator.md).