// Copyright (C) 2021 The Ontology Authors
// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package evm

import (
	"math/big"

	"github.com/laizy/web3"
)

// StateDB is an EVM database for full state querying.
type StateDB interface {
	CreateAccount(web3.Address)

	SubBalance(web3.Address, *big.Int)
	AddBalance(web3.Address, *big.Int)
	GetBalance(web3.Address) *big.Int

	GetNonce(web3.Address) uint64
	SetNonce(web3.Address, uint64)

	GetCodeHash(web3.Address) web3.Hash
	GetCode(web3.Address) []byte
	SetCode(web3.Address, []byte)
	GetCodeSize(web3.Address) int

	AddRefund(uint64)
	SubRefund(uint64)
	GetRefund() uint64

	GetCommittedState(web3.Address, web3.Hash) web3.Hash
	GetState(web3.Address, web3.Hash) web3.Hash
	SetState(web3.Address, web3.Hash, web3.Hash)

	Suicide(web3.Address) bool
	HasSuicided(web3.Address) bool

	// Exist reports whether the given account exists in state.
	// Notably this should also return true for suicided accounts.
	Exist(web3.Address) bool
	// Empty returns whether the given account is empty. Empty
	// is defined according to EIP161 (balance = nonce = code = 0).
	Empty(web3.Address) bool

	RevertToSnapshot(int)
	DiscardSnapshot(idx int)
	Snapshot() int

	AddLog(log *web3.StorageLog)
	AddPreimage(web3.Hash, []byte)

	ForEachStorage(web3.Address, func(web3.Hash, web3.Hash) bool) error
}

// CallContext provides a basic interface for the EVM calling conventions. The EVM
// depends on this context being implemented for doing subcalls and initialising new EVM contracts.
type CallContext interface {
	// Call another contract
	Call(env *EVM, me ContractRef, addr web3.Address, data []byte, gas, value *big.Int) ([]byte, error)
	// Take another's contract code and execute within our own context
	CallCode(env *EVM, me ContractRef, addr web3.Address, data []byte, gas, value *big.Int) ([]byte, error)
	// Same as CallCode except sender and value is propagated from parent to child scope
	DelegateCall(env *EVM, me ContractRef, addr web3.Address, data []byte, gas *big.Int) ([]byte, error)
	// Create a new contract
	Create(env *EVM, me ContractRef, data []byte, gas, value *big.Int) ([]byte, web3.Address, error)
}
