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

// ENS is a solidity contract
type ENS struct {
	c *contract.Contract
}

// DeployENS deploys a new ENS contract
func DeployENS(provider *jsonrpc.Client, from web3.Address, args ...interface{}) *contract.Txn {
	return contract.DeployContract(provider, from, abiENS, binENS, args...)
}

// NewENS creates a new instance of the contract at a specific address
func NewENS(addr web3.Address, provider *jsonrpc.Client) *ENS {
	return &ENS{c: contract.NewContract(addr, abiENS, provider)}
}

// Contract returns the contract object
func (_a *ENS) Contract() *contract.Contract {
	return _a.c
}

// calls

// Owner calls the owner method in the solidity contract
func (_a *ENS) Owner(node [32]byte, block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("owner", web3.EncodeBlock(block...), node)
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// Resolver calls the resolver method in the solidity contract
func (_a *ENS) Resolver(node [32]byte, block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("resolver", web3.EncodeBlock(block...), node)
	if err != nil {
		return
	}

	// decode outputs

	if err = mapstructure.Decode(out["0"], &retval0); err != nil {
		err = fmt.Errorf("failed to encode output at index 0")
	}

	return
}

// Ttl calls the ttl method in the solidity contract
func (_a *ENS) Ttl(node [32]byte, block ...web3.BlockNumber) (retval0 uint64, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("ttl", web3.EncodeBlock(block...), node)
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

// SetOwner sends a setOwner transaction in the solidity contract
func (_a *ENS) SetOwner(node [32]byte, owner web3.Address) *contract.Txn {
	return _a.c.Txn("setOwner", node, owner)
}

// SetResolver sends a setResolver transaction in the solidity contract
func (_a *ENS) SetResolver(node [32]byte, resolver web3.Address) *contract.Txn {
	return _a.c.Txn("setResolver", node, resolver)
}

// SetSubnodeOwner sends a setSubnodeOwner transaction in the solidity contract
func (_a *ENS) SetSubnodeOwner(node [32]byte, label [32]byte, owner web3.Address) *contract.Txn {
	return _a.c.Txn("setSubnodeOwner", node, label, owner)
}

// SetTTL sends a setTTL transaction in the solidity contract
func (_a *ENS) SetTTL(node [32]byte, ttl uint64) *contract.Txn {
	return _a.c.Txn("setTTL", node, ttl)
}

// events

var NewOwnerEventID = crypto.Keccak256Hash([]byte("NewOwner(bytes32,bytes32,address)"))

func (_a *ENS) NewOwnerTopicFilter(node [][32]byte, label [][32]byte) [][]web3.Hash {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	var labelRule []interface{}
	for _, labelItem := range label {
		labelRule = append(labelRule, labelItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{NewOwnerEventID}, nodeRule, labelRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *ENS) FilterNewOwnerEvent(node [][32]byte, label [][32]byte, startBlock uint64, endBlock ...uint64) ([]*NewOwnerEvent, error) {
	topic := _a.NewOwnerTopicFilter(node, label)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*NewOwnerEvent, 0)
	evts := _a.c.Abi.Events["NewOwner"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem NewOwnerEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

var NewResolverEventID = crypto.Keccak256Hash([]byte("NewResolver(bytes32,address)"))

func (_a *ENS) NewResolverTopicFilter(node [][32]byte) [][]web3.Hash {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{NewResolverEventID}, nodeRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *ENS) FilterNewResolverEvent(node [][32]byte, startBlock uint64, endBlock ...uint64) ([]*NewResolverEvent, error) {
	topic := _a.NewResolverTopicFilter(node)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*NewResolverEvent, 0)
	evts := _a.c.Abi.Events["NewResolver"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem NewResolverEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

var NewTTLEventID = crypto.Keccak256Hash([]byte("NewTTL(bytes32,uint64)"))

func (_a *ENS) NewTTLTopicFilter(node [][32]byte) [][]web3.Hash {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{NewTTLEventID}, nodeRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *ENS) FilterNewTTLEvent(node [][32]byte, startBlock uint64, endBlock ...uint64) ([]*NewTTLEvent, error) {
	topic := _a.NewTTLTopicFilter(node)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*NewTTLEvent, 0)
	evts := _a.c.Abi.Events["NewTTL"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem NewTTLEvent
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}

var TransferEventID = crypto.Keccak256Hash([]byte("Transfer(bytes32,address)"))

func (_a *ENS) TransferTopicFilter(node [][32]byte) [][]web3.Hash {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	var query [][]interface{}
	query = append(query, []interface{}{TransferEventID}, nodeRule)

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func (_a *ENS) FilterTransferEvent(node [][32]byte, startBlock uint64, endBlock ...uint64) ([]*TransferEvent, error) {
	topic := _a.TransferTopicFilter(node)

	logs, err := _a.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*TransferEvent, 0)
	evts := _a.c.Abi.Events["Transfer"]
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
