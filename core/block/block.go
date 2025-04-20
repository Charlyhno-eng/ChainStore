package block

import (
	"ChainStore/core/cryptography"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

type Block struct {
	ID           string
	Timestamp    time.Time
	Data         string
	Signature    string
	PublicKey    []byte
	PreviousHash string
	Version      int
}

func IsValidBlock(b Block) bool {
	message := []byte(b.Data)
	return ed25519.Verify(b.PublicKey, message, []byte(b.Signature))
}

func (b *Block) ComputeHash() string {
	data := b.ID + b.Timestamp.String() + b.Data + b.Signature + string(b.PublicKey) + b.PreviousHash
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func CreateNewBlock(data string, privKey ed25519.PrivateKey, previousHash string) Block {
	message := []byte(data)
	signature := cryptography.SignMessage(privKey, message)

	return Block{
		ID:           uuid.NewString(),
		Timestamp:    time.Now(),
		Data:         data,
		Signature:    string(signature),
		PublicKey:    privKey.Public().(ed25519.PublicKey),
		PreviousHash: previousHash,
		Version:      1,
	}
}