package main

import (
	"ChainStore/core/block"
	"ChainStore/core/cryptography"
	"ChainStore/core/network"
	"ChainStore/store/leveldb"
	"crypto/ed25519"
	"fmt"
	"log"
	"time"
)

func initStore() (*leveldb.BlockStore, error) {
	store, err := leveldb.NewBlockStore("data/blocks")
	if err != nil {
		return nil, err
	}
	return store, nil
}

func generateKeys() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	return cryptography.GenerateKeyPair()
}

func validateChain(store *leveldb.BlockStore) {
	valid, err := store.IsValidChain()
	if err != nil {
		panic(fmt.Sprintf("Chain validation failed: %v", err))
	}
	if !valid {
		panic("Blockchain integrity check failed")
	}
	fmt.Println("Blockchain integrity verified")
}

func createAndStoreBlockWithData(store *leveldb.BlockStore, privateKey ed25519.PrivateKey, node *network.Node, data string) {
	var previousHash string

	lastBlock, err := store.GetLastBlock()
	if err != nil {
		lastBlock = nil
	}

	if lastBlock != nil {
		previousHash = lastBlock.ComputeHash()
	}

	newBlock := block.CreateNewBlock(data, privateKey, previousHash)
	fmt.Printf("Generated Signature: %s\n", newBlock.Signature)

	if !block.IsValidBlock(newBlock) {
		log.Printf("Invalid block: block failed validity check")
		return
	}

	if err := store.AddBlock(newBlock); err != nil {
		log.Printf("Error adding block to store: %v", err)
		return
	}

	fmt.Println("Block added successfully")

	storedBlock, err := store.GetBlock(newBlock.ID)
	if err != nil {
		log.Printf("Error retrieving block from store: %v", err)
		return
	}

	fmt.Printf("Block read from the database: %+v\n", storedBlock)

	node.Broadcast(network.Message{
		Type: "new_block",
		Data: newBlock,
	})
}

func setupNetworkNode(store *leveldb.BlockStore) *network.Node {
	node := network.NewNode("localhost:3000", store)
	go node.Listen("3000")
	node.ConnectToPeer("localhost:3001")
	node.Broadcast(network.Message{Type: "ping"})
	return node
}

func main() {
	_, privateKey, err := generateKeys()
	if err != nil {
		panic(fmt.Sprintf("Error generating keys: %v", err))
	}

	store, err := initStore()
	if err != nil {
		panic(fmt.Sprintf("Error initializing store: %v", err))
	}
	defer store.Close()

	validateChain(store)

	node := setupNetworkNode(store)

	if err := node.SyncChain("localhost:3001"); err != nil {
		log.Printf("Chain sync failed: %v", err)
	}

	createAndStoreBlockWithData(store, privateKey, node, "Test blockchain!")

	// create a block every minute
	go func() {
		for {
			time.Sleep(60 * time.Second)
			timestamp := time.Now().Format(time.RFC3339)
			data := fmt.Sprintf("Auto block %s", timestamp)
			createAndStoreBlockWithData(store, privateKey, node, data)
		}
	}()

	select {}
}
