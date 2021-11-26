package erc20

import (
	"math/big"

	"github.com/laizy/web3"
	"github.com/laizy/web3/utils"
	"github.com/laizy/web3/utils/u256"
)

func (self *ERC20) AmountFloatWithDecimals(amount float64) *big.Int {
	const floatFactor = 10000000
	amt := u256.New(self.AmountWithDecimals(uint64(amount * floatFactor))).Div(floatFactor)
	return amt.ToBigInt()
}

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

func (a *ERC20) FilterApproval(opts *web3.FilterOpts, _owner []web3.Address, _spender []web3.Address) ([]*web3.Log, error) {
	return a.c.FilterLogs(opts, "Approval", eraseType(_owner), eraseType(_spender))
}

func (a *ERC20) FilterTransfer(opts *web3.FilterOpts, from []web3.Address, to []web3.Address) ([]*web3.Log, error) {
	return a.c.FilterLogs(opts, "Transfer", eraseType(from), eraseType(to))
}

func eraseType(addrs []web3.Address) []interface{} {
	var values []interface{}
	for _, addr := range addrs {
		values = append(values, addr)
	}

	return values
}
