package abigen

import (
	"go/format"
	"testing"

	"github.com/laizy/web3/abi"
	"github.com/laizy/web3/compiler"
	"github.com/laizy/web3/testutil"
	"github.com/laizy/web3/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var Artifact = func() *compiler.Artifact {
	if testutil.IsSolcInstalled() == false {
		return nil
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

    struct Transaction {
        uint256 timestamp;
        QueueOrigin l1QueueOrigin;
        address entrypoint;
        bytes data;
    }

    enum QueueOrigin {
        SEQUENCER_QUEUE,
        L1TOL2_QUEUE
    }

    function TestStruct(Transaction memory a,bytes memory b) public returns (bytes memory){
        return  b;
    }

    function getTxes(Transaction[] memory txes) external view returns (Transaction[] memory) {
        return txes;
    }
}
`
	solc := &compiler.Solidity{Path: "solc"}
	output, err := solc.CompileCode(code)
	utils.Ensure(err)
	return output["<stdin>:Sample"]
}()

func TestCodeGen(t *testing.T) {
	if testutil.IsSolcInstalled() == false {
		t.Skipf("skipping since solidity is not installed")
	}
	config := &Config{
		Package: "binding",
		Output:  "sample",
		Name:    "Sample",
	}

	artifacts := map[string]*compiler.Artifact{
		"Sample": Artifact,
	}
	res, err := NewGenerator(config, artifacts).Gen()
	assert.Nil(t, err)

	expected, _ := format.Source([]byte(`package binding

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/laizy/web3"
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/jsonrpc"
	"github.com/laizy/web3/utils"
	"github.com/mitchellh/mapstructure"
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
func (_a *Sample) Contract() *contract.Contract {
	return _a.c
}

// calls

// GetTxes calls the getTxes method in the solidity contract
func (_a *Sample) GetTxes(txes []Transaction, block ...web3.BlockNumber) (retval0 []Transaction, err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	out, err = _a.c.Call("getTxes", web3.EncodeBlock(block...), txes)
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

// TestStruct sends a TestStruct transaction in the solidity contract
func (_a *Sample) TestStruct(a Transaction, b []byte) *contract.Txn {
	return _a.c.Txn("TestStruct", a, b)
}

// events

//DepositEvent
type DepositEvent struct {
	From   web3.Address
	To     web3.Address
	Amount *big.Int
	Data   []byte
	Raw    *web3.Log
}

func (_a *Sample) FilterDepositEvent(opts *web3.FilterOpts, from []web3.Address, to []web3.Address) ([]*DepositEvent, error) {

	var _fromRule []interface{}
	for _, _fromItem := range from {
		_fromRule = append(_fromRule, _fromItem)
	}

	var _toRule []interface{}
	for _, _toItem := range to {
		_toRule = append(_toRule, _toItem)
	}

	logs, err := _a.c.FilterLogs(opts, "Deposit", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	res := make([]*DepositEvent, 0)
	evts := _a.c.Abi.Events["Deposit"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem DepositEvent
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
	From   web3.Address
	To     web3.Address
	Amount web3.Address
	Raw    *web3.Log
}

func (_a *Sample) FilterTransferEvent(opts *web3.FilterOpts, from []web3.Address, to []web3.Address, amount []web3.Address) ([]*TransferEvent, error) {

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

	logs, err := _a.c.FilterLogs(opts, "Transfer", fromRule, toRule, amountRule)
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
`))

	assert.Equal(t, string(expected), string(res.AbiFiles[0].Code))
}

func TestTupleStructs(t *testing.T) {
	if testutil.IsSolcInstalled() == false {
		t.Skipf("skipping since solidity is not installed")
	}
	assert := require.New(t)
	code, err := NewStructDefExtractor().ExtractFromAbi(abi.MustNewABI(Artifact.Abi)).RenderGoCode("binding")
	assert.Nil(err)
	expected, _ := format.Source([]byte(`package binding

import (
	"fmt"
	"math/big"

	"github.com/laizy/web3"
)

var (
	_ = big.NewInt
	_ = fmt.Printf
	_ = web3.HexToAddress
)

type Transaction struct {
	Timestamp     *big.Int
	L1QueueOrigin uint8
	Entrypoint    web3.Address
	Data          []byte
}
`))

	assert.Equal(string(expected), code)
}

func TestGenStruct(t *testing.T) {
	if testutil.IsSolcInstalled() == false {
		t.Skipf("skipping since solidity is not installed")
	}
	assert := require.New(t)

	defs := NewStructDefExtractor()
	abi1, err := abi.NewABI(Artifact.Abi)
	assert.Nil(err)

	defs.ExtractFromAbi(abi1)

	old := len(defs.Defs)
	defs.ExtractFromAbi(abi1)
	assert.Equal(old, len(defs.Defs)) // test dulplicate case

	var oldname string
	for name := range defs.Defs {
		oldname = name //read an struct from it
		break
	}
	defs.Defs[oldname] = &StructDef{Name: oldname}
	assert.PanicsWithError(ErrConflictDef.Error(), func() {
		defs.ExtractFromAbi(abi1)
	})
}

//
//var testFile = struct {
//	imports string
//	name    string
//	tester  string
//}{
//	`
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/ethereum/go-ethereum/crypto"
//	"github.com/laizy/web3"
//	"github.com/laizy/web3/jsonrpc"
//	"github.com/laizy/web3/testutil"
//	"github.com/stretchr/testify/require"
//	"testing"
//`,
//	"CallContract",
//	`
//
//`,
//}
//
//func TestBind(t *testing.T) {
//	if testutil.IsSolcInstalled() == false {
//		t.Skipf("skipping since solidity is not installed")
//	}
//	assert := require.New(t)
//	// Skip the test if no Go command can be found
//	gocmd := runtime.GOROOT() + "/bin/go"
//	if !common.FileExist(gocmd) {
//		t.Skip("go sdk not found for testing")
//	}
//	// Create a temporary workspace for the test suite
//	ws, err := ioutil.TempDir("", "binding-test")
//	if err != nil {
//		t.Fatalf("failed to create temporary workspace: %v", err)
//	}
//	//defer os.RemoveAll(ws)
//
//	pkg := filepath.Join(ws, "bindtest")
//	if err = os.MkdirAll(pkg, 0700); err != nil {
//		t.Fatalf("failed to create package: %v", err)
//	}
//
//	// Generate the test suite for all the contracts
//
//	artifacts := map[string]*compiler.Artifact{
//		"Sample": Artifact,
//	}
//	config := &Config{
//		Package: "bindtest",
//		Output:  pkg,
//		Name:    "Sample",
//	}
//	res, err := NewGenerator(config, artifacts).Gen()
//	assert.Nil(err)
//	for _, abif := range res.AbiFiles {
//		err = ioutil.WriteFile(filepath.Join(pkg, strings.ToLower(abif.FileName)+"_abi.go"), abif.Code, 0666)
//		assert.Nil(err)
//	}
//	for _, binf := range res.BinFiles {
//		err = ioutil.WriteFile(filepath.Join(pkg, strings.ToLower(binf.FileName)+"_bin.go"), binf.Code, 0666)
//		assert.Nil(err)
//	}
//
//	code := fmt.Sprintf(`
//			package bindtest
//
//			import (
//				"testing"
//				%s
//			)
//
//			func Test%s(t *testing.T) {
//				%s
//			}
//		`, testFile.imports, testFile.name, testFile.tester)
//	err = ioutil.WriteFile(filepath.Join(pkg, strings.ToLower(testFile.name+"_test.go")), []byte(code), 06666)
//	assert.Nil(err)
//}
