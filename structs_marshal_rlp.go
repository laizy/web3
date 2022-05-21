package web3

import (
	"fmt"
	"math/big"

	"golang.org/x/crypto/sha3"

	"github.com/umbracle/fastrlp"
)

func (t *Transaction) MarshalRLP() []byte {
	ar := fastrlp.DefaultArenaPool.Get()
	v := t.MarshalRLPWith(ar)
	data := v.MarshalTo(nil)
	fastrlp.DefaultArenaPool.Put(ar)
	return data
}

func (tx *Transaction) SignHash(chainId uint64) Hash {
	ar := fastrlp.DefaultArenaPool.Get()
	v := tx.MarshalRLPUnsignedWith(ar)
	// EIP155
	if chainId != 0 {
		v.Set(ar.NewUint(chainId))
		v.Set(ar.NewUint(0))
		v.Set(ar.NewUint(0))
	}
	hash := keccak256(v.MarshalTo(nil))
	fastrlp.DefaultArenaPool.Put(ar)
	return hash
}

func (t *Transaction) MarshalRLPUnsignedWith(arena *fastrlp.Arena) *fastrlp.Value {
	vv := arena.NewArray()

	vv.Set(arena.NewUint(t.Nonce))
	vv.Set(arena.NewUint(t.GasPrice))
	vv.Set(arena.NewUint(t.Gas))

	// Address may be empty
	if t.To != nil {
		vv.Set(arena.NewCopyBytes((*t.To)[:]))
	} else {
		vv.Set(arena.NewNull())
	}

	vv.Set(arena.NewBigInt(t.Value))
	vv.Set(arena.NewCopyBytes(t.Input))

	return vv
}

// MarshalRLPWith marshals the transaction to RLP with a specific fastrlp.Arena
func (t *Transaction) MarshalRLPWith(arena *fastrlp.Arena) *fastrlp.Value {
	vv := t.MarshalRLPUnsignedWith(arena)

	// signature values
	vv.Set(arena.NewCopyBytes(t.V))
	vv.Set(arena.NewCopyBytes(t.R))
	vv.Set(arena.NewCopyBytes(t.S))

	return vv
}

func TransactionFromRlp(data []byte) (*Transaction, error) {
	tx := &Transaction{}
	err := tx.UnmarshalRLP(data)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (t *Transaction) UnmarshalRLP(data []byte) error {
	p := &fastrlp.Parser{}
	v, err := p.Parse(data)
	if err != nil {
		return err
	}

	if v.Elems() != 9 {
		return fmt.Errorf("invalid data")
	}

	t.Nonce, err = v.Get(0).GetUint64()
	if err != nil {
		return err
	}
	t.GasPrice, err = v.Get(1).GetUint64()
	t.Gas, err = v.Get(2).GetUint64()
	to, err := v.Get(3).GetBytes(nil)
	if len(to) == 0 {
		t.To = nil
	} else {
		addr := BytesToAddress(to)
		t.To = &addr
	}
	t.Value = big.NewInt(0)
	err = v.Get(4).GetBigInt(t.Value)
	if err != nil {
		return err
	}
	t.Input, err = v.Get(5).GetBytes(nil)
	if err != nil {
		return err
	}
	t.R, err = v.Get(6).GetBytes(nil)
	if err != nil {
		return err
	}
	t.S, err = v.Get(7).GetBytes(nil)
	if err != nil {
		return err
	}
	t.V, err = v.Get(8).GetBytes(nil)
	if err != nil {
		return err
	}

	return nil
}

// to break circle deps with crypto
func keccak256(b []byte) (h Hash) {
	d := sha3.NewLegacyKeccak256()
	d.Write(b)
	d.Sum(h[:0])
	return
}
