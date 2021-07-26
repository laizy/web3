package crypto

import (
	"github.com/umbracle/fastrlp"
	"github.com/umbracle/go-web3"
	"golang.org/x/crypto/sha3"
)

func Keccak256(data ...[]byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	for _, d := range data {
		hash.Write(d)
	}
	return hash.Sum(nil)
}

func Keccak256Hash(code []byte) (result web3.Hash) {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(code)
	dst := hash.Sum(nil)
	copy(result[:], dst)
	return
}

// CreateAddress creates an ethereum address given the bytes and the nonce
func CreateAddress(b web3.Address, nonce uint64) web3.Address {
	a := &fastrlp.Arena{}
	v := a.NewArray()
	v.Set(a.NewBytes(b.Bytes()))
	v.Set(a.NewUint(nonce))
	data := v.MarshalTo(nil)
	return web3.BytesToAddress(Keccak256(data)[12:])
}

// CreateAddress2 creates an ethereum address given the address bytes, initial
// contract code hash and a salt.
func CreateAddress2(b web3.Address, salt [32]byte, inithash []byte) web3.Address {
	return web3.BytesToAddress(Keccak256([]byte{0xff}, b.Bytes(), salt[:], inithash)[12:])
}
