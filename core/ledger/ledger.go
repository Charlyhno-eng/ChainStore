package ledger

import (
	"ChainStore/core/block"
	"crypto/ed25519"
	"errors"
	"fmt"
)

type Ledger struct {
    Blocks []block.Block
}

func NewLedger() *Ledger {
    return &Ledger{
        Blocks: []block.Block{},
    }
}

func (l *Ledger) AddBlock(data string, privKey ed25519.PrivateKey) error {
	var previousHash string
	if len(l.Blocks) > 0 {
		previousHash = l.Blocks[len(l.Blocks)-1].ComputeHash()
	}

	newBlock := block.CreateNewBlock(data, privKey, previousHash)

	if !block.IsValidBlock(newBlock) {
		return errors.New("invalid block")
	}

	l.Blocks = append(l.Blocks, newBlock)
	return nil
}

func (l *Ledger) IsValidChain() bool {
	if len(l.Blocks) == 0 {
		return true
	}

	for i := 1; i < len(l.Blocks); i++ {
		currentBlock := l.Blocks[i]
		previousBlock := l.Blocks[i-1]

		if !block.IsValidBlock(currentBlock) {
			fmt.Printf("Invalid block %s (incorrect signature)\n", currentBlock.ID)
			return false
		}

		if currentBlock.PreviousHash != previousBlock.ComputeHash() {
			fmt.Printf("Incorrect chaining between blocks %s et %s\n", previousBlock.ID, currentBlock.ID)
			return false
		}
	}

	return true
}