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
	"github.com/laizy/web3"
	"github.com/laizy/web3/crypto"
	"github.com/laizy/web3/evm"
	"github.com/laizy/web3/evm/params"
	"github.com/laizy/web3/evm/storage"
	"github.com/laizy/web3/evm/storage/schema"
)

func applyTransaction(msg Message, statedb *storage.StateDB, usedGas *uint64, evm *evm.EVM, ctx Eip155Context, commitDB bool) (*web3.ExecutionResult, *web3.Receipt, error) {
	// Create a new context to be used in the EVM environment
	txContext := NewEVMTxContext(msg)

	// Update the evm with the new transaction context.
	evm.Reset(txContext, statedb)
	// Apply the transaction to the current state (included in the env)
	result, err := NewStateTransition(evm, msg, ctx.Coinbase).TransitionDb()
	if err != nil {
		return nil, nil, err
	}
	// flush changes to overlay db
	if commitDB {
		err = statedb.Commit()
		if err != nil {
			return nil, nil, err
		}
	}
	*usedGas += result.UsedGas

	// Create a new receipt for the transaction, storing the intermediate root and gas used by the tx
	// based on the eip phase, we're passing whether the root touch-delete accounts.
	status := uint64(1)
	if result.Failed() {
		status = 0
	}
	receipt := &web3.Receipt{
		Status:            status,
		TransactionHash:   msg.Hash(),
		TransactionIndex:  ctx.TxIndex,
		BlockHash:         ctx.BlockHash,
		From:              msg.From(),
		BlockNumber:       ctx.Height,
		GasUsed:           result.UsedGas,
		CumulativeGasUsed: *usedGas,
		LogsBloom:         nil,
		Logs:              nil,
	}
	// if the transaction created a contract, store the creation address in the receipt.
	if msg.To() == nil {
		nonce := statedb.GetNonce(evm.TxContext.Origin) - 1
		receipt.ContractAddress = crypto.CreateAddress(evm.TxContext.Origin, nonce)
	}
	// Set the receipt logs and create a bloom for filtering
	receipt.AddStorageLogs(statedb.GetLogs())

	return result, receipt, nil
}

// ApplyTransaction attempts to apply a transaction to the given state database
// and uses the input parameters for its environment. It returns the receipt
// for the transaction, gas used and an error if the transaction failed,
// indicating the block was invalid.
func ApplyTransaction(config *params.ChainConfig, bc schema.ChainDB, statedb *storage.StateDB,
	tx *web3.Transaction, ctx Eip155Context, usedGas *uint64,
	cfg evm.Config, checkNonce bool) (*web3.ExecutionResult, *web3.Receipt, error) {
	// Create a new context to be used in the EVM environment
	msg := MessageFromTx(tx, checkNonce)
	return ApplyMessage(config, bc, statedb, msg, ctx, usedGas, cfg, true)
}

func ApplyMessage(config *params.ChainConfig, bc schema.ChainDB, statedb *storage.StateDB, msg Message, ctx Eip155Context, usedGas *uint64, cfg evm.Config, commitDB bool) (*web3.ExecutionResult, *web3.Receipt, error) {
	// Create a new context to be used in the EVM environment
	blockContext := NewEVMBlockContext(ctx.Height, ctx.Timestamp, bc.GetBlockHash)
	vmenv := evm.NewEVM(blockContext, evm.TxContext{}, statedb, config, cfg)
	return applyTransaction(msg, statedb, usedGas, vmenv, ctx, commitDB)
}
