package ens

import (
	"fmt"
	"math/big"

	"github.com/laizy/web3"
)

var (
	_ = big.NewInt
	_ = fmt.Printf
	_ = web3.HexToAddress
)

type ABIChangedEvent struct {
	Node        [32]byte
	ContentType *big.Int

	Raw *web3.Log
}

type AddrChangedEvent struct {
	Node [32]byte
	A    web3.Address

	Raw *web3.Log
}

type ContentChangedEvent struct {
	Node [32]byte
	Hash [32]byte

	Raw *web3.Log
}

type NameChangedEvent struct {
	Node [32]byte
	Name string

	Raw *web3.Log
}

type NewOwnerEvent struct {
	Node  [32]byte
	Label [32]byte
	Owner web3.Address

	Raw *web3.Log
}

type NewResolverEvent struct {
	Node     [32]byte
	Resolver web3.Address

	Raw *web3.Log
}

type NewTTLEvent struct {
	Node [32]byte
	Ttl  uint64

	Raw *web3.Log
}

type PubkeyChangedEvent struct {
	Node [32]byte
	X    [32]byte
	Y    [32]byte

	Raw *web3.Log
}

type TransferEvent struct {
	Node  [32]byte
	Owner web3.Address

	Raw *web3.Log
}
