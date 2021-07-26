// Copyright (C) 2021 The Ontology Authors

package storage

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	comm "github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/store/leveldbstore"
	"github.com/stretchr/testify/require"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/evm/storage/overlaydb"
)

type dummy struct{}

func (d dummy) SubBalance(cache *CacheDB, addr ethgo.Address, val *big.Int) error {
	return nil
}
func (d dummy) AddBalance(cache *CacheDB, addr ethgo.Address, val *big.Int) error {
	return nil
}
func (d dummy) SetBalance(cache *CacheDB, addr ethgo.Address, val *big.Int) error {
	return nil
}
func (d dummy) GetBalance(cache *CacheDB, addr ethgo.Address) (*big.Int, error) {
	return big.NewInt(0), nil
}

var _ BalanceHandle = dummy{}

func TestEtherAccount(t *testing.T) {
	a := require.New(t)

	memback := leveldbstore.NewMemLevelDBStore()

	overlay := overlaydb.NewOverlayDB(memback)
	cache := NewCacheDB(overlay)

	// don't consider ong yet
	sd := NewStateDB(cache, ethgo.Hash{}, ethgo.Hash{})
	a.NotNil(sd, "fail")

	h := crypto.Keccak256Hash([]byte("hello"))

	ea := &EthAccount{
		Nonce:    1023,
		CodeHash: ethgo.Hash(h),
	}
	a.False(ea.IsEmpty(), "expect not empty")

	sink := comm.NewZeroCopySink(nil)
	ea.Serialization(sink)

	clone := &EthAccount{}
	source := comm.NewZeroCopySource(sink.Bytes())
	err := clone.Deserialization(source)
	a.Nil(err, "fail")
	a.Equal(clone, ea, "fail")

	pri, err := crypto.GenerateKey()
	a.Nil(err, "fail")
	ethAddr := crypto.PubkeyToAddress(pri.PublicKey)

	sd.cacheDB.PutEthAccount(ethgo.Address(ethAddr), *ea)

	getea := sd.getEthAccount(ethgo.Address(ethAddr))
	a.Equal(getea, *clone, "fail")

	a.Equal(sd.GetNonce(ethgo.Address(ethAddr)), ea.Nonce, "fail")
	sd.SetNonce(ethgo.Address(ethAddr), 1024)
	a.Equal(sd.GetNonce(ethgo.Address(ethAddr)), ea.Nonce+1, "fail")
	// don't effect code hash
	a.Equal(sd.getEthAccount(ethgo.Address(ethAddr)).CodeHash, ea.CodeHash, "fail")

	sd.SetCode(ethgo.Address(ethAddr), []byte("hello again"))
	a.Equal(sd.GetCodeHash(ethgo.Address(ethAddr)), crypto.Keccak256Hash([]byte("hello again")), "fail")
	a.Equal(sd.GetCode(ethgo.Address(ethAddr)), []byte("hello again"), "fail")

	a.False(sd.HasSuicided(ethgo.Address(ethAddr)), "fail")
	ret := sd.Suicide(ethgo.Address(ethAddr))
	a.True(ret, "fail")
	a.True(sd.HasSuicided(ethgo.Address(ethAddr)), "fail")

	// nonexist account get ==> default value
	pri2, _ := crypto.GenerateKey()
	anotherAddr := crypto.PubkeyToAddress(pri2.PublicKey)

	nonce := sd.GetNonce(ethgo.Address(anotherAddr))
	a.Equal(nonce, uint64(0), "fail")
	hash := sd.GetCodeHash(ethgo.Address(anotherAddr))
	a.Equal(hash, ethgo.Hash{}, "fail")

	sd.SetNonce(ethgo.Address(anotherAddr), 1)
	a.Equal(sd.GetNonce(ethgo.Address(anotherAddr)), uint64(1), "fail")
}
