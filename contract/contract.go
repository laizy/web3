package contract

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/abi"
	"github.com/umbracle/go-web3/jsonrpc"
	"github.com/umbracle/go-web3/registry"
	"github.com/umbracle/go-web3/utils"
)

// Contract is an Ethereum contract
type Contract struct {
	addr     web3.Address
	from     *web3.Address
	Abi      *abi.ABI
	Provider *jsonrpc.Client
}

// DeployContract deploys a contract
func DeployContract(provider *jsonrpc.Client, from web3.Address, abiVal *abi.ABI, bin []byte, args ...interface{}) *Txn {
	txn := &Txn{
		from:     from,
		provider: provider,
	}
	txn.Data = append(txn.Data, bin...)
	if method := abiVal.Constructor; method != nil {
		data, err := abi.Encode(args, method.Inputs)
		utils.Ensure(err)
		txn.Data = append(txn.Data, data...)
	}

	return txn
}

// NewContract creates a new contract instance
func NewContract(addr web3.Address, abi *abi.ABI, provider *jsonrpc.Client) *Contract {
	registry.Instance().RegisterFromAbi(abi)
	return &Contract{
		addr:     addr,
		Abi:      abi,
		Provider: provider,
	}
}

// Addr returns the address of the contract
func (c *Contract) Addr() web3.Address {
	return c.addr
}

// SetFrom sets the origin of the calls
func (c *Contract) SetFrom(addr web3.Address) {
	c.from = &addr
}

// EstimateGas estimates the gas for a contract call
func (c *Contract) EstimateGas(method string, args ...interface{}) (uint64, error) {
	return c.Txn(method, args).EstimateGas()
}

// Call calls a method in the contract
func (c *Contract) Call(method string, block web3.BlockNumber, args ...interface{}) (map[string]interface{}, error) {
	m, ok := c.Abi.Methods[method]
	if !ok {
		return nil, fmt.Errorf("method %s not found", method)
	}

	data := m.MustEncodeIDAndInput(args...)

	// Call function
	msg := &web3.CallMsg{
		To:   &c.addr,
		Data: data,
	}
	if c.from != nil {
		msg.From = *c.from
	}

	rawStr, err := c.Provider.Eth().Call(msg, block)
	if err != nil {
		return nil, err
	}

	// Decode output
	raw, err := hex.DecodeString(rawStr[2:])
	if err != nil {
		return nil, err
	}
	if len(raw) == 0 {
		return nil, fmt.Errorf("empty response")
	}
	respInterface, err := abi.Decode(m.Outputs, raw)
	if err != nil {
		return nil, err
	}

	resp := respInterface.(map[string]interface{})
	return resp, nil
}

// Txn creates a new transaction object
func (c *Contract) Txn(method string, args ...interface{}) *Txn {
	m, ok := c.Abi.Methods[method]
	if !ok {
		panic(fmt.Errorf("method %s not found", method))
	}
	data := m.MustEncodeIDAndInput(args...)

	return &Txn{
		from:     *c.from,
		to:       &c.addr,
		provider: c.Provider,
		Data:     data,
	}
}

// Txn is a transaction object
type Txn struct {
	provider *jsonrpc.Client

	from     web3.Address
	to       *web3.Address
	value    *big.Int
	nonce    uint64
	gasLimit uint64
	gasPrice uint64
	Data     []byte
	hash     web3.Hash
}

func (t *Txn) isContractDeployment() bool {
	return t.to == nil
}

// SetValue sets the value for the txn
func (t *Txn) SetValue(v *big.Int) *Txn {
	t.value = new(big.Int).Set(v)
	return t
}

func (t *Txn) SetNonce(nonce uint64) *Txn {
	t.nonce = nonce
	return t
}

// EstimateGas estimates the gas for the call
func (t *Txn) EstimateGas() (uint64, error) {
	if t.isContractDeployment() {
		return t.provider.Eth().EstimateGasContract(t.Data)
	}

	msg := &web3.CallMsg{
		From:  t.from,
		To:    t.to,
		Data:  t.Data,
		Value: t.value,
	}
	return t.provider.Eth().EstimateGas(msg)
}

// both Do and Wait functions
func (t *Txn) DoAndWait() (*web3.Receipt, error) {
	if err := t.Do(); err != nil {
		return nil, err
	}
	return t.Wait()
}

func (t *Txn) MustToTransaction() *web3.Transaction {
	tx, err := t.ToTransaction()
	utils.Ensure(err)

	return tx
}

func (t *Txn) Sign(signer *Signer) *SignedTx {
	t.from = signer.Address()
	if signer.Submit == false {
		t.nonce = signer.Nonce
		signer.Nonce += 1
		t.SetGasLimit(5000000)
	}
	tx := t.MustToTransaction()
	tx = signer.SignTx(tx)

	return &SignedTx{tx}
}

type SignedTx struct {
	*web3.Transaction
}

func (self *SignedTx) Execute(signer *Signer) (*web3.ExecutionResult, *web3.Receipt) {
	return signer.ExecuteTxn(self.Transaction)
}

func (self *SignedTx) Execute2(signer *Signer) *web3.Receipt {
	result, receipt := signer.ExecuteTxn(self.Transaction)
	if result.Err != nil {
		panic(fmt.Errorf("execution reverted: %s", result.RevertReson))
	}
	return receipt
}

func (self *SignedTx) SendTransaction(signer *Signer) *web3.Receipt {
	if signer.Submit {
		return signer.SendTransaction(self.Transaction)
	}

	return self.Execute2(signer)
}

func (t *Txn) ToTransaction() (*web3.Transaction, error) {
	var err error
	// estimate gas price
	if t.gasPrice == 0 {
		t.gasPrice, err = t.provider.Eth().GasPrice()
		if err != nil {
			return nil, err
		}
	}
	if t.nonce == 0 {
		nonce, err := t.provider.Eth().GetNonce(t.from, web3.Pending)
		if err != nil {
			return nil, err
		}

		t.nonce = nonce
	}
	utils.Ensure(err)
	// estimate gas limit
	if t.gasLimit == 0 {
		t.gasLimit, err = t.EstimateGas()
		if err != nil {
			return nil, err
		}
		t.gasLimit = t.gasLimit*130/100 + 500000
	}

	// send transaction
	txn := &web3.Transaction{
		From:     t.from,
		Input:    t.Data,
		GasPrice: t.gasPrice,
		Gas:      t.gasLimit,
		Value:    t.value,
		Nonce:    t.nonce,
	}
	if t.to != nil {
		txn.To = t.to
	}

	return txn, nil
}

// Do sends the transaction to the network
func (t *Txn) Do() error {
	// send transaction
	txn, err := t.ToTransaction()
	if err != nil {
		return err
	}
	t.hash, err = t.provider.Eth().SendTransaction(txn)
	if err != nil {
		return err
	}
	return nil
}

// SetGasPrice sets the gas price of the transaction
func (t *Txn) SetGasPrice(gasPrice uint64) *Txn {
	t.gasPrice = gasPrice
	return t
}

// SetGasLimit sets the gas limit of the transaction
func (t *Txn) SetGasLimit(gasLimit uint64) *Txn {
	t.gasLimit = gasLimit
	return t
}

// Wait waits till the transaction is mined
func (t *Txn) Wait() (receipt *web3.Receipt, err error) {
	if (t.hash == web3.Hash{}) {
		panic("transaction not executed")
	}

	for {
		receipt, err = t.provider.Eth().GetTransactionReceipt(t.hash)
		if err != nil {
			if err.Error() != "not found" {
				return nil, err
			}
		}
		if receipt != nil {
			break
		}
	}

	return
}
