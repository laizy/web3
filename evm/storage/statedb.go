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

	"github.com/holiman/uint256"
	comm "github.com/ontio/ontology/common"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/crypto"
	"github.com/umbracle/ethgo/evm/storage/overlaydb"
	"github.com/umbracle/ethgo/evm/storage/schema"
)

type BalanceHandle interface {
	SubBalance(cache *CacheDB, addr ethgo.Address, val *big.Int) error
	AddBalance(cache *CacheDB, addr ethgo.Address, val *big.Int) error
	SetBalance(cache *CacheDB, addr ethgo.Address, val *big.Int) error
	GetBalance(cache *CacheDB, addr ethgo.Address) (*big.Int, error)
}

type balanceHandle struct{}

func (self *balanceHandle) SubBalance(cache *CacheDB, addr ethgo.Address, val *big.Int) error {
	acct, err := cache.GetEthAccount(addr)
	if err != nil {
		return err
	}

	value, _ := uint256.FromBig(val)
	acct.Balance = acct.Balance.Sub(acct.Balance, value)
	cache.PutEthAccount(addr, acct)
	return nil
}

func (self *balanceHandle) AddBalance(cache *CacheDB, addr ethgo.Address, val *big.Int) error {
	acct, err := cache.GetEthAccount(addr)
	if err != nil {
		return err
	}

	value, _ := uint256.FromBig(val)
	acct.Balance = acct.Balance.Add(acct.Balance, value)
	cache.PutEthAccount(addr, acct)
	return nil
}

func (self *balanceHandle) SetBalance(cache *CacheDB, addr ethgo.Address, val *big.Int) error {
	acct, err := cache.GetEthAccount(addr)
	if err != nil {
		return err
	}

	value, _ := uint256.FromBig(val)
	acct.Balance = value
	cache.PutEthAccount(addr, acct)
	return nil
}

func (self *balanceHandle) GetBalance(cache *CacheDB, addr ethgo.Address) (*big.Int, error) {
	acct, err := cache.GetEthAccount(addr)
	if err != nil {
		return nil, err
	}

	return acct.Balance.ToBig(), nil
}

type StateDB struct {
	cacheDB       *CacheDB
	Suicided      map[ethgo.Address]bool
	logs          []*ethgo.StorageLog
	thash, bhash  ethgo.Hash
	txIndex       int
	refund        uint64
	snapshots     []*snapshot
	BalanceHandle BalanceHandle
}

func NewStateDB(cacheDB *CacheDB, thash, bhash ethgo.Hash) *StateDB {
	return &StateDB{
		cacheDB:       cacheDB,
		Suicided:      make(map[ethgo.Address]bool),
		logs:          nil,
		thash:         thash,
		bhash:         bhash,
		refund:        0,
		snapshots:     nil,
		BalanceHandle: &balanceHandle{},
	}
}

func (self *StateDB) Prepare(thash, bhash ethgo.Hash) {
	self.thash = thash
	self.bhash = bhash
	//	s.accessList = newAccessList()
}

func (self *StateDB) DbErr() error {
	return self.cacheDB.backend.Error()
}

func (self *StateDB) BlockHash() ethgo.Hash {
	return self.bhash
}

func (self *StateDB) GetLogs() []*ethgo.StorageLog {
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

	self.Suicided = make(map[ethgo.Address]bool)
	self.snapshots = self.snapshots[:0]

	return nil
}

