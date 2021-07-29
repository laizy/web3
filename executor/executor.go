package executor

import (
	"os"

	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/evm"
	"github.com/umbracle/go-web3/evm/params"
	"github.com/umbracle/go-web3/evm/storage"
	"github.com/umbracle/go-web3/evm/storage/overlaydb"
	"github.com/umbracle/go-web3/executor/remotedb"
	"github.com/umbracle/go-web3/jsonrpc"
)

type Executor struct {
	db        *remotedb.RemoteDB
	overlayDB *overlaydb.OverlayDB
	cacheDB   *storage.CacheDB
	chainID   uint64
	Trace     bool
}

func NewExecutor(client *jsonrpc.Client) *Executor {
	remote := remotedb.NewRemoteDB(client)
	overlay := overlaydb.NewOverlayDB(remote)
	cacheDB := storage.NewCacheDB(overlay)
	return &Executor{
		db:        remote,
		overlayDB: overlay,
		cacheDB:   cacheDB,
		chainID:   1234,
	}
}

func (self *Executor) ResetOverlay() {
	self.overlayDB = overlaydb.NewOverlayDB(self.db)
	self.cacheDB = storage.NewCacheDB(self.overlayDB)
}

type Eip155Context struct {
	BlockHash web3.Hash
	TxIndex   uint64
	Height    uint64
	Timestamp uint64
	Coinbase  web3.Address
}

func (self *Executor) ExecuteTransaction(tx *web3.Transaction, ctx Eip155Context) (*web3.ExecutionResult, *web3.Receipt, error) {
	usedGas := uint64(0)
	config := params.GetChainConfig(self.chainID)
	statedb := storage.NewStateDB(self.cacheDB, tx.Hash, ctx.BlockHash)
	evmConf := evm.Config{}
	if self.Trace {
		evmConf.Debug = true
		evmConf.Tracer = evm.NewJSONLogger(nil, os.Stdout)
	}
	result, receipt, err := ApplyTransaction(config, self.db, statedb, ctx.Height, ctx.Timestamp, tx, &usedGas,
		ctx.Coinbase, evmConf, false)

	if err != nil {
		return nil, nil, err
	}
	if err = statedb.DbErr(); err != nil {
		return nil, nil, err
	}
	receipt.TransactionIndex = ctx.TxIndex

	return result, receipt, nil
}
