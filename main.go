package main

import (
	"github.com/hickeroar/go-blockchain-gpg-playground/chain"
	"log"
)

func main() {
	blockchain := chain.InitBlockChain()
	blockchain.AddBlock("Foo")
	blockchain.AddBlock("Bar")
	blockchain.AddBlock("Baz")

	err := blockchain.ValidateChain()
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
}
