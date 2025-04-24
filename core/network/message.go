package network

import "ChainStore/core/block"

type BlocksRequestPayload struct {
    FromHeight int `json:"from_height"`
}

type BlocksResponsePayload struct {
    Blocks []*block.Block `json:"blocks"`
}