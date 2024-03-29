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

package web3

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/laizy/web3/evm/errors"
	"github.com/laizy/web3/utils"
	"github.com/laizy/web3/utils/common/hexutil"
)

// ExecutionResult includes all output after executing given evm
// message no matter the execution itself is successful or not.
type ExecutionResult struct {
	UsedGas      uint64        // Total used gas but include the refunded gas
	Err          error         // Any error encountered during the execution(listed in core/vm/errors.go)
	ReturnData   hexutil.Bytes // Returned data from evm(function result or data supplied with revert opcode)
	RevertReason string
}

// revert data signature is: Error(string) (0x08c379a0)
// https://ethereum.stackexchange.com/questions/83528/how-can-i-get-the-revert-reason-of-a-call-in-solidity-so-that-i-can-use-it-in-th
func DecodeRevert(ret []byte) (string, bool) {
	if len(ret) >= 68 && hex.EncodeToString(ret[:4]) == "08c379a0" {
		// data layout: sig(4bytes) + strpos(32bytes,should equal 2) + strlength(32bytes) + strdata
		data := ret[36:]
		length, err := readLength(data)
		utils.Ensure(err)
		return string(data[32 : 32+length]), true
	}

	return "", false
}

//copied from abi to avoid cycle dependencies
func readLength(data []byte) (int, error) {
	lengthBig := big.NewInt(0).SetBytes(data[0:32])
	if lengthBig.BitLen() > 63 {
		return 0, fmt.Errorf("length larger than int64: %v", lengthBig.Int64())
	}
	length := int(lengthBig.Uint64())
	if length > len(data) {
		return 0, fmt.Errorf("length insufficient %v require %v", len(data), length)
	}
	return length, nil
}

// Unwrap returns the internal evm error which allows us for further
// analysis outside.
func (result *ExecutionResult) Unwrap() error {
	return result.Err
}

// Failed returns the indicator whether the execution is successful or not
func (result *ExecutionResult) Failed() bool { return result.Err != nil }

// Return is a helper function to help caller distinguish between revert reason
// and function return. Return returns the data after execution if no error occurs.
func (result *ExecutionResult) Return() []byte {
	if result.Err != nil {
		return nil
	}
	return CopyBytes(result.ReturnData)
}

// Revert returns the concrete revert reason if the execution is aborted by `REVERT`
// opcode. Note the reason can be nil if no data supplied with revert opcode.
func (result *ExecutionResult) Revert() []byte {
	if result.Err != errors.ErrExecutionReverted {
		return nil
	}
	return CopyBytes(result.ReturnData)
}
