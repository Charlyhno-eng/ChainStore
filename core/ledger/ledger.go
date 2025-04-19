package ledger

import (
	"ChainStore/core/block"
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

func (l *Ledger) AddBlock(b block.Block) error {
    if !block.VerifyBlock(b) {
        return errors.New("bloc invalide : signature incorrecte")
    }

    l.Blocks = append(l.Blocks, b)
    return nil
}
