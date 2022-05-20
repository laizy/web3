package storage

import (
	"testing"

	"github.com/laizy/web3"
	"github.com/laizy/web3/evm/storage/overlaydb"
)

func TestSnapshot(t *testing.T) {
	caccheDB := NewCacheDB(overlaydb.NewOverlayDB(NewFakeDB()))
	statedb := NewStateDB(caccheDB)
	testAddr := web3.Address{1, 1, 1, 1}
	statedb.getEthAccount(testAddr)
	statedb.CreateAccount(testAddr)
	statedb.SetState(testAddr, web3.Hash{}, web3.Hash{})
	i := statedb.Snapshot()
	statedb.SetState(testAddr, web3.Hash{1}, web3.Hash{})
	statedb.RevertToSnapshot(i)
	statedb.getEthAccount(testAddr)
	i = statedb.Snapshot()
	statedb.DiscardSnapshot(i)
}
