// Copyright (C) 2021 The Ontology Authors
// Copyright 2015 The go-ethereum Authors
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

package executor

import (
	"fmt"

	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/crypto"
	"github.com/umbracle/go-web3/evm"
	"github.com/umbracle/go-web3/evm/params"
	"github.com/umbracle/go-web3/evm/storage"
	"github.com/umbracle/go-web3/executor/remotedb"
	"github.com/umbracle/go-web3/utils"
)

func applyTransaction(msg Message, statedb *storage.StateDB, tx *web3.Transaction, usedGas *uint64, evm *evm.EVM, feeReceiver web3.Address) (*web3.ExecutionResult, *web3.Receipt, error) {
	// Create a new context to be used in the EVM environment
	txContext := NewEVMTxContext(msg)

	// Update the evm with the new transaction context.
	evm.Reset(txContext, statedb)
	// Apply the transaction to the current state (included in the env)
	result, err := ApplyMessage(evm, msg, feeReceiver)
	if err != nil {
		return nil, nil, err
	}
	// flush changes to overlay db
	err = statedb.Commit()
	if err != nil {
		return nil, nil, err
	}
	*usedGas += result.UsedGas

	// Create a new receipt for the transaction, storing the intermediate root and gas used by the tx
	// based on the eip phase, we're passing whether the root touch-delete accounts.
	receipt := &web3.Receipt{
		TransactionHash:   tx.Hash,
		TransactionIndex:  0,
		BlockHash:         web3.Hash{},
		From:              msg.From(),
		BlockNumber:       0,
		GasUsed:           result.UsedGas,
		CumulativeGasUsed: 0,
		LogsBloom:         nil,
		Logs:              nil,
	}
	// if the transaction created a contract, store the creation address in the receipt.
	if msg.To() == nil {
		receipt.ContractAddress = crypto.CreateAddress(evm.TxContext.Origin, msg.Nonce())
	}
	// Set the receipt logs and create a bloom for filtering
	receipt.BlockHash = statedb.BlockHash()
	receipt.BlockNumber = evm.Context.BlockNumber.Uint64()
	receipt.AddStorageLogs(statedb.GetLogs())

	return result, receipt, err
}

// ApplyTransaction attempts to apply a transaction to the given state database
// and uses the input parameters for its environment. It returns the receipt
// for the transaction, gas used and an error if the transaction failed,
// indicating the block was invalid.
func ApplyTransaction(config *params.ChainConfig, bc *remotedb.RemoteDB, statedb *storage.StateDB, blockHeight, timestamp uint64, tx *web3.Transaction, usedGas *uint64, feeReceiver web3.Address, cfg evm.Config, checkNonce bool) (*web3.ExecutionResult, *web3.Receipt, error) {
	// Create a new context to be used in the EVM environment
	msg := MessageFromTx(tx, checkNonce)
	fmt.Println(utils.JsonString(msg))
	blockContext := NewEVMBlockContext(blockHeight, timestamp, bc.GetBlockHash)
	vmenv := evm.NewEVM(blockContext, evm.TxContext{}, statedb, config, cfg)
	return applyTransaction(msg, statedb, tx, usedGas, vmenv, feeReceiver)
}
