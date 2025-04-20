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

	// Get the previous block hash if it exists
	var previousHash string
	lastBlock, err := store.GetLastBlock()
	if err == nil && lastBlock != nil {
		previousHash = lastBlock.ComputeHash()
	}

	// Create a new block using the previous hash
	newBlock := block.CreateNewBlock("Test blockchain!", privateKey, previousHash)

	// Verify the cryptographic validity of the block
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