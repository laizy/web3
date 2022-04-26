/*
 * Copyright (C) 2021 The ontology Authors
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
	"fmt"
	"io"
	"math/big"

	"github.com/laizy/web3"
	"github.com/laizy/web3/crypto"
	"github.com/laizy/web3/evm/storage/schema"
	"github.com/laizy/web3/utils/codec"
	"github.com/laizy/web3/utils/common/hexutil"
	"github.com/laizy/web3/utils/common/uint256"
)

type BalanceHandle interface {
	SubBalance(cache *CacheDB, addr web3.Address, val *big.Int) error
	AddBalance(cache *CacheDB, addr web3.Address, val *big.Int) error
	SetBalance(cache *CacheDB, addr web3.Address, val *big.Int) error
	GetBalance(cache *CacheDB, addr web3.Address) (*big.Int, error)
}

type balanceHandle struct{}

func (self *balanceHandle) SubBalance(cache *CacheDB, addr web3.Address, val *big.Int) error {
	acct, err := cache.GetEthAccount(addr)
	if err != nil {
		return err
	}

	value, _ := uint256.FromBig(val)
	acct.Balance = acct.Balance.Sub(acct.Balance, value)
	cache.PutEthAccount(addr, acct)
	return nil
}

func (self *balanceHandle) AddBalance(cache *CacheDB, addr web3.Address, val *big.Int) error {
	acct, err := cache.GetEthAccount(addr)
	if err != nil {
		return err
	}

	value, _ := uint256.FromBig(val)
	acct.Balance = acct.Balance.Add(acct.Balance, value)
	cache.PutEthAccount(addr, acct)
	return nil
}

func (self *balanceHandle) SetBalance(cache *CacheDB, addr web3.Address, val *big.Int) error {
	acct, err := cache.GetEthAccount(addr)
	if err != nil {
		return err
	}

	value, _ := uint256.FromBig(val)
	acct.Balance = value
	cache.PutEthAccount(addr, acct)
	return nil
}

func (self *balanceHandle) GetBalance(cache *CacheDB, addr web3.Address) (*big.Int, error) {
	acct, err := cache.GetEthAccount(addr)
	if err != nil {
		return nil, err
	}

	return acct.Balance.ToBig(), nil
}

type StateDB struct {
	cacheDB       *CacheDB
	Suicided      map[web3.Address]bool
	logs          []*web3.StorageLog
	thash, bhash  web3.Hash
	txIndex       int
	refund        uint64
	snapshots     []*snapshot
	BalanceHandle BalanceHandle
}

func NewStateDB(cacheDB *CacheDB, thash, bhash web3.Hash) *StateDB {
	return &StateDB{
		cacheDB:       cacheDB,
		Suicided:      make(map[web3.Address]bool),
		logs:          nil,
		thash:         thash,
		bhash:         bhash,
		refund:        0,
		snapshots:     nil,
		BalanceHandle: &balanceHandle{},
	}
}

func (self *StateDB) Prepare(thash, bhash web3.Hash) {
	self.thash = thash
	self.bhash = bhash
	//	s.accessList = newAccessList()
}

func (self *StateDB) DbErr() error {
	return self.cacheDB.backend.Error()
}

func (self *StateDB) BlockHash() web3.Hash {
	return self.bhash
}

func (self *StateDB) GetLogs() []*web3.StorageLog {
	return self.logs
}

func (self *StateDB) Commit() error {
	err := self.CommitToCacheDB()
	if err != nil {
		return err
	}
	self.cacheDB.Commit()
	return nil
}

func (self *StateDB) CommitToCacheDB() error {
	for addr := range self.Suicided {
		self.cacheDB.DelEthAccount(addr) //todo : check consistence with ethereum
		err := self.cacheDB.CleanContractStorageData(addr)
		if err != nil {
			return err
		}
	}

	self.Suicided = make(map[web3.Address]bool)
	self.snapshots = self.snapshots[:0]

	return nil
}

type snapshot struct {
	//changes  overlaydb.IMemoryDB
	recover  map[string][]byte
	deletes  map[string]bool
	suicided map[web3.Address]bool
	logsSize int
	refund   uint64
}

func (self *StateDB) AddRefund(gas uint64) {
	self.refund += gas
}

// SubRefund removes gas from the refund counter.
// This method will panic if the refund counter goes below zero
func (self *StateDB) SubRefund(gas uint64) {
	if gas > self.refund {
		panic(fmt.Sprintf("Refund counter below zero (gas: %d > refund: %d)", gas, self.refund))
	}

	self.refund -= gas
}

func genKey(contract web3.Address, key web3.Hash) []byte {
	var result []byte
	result = append(result, contract.Bytes()...)
	result = append(result, key.Bytes()...)
	return result
}

func (self *StateDB) GetState(contract web3.Address, key web3.Hash) web3.Hash {
	val, err := self.cacheDB.Get(genKey(contract, key))
	if err != nil {
		self.cacheDB.SetDbErr(err)
	}

	return web3.BytesToHash(val)
}

// GetRefund returns the current value of the refund counter.
func (self *StateDB) GetRefund() uint64 {
	return self.refund
}

func (self *StateDB) SetState(contract web3.Address, key, value web3.Hash) {
	self.cacheDB.Put(genKey(contract, key), value[:])
}

func (self *StateDB) GetCommittedState(addr web3.Address, key web3.Hash) web3.Hash {
	k := self.cacheDB.GenAccountStateKey(addr, key[:])
	val, err := self.cacheDB.backend.Get(k)
	if err != nil {
		self.cacheDB.SetDbErr(err)
	}

	return web3.BytesToHash(val)
}

type EthAccount struct {
	Nonce    uint64
	Balance  *uint256.Int
	Code     hexutil.Bytes
	CodeHash web3.Hash
}

func (self *EthAccount) IsEmpty() bool {
	return self.Nonce == 0 && self.CodeHash == web3.Hash{}
}

func (self *EthAccount) Serialization(sink *codec.ZeroCopySink) {
	sink.WriteUint64(self.Nonce)
	var balance [32]byte
	if self.Balance != nil {
		balance = self.Balance.Bytes32()
	}
	sink.WriteBytes(balance[:])
	sink.WriteVarBytes(self.Code)
	sink.WriteHash(self.CodeHash)
}

func (self *EthAccount) Deserialization(source *codec.ZeroCopySource) error {
	self.Nonce, _ = source.NextUint64()
	balance, _ := source.NextBytes(32)
	self.Balance = uint256.NewInt().SetBytes32(balance)
	self.Code, _ = source.ReadVarBytes()
	hash, eof := source.NextHash()
	if eof {
		return io.ErrUnexpectedEOF
	}
	self.CodeHash = web3.Hash(hash)

	return nil
}

func (self *CacheDB) GetEthAccount(addr web3.Address) (val EthAccount, err error) {
	value, err := self.get(schema.ST_ETH_ACCOUNT, addr[:])
	if err != nil {
		return val, err
	}

	if len(value) == 0 {
		return val, nil
	}

	err = val.Deserialization(codec.NewZeroCopySource(value))

	return val, err
}

func (self *CacheDB) PutEthAccount(addr web3.Address, val EthAccount) {
	var raw []byte
	if !val.IsEmpty() {
		raw = codec.SerializeToBytes(&val)
	}

	self.put(schema.ST_ETH_ACCOUNT, addr[:], raw)
}

func (self *CacheDB) DelEthAccount(addr web3.Address) {
	self.put(schema.ST_ETH_ACCOUNT, addr[:], nil)
}

func (self *StateDB) getEthAccount(addr web3.Address) (val EthAccount) {
	account, err := self.cacheDB.GetEthAccount(addr)
	if err != nil {
		self.cacheDB.SetDbErr(err)
		return val
	}

	return account
}

func (self *StateDB) GetNonce(addr web3.Address) uint64 {
	return self.getEthAccount(addr).Nonce
}

func (self *StateDB) SetNonce(addr web3.Address, nonce uint64) {
	account := self.getEthAccount(addr)
	account.Nonce = nonce
	self.cacheDB.PutEthAccount(addr, account)
}

func (self *StateDB) GetCodeHash(addr web3.Address) (hash web3.Hash) {
	return self.getEthAccount(addr).CodeHash
}

func (self *StateDB) GetCode(addr web3.Address) []byte {
	return self.getEthAccount(addr).Code
}

func (self *StateDB) SetCode(addr web3.Address, code []byte) {
	codeHash := crypto.Keccak256Hash(code)
	account := self.getEthAccount(addr)
	account.CodeHash = codeHash
	account.Code = code
	self.cacheDB.PutEthAccount(addr, account)
}

func (self *StateDB) GetCodeSize(addr web3.Address) int {
	// todo : add cache to speed up
	return len(self.GetCode(addr))
}

func (self *StateDB) Suicide(addr web3.Address) bool {
	acct := self.getEthAccount(addr)
	if acct.IsEmpty() {
		return false
	}
	self.Suicided[addr] = true
	err := self.BalanceHandle.SetBalance(self.cacheDB, addr, big.NewInt(0))
	if err != nil {
		self.cacheDB.SetDbErr(err)
	}
	return true
}

func (self *StateDB) HasSuicided(addr web3.Address) bool {
	return self.Suicided[addr]
}

func (self *StateDB) Exist(addr web3.Address) bool {
	if self.Suicided[addr] {
		return true
	}
	acct := self.getEthAccount(addr)
	balance, err := self.BalanceHandle.GetBalance(self.cacheDB, addr)
	if err != nil {
		self.cacheDB.SetDbErr(err)
		return false
	}
	if !acct.IsEmpty() || balance.Sign() > 0 {
		return true
	}

	return false
}

func (self *StateDB) Empty(addr web3.Address) bool {
	acct := self.getEthAccount(addr)

	balance, err := self.BalanceHandle.GetBalance(self.cacheDB, addr)
	if err != nil {
		self.cacheDB.SetDbErr(err)
		return false
	}

	return acct.IsEmpty() && balance.Sign() == 0
}

func (self *StateDB) AddLog(log *web3.StorageLog) {
	self.logs = append(self.logs, log)
}

func (self *StateDB) AddPreimage(web3.Hash, []byte) {
	// todo
}

func (self *StateDB) ForEachStorage(web3.Address, func(web3.Hash, web3.Hash) bool) error {
	panic("todo")
}

func (self *StateDB) CreateAccount(web3.Address) {
	return
}

func (self *StateDB) Snapshot() int {
	recover, deletes := self.cacheDB.memdb.Changes()
	suicided := make(map[web3.Address]bool)
	for k, v := range self.Suicided {
		suicided[k] = v
	}

	sn := &snapshot{
		recover:  recover,
		deletes:  deletes,
		suicided: suicided,
		logsSize: len(self.logs),
		refund:   self.refund,
	}

	self.snapshots = append(self.snapshots, sn)

	return len(self.snapshots) - 1
}

func (self *StateDB) RevertToSnapshot(idx int) {
	if idx >= len(self.snapshots) {
		panic("can not to revert snapshot")
	}
	self.cacheDB.memdb.Revert(self.cacheDB.memdb.Changes())
	for i := len(self.snapshots) - 1; i > idx; i-- {
		self.cacheDB.memdb.Revert(self.snapshots[i].recover, self.snapshots[i].deletes)
	}

	sn := self.snapshots[idx]
	self.snapshots = self.snapshots[:idx]
	self.Suicided = sn.suicided
	self.refund = sn.refund
	self.logs = self.logs[:sn.logsSize]
}

func (self *StateDB) SubBalance(addr web3.Address, val *big.Int) {
	err := self.BalanceHandle.SubBalance(self.cacheDB, addr, val)
	if err != nil {
		self.cacheDB.SetDbErr(err)
		return
	}
}

func (self *StateDB) AddBalance(addr web3.Address, val *big.Int) {
	err := self.BalanceHandle.AddBalance(self.cacheDB, addr, val)
	if err != nil {
		self.cacheDB.SetDbErr(err)
		return
	}
}

func (self *StateDB) GetBalance(addr web3.Address) *big.Int {
	balance, err := self.BalanceHandle.GetBalance(self.cacheDB, addr)
	if err != nil {
		self.cacheDB.SetDbErr(err)
		return big.NewInt(0)
	}

	return balance
}
