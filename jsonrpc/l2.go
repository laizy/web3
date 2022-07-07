package jsonrpc

import (
	"github.com/laizy/web3"
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

func (l *L2) GetRollupStateHash(batchIndex uint64) (web3.Hash, error) {
	var out web3.Hash
	err := l.c.Call("l2_getState", &out, batchIndex)
	return out, err
}

type InputChainInfo struct {
	PendingQueueIndex hexutil.Uint64
	TotalBatches      hexutil.Uint64
	QueueSize         hexutil.Uint64
}

type GlobalInfo struct {
	//total batch num in l1 RollupInputChain contract
	L1InputInfo InputChainInfo
	//l2 client have checked tx batch num
	L2CheckedBatchNum hexutil.Uint64
	//the total block num l2 already checked,start from 1, because genesis block do not need to check
	L2CheckedBlockNum hexutil.Uint64
	//l2 client head block num
	L2HeadBlockNumber   hexutil.Uint64
	L1SyncedBlockNumber hexutil.Uint64
	L1SyncedTimestamp   *hexutil.Uint64
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
	err := l.c.Call("l2_inputBatchNumber", &out)
	return out, err
}

func (l *L2) StateBatchNumber() (uint64, error) {
	out := uint64(0)
	err := l.c.Call("l2_stateBatchNumber", &out)
	return out, err
}

type RPCBatch struct {
	Sequencer    web3.Address        `json:"sequencer"`
	BatchNumber  hexutil.Uint64      `json:"batchNumber"`
	BatchHash    web3.Hash           `json:"batchHash"`
	QueueStart   hexutil.Uint64      `json:"queueStart"`
	QueueNum     hexutil.Uint64      `json:"queueNum"`
	Transactions []*web3.Transaction `json:"transactions"`
}

func (l *L2) GetBatch(batchNumber uint64, useDetail bool) (*RPCBatch, error) {
	out := RPCBatch{}
	err := l.c.Call("l2_getBatch", &out, batchNumber, useDetail)
	return &out, err
}

type RPCEnqueuedTx struct {
	QueueIndex hexutil.Uint64 `json:"queueIndex"`
	From       web3.Address   `json:"from"`
	To         web3.Address   `json:"to"`
	RlpTx      hexutil.Bytes  `json:"rlpTx"`
	Timestamp  hexutil.Uint64 `json:"timestamp"`
}

func (l *L2) GetEnqueuedTxs(queueStart, queueNum uint64) ([]*RPCEnqueuedTx, error) {
	out := make([]*RPCEnqueuedTx, 0)
	err := l.c.Call("l2_getEnqueuedTxs", &out, queueStart, queueNum)
	return out, err
}

type RPCBatchState struct {
	Index     hexutil.Uint64
	Proposer  web3.Address
	Timestamp hexutil.Uint64
	BlockHash web3.Hash
}

func (l *L2) GetBatchState(batchNumber uint64) (*RPCBatchState, error) {
	out := RPCBatchState{}
	err := l.c.Call("l2_getBatchState", &out, batchNumber)
	return &out, err
}

type BlockNumberOrHash struct {
	BlockNumber      *uint64    `json:"blockNumber,omitempty"`
	BlockHash        *web3.Hash `json:"blockHash,omitempty"`
	RequireCanonical bool       `json:"requireCanonical,omitempty"`
}

func (l *L2) GetReadStorageProof(blockNumOrHash *BlockNumberOrHash) ([]string, error) {
	result := make([]string, 0)
	err := l.c.Call("debug_getReadStorageProofAtBlock", &result, blockNumOrHash)
	return result, err
}

type L1RelayMsgParams struct {
	Target       web3.Address   `json:"target"`
	Sender       web3.Address   `json:"sender"`
	Message      hexutil.Bytes  `json:"message"`
	MessageIndex hexutil.Uint64 `json:"messageIndex"`
	RLPHeader    hexutil.Bytes  `json:"rlpHeader"`
	StateInfo    *RPCBatchState `json:"stateInfo"`
	Proof        []web3.Hash    `json:"proof"`
}

type L2RelayMsgParams struct {
	Target       web3.Address   `json:"target"`
	Sender       web3.Address   `json:"sender"`
	Message      hexutil.Bytes  `json:"message"`
	MessageIndex hexutil.Uint64 `json:"messageIndex"`
	MMRSize      hexutil.Uint64 `json:"mmrSize"`
	Proof        []web3.Hash    `json:"proof"`
}

func (l *L2) GetL2MMRProof(msgIndex, size uint64) ([]web3.Hash, error) {
	result := make([]web3.Hash, 0)
	err := l.c.Call("l2_getL2MMRProof", &result, msgIndex, size)
	return result, err
}

func (l *L2) GetL1RelayMsgParams(msgIndex uint64) (*L1RelayMsgParams, error) {
	result := &L1RelayMsgParams{}
	err := l.c.Call("l2_getL1RelayMsgParams", &result, msgIndex)
	return result, err
}

func (l *L2) GetL2RelayMsgParams(msgIndex uint64) (*L2RelayMsgParams, error) {
	result := &L2RelayMsgParams{}
	err := l.c.Call("l2_getL1RelayMsgParams", &result, msgIndex)
	return result, err
}
