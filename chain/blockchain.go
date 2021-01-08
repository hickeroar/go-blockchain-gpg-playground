package chain

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hickeroar/go-blockchain-gpg-playground/sign"
	"os"
)

type BlockChain struct {
	Blocks []*Block
}

// This is essentially a starting point for a log entry table in a database.
type Block struct {
	ID        int64
	Data      string
	Hash      []byte
	Signature []byte
}

func (b *Block) DeriveHash(prevHash []byte, prevSignature []byte) {
	data := bytes.Join([][]byte{[]byte(b.Data), prevHash, prevSignature}, []byte{})
	hash := sha256.Sum256(data)
	b.Hash = hash[:]
}

func (b *Block) ValidateHash(prevBlock *Block) bool {
	testBlock := &Block{b.ID, b.Data, []byte{}, []byte{}}
	testBlock.DeriveHash(prevBlock.Hash, prevBlock.Signature)
	return bytes.Compare(testBlock.Hash, b.Hash) == 0
}

func (b *Block) DeriveSignature() {
	signature := sign.CreateSignature(b.Data, b.Hash)
	b.Signature = signature
}

func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]

	newBlock := &Block{prevBlock.ID +1, data, []byte{}, []byte{}}
	newBlock.DeriveHash(prevBlock.Hash, prevBlock.Signature)
	newBlock.DeriveSignature()

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
	genesisBlock := &Block{0, "Genesis", []byte{}, []byte{}}
	genesisBlock.DeriveHash([]byte{}, []byte{})
	genesisBlock.DeriveSignature()
	chain.Blocks = append(chain.Blocks, genesisBlock)
	chain.WriteChain()
}

func (chain *BlockChain) PreviousBlock(block *Block) *Block {
	if block.ID > 0 {
		return chain.Blocks[block.ID-1]
	} else {
		return &Block{-1, "", []byte{}, []byte{}}
	}
}

func (chain *BlockChain) ValidateChain() error {
	for i := 0; i < len(chain.Blocks); i++ {
		block := chain.Blocks[i]
		prevBlock := chain.PreviousBlock(block)

		if !block.ValidateHash(prevBlock) {
			return errors.New(fmt.Sprintf("Block %d's hash could not be validated.", block.ID))
		}

		if !sign.VerifySignature(block.Data, block.Hash, block.Signature) {
			return errors.New(fmt.Sprintf("Block %d's signature could not be validated.", block.ID))
		}
	}

	return nil
}

func InitBlockChain() *BlockChain {
	chain := &BlockChain{[]*Block{}}
	chain.Genesis()
	return chain
}