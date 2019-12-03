# QOS Daemon server

`qosd`: commands to create, initialize, and start QOS network.

* `init`                  [Initialization](#initialization)
* `add-genesis-accounts`  [Genesis accounts](#genesis-accounts)
* `add-guardian`          [Guardian](#guardian)
* `gentx`                 [Gentx](#gentx)
* `collect-gentxs`        [Collect gentxs](#collect-gentxs)
* `config-root-ca`        [Root CA](#root-ca)
* `start`                 [Start](#start)
* `export        `        [Export](#export)
* `testnet`               [Testnet](#testnet)
* `unsafe-reset-all`      [Reset](#reset)
* `tendermint`            [Tendermint](#tendermint)
* `version`               [Version](#version)

global parameters:

| parameter | default value | description |
| :--- | :---: | :--- |
|--home string        | "$HOME/.qosd" |directory for config and data (default "$HOME/.qosd")|
|--log_level string   | "main:info,state:info,*:error" |Log level (default "main:info,state:info,*:error")|
|--trace              |  |print out full stack trace on errors|


## Initialization

`qosd init --moniker <your_custom_moniker> --chain-id <chain_id> --overwrite <overwrite>`

main parameters:

- `--moniker`   name in the P2P network, corresponds to `moniker` configuration item in `config.toml`
- `--chain-id`  chain ID
- `--overwrite` whether to overwrite the existing initial file

initialize `genesis`、`priv-validator`、`p2p-node` files:

```bash
$ qosd init --moniker capricorn-1000
```

result:
```bash
{
  "chain_id": "test-chain-9nlhQS",
  "node_id": "c427167c8d2838b00a46e33c4b325a7f05bd2c16",
  "app_message": null
}
```

two directories `data` and `config` will be created under `$HOME/.qosd/`:
`data` directory for storing data saved after network startup.
`config` holds `config.toml`，`genesis.json`，`node_key.json`，`priv_validator.json`.

## Genesis accounts

`qosd add-genesis-accounts <account_coin_s>`

`<account_coin_s>` account currency list, eg:[address1],[coin1],[coin2];[address2],[coin1],[coin2]

add genesis accounts to `genesis.json` file:
```bash
$ qosd add-genesis-accounts qosacc1c7nh7qquvjm3p28xpsnfn420437ztvzy2hwdtk,10000QOS
```

## Guardian

`qosd add-guardian --address <address> --description <description>`

main parameters:

- `--address`     address of guardian
- `--description` description

add guardian to `genesis.json` file:
```bash
$ qosd add-guardian --address qosacc1c7nh7qquvjm3p28xpsnfn420437ztvzy2hwdtk --description "this is the description"
```

## Gentx

Generate a create validator([TxCreateValidator](../spec/staking.md#TxCreateValidator)) transaction.

`qosd gentx --moniker <validator_name> --owner <account_address> --tokens <tokens>`

visit [create validator](qoscli.md#成为验证节点) for more information.

```bash
$ qosd gentx --moniker "Arya's node" --owner qosacc1c7nh7qquvjm3p28xpsnfn420437ztvzy2hwdtk --tokens 1000
```

the signed transaction data file with the file name `node ID@IP` will be generated in the `$HOME/.qosd/config/gentx` directory.

## Collect gentxs

`qosd collect-gentxs`

Collect the transaction data in the `gentx` directory and fill it into `genesis.json`.

## Root CA

`qosd config-root-ca --qcp <qcp_root.pub> --qsc <qsc_root.pub>`

`<qcp_root.pub>`、`<qsc_root.pub>` is the file path of root CA public key

setting the root CA information for [QSC](qoscli.md#qsc) and [QCP](qoscli.md#qcp).

## Start

`qosd start`

| parameter | default value | description |
| :--- | :---: | :--- |
|--abci string                     | "socket" |Specify abci transport (socket | grpc) (default "socket")|
|--address string                  | "tcp://0.0.0.0:26658") |Listen address (default "tcp://0.0.0.0:26658")|
|--consensus.create_empty_blocks   | true |Set this to false to only produce blocks when there are txs or when the AppHash changes (default true)|
|--fast_sync                       | true |Fast blockchain syncing (default true)|
|--moniker string                  | <your_computer_name> |Node Name|
|--p2p.laddr string                | "tcp://0.0.0.0:26656" |Node listen address. (0.0.0.0:0 means any interface, any port) (default "tcp://0.0.0.0:26656")|
|--p2p.persistent_peers string     | "" |Comma-delimited ID@host:port persistent peers|
|--p2p.pex                         | true |Enable/disable Peer-Exchange (default true)|
|--p2p.private_peer_ids string     | "" |Comma-delimited private peer IDs|
|--p2p.seed_mode                   | false |Enable/disable seed mode|
|--p2p.seeds string                | "" |Comma-delimited ID@host:port seed nodes|
|--p2p.upnp                        | false |Enable/disable UPNP port forwarding|
|--priv_validator_laddr string     | "" |Socket address to listen on for connections from external priv_validator process|
|--proxy_app string                | "tcp://127.0.0.1:26658" |Proxy app address, or 'nilapp' or 'kvstore' for local testing. (default "tcp://127.0.0.1:26658")|
|--pruning string                  | "syncable" |Pruning strategy: syncable, nothing, everything (default "syncable")|
|--rpc.grpc_laddr string           | "" |GRPC listen address (BroadcastTx only). Port required|
|--rpc.laddr string                | "tcp://0.0.0.0:26657" |RPC listen address. Port required (default "tcp://0.0.0.0:26657")|
|--rpc.unsafe                      | false |Enabled unsafe rpc methods|
|--trace-store string              | false |Enable KVStore tracing to an output file|
|--with-tendermint                 | true |Run abci app embedded in-process with tendermint|

start the QOS network:
```bash
$ qosd start
```

## Export

`qosd export --height <block_height> --for-zero-height <export_state_to_start_at_height_zero> -o <directory for exported json file>`

main parameters:

- `--height`            export block height
- `--for-zero-height`   whether to export the state for restarting the network from zero height
- `--o`                 export file location

export state with block height 4:
```bash
qosd export --height 4
```

the file named `genesis-<height>-<timestamp>.json` will be generated under `$HOME/.qosd`.

## Testnet

`qosd testnet` batch generation of cluster configuration files:
```bash
testnet will create "v" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).

Note, strict routability for addresses is turned off in the config file.

Example:

	qosd testnet --chain-id=qostest --v=4 --o=./output --starting-ip-address=192.168.1.2 --genesis-accounts=qosacc1c7nh7qquvjm3p28xpsnfn420437ztvzy2hwdtk,1000000qos

Usage:
  qosd testnet [flags]

Flags:
      --chain-id string              Chain ID
      --compound                     whether the validator's income is calculated as compound interest, default: true (default true)
      --genesis-accounts string      Add genesis accounts to genesis.json, eg: qosacc1c7nh7qquvjm3p28xpsnfn420437ztvzy2hwdtk,1000000qos,1000000qstars. Multiple accounts separated by ';'
      --guardians string             addresses for guardian. Multiple addresses separated by ','
  -h, --help                         help for testnet
      --home-client string           directory for keybase (default "$HOME/.qoscli")
      --hostname-prefix string       Hostname prefix (node results in persistent peers list ID0@node0:26656, ID1@node1:26656, ...) (default "node")
      --node-dir-prefix string       Prefix the directory name for each node with (node results in node0, node1, ...) (default "node")
      --o string                     Directory to store initialization data for the testnet (default "./mytestnet")
      --qcp-root-ca string           Config pubKey of root CA for QSC
      --qsc-root-ca string           Config pubKey of root CA for QCP
      --starting-ip-address string   Starting IP address (192.168.0.1 results in persistent peers list ID0@192.168.0.1:26656, ID1@192.168.0.2:26656, ...)
      --v int                        Number of validators to initialize the testnet with (default 4)

Global Flags:
      --home string        directory for config and data (default "$HOME/.qosd")
      --log_level string   Log level (default "main:info,state:info,*:error")
      --trace              print out full stack trace on errors


```

main parameters:
- chain-id            chain ID
- genesis-accounts    genesis accounts
- hostname-prefix     host name prefix
- moniker             moniker
- qcp-root-ca         pubKey of root CA for QCP
- qsc-root-ca         pubKey of root CA for QSC
- compound            delegation `compound`, default `true`
- starting-ip-address start IP

## Reset

`qosd unsafe-reset-all`

Reset the blockchain database, delete the address book file, and reset the state to the initial state.

## Tendermint

tendermint subcommands:

- `qosd tendermint show-address`    Show this node's tendermint validator address
- `qosd tendermint show-node-id`    Show this node's ID
- `qosd tendermint show-validator`  Show this node's tendermint validator info

## Version
view [qoscli version](qoscli.md#version)
