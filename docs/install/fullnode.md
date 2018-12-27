# Setup A Full-node

Before setting up your validator node, make sure you already had QOS installed by following this [guide](installation.md).

## Init Your Node

These instructions are for setting up a brand new full node from scratch.

First, initialize the node and create the necessary config files:

```bash
$ qosd init --name <your_custom_moniker>
```
The default directory for config and data is *$HOME/.qosd*, you can change it by adding `--home` flag. 

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

Replace your local `genesis.json`:
```bash
$ curl https://raw.githubusercontent.com/QOSGroup/qos-testnets/master/capricorn-1000/genesis.json > ~/.qosd/config/genesis.json
```

Edit the `config.toml`:
```toml
# Comma separated list of seed nodes to connect to
seeds = "5d9fcba29ce9a066cdd6e4c45001567a4bd1dbf4@47.100.231.9:26656"
```

Note we use the `capricorn-1000` directory, that may not be the latest. Search your version In the [testnets repo](https://github.com/QOSGroup/qos-testnets).

## Run a Full Node

Start the full node with start command:

```bash
$ qosd start --with-tendermint
```

You can see the node is running, your node is not a validator, and your node is synchronizing blocks from the QOS testnet.
```bash
Starting ABCI with Tendermint                module=main 
Starting multiAppConn                        module=proxy impl=multiAppConn
Starting localClient                         module=abci-client connection=query impl=localClient
Starting localClient                         module=abci-client connection=mempool impl=localClient
Starting localClient                         module=abci-client connection=consensus impl=localClient
ABCI Handshake                               module=consensus appHeight=0 appHash=
ABCI Replay Blocks                           module=consensus appHeight=0 storeHeight=0 stateHeight=0
update Validators                            module=main len=4
Completed ABCI Handshake - Tendermint and App are synced module=consensus appHeight=0 appHash=
This node is not a validator                 module=consensus addr=666A495A6B05C975B241880785665417B5CEA2A6 pubKey=PubKeyEd25519{36BA673E7CC36F09C353720441C439A96E81B54689BAC219F0D24C52C3D23E65}
Starting Node                                module=node impl=Node
Starting EventBus                            module=events impl=EventBus
Local listener                               module=p2p ip=0.0.0.0 port=26656
Starting DefaultListener                     module=p2p impl=Listener(@172.31.230.212:26656)
P2P Node ID                                  module=node ID=db49a8d5a902910e0f8aee19e1b4889d6a235a91 file=/root/.qosd/config/node_key.json
Add our address to book                      module=p2p book=/root/.qosd/config/addrbook.json addr=db49a8d5a902910e0f8aee19e1b4889d6a235a91@172.31.230.212:26656
Starting RPC HTTP server on tcp://0.0.0.0:26657 module=rpc-server 
Starting P2P Switch                          module=p2p impl="P2P Switch"
Starting EvidenceReactor                     module=evidence impl=EvidenceReactor
Starting PEXReactor                          module=p2p impl=PEXReactor
Starting AddrBook                            module=p2p book=/root/.qosd/config/addrbook.json impl=AddrBook
Starting MempoolReactor                      module=mempool impl=MempoolReactor
Starting BlockchainReactor                   module=blockchain impl=BlockchainReactor
Starting BlockPool                           module=blockchain impl=BlockPool
Starting ConsensusReactor                    module=consensus impl=ConsensusReactor
ConsensusReactor                             module=consensus fastSync=true
Saving AddrBook to file                      module=p2p book=/root/.qosd/config/addrbook.json size=1
Starting IndexerService                      module=txindex impl=IndexerService
Ensure peers                                 module=p2p numOutPeers=0 numInPeers=0 numDialing=0 numToDial=10
Will dial address                            module=p2p addr=5d9fcba29ce9a066cdd6e4c45001567a4bd1dbf4@47.100.231.9:26656
Dialing peer                                 module=p2p address=5d9fcba29ce9a066cdd6e4c45001567a4bd1dbf4@47.100.231.9:26656
Successful handshake with peer               module=p2p peer=47.100.231.9:26656 peerNodeInfo="NodeInfo{id: 5d9fcba29ce9a066cdd6e4c45001567a4bd1dbf4, moniker: qos0, network: capricorn-1000 [listen 172.19.222.64:26656], version: 0.23.1 ([amino_version=0.12.0 p2p_version=0.5.0 consensus_version=v1/0.2.2 rpc_version=0.7.0/3 tx_index=on rpc_addr=tcp://0.0.0.0:26657])}"
Starting Peer                                module=p2p peer=47.100.231.9:26656 impl="Peer{MConn{47.100.231.9:26656} 5d9fcba29ce9a066cdd6e4c45001567a4bd1dbf4 out}"
Starting MConnection                         module=p2p peer=47.100.231.9:26656 impl=MConn{47.100.231.9:26656}
Added peer                                   module=p2p peer="Peer{MConn{47.100.231.9:26656} 5d9fcba29ce9a066cdd6e4c45001567a4bd1dbf4 out}"
Dialing peer                                 module=p2p address=5d9fcba29ce9a066cdd6e4c45001567a4bd1dbf4@47.100.231.9:26656
update Validators                            module=main len=0
Executed block                               module=state height=1 validTxs=0 invalidTxs=0
Committed state                              module=state height=1 txs=0 appHash=E5C9EABCC5C3ACB7EA6D8ED4D17B997BFCDD6F4F
Recheck txs                                  module=mempool numtxs=0 height=1
Indexed block                                module=txindex height=1
mint reward                                  module=main predict=8085999 actual=8085999
validatorVoteInfo                            module=main height=2 address1nfsgxj0l4gtgje0ydmjg6harsfvmduxtq8fdwa="not vote"
update Validators                            module=main len=0
Executed block                               module=state height=2 validTxs=0 invalidTxs=0
Committed state                              module=state height=2 txs=0 appHash=F70CB6559B9DA8015A63547696DC011032B7161F
Recheck txs                                  module=mempool numtxs=0 height=2
Indexed block                                module=txindex height=2
mint reward                                  module=main predict=8085999 actual=8085999
validatorVoteInfo                            module=main height=3 address1nfsgxj0l4gtgje0ydmjg6harsfvmduxtq8fdwa="not vote"
...
```

Check that everything is running smoothly:

```bash
$ qoscli tendermint status
```

If you see the catching_up is false, it means your node is fully synced with the network, otherwise your node is still downloading blocks. 

If your node is fully synced, follow [be a validator](validator.md) guild to become a testnet validator.