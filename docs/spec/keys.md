# Keys
公私钥、地址存储、找回工具

## Add
添加账户地址信息</br>
Add Sansa</br>
passphrase：12345678
```
qoscli keys add Sansa
Enter a passphrase for your key:
Repeat the passphrase:
NAME:	TYPE:	ADDRESS:						PUBKEY:
Sansa	local	address1spdn868fzcpah8zd74tjck0e5akacgt2gmccnq	PubKeyEd25519{0DCE1804ED4932556D78EE646B6D4535B85106A09E29A70D62CEB66011E115CC}
**Important** write this seed phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

visual shop light arrive broken patch chest juice treat witness leopard mansion pitch strong crucial decade toward salad material youth slide seed crumble reopen
```

## List
```
qoscli keys list
NAME:	TYPE:	ADDRESS:						PUBKEY:
Arya	local	address1cnfqru6rts4nz224mvrf58ne427uthmcut4kc3	PubKeyEd25519{B495C3CC88D6D4D669DD7882C560C506ACA231BD94A72A14F2FA76E340EAB236}
Sansa	local	address1spdn868fzcpah8zd74tjck0e5akacgt2gmccnq	PubKeyEd25519{0DCE1804ED4932556D78EE646B6D4535B85106A09E29A70D62CEB66011E115CC}
```

## More
```
qoscli keys --help
Keys allows you to manage your local keystore for tendermint.

    These keys may be in any format supported by go-crypto and can be
    used by light-clients, full nodes, or any other application that
    needs to sign with a private key.

Usage:
  basecli keys [command]

Available Commands:
  mnemonic    Compute the bip39 mnemonic
  new         Interactive command to derive a new private key, encrypt it, and save to disk
  add         Create a new key, or import from seed
  list        List all keys
                   
  delete      Delete the given key
  update      Change the password used to protect private key

Flags:
      --chain-id string   Chain ID of tendermint node
      --height int        block height to query, omit to get most recent provable block
  -h, --help              help for keys
      --node string       <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")
      --trust-node        Trust connected full node (don't verify proofs for responses)

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "/home/imuge/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors

Use "basecli keys [command] --help" for more information about a command.
```