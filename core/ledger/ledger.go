package ledger

import (
	"ChainStore/core/block"
	"crypto/ed25519"
	"errors"
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
	var prevHash string
	if len(l.Blocks) > 0 {
		prevHash = l.Blocks[len(l.Blocks)-1].ComputeHash()
	}

	newBlock := block.CreateNewBlock(data, privKey, prevHash)

	if !block.IsValidBlock(newBlock) {
		return errors.New("bloc invalide")
	}

	l.Blocks = append(l.Blocks, newBlock)
	return nil
}