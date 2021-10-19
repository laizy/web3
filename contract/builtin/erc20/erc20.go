package erc20

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/laizy/web3"
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/crypto"
	"github.com/laizy/web3/jsonrpc"
	"github.com/laizy/web3/utils"
)

var (
	_ = big.NewInt
	_ = fmt.Printf
	_ = utils.JsonStr
)

// ERC20 is a solidity contract
type ERC20 struct {
	c *contract.Contract
}

// NewERC20 creates a new instance of the contract at a specific address
func NewERC20(addr web3.Address, provider *jsonrpc.Client) *ERC20 {
	return &ERC20{c: contract.NewContract(addr, abiERC20, provider)}
}

// Contract returns the contract object
func (a *ERC20) Contract() *contract.Contract {
	return a.c
}

// calls

// Allowance calls the allowance method in the solidity contract
func (a *ERC20) Allowance(owner web3.Address, spender web3.Address, block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	var ok bool

	out, err = a.c.Call("allowance", web3.EncodeBlock(block...), owner, spender)
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].(*big.Int)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}

	return
}

// BalanceOf calls the balanceOf method in the solidity contract
func (a *ERC20) BalanceOf(owner web3.Address, block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	var ok bool

	out, err = a.c.Call("balanceOf", web3.EncodeBlock(block...), owner)
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["balance"].(*big.Int)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}

	return
}

// Decimals calls the decimals method in the solidity contract
func (a *ERC20) Decimals(block ...web3.BlockNumber) (retval0 uint8, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	var ok bool

	out, err = a.c.Call("decimals", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].(uint8)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}

	return
}

// Name calls the name method in the solidity contract
func (a *ERC20) Name(block ...web3.BlockNumber) (retval0 string, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	var ok bool

	out, err = a.c.Call("name", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].(string)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}

	return
}

// Symbol calls the symbol method in the solidity contract
func (a *ERC20) Symbol(block ...web3.BlockNumber) (retval0 string, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	var ok bool

	out, err = a.c.Call("symbol", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].(string)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}

	return
}

// TotalSupply calls the totalSupply method in the solidity contract
func (a *ERC20) TotalSupply(block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	var ok bool

	out, err = a.c.Call("totalSupply", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].(*big.Int)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}

	return
}

// txns

// Approve sends a approve transaction in the solidity contract
func (a *ERC20) Approve(spender web3.Address, value *big.Int) *contract.Txn {
	return a.c.Txn("approve", spender, value)
}

// Transfer sends a transfer transaction in the solidity contract
func (a *ERC20) Transfer(to web3.Address, value *big.Int) *contract.Txn {
	return a.c.Txn("transfer", to, value)
}

// TransferFrom sends a transferFrom transaction in the solidity contract
func (a *ERC20) TransferFrom(from web3.Address, to web3.Address, value *big.Int) *contract.Txn {
	return a.c.Txn("transferFrom", from, to, value)
}

// events

//ApprovalEvent
type ApprovalEvent struct {
	Owner   web3.Address
	Spender web3.Address
	Value   *big.Int
	Raw     *web3.Log
}

var ApprovalEventID = crypto.Keccak256Hash([]byte("Approval(address,address,uint256)"))

func (a *ERC20) ApprovalTopicFilter(owner []web3.Address, spender []web3.Address) [][]web3.Hash {
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}
	query := append([][]interface{}{{ApprovalEventID}}, ownerRule, spenderRule)
	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (a *ERC20) FilterApprovalEvent(owner []web3.Address, spender []web3.Address, startBlock uint64, endBlock ...uint64) ([]*ApprovalEvent, error) {
	topic := a.ApprovalTopicFilter(owner, spender)
	logs, err := a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*ApprovalEvent, 0)
	evts := a.c.Abi.Events["Approval"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem ApprovalEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

//TransferEvent
type TransferEvent struct {
	From  web3.Address
	To    web3.Address
	Value *big.Int
	Raw   *web3.Log
}

func (a *ERC20) FilterTransferEvent(opts *web3.FilterOpts, from []web3.Address, to []web3.Address) ([]*TransferEvent, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, err := a.c.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	res := make([]*TransferEvent, 0)
	evts := a.c.Abi.Events["Transfer"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem TransferEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}
