/*
 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */

package storage

import (
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/evm/storage/overlaydb"
	"github.com/umbracle/ethgo/evm/storage/schema"
)

// CacheDB is smart contract execute cache, it contain transaction cache and block cache
// When smart contract execute finish, need to commit transaction cache to block cache
type CacheDB struct {
	memdb      *overlaydb.MemDB
	backend    *overlaydb.OverlayDB
	keyScratch []byte
}

const initCap = 1024
const initKvNum = 16

// NewCacheDB return a new contract cache
func NewCacheDB(store *overlaydb.OverlayDB) *CacheDB {
	return &CacheDB{
		backend: store,
		memdb:   overlaydb.NewMemDB(initCap, initKvNum),
	}
}

func (self *CacheDB) SetDbErr(err error) {
	self.backend.SetError(err)
}

func (self *CacheDB) Reset() {
	self.memdb.Reset()
}

func ensureBuffer(b []byte, n int) []byte {
	if cap(b) < n {
		return make([]byte, n)
	}
	return b[:n]
}

func makePrefixedKey(dst []byte, prefix byte, key []byte) []byte {
	dst = ensureBuffer(dst, len(key)+1)
	dst[0] = prefix
	copy(dst[1:], key)
	return dst
}

// Commit current transaction cache to block cache
func (self *CacheDB) Commit() {
	self.memdb.ForEach(func(key, val []byte) {
		if len(val) == 0 {
			self.backend.Delete(key)
		} else {
			self.backend.Put(key, val)
		}
	})
	self.memdb.Reset()
}

func (self *CacheDB) Put(key []byte, value []byte) {
	self.put(schema.ST_STORAGE, key, value)
}

func (self *CacheDB) GenAccountStateKey(addr ethgo.Address, key []byte) []byte {
	k := make([]byte, 0, 1+20+32)
	k = append(k, byte(schema.ST_STORAGE))
	k = append(k, addr[:]...)
	k = append(k, key[:]...)

	return k
}

func (self *CacheDB) put(prefix schema.DataEntryPrefix, key []byte, value []byte) {
	self.keyScratch = makePrefixedKey(self.keyScratch, byte(prefix), key)
	self.memdb.Put(self.keyScratch, value)
}

func (self *CacheDB) Get(key []byte) ([]byte, error) {
	return self.get(schema.ST_STORAGE, key)
}

func (self *CacheDB) get(prefix schema.DataEntryPrefix, key []byte) ([]byte, error) {
	self.keyScratch = makePrefixedKey(self.keyScratch, byte(prefix), key)
	value, unknown := self.memdb.Get(self.keyScratch)
	if unknown {
		v, err := self.backend.Get(self.keyScratch)
		if err != nil {
			return nil, err
		}
		value = v
	}

	return value, nil
}

func (self *CacheDB) Delete(key []byte) {
	self.delete(schema.ST_STORAGE, key)
}

// Delete item from cache
func (self *CacheDB) delete(prefix schema.DataEntryPrefix, key []byte) {
	self.keyScratch = makePrefixedKey(self.keyScratch, byte(prefix), key)
	self.memdb.Delete(self.keyScratch)
}

func (self *CacheDB) NewIterator(key []byte) schema.StoreIterator {
	pkey := make([]byte, 1+len(key))
	pkey[0] = byte(schema.ST_STORAGE)
	copy(pkey[1:], key)
	prefixRange := util.BytesPrefix(pkey)
	backIter := self.backend.NewIterator(pkey)
	memIter := self.memdb.NewIterator(prefixRange)

	return &Iter{overlaydb.NewJoinIter(memIter, backIter)}
}

type Iter struct {
	*overlaydb.JoinIter
}

func (self *Iter) Key() []byte {
	key := self.JoinIter.Key()
	if len(key) != 0 {
		key = key[1:] // remove the first prefix
	}
	return key
}

func (self *CacheDB) CleanContractStorageData(addr ethgo.Address) error {
	iter := self.NewIterator(addr[:])

	for has := iter.First(); has; has = iter.Next() {
		self.Delete(iter.Key())
	}
	iter.Release()

	return iter.Error()
}
