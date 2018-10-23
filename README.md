# ring-leth
testing implementation of ring-verify precompile in geth.

the contracts in this are the same as noot/ring-mixer. javascript deployment and testing will live there, while golang deployment and testing will live here. they should work the same but sometimes it's fun to interact with solidity in go :)
 
### dependencies
go version 1.10.3

### instructions
please see instructions for github.com/ChainSafeSystems/leth for golang tools install

then, 

```
go get github.com/noot/ring-leth
cd $GOPATH/github.com/noot/ring-leth
leth compile
```

this clones the needed repos and compiles the contract.

as well, we will need our forked go-ethereum with the ring-verify precompile.

```
go get github.com/noot/go-ethereum
cd $GOPATH/github.com/noot/go-ethereum
make geth
```

geth may need other dependencies; please see the geth README.

we also need to create a custom genesis block for our ring testnet. copy the following into `genesis.json`
```
{
  "difficulty" : "0x200",
  "extraData"  : "",
  "gasLimit"   : "0x8000000",
  "alloc": {
     "0x87C03Cc1F1a0BfABa24c23e929fFc61b6BB0d580": { "balance": "20000000000000000000" },
     "0x8f9b540b19520f8259115a90e4b4ffaeac642a30": { "balance": "20000000000000000000" }
  },
  "config": {
        "chainId": 15,
        "homesteadBlock": 0,
        "eip155Block": 0,
        "eip158Block": 0
    }
}
```
now let's run geth:

```
cd ~ && mkdir ringchaindata
$GOPATH/github.com/ethereum/go-ethereum/build/bin/geth init genesis.json --datadir ./ringchaindata
$GOPATH/github.com/ethereum/go-ethereum/build/bin/geth --rpc --datadir=~/ringchaindata --rpcport "8545"  --nodiscover --rpcapi="eth,web3,personal,net,miner,txpool"
```

in another terminal:

``` 
geth attach --datadir ./ringchaindata 
miner.start()
```

back to ring-leth, let's deploy and test!

```
cd $GOPATH/github.com/noot/ring-leth
leth deploy && leth test
```
