package l2

import "github.com/laizy/web3"

// format: queueNum(uint64) + queueStart(uint64) + batchNum(uint64) + batch0Time(uint64) +
// batchLeftTimeDiff([]uint32) + batchesData
// batchesData: version(0) + rlp([][]transaction)
type RollupInputBatches struct {
	QueueNum   uint64
	QueueStart uint64
	SubBatches []*SubBatch
}

type SubBatch struct {
	Timestamp uint64
	Txs       []*web3.Transaction
}
