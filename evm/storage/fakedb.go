package storage

import (
	"github.com/laizy/web3"

	"github.com/laizy/web3/evm/storage/schema"
)

type FakeDB struct {
	schema.PersistStore
}

func (self *FakeDB) Get(key []byte) ([]byte, error) {
	return nil, schema.ErrNotFound
}

func NewFakeDB() *FakeDB {
	return &FakeDB{}
}

func (self *FakeDB) GetBlockHash(height uint64) web3.Hash {
	return web3.Hash{}
}
