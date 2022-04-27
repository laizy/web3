package storage

import (
	"fmt"

	"github.com/laizy/web3/evm/storage/schema"
)

type FakeDB struct {
	schema.PersistStore
}

func (self *FakeDB) Get(key []byte) ([]byte, error) {
	return nil, fmt.Errorf("no presistence")
}

func NewFakeDB() *FakeDB {
	return &FakeDB{}
}
