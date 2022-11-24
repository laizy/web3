package jsonrpc

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/laizy/web3"
	"github.com/laizy/web3/utils/codec"
	"github.com/stretchr/testify/assert"
)

const url = "http://192.168.6.237:23333"

var l2 *L2

func getL2Client(t *testing.T) *L2 {
	if l2 == nil {
		c, err := NewClient(url)
		assert.NoError(t, err)
		l2 = c.L2()
	}
	return l2
}

func TestL2_GlobalInfo(t *testing.T) {
	info, err := getL2Client(t).GlobalInfo()
	if err != nil {
		t.Skipf("skipping since client is not available")
	}
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
	if err != nil {
		t.Skipf("skipping since client is not available")
	}
	if len(batch) > 0 { // try to decode
		err := DecodeBatch(batch)
		assert.NoError(t, err)
	}
}

func TestL2_GetRollupStateHash(t *testing.T) {
	info, err := getL2Client(t).GlobalInfo()
	if err != nil {
		t.Skipf("skipping since client is not available")
	}
	stateHash, err := getL2Client(t).GetRollupStateHash(uint64(info.L1InputInfo.TotalBatches) / 2)
	assert.NoError(t, err)
	assert.False(t, stateHash.IsEmpty())
	t.Log(stateHash.String())
}

func TestL2_InputBatchNumber(t *testing.T) {
	num, err := getL2Client(t).InputBatchNumber()
	if err != nil {
		t.Skipf("skipping since client is not available")
	}
	assert.True(t, num > 0)
	t.Log(num)
}

func TestL2_StateBatchNumber(t *testing.T) {
	num, err := getL2Client(t).StateBatchNumber()
	if err != nil {
		t.Skipf("skipping since client is not available")
	}
	assert.True(t, num > 0)
	t.Log(num)
}

func TestL2_GetBatch(t *testing.T) {
	batch, err := getL2Client(t).GetBatch(10, true)
	if err != nil {
		t.Skipf("skipping since client is not available")
	}
	jsonBatch, err := json.MarshalIndent(batch, "", "	")
	assert.NoError(t, err)
	t.Log(string(jsonBatch))
}

func TestL2_GetEnqueuedTxs(t *testing.T) {
	txs, err := getL2Client(t).GetEnqueuedTxs(100, 1)
	if err != nil {
		t.Skipf("skipping since client is not available")
	}
	assert.Equal(t, 1, len(txs))
	jsonTxs, err := json.MarshalIndent(txs, "", "	")
	assert.NoError(t, err)
	t.Log(string(jsonTxs))
}

func TestL2_GetBatchState(t *testing.T) {
	info, err := getL2Client(t).GlobalInfo()
	if err != nil {
		t.Skipf("skipping since client is not available")
	}
	state, err := getL2Client(t).GetBatchState(uint64(info.L1InputInfo.TotalBatches) / 2)
	assert.NoError(t, err)
	jsonState, err := json.MarshalIndent(state, "", "	")
	assert.NoError(t, err)
	t.Log(string(jsonState))
}

func TestL2_GetReadStorageProof(t *testing.T) {
	//totalBlock, err := getL2Client(t).c.Eth().BlockNumber()
	//if err != nil {
	//	t.Skipf("skipping since client is not available")
	//}
	// 986 0x3da, 958 0x3be
	batchIndex := uint64(1)
	proofs, err := getL2Client(t).GetReadStorageProof(nil, web3.Hash{}, batchIndex)
	if err != nil {
		t.Skipf("skipping since client is not available")
	}
	for _, proof := range proofs {
		t.Log(proof)
	}
}

func TestL2_GetL2MMRProof(t *testing.T) {
	proofs, err := getL2Client(t).GetL2MMRProof(1, 5)
	if err != nil {
		t.Skipf("skipping since client is not available")
	}
	for _, proof := range proofs {
		t.Log(proof.String())
	}
}

func TestL2_GetL1RelayMsgParams(t *testing.T) {
	params, err := getL2Client(t).GetL1RelayMsgParams(0)
	if err != nil {
		t.Skipf("skipping since client is not available")
	}
	jsonParams, err := json.MarshalIndent(params, "", "	")
	assert.NoError(t, err)
	t.Log(string(jsonParams))
}

func TestL2_GetL2RelayMsgParams(t *testing.T) {
	params, err := getL2Client(t).GetL2RelayMsgParams(0)
	if err != nil {
		t.Skipf("skipping since client is not available")
	}
	jsonParams, err := json.MarshalIndent(params, "", "	")
	assert.NoError(t, err)
	t.Log(string(jsonParams))
}
