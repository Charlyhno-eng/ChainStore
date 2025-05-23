# ChainStore

**ChainStore** is a simplified blockchain foundation written in **Go**.  
This project is intended for **educational purposes** and **does not include a consensus mechanism** (like Proof of Work or Proof of Stake) or any form of **cryptocurrency**. It serves as a **basic model** for understanding the architecture of a blockchain system and can be extended to support more complex features.

---

## Features

- Ed25519 public/private key generation
- Block creation and digital signature
- Persistent block storage using LevelDB
- Basic peer-to-peer networking (connections, handshake, message broadcasting)
- Blockchain integrity verification at startup
- Automatic block creation every 60 seconds

---

## Project Structure

```
ChainStore
├── cmd
│   └── blockchain-node
│       └── main.go              # Main entry point
├── core
│   ├── block
│   │   └── block.go             # Block structure and validation
│   ├── cryptography
│   │   └── keys.go              # Ed25519 key generation
│   ├── ledger
│   │   └── ledger.go            # Reserved for future ledger/transaction management
│   └── network
│       ├── message.go           # Network message structures
│       └── network.go           # Peer-to-peer networking logic
├── data                        # Persistent storage directory for blocks
├── store
│   └── leveldb
│       └── store.go             # LevelDB block store implementation
├── go.mod                      # Go module file
├── go.sum                      # Go module checksum file
└── README.md                   # Project documentation
```

---

## Running the Node

To start a node, use the following command:

```bash
go run cmd/blockchain-node/main.go
```

To test the blockchain, make two clones of the project. One on port 3000 listening to 3001, and one on port 3001 listening to 3000.

A new block is automatically generated every **60 seconds**.

---

## Limitations

This is a minimal project and does not include:
- Any consensus algorithm
- Fork resolution or reorganization
- Transactions or smart contract support
- User interface

---

## Use Cases

- Prototype for a custom blockchain implementation
- Experimentation with peer-to-peer network design
- Educational resource for understanding blockchain basics
- Base layer for adding consensus or tokenization mechanisms

---

## Dependencies

- [Go](https://golang.org/)
- [LevelDB](https://github.com/syndtr/goleveldb)

---

## License

This project is provided as-is for personal, educational, or experimental use.

---

## Contributions

Contributions are welcome.  
Fork the repository, test it locally with multiple nodes, and extend it with new features.
