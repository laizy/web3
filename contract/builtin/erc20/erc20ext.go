package erc20

import (
	"encoding/json"
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

// Erc20Approval represents a Approval event raised by the Erc20 contract.
type Erc20Approval struct {
	Owner   web3.Address
	Spender web3.Address
	Value   *big.Int
	Raw     web3.Log // Blockchain specific contextual infos
}

func (a *ERC20) FilterApproval(opts *web3.FilterOpts, _owner []web3.Address, _spender []web3.Address) ([]*Erc20Approval, error) {
	var _ownerRule []interface{}
	for _, _fromItem := range _owner {
		_ownerRule = append(_ownerRule, _fromItem)
	}
	var _spenderRule []interface{}
	for _, _toItem := range _spender {
		_spenderRule = append(_spenderRule, _toItem)
	}
	logs, err := a.c.FilterLogs(opts, "Approval", _ownerRule, _spenderRule)
	if err != nil {
		return nil, err
	}
	res := make([]*Erc20Approval, 0)
	approveEvent := a.c.Abi.Events["Approval"]
	for _, log := range logs {
		args, err := approveEvent.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var appr Erc20Approval
		err = mapToStruct(args, &appr)
		if err != nil {
			return nil, err
		}
		res = append(res, &appr)
	}
	return res, nil
}

type Erc20Transfer struct {
	From  web3.Address
	To    web3.Address
	Value *big.Int
}

func (a *ERC20) FilterTransfer(opts *web3.FilterOpts, from []web3.Address, to []web3.Address) ([]*Erc20Transfer, error) {
	var _fromRule []interface{}
	for _, _fromItem := range from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range to {
		_toRule = append(_toRule, _toItem)
	}

	logs, err := a.c.FilterLogs(opts, "Transfer", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	res := make([]*Erc20Transfer, 0)
	evts := a.c.Abi.Events["Transfer"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var tf Erc20Transfer
		err = mapToStruct(args, &tf)
		if err != nil {
			return nil, err
		}
		res = append(res, &tf)
	}
	return res, nil
}

func mapToStruct(m map[string]interface{}, evt interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, evt)
}
