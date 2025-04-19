package main

import (
	"ChainStore/core/block"
	"ChainStore/core/cryptography"
	"ChainStore/store/leveldb"
	"fmt"
)

func main() {
	// Generate a key pair 
	_, privateKey, err := cryptography.GenerateKeyPair()
	if err != nil {
		panic(err)
	}

	// Open the LevelDB db
	store, err := leveldb.NewBlockStore("data/blocks")
	if err != nil {
		panic(err)
	}
	defer store.Close()

	// Determine the hash of the previous block
	var previousHash string
	lastBlock, err := store.GetLastBlock()
	if err == nil && lastBlock != nil {
		previousHash = lastBlock.Hash()
	}

	// Create a new block with the previous hash
	newBlock := block.CreateNewBlock("Hello blockchain!", privateKey, previousHash)

	// Check the cryptographic validity
	isValid := block.VerifyBlock(newBlock)
	if !isValid {
		panic("Bloc invalide !")
	}

	// Add the block to the db
	err = store.AddBlock(newBlock)
	if err != nil {
		panic(err)
	}

	fmt.Println("Bloc ajouté avec succès !")

	// Read the block from the base to confirm the recording
	storedBlock, err := store.GetBlock(newBlock.ID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Block read from the database : %+v\n", storedBlock)
}


// go run cmd/blockchain-node/main.go