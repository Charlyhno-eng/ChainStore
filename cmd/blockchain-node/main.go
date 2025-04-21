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

	// Open the LevelDB database
	store, err := leveldb.NewBlockStore("data/blocks")
	if err != nil {
		panic(err)
	}
	defer store.Close()

	// Validate the integrity of the chain at startup
	valid, err := store.IsValidChain()
	if err != nil {
		panic(fmt.Sprintf("Chain validation failed: %v", err))
	}
	if !valid {
		panic("Blockchain integrity check failed")
	}
	fmt.Println("Blockchain integrity verified")

	// Get the previous block hash if it exists
	var previousHash string
	lastBlock, err := store.GetLastBlock()
	if err == nil && lastBlock != nil {
		previousHash = lastBlock.ComputeHash()
	}

	// Create a new block using the previous hash
	newBlock := block.CreateNewBlock("Test blockchain!", privateKey, previousHash)

	// Print the generated signature for debugging purposes
	fmt.Printf("Generated Signature: %s\n", newBlock.Signature)

	// Verify the cryptographic validity of the block before adding it to the store
	if !block.IsValidBlock(newBlock) {
		panic("Invalid block")
	}

	// Add the block to the database
	err = store.AddBlock(newBlock)
	if err != nil {
		panic(err)
	}

	fmt.Println("Block added successfully")

	// Retrieve and display the block from the database
	storedBlock, err := store.GetBlock(newBlock.ID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Block read from the database: %+v\n", storedBlock)
}



// go run cmd/blockchain-node/main.go