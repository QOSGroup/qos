# Setup A Full-node

Before setting up your validator node, make sure you already had QOS installed by following this [guide](installation.md).

## Init Your Node

These instructions are for setting up a brand new full node from scratch.

First, initialize the node and create the necessary config files:

```bash
qosd init --name <your_custom_moniker>
```

::: warning Note
Only ASCII characters are supported for the `--name`. Using Unicode characters will render your node unreachable.
:::

You can edit `moniker` in the `~/.qosd/config/config.toml` file:

```toml
# A custom human readable name for this node
moniker = "<your_custom_moniker>"
```

Your full node has been initialized!

## Get Configuration Files

Replace your local `genesis.json` to testnet's:
```bash
curl https://raw.githubusercontent.com/QOSGroup/qos-testnets/master/capricorn-1000/genesis.json > ~/.qosd/config/genesis.json
```

Edit the following config in your local `config.toml` according to [`config.toml`](https://raw.githubusercontent.com/QOSGroup/qos-testnets/master/capricorn-1000/config.toml)
```toml
# Comma separated list of nodes to keep persistent connections to
persistent_peers = "<according_to_testnet_config>"
```

Note we use the `capricorn-1000` directory, that may not be the latest. Search your version In the [testnets repo](https://github.com/QOSGroup/qos-testnets).

## Run a Full Node

Start the full node with this command:

```bash
qosd start --with-tendermint
```

Check that everything is running smoothly:

```bash
qoscli tendermint status
```

If you see the catching_up is false, it means your node is fully synced with the network, otherwise your node is still downloading blocks. Once fully synced, you could upgrade your node to a validator node. The instructions is in [here](validator.md).