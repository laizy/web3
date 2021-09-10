package abigen

import (
	"bytes"
	"fmt"
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
        address indexed _from,
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

	fmt.Println(b.String())
}

func TestGenCode(t *testing.T) {
	artifact := &compiler.Artifact{
		Abi: `[{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}]`,
	}
	config := &Config{
		Package: "testPkg",
		Output:  "testOutput",
		Name:    "ERC20",
	}

	b := bytes.NewBuffer(nil)
	err := GenCodeToWriter(config.Name, artifact, config, b, nil)
	assert.Nil(t, err)
	expected := `package testPkg

import (
	"fmt"
	"math/big"

	"github.com/laizy/web3"
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/jsonrpc"
)

var (
	_ = big.NewInt
	_ = fmt.Printf
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

// txns


//Approval
type Approval struct { 
    owner  web3.Address
    spender  web3.Address
    value  *big.Int
}

func (a *ERC20) FilterApproval(opts *web3.FilterOpts, owner []web3.Address, spender []web3.Address)([]*Approval, error){
	
    var ownerRule []interface{}
    for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
    var spenderRule []interface{}
    for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}
    
    logs, err := a.c.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	res := make([]*Approval, 0)
	evts := a.c.Abi.Events["Approval"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem Approval
		err = mapToStruct(args, &evtItem)
		if err != nil {
			return nil, err
		}
		res = append(res, &evtItem)
	}
	return res, nil
}

//Transfer
type Transfer struct { 
    from  web3.Address
    to  web3.Address
    value  *big.Int
}

func (a *ERC20) FilterTransfer(opts *web3.FilterOpts, from []web3.Address, to []web3.Address)([]*Transfer, error){
	
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
	res := make([]*Transfer, 0)
	evts := a.c.Abi.Events["Transfer"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem Transfer
		err = mapToStruct(args, &evtItem)
		if err != nil {
			return nil, err
		}
		res = append(res, &evtItem)
	}
	return res, nil
}

func mapToStruct(m map[string]interface{}, evt interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, evt)
}`
	assert.Equal(t, expected, b.String())
}
