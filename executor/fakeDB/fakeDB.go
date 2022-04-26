package fakeDB

import (
	"fmt"

	"github.com/laizy/web3/evm/storage/schema"
)

type FakeDB struct {
}

func (self *FakeDB) Get(key []byte) ([]byte, error) {
	return nil, fmt.Errorf("no presistence")
}

func (self *FakeDB) BatchPut(key []byte, value []byte) {
	panic("todo ")
} //Put a key-value pair to batch
func (self *FakeDB) BatchDelete(key []byte) {
	panic("todo")

} //Delete the key in batch
func (self *FakeDB) NewIterator(prefix []byte) schema.StoreIterator {
	panic("todo")
} //Return the iterator of store
