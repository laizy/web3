package abigen

import (
	"bytes"
	"strings"
	"testing"
	"text/template"

	"github.com/laizy/web3/abi"
	"github.com/laizy/web3/compiler"
	"github.com/stretchr/testify/assert"
)

func TestGenCode(t *testing.T) {

	artifact := &compiler.Artifact{
		Abi: `[{"anonymous":false,"inputs":[{"indexed":true,"name":"_owner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}]`,
	}
	config := &Config{
		Package: "testPkg",
		Output:  "testOutput",
		Name:    "testName",
	}
	funcMap := template.FuncMap{
		"title":       strings.Title,
		"clean":       cleanName,
		"arg":         encodeArg,
		"outputArg":   outputArg,
		"funcName":    funcName,
		"tupleElems":  tupleElems,
		"tupleLen":    tupleLen,
		"toCamelCase": toCamelCase,
	}
	tmplAbi, err := template.New("eth-abi").Funcs(funcMap).Parse(templateAbiStr)
	assert.Nil(t, err)
	// parse abi
	abi, err := abi.NewABI(artifact.Abi)
	assert.Nil(t, err)
	input := map[string]interface{}{
		"Ptr":      "a",
		"Config":   config,
		"Contract": artifact,
		"Abi":      abi,
		"Name":     "ERC20",
	}

	var b bytes.Buffer
	if err := tmplAbi.Execute(&b, input); err != nil {
		assert.Nil(t, err)
	}

	expected := `package testPkg

import (
    "encoding/json"
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
    Owner  web3.Address
    Spender  web3.Address
    Value  *big.Int
}

func (a *ERC20) FilterApproval(opts *web3.FilterOpts, owner []web3.Address, spender []web3.Address)([]*Approval, error){
	
    var _ownerRule []interface{}
    for _, _ownerItem := range owner {
		_ownerRule = append(_ownerRule, _ownerItem)
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
    From  web3.Address
    To  web3.Address
    Value  *big.Int
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
