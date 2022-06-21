package jsonrpc

import (
	"github.com/laizy/web3"
	"github.com/laizy/web3/utils/common"
	"github.com/laizy/web3/utils/common/hexutil"
)

// L2 is the l2 client namespace
type L2 struct {
	c *Client
}

// L2 returns the reference to the l2 namespace
func (c *Client) L2() *L2 {
	return c.endpoints.l
}

//tx batch data is already encoded as params of AppendBatch in RollupInputChain.sol, just add a func selector beyond it
//to invoke the AppendBatch is fine.
func (l *L2) GetPendingTxBatches() ([]byte, error) {
	var out []byte
	err := l.c.Call("l2_getPendingTxBatches", &out)
	return out, err
}

type InputChainInfo struct {
	PendingQueueIndex uint64
	TotalBatches      uint64
	QueueSize         uint64
}
type GlobalInfo struct {
	//total batch num in l1 RollupInputChain contract
	L1InputInfo InputChainInfo
	//l2 client have checked tx batch num
	L2CheckedBatchNum uint64
	//the total block num l2 already checked,start from 1, because genesis block do not need to check
	L2CheckedBlockNum uint64
	//l2 client head block num
	L2HeadBlockNumber   uint64
	L1SyncedBlockNumber uint64
	L1SyncedTimestamp   *uint64
}

//tx batch data is already encoded as params of AppendBatch in RollupInputChain.sol, just add a func selector beyond it
//to invoke the AppendBatch is fine.
func (l *L2) GlobalInfo() (*GlobalInfo, error) {
	var out GlobalInfo
	err := l.c.Call("l2_globalInfo", &out)
	return &out, err
}

func (l *L2) InputBatchNumber() (uint64, error) {
	out := uint64(0)
	err := l.c.Call("rollup_inputBatchNumber", &out)
	return out, err
}

func (l *L2) StateBatchNumber() (uint64, error) {
	out := uint64(0)
	err := l.c.Call("rollup_stateBatchNumber", &out)
	return out, err
}

type RPCBatch struct {
	Sequencer    common.Address      `json:"sequencer"`
	BatchNumber  uint64              `json:"batchNumber"`
	BatchHash    uint64              `json:"batchHash"`
	Transactions []*web3.Transaction `json:"transactions"`
}

func (l *L2) GetBatch() (*RPCBatch, error) {
	out := RPCBatch{}
	err := l.c.Call("rollup_getBatch", &out)
	return &out, err
}

type RPCBatchState struct {
	Index     hexutil.Uint64
	Proposer  common.Address
	Timestamp hexutil.Uint64
	BlockHash common.Hash
}

func (l *L2) GetBatchState() (*RPCBatchState, error) {
	out := RPCBatchState{}
	err := l.c.Call("rollup_getBatchState", &out)
	return &out, err
}
