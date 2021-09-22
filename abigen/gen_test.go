package abigen

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/laizy/web3/utils/common"

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
	
event NoName(
address indexed,
address
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

	 constructor(){ emit Deposit(msg.sender,msg.sender,100000,bytes("test")); emit NoName(msg.sender,msg.sender);}

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

	var fromRule []interface{}
	for _, _fromItem := range from {
		fromRule = append(fromRule, _fromItem)
	}

	var toRule []interface{}
	for _, _toItem := range to {
		toRule = append(toRule, _toItem)
	}

	logs, err := _a.c.FilterLogs(opts, "Deposit", fromRule, toRule)
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

//NoNameEvent
type NoNameEvent struct {
	Arg0 web3.Address
	Arg1 web3.Address
	Raw  *web3.Log
}

func (_a *Sample) FilterNoNameEvent(opts *web3.FilterOpts, arg0 []web3.Address) ([]*NoNameEvent, error) {

	var arg0Rule []interface{}
	for _, arg0Item := range arg0 {
		arg0Rule = append(arg0Rule, arg0Item)
	}

	logs, err := _a.c.FilterLogs(opts, "NoName", arg0Rule)
	if err != nil {
		return nil, err
	}
	res := make([]*NoNameEvent, 0)
	evts := _a.c.Abi.Events["NoName"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem NoNameEvent
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

var testFile = struct {
	imports string
	name    string
	tester  string
}{`
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/laizy/web3"
	"github.com/laizy/web3/jsonrpc"
	"github.com/laizy/web3/testutil"
	"github.com/stretchr/testify/require"
`,
	"ParseEvent",

	`
assert := require.New(t)

	server := testutil.NewTestServer(t, nil)
	defer server.Close()

	client, err := jsonrpc.NewClient(server.HTTPAddr())
	assert.Nil(err)

	sender := server.Account(0)
	nonce, err := client.Eth().GetNonce(sender, web3.Latest)
	assert.Nil(err)
	txn := DeploySample(client, server.Account(0))
	err = txn.Do()
	assert.Nil(err)

	creationAddress := crypto.CreateAddress(common.BytesToAddress(sender.Bytes()), nonce)

	sss := NewSample(web3.BytesToAddress(creationAddress.Bytes()), client)
	sss.Contract().SetFrom(sender)

	time.Sleep(5 * time.Second)
	events, err := sss.FilterDepositEvent(&web3.FilterOpts{Start: 0}, nil, nil)
	assert.Nil(err)
	eveent2, err := sss.FilterNoNameEvent(&web3.FilterOpts{Start: 0}, nil)
	assert.Nil(err)

	want1 := DepositEvent{From: sender, To: sender, Amount: big.NewInt(100000), Data: []byte("test")}
	for _, eve := range events {
		assert.Equal(want1.From, eve.From)
		assert.Equal(want1.To, eve.To)
		assert.Equal(want1.Amount, eve.Amount)
		assert.Equal(want1.Data, eve.Data)
	}

	want2 := NoNameEvent{Arg0: sender, Arg1: sender}
	for _, eve := range eveent2 {
		assert.Equal(want2.Arg0, eve.Arg0)
		assert.Equal(want2.Arg1, eve.Arg1)

	}
`,
}

func TestBind(t *testing.T) {
	if testutil.IsSolcInstalled() == false {
		t.Skipf("skipping since solidity is not installed")
	}
	assert := require.New(t)
	// Skip the test if no Go command can be found
	gocmd := runtime.GOROOT() + "/bin/go"
	if !common.FileExist(gocmd) {
		t.Skip("go sdk not found for testing")
	}
	// Create a temporary workspace for the test suite
	ws, err := ioutil.TempDir("", "binding-test")
	if err != nil {
		t.Fatalf("failed to create temporary workspace: %v", err)
	}
	//defer os.RemoveAll(ws)

	pkg := filepath.Join(ws, "bindtest")
	if err = os.MkdirAll(pkg, 0700); err != nil {
		t.Fatalf("failed to create package: %v", err)
	}

	// Generate the test suite for all the contracts

	artifacts := map[string]*compiler.Artifact{
		"Sample": Artifact,
	}
	config := &Config{
		Package: "bindtest",
		Output:  pkg,
		Name:    "Sample",
	}
	if err := GenCode(artifacts, config); err != nil {
		t.Fatalf("genCode: %v", err)
	}

	code := fmt.Sprintf(`
			package bindtest

			import (
				"testing"
				%s
			)

			func Test%s(t *testing.T) {
				%s
			}
		`, testFile.imports, testFile.name, testFile.tester)
	err = ioutil.WriteFile(filepath.Join(pkg, strings.ToLower(testFile.name+"_test.go")), []byte(code), 06666)
	assert.Nil(err)

	// Convert the package to go modules and use the current source for go-ethereum
	moder := exec.Command(gocmd, "mod", "init", "bindtest")
	moder.Dir = pkg
	if out, err := moder.CombinedOutput(); err != nil {
		t.Fatalf("failed to convert binding test to modules: %v\n%s", err, out)
	}
	pwd, _ := os.Getwd()
	replacer := exec.Command(gocmd, "mod", "edit", "-x", "-require", "github.com/laizy/web3@v0.0.0", "-replace", "github.com/laizy/web3="+filepath.Join(pwd, "..")) // Repo root
	replacer.Dir = pkg
	if out, err := replacer.CombinedOutput(); err != nil {
		t.Fatalf("failed to replace binding test dependency to current source tree: %v\n%s", err, out)
	}
	tidier := exec.Command(gocmd, "mod", "tidy")
	tidier.Dir = pkg
	if out, err := tidier.CombinedOutput(); err != nil {
		t.Fatalf("failed to tidy Go module file: %v\n%s", err, out)
	}
	// Test the entire package and report any failures
	cmd := exec.Command(gocmd, "test", "-v", "-count", "1")
	cmd.Dir = pkg
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to run binding test: %v\n%s", err, out)
	}
}
