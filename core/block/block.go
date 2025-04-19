package block

import (
	"ChainStore/core/cryptography"
	"crypto/ed25519"
	"time"

	"github.com/google/uuid"
)

type Block struct {
    ID        string
    Timestamp time.Time
    Data      string
    Signature string
    PublicKey []byte
}

func CreateNewBlock(data string, privKey ed25519.PrivateKey) Block {
	message := []byte(data)
    signature := cryptography.SignMessage(privKey, message)

    return Block {
		ID:        uuid.NewString(),
        Timestamp: time.Now(),
        Data:      data,
        Signature: string(signature),
        PublicKey: privKey.Public().(ed25519.PublicKey),
	}
}

func VerifyBlock(b Block) bool {
	message := []byte(b.Data)
	return ed25519.Verify(b.PublicKey, message, []byte(b.Signature))
}