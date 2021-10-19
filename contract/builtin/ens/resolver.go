package ens

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/laizy/web3"
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/crypto"
	"github.com/laizy/web3/jsonrpc"
	"github.com/laizy/web3/utils"
	"github.com/mitchellh/mapstructure"
)

var (
	_ = json.Unmarshal
	_ = big.NewInt
	_ = fmt.Printf
	_ = utils.JsonStr
	_ = mapstructure.Decode
	_ = crypto.Keccak256Hash
)

// Resolver is a solidity contract
type Resolver struct {
	c *contract.Contract
}

// DeployResolver deploys a new Resolver contract
func DeployResolver(provider *jsonrpc.Client, from web3.Address, args ...interface{}) *contract.Txn {
	return contract.DeployContract(provider, from, abiResolver, binResolver, args...)
}

// NewResolver creates a new instance of the contract at a specific address
func NewResolver(addr web3.Address, provider *jsonrpc.Client) *Resolver {
	return &Resolver{c: contract.NewContract(addr, abiResolver, provider)}
}

// Contract returns the contract object
func (_a *Resolver) Contract() *contract.Contract {
	return _a.c
}

// calls

// ABI calls the ABI method in the solidity contract
func (_a *Resolver) ABI(node [32]byte, contentTypes *big.Int, block ...web3.BlockNumber) (retval0 *big.Int, retval1 []byte, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("ABI", web3.EncodeBlock(block...), node, contentTypes)
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["contentType"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}
	if err = mapstructure.Decode(out["data"], &retval1); err != nil {
		err = fmt.Errorf("failed to encode output at index 1")
	}

	return
}

// Addr calls the addr method in the solidity contract
func (_a *Resolver) Addr(node [32]byte, block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("addr", web3.EncodeBlock(block...), node)
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["ret"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// Content calls the content method in the solidity contract
func (_a *Resolver) Content(node [32]byte, block ...web3.BlockNumber) (retval0 [32]byte, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("content", web3.EncodeBlock(block...), node)
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["ret"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// Name calls the name method in the solidity contract
func (_a *Resolver) Name(node [32]byte, block ...web3.BlockNumber) (retval0 string, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("name", web3.EncodeBlock(block...), node)
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["ret"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// Pubkey calls the pubkey method in the solidity contract
func (_a *Resolver) Pubkey(node [32]byte, block ...web3.BlockNumber) (retval0 [32]byte, retval1 [32]byte, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("pubkey", web3.EncodeBlock(block...), node)
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["x"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}
	if err = mapstructure.Decode(out["y"], &retval1); err != nil {
		err = fmt.Errorf("failed to encode output at index 1")
	}

	return
}

// SupportsInterface calls the supportsInterface method in the solidity contract
func (_a *Resolver) SupportsInterface(interfaceID [4]byte, block ...web3.BlockNumber) (retval0 bool, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("supportsInterface", web3.EncodeBlock(block...), interfaceID)
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// txns

// SetABI sends a setABI transaction in the solidity contract
func (_a *Resolver) SetABI(node [32]byte, contentType *big.Int, data []byte) *contract.Txn {
	return _a.c.Txn("setABI", node, contentType, data)
}

// SetAddr sends a setAddr transaction in the solidity contract
func (_a *Resolver) SetAddr(node [32]byte, addr web3.Address) *contract.Txn {
	return _a.c.Txn("setAddr", node, addr)
}

// SetContent sends a setContent transaction in the solidity contract
func (_a *Resolver) SetContent(node [32]byte, hash [32]byte) *contract.Txn {
	return _a.c.Txn("setContent", node, hash)
}

// SetName sends a setName transaction in the solidity contract
func (_a *Resolver) SetName(node [32]byte, name string) *contract.Txn {
	return _a.c.Txn("setName", node, name)
}

// SetPubkey sends a setPubkey transaction in the solidity contract
func (_a *Resolver) SetPubkey(node [32]byte, x [32]byte, y [32]byte) *contract.Txn {
	return _a.c.Txn("setPubkey", node, x, y)
}

// events

var ABIChangedEventID = crypto.Keccak256Hash([]byte("ABIChanged(bytes32,uint256)"))

func (_a *Resolver) ABIChangedTopicFilter(node [][32]byte, contentType []*big.Int) [][]web3.Hash {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	var contentTypeRule []interface{}
	for _, contentTypeItem := range contentType {
		contentTypeRule = append(contentTypeRule, contentTypeItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{ABIChangedEventID}, nodeRule, contentTypeRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *Resolver) FilterABIChangedEvent(node [][32]byte, contentType []*big.Int, startBlock uint64, endBlock ...uint64) ([]*ABIChangedEvent, error) {
	topic := _a.ABIChangedTopicFilter(node, contentType)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*ABIChangedEvent, 0)
	evts := _a.c.Abi.Events["ABIChanged"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem ABIChangedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

var AddrChangedEventID = crypto.Keccak256Hash([]byte("AddrChanged(bytes32,address)"))

func (_a *Resolver) AddrChangedTopicFilter(node [][32]byte) [][]web3.Hash {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{AddrChangedEventID}, nodeRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *Resolver) FilterAddrChangedEvent(node [][32]byte, startBlock uint64, endBlock ...uint64) ([]*AddrChangedEvent, error) {
	topic := _a.AddrChangedTopicFilter(node)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*AddrChangedEvent, 0)
	evts := _a.c.Abi.Events["AddrChanged"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem AddrChangedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

var ContentChangedEventID = crypto.Keccak256Hash([]byte("ContentChanged(bytes32,bytes32)"))

func (_a *Resolver) ContentChangedTopicFilter(node [][32]byte) [][]web3.Hash {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{ContentChangedEventID}, nodeRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *Resolver) FilterContentChangedEvent(node [][32]byte, startBlock uint64, endBlock ...uint64) ([]*ContentChangedEvent, error) {
	topic := _a.ContentChangedTopicFilter(node)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*ContentChangedEvent, 0)
	evts := _a.c.Abi.Events["ContentChanged"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem ContentChangedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

var NameChangedEventID = crypto.Keccak256Hash([]byte("NameChanged(bytes32,string)"))

func (_a *Resolver) NameChangedTopicFilter(node [][32]byte) [][]web3.Hash {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{NameChangedEventID}, nodeRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *Resolver) FilterNameChangedEvent(node [][32]byte, startBlock uint64, endBlock ...uint64) ([]*NameChangedEvent, error) {
	topic := _a.NameChangedTopicFilter(node)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*NameChangedEvent, 0)
	evts := _a.c.Abi.Events["NameChanged"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem NameChangedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

var PubkeyChangedEventID = crypto.Keccak256Hash([]byte("PubkeyChanged(bytes32,bytes32,bytes32)"))

func (_a *Resolver) PubkeyChangedTopicFilter(node [][32]byte) [][]web3.Hash {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{PubkeyChangedEventID}, nodeRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *Resolver) FilterPubkeyChangedEvent(node [][32]byte, startBlock uint64, endBlock ...uint64) ([]*PubkeyChangedEvent, error) {
	topic := _a.PubkeyChangedTopicFilter(node)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*PubkeyChangedEvent, 0)
	evts := _a.c.Abi.Events["PubkeyChanged"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem PubkeyChangedEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}
