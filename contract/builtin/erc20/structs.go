package erc20

import (
	"fmt"
	"math/big"

	"github.com/laizy/web3"
)

var (
	_ = big.NewInt
	_ = fmt.Printf
	_ = web3.HexToAddress
)

type ApprovalEvent struct {
	Owner   web3.Address
	Spender web3.Address
	Value   *big.Int

	Raw *web3.Log
}

type TransferEvent struct {
	From  web3.Address
	To    web3.Address
	Value *big.Int

	Raw *web3.Log
}
