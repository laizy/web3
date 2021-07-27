package erc20

import (
	"math/big"

	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/utils"
	"github.com/umbracle/go-web3/utils/u256"
)

func (self *ERC20) AmountWithDecimals(amount uint64) *big.Int {
	decimals, err := self.Decimals(web3.Latest)
	utils.Ensure(err)

	return u256.New(10).ExpUint8(decimals).MulUint64(amount).ToBigInt()
}
