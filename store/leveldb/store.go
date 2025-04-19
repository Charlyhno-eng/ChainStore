package leveldb

import (
	"ChainStore/core/block"
	"encoding/json"

	"github.com/syndtr/goleveldb/leveldb"
)

type BlockStore struct {
	db *leveldb.DB
}

func NewBlockStore(path string) (*BlockStore, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	return &BlockStore{db: db}, nil
}

func (s *BlockStore) AddBlock(b block.Block) error {
	data, err := json.Marshal(b)
	if err != nil {
		return err
	}

	err = s.db.Put([]byte(b.ID), data, nil)
	if err != nil {
		return err
	}

	err = s.db.Put([]byte("last"), []byte(b.ID), nil)
	if err != nil {
		return err
	}

	return s.db.Put([]byte(b.ID), data, nil)
}

func (s *BlockStore) GetBlock(id string) (*block.Block, error) {
	data, err := s.db.Get([]byte(id), nil)
	if err != nil {
		return nil, err
	}

	var b block.Block
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, err
	}

	return &b, nil
}

func (s *BlockStore) GetLastBlock() (*block.Block, error) {
	lastID, err := s.db.Get([]byte("last"), nil)
	if err != nil {
		return nil, err
	}

	return s.GetBlock(string(lastID))
}

func (s *BlockStore) Close() {
	s.db.Close()
}