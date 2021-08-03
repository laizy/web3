package remotedb

import (
	"encoding/hex"
	"fmt"

	"github.com/holiman/uint256"
	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/crypto"
	"github.com/umbracle/go-web3/evm/storage"
	"github.com/umbracle/go-web3/evm/storage/schema"
	"github.com/umbracle/go-web3/jsonrpc"
	"github.com/umbracle/go-web3/utils"
	"github.com/umbracle/go-web3/utils/codec"
)

type RemoteDB struct {
	Trace    bool
	client   *jsonrpc.Client
	Accounts map[web3.Address]*storage.EthAccount
	Storage  map[storageKey]web3.Hash
}

func NewRemoteDB(client *jsonrpc.Client) *RemoteDB {
	return &RemoteDB{
		client:   client,
		Accounts: make(map[web3.Address]*storage.EthAccount),
		Storage:  make(map[storageKey]web3.Hash),
	}
}

type storageKey struct {
	Addr web3.Address
	Key  web3.Hash
}

func (self *RemoteDB) GetAccount(addr web3.Address) *storage.EthAccount {
	if acc := self.Accounts[addr]; acc != nil {
		return acc
	}

	nonce, err := self.client.Eth().GetNonce(addr, web3.Latest)
	utils.Ensure(err)
	balance, err := self.client.Eth().GetBalance(addr, web3.Latest)
	utils.Ensure(err)
	code, err := self.client.Eth().GetCode(addr)
	utils.Ensure(err)
	codeRaw, err := hex.DecodeString(code[2:])
	utils.Ensure(err)
	hash := crypto.Keccak256Hash(codeRaw)
	bal, _ := uint256.FromBig(balance)
	acct := &storage.EthAccount{
		Nonce:    nonce,
		Balance:  bal,
		Code:     codeRaw,
		CodeHash: hash,
	}

	if self.Trace {
		fmt.Printf("[remotedb] get account %s: %s\n", addr, utils.JsonString(acct))
	}

	self.Accounts[addr] = acct
	return acct
}

func (self *RemoteDB) GetStorage(addr web3.Address, key web3.Hash) web3.Hash {
	skey := storageKey{Addr: addr, Key: key}
	if val, ok := self.Storage[skey]; ok {
		return val
	}

	val, err := self.client.Eth().GetStorage(addr, key, web3.Latest)
	utils.Ensure(err)
	self.Storage[skey] = val

	if self.Trace {
		fmt.Printf("[remotedb] get storage, contract: %s, key: %s, value:%s\n", addr, key, val)
	}

	return val
}

func (self *RemoteDB) Get(key []byte) ([]byte, error) {
	switch schema.DataEntryPrefix(key[0]) {
	case schema.ST_ETH_ACCOUNT:
		addr := web3.BytesToAddress(key[1:])
		acct := self.GetAccount(addr)
		return codec.SerializeToBytes(acct), nil
	case schema.ST_STORAGE:
		addr := web3.BytesToAddress(key[1:21])
		key := web3.BytesToHash(key[21:])

		val := self.GetStorage(addr, key)
		return val.Bytes(), nil
	default:
		panic("unkown key prefix")
	}
}

func (self *RemoteDB) BatchPut(key []byte, value []byte) {
	panic("todo")
}

func (self *RemoteDB) BatchDelete(key []byte) {
	panic("todo")
}

func (self *RemoteDB) NewIterator(prefix []byte) schema.StoreIterator {
	panic("todo")
}

func (self *RemoteDB) GetBlockHash(height uint64) web3.Hash {
	block, err := self.client.Eth().GetBlockByNumber(web3.BlockNumber(height), false)
	utils.Ensure(err)
	return block.Hash
}