type snapshot struct {
	changes  *overlaydb.MemDB
	suicided map[ethgo.Address]bool
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

func genKey(contract ethgo.Address, key ethgo.Hash) []byte {
	var result []byte
	result = append(result, contract.Bytes()...)
	result = append(result, key.Bytes()...)
	return result
}

func (self *StateDB) GetState(contract ethgo.Address, key ethgo.Hash) ethgo.Hash {
	val, err := self.cacheDB.Get(genKey(contract, key))
	if err != nil {
		self.cacheDB.SetDbErr(err)
	}

	return ethgo.BytesToHash(val)
}

// GetRefund returns the current value of the refund counter.
func (self *StateDB) GetRefund() uint64 {
	return self.refund
}

func (self *StateDB) SetState(contract ethgo.Address, key, value ethgo.Hash) {
	self.cacheDB.Put(genKey(contract, key), value[:])
}

func (self *StateDB) GetCommittedState(addr ethgo.Address, key ethgo.Hash) ethgo.Hash {
	k := self.cacheDB.GenAccountStateKey(addr, key[:])
	val, err := self.cacheDB.backend.Get(k)
	if err != nil {
		self.cacheDB.SetDbErr(err)
	}

	return ethgo.BytesToHash(val)
}

type EthAccount struct {
	Nonce    uint64
	Balance  *uint256.Int
	Code     []byte
	CodeHash ethgo.Hash
}

func (self *EthAccount) IsEmpty() bool {
	return self.Nonce == 0 && self.CodeHash == ethgo.Hash{}
}

func (self *EthAccount) Serialization(sink *comm.ZeroCopySink) {
	sink.WriteUint64(self.Nonce)
	balance := self.Balance.Bytes32()
	sink.WriteBytes(balance[:])
	sink.WriteVarBytes(self.Code)
	sink.WriteHash(comm.Uint256(self.CodeHash))
}

func (self *EthAccount) Deserialization(source *comm.ZeroCopySource) error {
	self.Nonce, _ = source.NextUint64()
	balance, _ := source.NextBytes(32)
	self.Balance = uint256.NewInt(0).SetBytes32(balance)
	self.Code, _ = source.ReadVarBytes()
	hash, eof := source.NextHash()
	if eof {
		return io.ErrUnexpectedEOF
	}
	self.CodeHash = ethgo.Hash(hash)

	return nil
}

func (self *CacheDB) GetEthAccount(addr ethgo.Address) (val EthAccount, err error) {
	value, err := self.get(schema.ST_ETH_ACCOUNT, addr[:])
	if err != nil {
		return val, err
	}

	if len(value) == 0 {
		return val, nil
	}

	err = val.Deserialization(comm.NewZeroCopySource(value))

	return val, err
}

func (self *CacheDB) PutEthAccount(addr ethgo.Address, val EthAccount) {
	var raw []byte
	if !val.IsEmpty() {
		raw = comm.SerializeToBytes(&val)
	}

	self.put(schema.ST_ETH_ACCOUNT, addr[:], raw)
}

func (self *CacheDB) DelEthAccount(addr ethgo.Address) {
	self.put(schema.ST_ETH_ACCOUNT, addr[:], nil)
}

func (self *StateDB) getEthAccount(addr ethgo.Address) (val EthAccount) {
	account, err := self.cacheDB.GetEthAccount(addr)
	if err != nil {
		self.cacheDB.SetDbErr(err)
		return val
	}

	return account
}

func (self *StateDB) GetNonce(addr ethgo.Address) uint64 {
	return self.getEthAccount(addr).Nonce
}

func (self *StateDB) SetNonce(addr ethgo.Address, nonce uint64) {
	account := self.getEthAccount(addr)
	account.Nonce = nonce
	self.cacheDB.PutEthAccount(addr, account)
}

func (self *StateDB) GetCodeHash(addr ethgo.Address) (hash ethgo.Hash) {
	return self.getEthAccount(addr).CodeHash
}

func (self *StateDB) GetCode(addr ethgo.Address) []byte {
	return self.getEthAccount(addr).Code
}

func (self *StateDB) SetCode(addr ethgo.Address, code []byte) {
	codeHash := crypto.Keccak256Hash(code)
	account := self.getEthAccount(addr)
	account.CodeHash = codeHash
	account.Code = code
	self.cacheDB.PutEthAccount(addr, account)
}

func (self *StateDB) GetCodeSize(addr ethgo.Address) int {
	// todo : add cache to speed up
	return len(self.GetCode(addr))
}

func (self *StateDB) Suicide(addr ethgo.Address) bool {
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

func (self *StateDB) HasSuicided(addr ethgo.Address) bool {
	return self.Suicided[addr]
}

func (self *StateDB) Exist(addr ethgo.Address) bool {
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

func (self *StateDB) Empty(addr ethgo.Address) bool {
	acct := self.getEthAccount(addr)

	balance, err := self.BalanceHandle.GetBalance(self.cacheDB, addr)
	if err != nil {
		self.cacheDB.SetDbErr(err)
		return false
	}

	return acct.IsEmpty() && balance.Sign() == 0
}

func (self *StateDB) AddLog(log *ethgo.StorageLog) {
	self.logs = append(self.logs, log)
}

func (self *StateDB) AddPreimage(ethgo.Hash, []byte) {
	// todo
}

func (self *StateDB) ForEachStorage(ethgo.Address, func(ethgo.Hash, ethgo.Hash) bool) error {
	panic("todo")
}

func (self *StateDB) CreateAccount(ethgo.Address) {
	return
}

func (self *StateDB) Snapshot() int {
	changes := self.cacheDB.memdb.DeepClone()
	suicided := make(map[ethgo.Address]bool)
	for k, v := range self.Suicided {
		suicided[k] = v
	}

	sn := &snapshot{
		changes:  changes,
		suicided: suicided,
		logsSize: len(self.logs),
		refund:   self.refund,
	}

	self.snapshots = append(self.snapshots, sn)

	return len(self.snapshots) - 1
}

func (self *StateDB) RevertToSnapshot(idx int) {
	if idx+1 > len(self.snapshots) {
		panic("can not to revert snapshot")
	}

	sn := self.snapshots[idx]

	self.snapshots = self.snapshots[:idx]
	self.cacheDB.memdb = sn.changes
	self.Suicided = sn.suicided
	self.refund = sn.refund
	self.logs = self.logs[:sn.logsSize]
}

func (self *StateDB) SubBalance(addr ethgo.Address, val *big.Int) {
	err := self.BalanceHandle.SubBalance(self.cacheDB, addr, val)
	if err != nil {
		self.cacheDB.SetDbErr(err)
		return
	}
}

func (self *StateDB) AddBalance(addr ethgo.Address, val *big.Int) {
	err := self.BalanceHandle.AddBalance(self.cacheDB, addr, val)
	if err != nil {
		self.cacheDB.SetDbErr(err)
		return
	}
}

func (self *StateDB) GetBalance(addr ethgo.Address) *big.Int {
	balance, err := self.BalanceHandle.GetBalance(self.cacheDB, addr)
	if err != nil {
		self.cacheDB.SetDbErr(err)
		return big.NewInt(0)
	}

	return balance
}
