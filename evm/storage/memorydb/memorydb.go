package memorydb

import (
	"sort"
	"sync"

	"github.com/laizy/web3/evm/storage/overlaydb"
	"github.com/laizy/web3/evm/storage/schema"
	"github.com/laizy/web3/utils/common"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// Database is an ephemeral key-value store. Apart from basic data storage
// functionality it also supports batch writes and iterating over the keyspace in
// binary-alphabetical order.
type Database struct {
	db      map[string][]byte
	recover map[string][]byte
	deletes map[string]bool
	lock    sync.RWMutex
}

// New returns a wrapped map with all the required database interface methods
// implemented.
func New() *Database {
	return &Database{
		db:      make(map[string][]byte),
		deletes: make(map[string]bool),
		recover: make(map[string][]byte),
	}
}

func NewWithDB(db map[string][]byte) *Database {
	return &Database{
		//reuse
		db:      db,
		deletes: make(map[string]bool),
		recover: make(map[string][]byte),
	}
}

// Get retrieves the given key if it's present in the key-value store.
func (db *Database) Get(key []byte) (value []byte, unknow bool) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if db.db == nil {
		panic("no db")
	}
	if entry, ok := db.db[string(key)]; ok {
		return common.CopyBytes(entry), false
	}
	return nil, true
}

// Put inserts the given value into the key-value store.
func (db *Database) Put(key []byte, value []byte) {
	db.lock.Lock()
	defer db.lock.Unlock()

	if db.db == nil {
		panic("no db")
	}
	if _, exist := db.recover[string(key)]; !exist {
		if len(db.db[string(key)]) == 0 {
			db.deletes[string(key)] = true
		} else {
			db.recover[string(key)] = common.CopyBytes(db.db[string(key)])
		}
	}
	db.db[string(key)] = common.CopyBytes(value)
}

// Delete removes the key from the key-value store.
func (db *Database) Delete(key []byte) {
	db.lock.Lock()
	defer db.lock.Unlock()

	if db.db == nil {
		panic("no db")
	}
	if _, exist := db.recover[string(key)]; !exist {
		if len(db.db[string(key)]) == 0 {
			db.deletes[string(key)] = true
		} else {
			db.recover[string(key)] = common.CopyBytes(db.db[string(key)])
		}
	}
	delete(db.db, string(key))
}

func (db *Database) ForEach(f func(key, val []byte)) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	for k, v := range db.db {
		f([]byte(k), common.CopyBytes(v))
	}
}
func (db *Database) NewIterator(slice *util.Range) schema.StoreIterator {

	db.lock.RLock()
	defer db.lock.RUnlock()

	var (
		st     = string(slice.Start)
		keys   = make([]string, 0, len(db.db))
		values = make([][]byte, 0, len(db.db))
	)
	// Collect the keys from the memory database corresponding to the given prefix
	// and start
	for key := range db.db {
		if key >= st {
			keys = append(keys, key)
		}
	}
	// Sort the items and retrieve the associated values
	sort.Strings(keys)
	for _, key := range keys {
		values = append(values, db.db[key])
	}
	return &iterator{
		keys:   keys,
		values: values,
	}
}
func (db *Database) Reset() {
	db.lock.Lock()
	defer db.lock.Unlock()

	if db.db == nil {
		panic("no db")
	}
	db.db = make(map[string][]byte)
	db.recover = make(map[string][]byte)
	db.deletes = make(map[string]bool)
}

//Now just simply copy data
func (db *Database) DeepClone() overlaydb.IMemoryDB {
	db.lock.Lock()
	defer db.lock.Unlock()

	d := &Database{db: make(map[string][]byte)}
	for k, v := range db.db {
		d.db[k] = common.CopyBytes(v)
	}
	return d
}

func (db *Database) Changes() (recover map[string][]byte, delete map[string]bool) {
	db.lock.Lock()
	defer db.lock.Unlock()

	recover = db.recover
	delete = db.deletes
	db.recover = make(map[string][]byte)
	db.deletes = make(map[string]bool)
	return
}

func (db *Database) Revert(recover map[string][]byte, deletes map[string]bool) {
	db.lock.Lock()
	defer db.lock.Unlock()
	for k, v := range recover {
		db.db[k] = v
	}
	for k, _ := range deletes {
		delete(db.db, k)
	}
}
