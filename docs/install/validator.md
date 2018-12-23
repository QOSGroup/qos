# Running a Validator Node

Before setting up your validator node, make sure you've already installed QOS by this [guide](fullnode.md)

Validators are responsible for committing new blocks to the blockchain through consensus.

## Get Testnet Token 

### Create Owner Account

You need to get `qosd` and `qoscli` installed first. Then, if you don't create an account before, follow the instructions below to create a new account.
And you could also use your account you created before. 

```
$ qoscli keys add <NAME_OF_KEY>
```

Then, you should set a password of at least 8 characters.

The output will look like the following:
```
NAME:   TYPE:   ADDRESS:                                                PUBKEY:
Peter local   address1epvxmtxx99gy5xv7k7sl55994pehxgqt03va2s  D+pHqEJVjQMiRzl5PbL8FraVZqWqxrxcTF7akcCIDfo=
**Important** write this seed phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

broom resource trash summer crop embrace stadium fish brief dolphin run decrease brief heart upgrade icon toe lift dawn regret dumb indoor drop glide
```

You could see the address and public key of this account. //TODO 是否要区分测试网和正式网地址前缀？

The seed phrase of this account will also be displayed. You could use these 24 phrases to recover this account in another server. The recover command is:
```
$ qoscli keys add <NAME_OF_KEY> --recover
```

### Claim tokens

You can always get some `QOS` //TODO ?  by using the [Faucet TODO](). The faucet will send you 10QOS every request, Please don't abuse it.

You can use the following command to check the balance of you account.
```
$ qoscli query account <NAME_OF_KEY or ADDRESS_OF_ACCOUNT>
```

Once you have created your own address, you can use this　account to stake as a validator. 


## Create Validator

### Confirm Your Validator is Synced

Your validator is active if the following command returns anything:

```
$ qoscli status
```

You should also be able to see `catching_up` is `false`. 

You need to get the public key of your node before upgrade your node to a validator node, it can be used to create a new validator by staking tokens. 

You can find your validator's pubkey by running:

```
$ qosd tendermint show-validator
{"type":"tendermint/PubKeyEd25519","value":"4X3GGmx2/D9UrQ9nKeB86zr+3SfI+QF4GI8t0QKS7CE="}
```

The use the [validator commands](../client/validator.md) to create the validator: 
```
qoscli tx create-validator --owner owner --name name --pubkey PJ58L4OuZp20opx2YhnMhkcTzdEWI+UayicuckdKaTo= --tokens 10 --description "I am a validator."
```

To read more about stake mechanism in QOS, go to this [doc](../spec/staking.md)

### View Validator Info

View the validator's information with this command:

```
// TODO
```


### Confirm Your Validator is Running

Your validator is active if the following command returns anything:

```
$ qosd status
```

You should also be able to see your power is above 0.


### Use QOS Testnet Explorer

You should also be able to see your validator on the [Explorer TODO]() If your bonded token is in top 100. The `bech32` encoded `address` of you validator you can find in the `~/.qosd/config/priv_validator.json` file.

Also, you can see all the other information of the testnet on the explorer.