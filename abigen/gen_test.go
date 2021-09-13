package abigen

import (
	"bytes"
	"go/format"
	"testing"

	"github.com/laizy/web3/compiler"
	"github.com/laizy/web3/testutil"
	"github.com/laizy/web3/utils"
	"github.com/stretchr/testify/assert"
)

func TestEventGen(t *testing.T) {
	if testutil.IsSolcInstalled() == false {
		t.Skipf("skipping since solidity is not installed")
	}
	code := `
pragma experimental ABIEncoderV2;
contract Sample {
    event Deposit (
        address indexed _from, // test name with _ will translate to From
        address indexed _to,
        uint256 _amount,
        bytes _data
    );

	event Transfer (
		address indexed from,
		address indexed to,
		address indexed amount
	);
}
`
	solc := &compiler.Solidity{Path: "solc"}
	output, err := solc.CompileCode(code)
	utils.Ensure(err)
	artifact := output["<stdin>:Sample"]
	config := &Config{
		Package: "binding",
		Output:  "sample",
		Name:    "Sample",
	}

	b := bytes.NewBuffer(nil)
	err = GenCodeToWriter(config.Name, artifact, config, b, nil)
	assert.Nil(t, err)

	expected, _ := format.Source([]byte(`package binding

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/laizy/web3"
	"github.com/laizy/web3/utils"
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/jsonrpc"
)

var (
	_ = big.NewInt
	_ = fmt.Printf
	_ = utils.JsonStr
)

// Sample is a solidity contract
type Sample struct {
	c *contract.Contract
}

// DeploySample deploys a new Sample contract
func DeploySample(provider *jsonrpc.Client, from web3.Address, args ...interface{}) *contract.Txn {
	return contract.DeployContract(provider, from, abiSample, binSample, args...)
}

// NewSample creates a new instance of the contract at a specific address
func NewSample(addr web3.Address, provider *jsonrpc.Client) *Sample {
	return &Sample{c: contract.NewContract(addr, abiSample, provider)}
}

// Contract returns the contract object
func (a *Sample) Contract() *contract.Contract {
	return a.c
}

// calls

// txns

//Deposit
type Deposit struct {
	From   web3.Address
	To     web3.Address
	Amount *big.Int
	Data   []byte
	Raw    *web3.Log
}

func (a *Sample) FilterDeposit(opts *web3.FilterOpts, from []web3.Address, to []web3.Address) ([]*Deposit, error) {

	var _fromRule []interface{}
	for _, _fromItem := range from {
		_fromRule = append(_fromRule, _fromItem)
	}

	var _toRule []interface{}
	for _, _toItem := range to {
		_toRule = append(_toRule, _toItem)
	}

	logs, err := a.c.FilterLogs(opts, "Deposit", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	res := make([]*Deposit, 0)
	evts := a.c.Abi.Events["Deposit"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem Deposit
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

//Transfer
type Transfer struct {
	From   web3.Address
	To     web3.Address
	Amount web3.Address
	Raw    *web3.Log
}

func (a *Sample) FilterTransfer(opts *web3.FilterOpts, from []web3.Address, to []web3.Address, amount []web3.Address) ([]*Transfer, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, err := a.c.FilterLogs(opts, "Transfer", fromRule, toRule, amountRule)
	if err != nil {
		return nil, err
	}
	res := make([]*Transfer, 0)
	evts := a.c.Abi.Events["Transfer"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem Transfer
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}`))

	assert.Equal(t, string(expected), b.String())
}
