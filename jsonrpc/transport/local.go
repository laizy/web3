package transport

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sync/atomic"

	"github.com/laizy/web3"
	"github.com/laizy/web3/evm/storage"
	"github.com/laizy/web3/evm/storage/schema"
	"github.com/laizy/web3/executor"
	"github.com/laizy/web3/utils"
	"github.com/laizy/web3/utils/common/hexutil"
)

type Local struct {
	db          schema.ChainDB
	exec        *executor.Executor
	BlockNumber uint64
	BlockHashes map[uint64]web3.Hash
	Receipts    map[web3.Hash]*web3.Receipt
	nextId      uint64
}

func NewLocal(db schema.ChainDB) *Local {
	return &Local{
		db:          db,
		exec:        executor.NewExecutor(db),
		BlockNumber: 0,
		BlockHashes: make(map[uint64]web3.Hash),
		Receipts:    make(map[web3.Hash]*web3.Receipt),
		nextId:      0,
	}
}

// Close implements the transport interface
func (self *Local) Close() error {
	return nil
}

func (self *Local) nextID() uint64 {
	id := atomic.AddUint64(&self.nextId, 1)
	return id
}

// Call implements the transport interface
func (self *Local) Call(method string, out interface{}, params ...interface{}) error {
	var result []byte
	switch method {
	case "eth_getCode":
		addr := params[0].(web3.Address)
		cacheDB := storage.NewCacheDB(self.exec.OverlayDB)
		val, err := cacheDB.GetEthAccount(addr)
		if err != nil {
			return err
		}
		result = utils.JsonBytes(val.Code)
	case "eth_blockNumber":
		result = utils.JsonBytes(hexutil.Uint64(self.BlockNumber))
	case "eth_call":
		msg := params[0].(*web3.CallMsg)
		// blockNum := params[0].(string)
		res, err := self.CallEvm(msg)
		if err != nil {
			return err
		}
		result = utils.JsonBytes(res.ReturnData)
	case "eth_estimateGas":
		msg := params[0].(*web3.CallMsg)
		res, err := self.CallEvm(msg)
		if err != nil {
			return err
		}
		result = utils.JsonBytes(hexutil.Uint64(res.UsedGas))
	case "eth_sendTransaction":
		txn := params[0].(*web3.Transaction)
		_, receipt, err := self.exec.ExecuteTransaction(txn, executor.Eip155Context{
			BlockHash: web3.Hash{},
			Height:    self.BlockNumber,
		})
		if err != nil {
			return err
		}
		self.Receipts[txn.Hash()] = receipt
		self.BlockNumber += 1
		result = []byte(txn.Hash().String())
	default:
		panic(fmt.Errorf("unimplemented method: %s", method))
	}

	return json.Unmarshal(result, out)
}

func (self *Local) CallEvm(msg *web3.CallMsg) (*web3.ExecutionResult, error) {
	res, _, err := self.exec.Call(CallMsg{msg}, executor.Eip155Context{
		BlockHash: web3.Hash{},
		Height:    self.BlockNumber,
	})
	if err != nil {
		return nil, err
	}
	if res.Failed() {
		return nil, fmt.Errorf(res.RevertReason)
	}

	return res, nil
}

type CallMsg struct {
	msg *web3.CallMsg
}

func (self CallMsg) From() web3.Address {
	return self.msg.From
}

func (self CallMsg) To() *web3.Address {
	return self.msg.To
}

func (self CallMsg) GasPrice() *big.Int {
	return big.NewInt(0).SetUint64(self.msg.GasPrice)
}

func (self CallMsg) Gas() uint64 {
	return 20000000
}

func (self CallMsg) Value() *big.Int {
	return self.msg.Value
}

func (self CallMsg) Nonce() uint64 {
	return 0
}

func (self CallMsg) CheckNonce() bool {
	return false
}

func (self CallMsg) Data() []byte {
	return self.msg.Data
}

func (self CallMsg) Hash() web3.Hash {
	return web3.Hash{}
}
