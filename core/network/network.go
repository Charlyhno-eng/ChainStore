package network

import (
	"ChainStore/store/leveldb"
	"encoding/json"
	"log"
	"net"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Peer struct {
	Conn   net.Conn
	Addr   string
	Writer *json.Encoder
}

type Node struct {
	Address     string
	Peers       map[string]*Peer
	Blockchain  *leveldb.BlockStore
}

func (n *Node) handleMessage(msg Message, peer *Peer) {
	switch msg.Type {
	case "ping":
		log.Printf("Ping received from %s\n", peer.Addr)
		peer.Writer.Encode(Message{Type: "pong", Data: nil})

	case "new_block":
		log.Printf("New block received from %s\n", peer.Addr)
		// À compléter dans les prochaines étapes

	default:
		log.Printf("Unknown message type: %s\n", msg.Type)
	}
}

func (n *Node) ConnectToPeer(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Printf("Failed to connect to %s: %v", address, err)
		return
	}

	log.Printf("Connected to peer %s\n", address)

	peer := &Peer{
		Conn:   conn,
		Addr:   address,
		Writer: json.NewEncoder(conn),
	}

	n.Peers[address] = peer

	go n.handleConnection(conn)
}

func (n *Node) handleConnection(conn net.Conn) {
	addr := conn.RemoteAddr().String()
	log.Printf("New connection from %s\n", addr)

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	peer := &Peer{
		Conn:   conn,
		Addr:   addr,
		Writer: encoder,
	}

	n.Peers[addr] = peer

	for {
		var msg Message
		if err := decoder.Decode(&msg); err != nil {
			log.Printf("Connection lost with %s: %v\n", addr, err)
			delete(n.Peers, addr)
			return
		}
		n.handleMessage(msg, peer)
	}
}


func (n *Node) Listen(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	log.Printf("Listening on port %s\n", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		go n.handleConnection(conn)
	}
}

func (n *Node) Broadcast(msg Message) {
	for _, peer := range n.Peers {
		err := peer.Writer.Encode(msg)
		if err != nil {
			log.Printf("Failed to send to peer %s: %v", peer.Addr, err)
		}
	}
}

func NewNode(address string, blockchain *leveldb.BlockStore) *Node {
	return &Node{
		Address:    address,
		Peers:      make(map[string]*Peer),
		Blockchain: blockchain,
	}
}
