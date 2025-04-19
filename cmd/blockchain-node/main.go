package main

import (
	"ChainStore/core/block"
	"ChainStore/core/cryptography"
	"ChainStore/store/leveldb"
	"fmt"
)

func main() {
	_, privateKey, err := cryptography.GenerateKeyPair()
	if err != nil {
		panic(err)
	}

	newBlock := block.CreateNewBlock("Hello blockchain!", privateKey)
	isValid := block.VerifyBlock(newBlock)

	if !isValid {
		panic("Bloc invalide !")
	}

	store, err := leveldb.NewBlockStore("data/blocks")
	if err != nil {
		panic(err)
	}
	defer store.Close()

	err = store.AddBlock(newBlock)
	if err != nil {
		panic(err)
	}

	fmt.Println("Bloc ajouté avec succès !")

	storedBlock, err := store.GetBlock(newBlock.ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Bloc lu depuis la DB : %+v\n", storedBlock)
}


// go run cmd/blockchain-node/main.go