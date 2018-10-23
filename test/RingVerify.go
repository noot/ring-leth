package test

import (
	"fmt"
	"log"
	"path/filepath"
	"io/ioutil"
	"strings"
	"context"
	"math/big"

	"github.com/noot/ring-mixer/bindings"

	"github.com/ChainSafeSystems/leth/core"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	//"github.com/ethereum/go-ethereum/common"
	//"github.com/ethereum/go-ethereum/core/types"
)

// type Receipt struct {
//     // Consensus fields
//     PostState         []byte `json:"root"`
//     Status            uint64 `json:"status"`
//     CumulativeGasUsed uint64 `json:"cumulativeGasUsed" gencodec:"required"`
//     Bloom             Bloom  `json:"logsBloom"         gencodec:"required"`
//     Logs              []*Log `json:"logs"              gencodec:"required"`

//     // Implementation fields (don't reorder!)
//     TxHash          common.Hash    `json:"transactionHash" gencodec:"required"`
//     ContractAddress common.Address `json:"contractAddress"`
//     GasUsed         uint64         `json:"gasUsed" gencodec:"required"`
// }

// type Log struct {
//     // Consensus fields:
//     // address of the contract that generated the event
//     Address common.Address `json:"address" gencodec:"required"`
//     // list of topics provided by the contract.
//     Topics []common.Hash `json:"topics" gencodec:"required"`
//     // supplied by the contract, usually ABI-encoded
//     Data []byte `json:"data" gencodec:"required"`

//     // Derived fields. These fields are filled in by the node
//     // but not secured by consensus.
//     // block in which the transaction was included
//     BlockNumber uint64 `json:"blockNumber"`
//     // hash of the transaction
//     TxHash common.Hash `json:"transactionHash" gencodec:"required"`
//     // index of the transaction in the block
//     TxIndex uint `json:"transactionIndex" gencodec:"required"`
//     // hash of the block in which the transaction was included
//     BlockHash common.Hash `json:"blockHash"`
//     // index of the log in the block
//     Index uint `json:"logIndex" gencodec:"required"`

//     // The Removed field is true if this log was reverted due to a chain reorganisation.
//     // You must pay attention to this field if you receive logs through a filter query.
//     Removed bool `json:"removed"`
// }

func Test() {
	conn, err := core.NewConnection("default")
	if err != nil {
		log.Fatal(err)
	}

	address, err := core.ContractAddress("RingVerify", "default")
	if err != nil {
		log.Fatal(err)
	}

	path, _ := filepath.Abs("./keystore/UTC--2018-05-17T21-58-52.188632298Z--8f9b540b19520f8259115a90e4b4ffaeac642a30")
	key, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("could not find keystore json file", err)
	}

	ringVerify, err := bindings.NewRingVerify(address, conn)
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewTransactor(strings.NewReader(string(key)), "password")
	if err != nil {
		log.Fatalf("could not unlock accounst")
	}

	tx, err := ringVerify.Hash(auth, big.NewInt(77))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tx.Hash().Hex())

	done := make(chan bool)
	go func() {
		core.AwaitTx(conn, tx.Hash(), done)
	}()
	<-done

	receipt, err := conn.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("tx logs: %x\n", receipt.Logs[0].Topics)

	tx, err = ringVerify.Verify(auth, []byte{0xff})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tx.Hash().Hex())

	done = make(chan bool)
	go func() {
		core.AwaitTx(conn, tx.Hash(), done)
	}()
	<-done

	receipt, err = conn.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx logs: %x\n", receipt.Logs[0].Topics)
}