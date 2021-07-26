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

	"github.com/umbracle/ethgo"
)

// StateDB is an EVM database for full state querying.
type StateDB interface {
	CreateAccount(ethgo.Address)

	SubBalance(ethgo.Address, *big.Int)
	AddBalance(ethgo.Address, *big.Int)
	GetBalance(ethgo.Address) *big.Int

	GetNonce(ethgo.Address) uint64
	SetNonce(ethgo.Address, uint64)

	GetCodeHash(ethgo.Address) ethgo.Hash
	GetCode(ethgo.Address) []byte
	SetCode(ethgo.Address, []byte)
	GetCodeSize(ethgo.Address) int

	AddRefund(uint64)
	SubRefund(uint64)
	GetRefund() uint64

	GetCommittedState(ethgo.Address, ethgo.Hash) ethgo.Hash
	GetState(ethgo.Address, ethgo.Hash) ethgo.Hash
	SetState(ethgo.Address, ethgo.Hash, ethgo.Hash)

	Suicide(ethgo.Address) bool
	HasSuicided(ethgo.Address) bool

	// Exist reports whether the given account exists in state.
	// Notably this should also return true for suicided accounts.
	Exist(ethgo.Address) bool
	// Empty returns whether the given account is empty. Empty
	// is defined according to EIP161 (balance = nonce = code = 0).
	Empty(ethgo.Address) bool

	RevertToSnapshot(int)
	Snapshot() int

	AddLog(log *ethgo.StorageLog)
	AddPreimage(ethgo.Hash, []byte)

	ForEachStorage(ethgo.Address, func(ethgo.Hash, ethgo.Hash) bool) error
}

// CallContext provides a basic interface for the EVM calling conventions. The EVM
// depends on this context being implemented for doing subcalls and initialising new EVM contracts.
type CallContext interface {
	// Call another contract
	Call(env *EVM, me ContractRef, addr ethgo.Address, data []byte, gas, value *big.Int) ([]byte, error)
	// Take another's contract code and execute within our own context
	CallCode(env *EVM, me ContractRef, addr ethgo.Address, data []byte, gas, value *big.Int) ([]byte, error)
	// Same as CallCode except sender and value is propagated from parent to child scope
	DelegateCall(env *EVM, me ContractRef, addr ethgo.Address, data []byte, gas *big.Int) ([]byte, error)
	// Create a new contract
	Create(env *EVM, me ContractRef, data []byte, gas, value *big.Int) ([]byte, ethgo.Address, error)
}
