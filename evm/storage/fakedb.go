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

func (self *FakeDB) NewIterator(prefix []byte) schema.StoreIterator {
	return &fakeIter{}
}

type fakeIter struct{}

func (self *fakeIter) Next() bool {
	return false
}

func (self *fakeIter) First() bool {
	return false
}

func (self *fakeIter) Key() []byte {
	return nil
}
func (self *fakeIter) Value() []byte {
	return nil
}

func (self *fakeIter) Release() {
	return
}

func (self *fakeIter) Error() error {
	return nil
}
