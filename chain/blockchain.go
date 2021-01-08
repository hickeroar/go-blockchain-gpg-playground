package chain

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/hickeroar/go-blockchain-gpg-playground/sign"
	"os"
)

type BlockChain struct {
	Blocks []*Block
}

// This is essentially a starting point for a log entry table in a database.
type Block struct {
	BlockIndex int64
	Data       string
	Timestamp  int64
	Signature  []byte
}

func (b *Block) DeriveSignature(prevSignature []byte) {
	// The signature produced is derived from the payload concatenated with the previous signature.
	// This allows the signature itself to function as the blockchain "hash." Two birds, one stone.
	signature := sign.CreateSignature(b.Data, b.Timestamp, prevSignature)
	b.Signature = signature
}

func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]

	newBlock := &Block{prevBlock.BlockIndex + 1, data, crypto.GetUnixTime(), []byte{}}
	newBlock.DeriveSignature(prevBlock.Signature)
	chain.Blocks = append(chain.Blocks, newBlock)

	chain.WriteChain()
}

func (chain *BlockChain) WriteChain() {
	jsonBytes, _ := json.Marshal(chain.Blocks)
	jsonFile, _ := os.Create("chain.json")
	_, _ = jsonFile.WriteString(string(jsonBytes))
	_ = jsonFile.Close()
}

func (chain *BlockChain) Genesis() {
	genesisBlock := &Block{0, "Genesis", crypto.GetUnixTime(), []byte{}}
	genesisBlock.DeriveSignature([]byte{})
	chain.Blocks = append(chain.Blocks, genesisBlock)

	chain.WriteChain()
}

func (chain *BlockChain) PreviousBlock(block *Block) *Block {
	if block.BlockIndex > 0 {
		return chain.Blocks[block.BlockIndex-1]
	} else {
		return &Block{-1, "", crypto.GetUnixTime(), []byte{}}
	}
}

func (chain *BlockChain) ValidateChain() error {
	for i := 0; i < len(chain.Blocks); i++ {
		block := chain.Blocks[i]
		prevBlock := chain.PreviousBlock(block)

		if !sign.VerifySignature(block.Data, block.Timestamp, prevBlock.Signature, block.Signature) {
			return errors.New(fmt.Sprintf("Block %d's signature could not be validated.", block.BlockIndex))
		}
	}

	return nil
}

func InitBlockChain() *BlockChain {
	chain := &BlockChain{[]*Block{}}
	chain.Genesis()
	return chain
}
