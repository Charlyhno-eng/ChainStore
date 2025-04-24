package network

import (
	"ChainStore/core/block"
	"ChainStore/store/leveldb"
	"encoding/json"
	"log"
	"net"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type HandshakePayload struct {
	Address string `json:"address"`
	Height  int    `json:"height"`
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

func (n *Node) handleMessage(conn net.Conn, msg Message) {
	switch msg.Type {
	case "ping":
		n.handlePing(conn)
	case "handshake":
		n.handleHandshake(conn, msg)
	case "blocks_request":
		n.handleBlocksRequest(conn, msg)
	case "blocks_response":
		n.handleBlocksResponse(conn, msg)
	case "new_block":
		n.handleNewBlock(conn, msg)
	}
}

func (n *Node) handlePing(conn net.Conn) {
	log.Printf("Ping received from %s", conn.RemoteAddr())
}

func (n *Node) handleHandshake(_ net.Conn, msg Message) {
	var payload HandshakePayload
	raw, _ := json.Marshal(msg.Data)
	if err := json.Unmarshal(raw, &payload); err != nil {
		log.Printf("Invalid handshake payload: %v", err)
		return
	}

	log.Printf("Handshake received from %s - Height: %d", payload.Address, payload.Height)
}

func (n *Node) handleBlocksRequest(conn net.Conn, msg Message) {
	var req BlocksRequestPayload
	raw, _ := json.Marshal(msg.Data)
	if err := json.Unmarshal(raw, &req); err != nil {
		log.Printf("Invalid blocks_request payload: %v", err)
		return
	}

	blocks, err := n.Blockchain.GetBlocksFromHeight(req.FromHeight)
	if err != nil {
		log.Printf("Error retrieving blocks: %v", err)
		return
	}

	peer := n.Peers[conn.RemoteAddr().String()]
	if peer != nil {
		resp := Message{
			Type: "blocks_response",
			Data: BlocksResponsePayload{Blocks: blocks},
		}
		peer.Writer.Encode(resp)
	}
}

func (n *Node) handleBlocksResponse(_ net.Conn, msg Message) {
	var payload BlocksResponsePayload
	raw, _ := json.Marshal(msg.Data)
	if err := json.Unmarshal(raw, &payload); err != nil {
		log.Printf("Invalid blocks_response payload: %v", err)
		return
	}

	for _, pb := range payload.Blocks {
		b := *pb
		if block.IsValidBlock(b) {
			if err := n.Blockchain.AddBlock(b); err != nil {
				log.Printf("Error adding block: %v", err)
			}
		}
	}
}

func (n *Node) handleNewBlock(_ net.Conn, msg Message) {
	var b block.Block
	raw, _ := json.Marshal(msg.Data)
	if err := json.Unmarshal(raw, &b); err != nil {
		log.Printf("Invalid new_block payload: %v", err)
		return
	}

	log.Printf("New block received: %s", b.ID)

	if !block.IsValidBlock(b) {
		log.Printf("Invalid block received: %s", b.ID)
		return
	}

	if err := n.Blockchain.AddBlock(b); err != nil {
		log.Printf("Error adding block %s: %v", b.ID, err)
		return
	}

	log.Printf("Block %s successfully added to the blockchain", b.ID)

	n.Broadcast(Message{
		Type: "new_block",
		Data: b,
	})

	log.Printf("Block %s broadcasted to peers", b.ID)
}

func (n *Node) SendHandshake(peer *Peer) {
	height := n.Blockchain.GetHeight()
	payload := HandshakePayload{
		Address: n.Address,
		Height:  height,
	}

	msg := Message{
		Type: "handshake",
		Data: payload,
	}

	err := peer.Writer.Encode(msg)
	if err != nil {
		log.Printf("Failed to send handshake to %s: %v", peer.Addr, err)
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
	n.SendHandshake(peer)

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
			log.Printf("Error decoding message: %v", err)
			return
		}

		n.handleMessage(conn, msg)
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