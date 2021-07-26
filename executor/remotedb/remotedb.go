package remotedb

import (
	"encoding/hex"

	"github.com/holiman/uint256"
	common2 "github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/store/common"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/crypto"
	"github.com/umbracle/ethgo/evm/storage"
	"github.com/umbracle/ethgo/evm/storage/schema"
	"github.com/umbracle/ethgo/jsonrpc"
	"github.com/umbracle/ethgo/utils"
)

type RemoteDB struct {
	client   *jsonrpc.Client
	Accounts map[ethgo.Address]*storage.EthAccount
	Storage  map[storageKey]ethgo.Hash
}

func NewRemoteDB(url string) *RemoteDB {
	client, err := jsonrpc.NewClient(url)
	utils.Ensure(err)
	return &RemoteDB{
		client:   client,
		Accounts: make(map[ethgo.Address]*storage.EthAccount),
		Storage:  make(map[storageKey]ethgo.Hash),
	}
}

type storageKey struct {
	Addr ethgo.Address
	Key  ethgo.Hash
}

func (self *RemoteDB) GetAccount(addr ethgo.Address) *storage.EthAccount {
	if acc := self.Accounts[addr]; acc != nil {
		return acc
	}

	nonce, err := self.client.Eth().GetNonce(addr, ethgo.Latest)
	utils.Ensure(err)
	balance, err := self.client.Eth().GetBalance(addr, ethgo.Latest)
	utils.Ensure(err)
	code, err := self.client.Eth().GetCode(addr, ethgo.Latest)
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

	self.Accounts[addr] = acct
	return acct
}

func (self *RemoteDB) GetStorage(addr ethgo.Address, key ethgo.Hash) ethgo.Hash {
	skey := storageKey{Addr: addr, Key: key}
	if val, ok := self.Storage[skey]; ok {
		return val
	}

	val, err := self.client.Eth().GetStorage(addr, key, ethgo.Latest)
	utils.Ensure(err)
	self.Storage[skey] = val

	return val
}

func (self *RemoteDB) Get(key []byte) ([]byte, error) {
	switch schema.DataEntryPrefix(key[0]) {
	case schema.ST_ETH_ACCOUNT:
		addr := ethgo.BytesToAddress(key[1:])
		acct := self.GetAccount(addr)
		return common2.SerializeToBytes(acct), nil
	case schema.ST_STORAGE:
		addr := ethgo.BytesToAddress(key[1:21])
		key := ethgo.BytesToHash(key[21:])

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

func (self *RemoteDB) NewIterator(prefix []byte) common.StoreIterator {
	panic("todo")
}

func (self *RemoteDB) GetBlockHash(height uint64) ethgo.Hash {
	block, err := self.client.Eth().GetBlockByNumber(ethgo.BlockNumber(height), false)
	utils.Ensure(err)
	return block.Hash
}
