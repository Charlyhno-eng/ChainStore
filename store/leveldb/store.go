package leveldb

import (
	"ChainStore/core/block"
	"encoding/json"
	"fmt"
	"sort"

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

func (s *BlockStore) GetAllBlocks() ([]block.Block, error) {
    var blocks []block.Block
    iter := s.db.NewIterator(nil, nil)
    defer iter.Release()

    for iter.Next() {
        key := string(iter.Key())
        if key == "last" {
            continue
        }

        var b block.Block
        if err := json.Unmarshal(iter.Value(), &b); err != nil {
            continue
        }

        blocks = append(blocks, b)
    }

    if err := iter.Error(); err != nil {
        return nil, err
    }

    return blocks, nil
}

func (bs *BlockStore) GetHeight() int {
	blocks, err := bs.GetAllBlocks()
	if err != nil {
		return 0
	}
	return len(blocks)
}

func (s *BlockStore) GetBlocksFromHeight(height int) ([]*block.Block, error) {
	all, err := s.GetAllBlocks()
	if err != nil {
		return nil, err
	}

	sort.Slice(all, func(i, j int) bool {
		return all[i].Timestamp.Before(all[j].Timestamp)
	})

	if height >= len(all) {
		return []*block.Block{}, nil
	}

	var result []*block.Block
	for i := height; i < len(all); i++ {
		b := all[i]
		result = append(result, &b)
	}

	return result, nil
}

func (s *BlockStore) IsValidChain() (bool, error) {
    blocks, err := s.GetAllBlocks()
    if err != nil {
        return false, err
    }

	if len(blocks) < 2 {
        return true, nil
    }

    sort.Slice(blocks, func(i, j int) bool {
        return blocks[i].Timestamp.Before(blocks[j].Timestamp)
    })

    for i := 1; i < len(blocks); i++ {
        curr := blocks[i]
        prev := blocks[i-1]

        if !block.IsValidBlock(curr) {
            return false, fmt.Errorf("invalid signature on block %s", curr.ID)
        }
		
        if curr.PreviousHash != prev.ComputeHash() {
            return false, fmt.Errorf("invalid linkage: %s -> %s", prev.ID, curr.ID)
        }
    }
    return true, nil
}

func (s *BlockStore) Close() {
	s.db.Close()
}