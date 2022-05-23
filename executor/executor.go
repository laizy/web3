package executor

import (
	"os"

	"github.com/laizy/web3"
	"github.com/laizy/web3/evm"
	"github.com/laizy/web3/evm/params"
	"github.com/laizy/web3/evm/storage"
	"github.com/laizy/web3/evm/storage/overlaydb"
	"github.com/laizy/web3/evm/storage/schema"
)

type Executor struct {
	db        schema.ChainDB
	OverlayDB *overlaydb.OverlayDB
	ChainID   uint64
	Trace     bool
}

func NewExecutor(db schema.ChainDB, chainID uint64) *Executor {
	overlay := overlaydb.NewOverlayDB(db)
	//remote.Trace = true
	return &Executor{
		db:        db,
		OverlayDB: overlay,
		ChainID:   chainID,
	}
}

func (self *Executor) ResetOverlay() {
	self.OverlayDB = overlaydb.NewOverlayDB(self.db)
}

type Eip155Context struct {
	BlockHash web3.Hash
	TxIndex   uint64
	Height    uint64
	Timestamp uint64
	Coinbase  web3.Address
}

func (self *Executor) Call(msg Message, ctx Eip155Context) (*web3.ExecutionResult, *web3.Receipt, error) {
	usedGas := uint64(0)
	config := params.GetChainConfig(self.ChainID)
	statedb := storage.NewStateDB(storage.NewCacheDB(self.OverlayDB))
	evmConf := evm.Config{}
	if self.Trace {
		evmConf.Debug = true
		evmConf.Tracer = evm.NewJSONLogger(nil, os.Stdout)
	}
	result, receipt, err := ApplyMessage(config, self.db, statedb, msg, ctx, &usedGas, evmConf, false)

	if err != nil {
		return nil, nil, err
	}
	if err = statedb.DbErr(); err != nil {
		return nil, nil, err
	}
	receipt.TransactionIndex = ctx.TxIndex

	return result, receipt, nil

}

func (self *Executor) ExecuteTransaction(tx *web3.Transaction, ctx Eip155Context) (*web3.ExecutionResult, *web3.Receipt, error) {
	usedGas := uint64(0)
	config := params.GetChainConfig(self.ChainID)
	cacheDB := storage.NewCacheDB(self.OverlayDB)
	statedb := storage.NewStateDB(cacheDB)
	evmConf := evm.Config{}
	if self.Trace {
		evmConf.Debug = true
		evmConf.Tracer = evm.NewJSONLogger(nil, os.Stdout)
	}
	result, receipt, err := ApplyTransaction(config, self.db, statedb, tx, ctx, &usedGas,
		evmConf, false)

	if err != nil {
		return nil, nil, err
	}
	if err = statedb.DbErr(); err != nil {
		return nil, nil, err
	}
	receipt.TransactionIndex = ctx.TxIndex

	return result, receipt, nil
}
