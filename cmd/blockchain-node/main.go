package main

import (
	"ChainStore/core/block"
	"ChainStore/core/cryptography"
	"fmt"
)

func main() {
    publicKey, privateKey, err := cryptography.GenerateKeyPair()

    if err != nil {
        panic(err)
    }

	newBlock := block.CreateNewBlock("Hello blockchain!", privateKey)
	fmt.Printf("Bloc créé : %+v\n", newBlock)
    fmt.Printf("Clé publique : %+v\n", publicKey)
}

// go run cmd/blockchain-node/main.go