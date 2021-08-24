package erc20

import (
	"math/big"

	"github.com/laizy/web3"
	"github.com/laizy/web3/utils"
	"github.com/laizy/web3/utils/u256"
)

func (self *ERC20) AmountWithDecimals(amount uint64) *big.Int {
	decimals, err := self.Decimals(web3.Latest)
	utils.Ensure(err)

	return u256.New(10).ExpUint8(decimals).MulUint64(amount).ToBigInt()
}

func (self *ERC20) AmountWithoutDecimals(amount *big.Int) uint64 {
	decimals, err := self.Decimals(web3.Latest)
	utils.Ensure(err)

	return u256.New(amount).Div(u256.New(10).ExpUint8(decimals)).ToBigInt().Uint64()
}
