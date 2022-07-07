package jsonrpc

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/laizy/web3/utils/codec"
	"github.com/stretchr/testify/assert"
)

var l2 *L2

func getL2Client(t *testing.T) *L2 {
	if l2 == nil {
		c, err := NewClient("http://172.168.3.73:8545")
		assert.NoError(t, err)
		l2 = c.L2()
	}
	return l2
}

func TestL2_GlobalInfo(t *testing.T) {
	info, err := getL2Client(t).GlobalInfo()
	assert.NoError(t, err)
	jsonInfo, err := json.MarshalIndent(info, "", "	")
	assert.NoError(t, err)
	t.Log(string(jsonInfo))
}

type RollupInputBatches struct {
	//BatchIndex ignored when calc hash, because its useless in l2 system
	BatchIndex uint64
	QueueNum   uint64
	QueueStart uint64
	SubBatches []*SubBatch
}

type SubBatch struct {
	Timestamp uint64
	//Txs       [][]byte doesn't decode
}

func DecodeBatch(b []byte) error {
	reader := codec.NewZeroCopyReader(b)
	self := &RollupInputBatches{}
	self.BatchIndex = reader.ReadUint64BE()
	self.QueueNum = reader.ReadUint64BE()
	self.QueueStart = reader.ReadUint64BE()
	batchNum := reader.ReadUint64BE()
	if batchNum == 0 {
		//check length
		if reader.Len() != 0 {
			return fmt.Errorf("wrong b length")
		}
		return reader.Error()
	}
	batchTime := reader.ReadUint64BE()
	batchesTime := []uint64{batchTime}
	for i := uint64(0); i < batchNum-1; i++ {
		batchTime = batchTime + uint64(reader.ReadUint32BE())
		if reader.Error() != nil {
			return reader.Error()
		}
		batchesTime = append(batchesTime, batchTime)
	}

	version := reader.ReadUint8()
	if version != 0 {
		return fmt.Errorf("unknown batch version: %d", version)
	}
	return nil
}

func TestL2_GetPendingTxBatches(t *testing.T) {
	batch, err := getL2Client(t).GetPendingTxBatches()
	assert.NoError(t, err)
	if len(batch) > 0 { // try to decode
		err := DecodeBatch(batch)
		assert.NoError(t, err)
	}
}

func TestL2_GetRollupStateHash(t *testing.T) {
	info, err := getL2Client(t).GlobalInfo()
	assert.NoError(t, err)
	stateHash, err := getL2Client(t).GetRollupStateHash(uint64(info.L1InputInfo.TotalBatches) / 2)
	assert.NoError(t, err)
	assert.False(t, stateHash.IsEmpty())
	t.Log(stateHash.String())
}

func TestL2_InputBatchNumber(t *testing.T) {
	num, err := getL2Client(t).InputBatchNumber()
	assert.NoError(t, err)
	assert.True(t, num > 0)
	t.Log(num)
}

func TestL2_StateBatchNumber(t *testing.T) {
	num, err := getL2Client(t).StateBatchNumber()
	assert.NoError(t, err)
	assert.True(t, num > 0)
	t.Log(num)
}

func TestL2_GetBatch(t *testing.T) {
	info, err := getL2Client(t).GlobalInfo()
	assert.NoError(t, err)
	batch, err := getL2Client(t).GetBatch(uint64(info.L1InputInfo.TotalBatches)/2, true)
	assert.NoError(t, err)
	jsonBatch, err := json.MarshalIndent(batch, "", "	")
	assert.NoError(t, err)
	t.Log(string(jsonBatch))
}

func TestL2_GetEnqueuedTxs(t *testing.T) {
	txs, err := getL2Client(t).GetEnqueuedTxs(100, 200)
	assert.NoError(t, err)
	assert.Equal(t, 200, len(txs))
	jsonTxs, err := json.MarshalIndent(txs, "", "	")
	assert.NoError(t, err)
	t.Log(string(jsonTxs))
}

func TestL2_GetBatchState(t *testing.T) {
	info, err := getL2Client(t).GlobalInfo()
	assert.NoError(t, err)
	state, err := getL2Client(t).GetBatchState(uint64(info.L1InputInfo.TotalBatches) / 2)
	assert.NoError(t, err)
	jsonState, err := json.MarshalIndent(state, "", "	")
	assert.NoError(t, err)
	t.Log(string(jsonState))
}

func TestL2_GetReadStorageProof(t *testing.T) {
	info, err := getL2Client(t).GlobalInfo()
	assert.NoError(t, err)
	blk := uint64(info.L2CheckedBlockNum)
	proofs, err := getL2Client(t).GetReadStorageProof(&BlockNumberOrHash{BlockNumber: &blk})
	assert.NoError(t, err)
	for _, proof := range proofs {
		t.Log(proof)
	}
}

func TestL2_GetL2MMRProof(t *testing.T) {
	proofs, err := getL2Client(t).GetL2MMRProof(0, 1)
	assert.NoError(t, err)
	for _, proof := range proofs {
		t.Log(proof.String())
	}
}

func TestL2_GetL1RelayMsgParams(t *testing.T) {
	params, err := getL2Client(t).GetL1RelayMsgParams(0)
	assert.NoError(t, err)
	jsonParams, err := json.MarshalIndent(params, "", "	")
	assert.NoError(t, err)
	t.Log(string(jsonParams))
}

func TestL2_GetL2RelayMsgParams(t *testing.T) {
	params, err := getL2Client(t).GetL2RelayMsgParams(0)
	assert.NoError(t, err)
	jsonParams, err := json.MarshalIndent(params, "", "	")
	assert.NoError(t, err)
	t.Log(string(jsonParams))
}
