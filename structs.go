package web3

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/laizy/web3/utils"
)

// Lengths of hashes and addresses in bytes.
const (
	// HashLength is the expected length of the hash
	HashLength = 32
	// AddressLength is the expected length of the address
	AddressLength = 20
)

// Address is an Ethereum address
type Address [20]byte

// HexToAddress converts an hex string value to an address object
func HexToAddress(str string) Address {
	a := Address{}
	err := a.UnmarshalText([]byte(str))
	utils.Ensure(err)
	return a
}

// UnmarshalText implements the unmarshal interface
func (a *Address) UnmarshalText(b []byte) error {
	return unmarshalTextByte(a[:], b, 20)
}

func (a *Address) IsZero() bool {
	var zero Address
	return *a == zero
}

// MarshalText implements the marshal interface
func (a Address) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

func (a Address) String() string {
	return "0x" + hex.EncodeToString(a[:])
}

func (a Address) Bytes() []byte {
	return a[:]
}

// Hash is an Ethereum hash
type Hash [32]byte

// HexToHash converts an hex string value to a hash object
func HexToHash(str string) Hash {
	h := Hash{}
	h.UnmarshalText([]byte(str))
	return h
}

// UnmarshalText implements the unmarshal interface
func (h *Hash) UnmarshalText(b []byte) error {
	return unmarshalTextByte(h[:], b, 32)
}

// MarshalText implements the marshal interface
func (h Hash) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

func (h Hash) String() string {
	return "0x" + hex.EncodeToString(h[:])
}

func (h Hash) Bytes() []byte {
	return h[:]
}

type Block struct {
	Number             uint64
	Hash               Hash
	ParentHash         Hash
	Sha3Uncles         Hash
	TransactionsRoot   Hash
	StateRoot          Hash
	ReceiptsRoot       Hash
	Miner              Address
	Difficulty         *big.Int
	ExtraData          []byte
	GasLimit           uint64
	GasUsed            uint64
	Timestamp          uint64
	Transactions       []*Transaction
	TransactionsHashes []Hash
	Uncles             []Hash
}

type Transaction struct {
	Hash        Hash
	From        Address
	To          *Address
	Input       []byte
	GasPrice    uint64
	Gas         uint64
	Value       *big.Int
	Nonce       uint64
	V           []byte
	R           []byte
	S           []byte
	BlockHash   Hash
	BlockNumber uint64
	TxnIndex    uint64
}

func (t *Transaction) ToCallMsg() *CallMsg {
	return &CallMsg{
		From:     t.From,
		To:       t.To,
		Data:     t.Input,
		Value:    t.Value,
		GasPrice: t.GasPrice,
	}
}

type CallMsg struct {
	From     Address
	To       *Address
	Data     []byte
	GasPrice uint64
	Value    *big.Int
}

type LogFilter struct {
	Address   []Address
	Topics    []*Hash
	BlockHash *Hash
	From      *BlockNumber
	To        *BlockNumber
}

func (l *LogFilter) SetFromUint64(num uint64) {
	b := BlockNumber(num)
	l.From = &b
}

func (l *LogFilter) SetToUint64(num uint64) {
	b := BlockNumber(num)
	l.To = &b
}

func (l *LogFilter) SetTo(b BlockNumber) {
	l.To = &b
}

type Receipt struct {
	TransactionHash   Hash
	TransactionIndex  uint64
	ContractAddress   Address
	BlockHash         Hash
	From              Address
	BlockNumber       uint64
	GasUsed           uint64
	CumulativeGasUsed uint64
	LogsBloom         []byte
	Logs              []*Log
}

func (self *Receipt) AddStorageLogs(logs []*StorageLog) {
	for ind, log := range logs {
		l := &Log{
			Removed:          false,
			LogIndex:         uint64(ind),
			TransactionIndex: self.TransactionIndex,
			TransactionHash:  self.TransactionHash,
			BlockHash:        self.BlockHash,
			BlockNumber:      self.BlockNumber,
			Address:          log.Address,
			Topics:           log.Topics,
			Data:             log.Data,
		}
		l.ParseEvent()
		self.Logs = append(self.Logs, l)
	}
}

type Log struct {
	Removed          bool
	LogIndex         uint64
	TransactionIndex uint64
	TransactionHash  Hash
	BlockHash        Hash
	BlockNumber      uint64
	Address          Address
	Topics           []Hash
	Data             []byte
	Event            *ParsedEvent
}

func (self *Log) ParseEvent() {
	parsed, err := GetParser().ParseLog(self)
	if err == nil {
		self.Event = parsed
	}
}

type StorageLog struct {
	Address Address
	Topics  []Hash
	Data    []byte
}

type BlockNumber int

const (
	Latest   BlockNumber = -1
	Earliest             = -2
	Pending              = -3
)

func (b BlockNumber) String() string {
	switch b {
	case Latest:
		return "latest"
	case Earliest:
		return "earliest"
	case Pending:
		return "pending"
	}
	if b < 0 {
		panic("internal. blocknumber is negative")
	}
	return fmt.Sprintf("0x%x", uint64(b))
}

func EncodeBlock(block ...BlockNumber) BlockNumber {
	if len(block) != 1 {
		return Latest
	}
	return block[0]
}

type ParsedEvent struct {
	Contract string
	Sig      string
	Values   map[string]interface{}
}

// BytesToHash sets b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

// SetBytes sets the hash to the value of b.
// If b is larger than len(h), b will be cropped from the left.
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}

func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

func CopyBytes(b []byte) (copiedBytes []byte) {
	if b == nil {
		return nil
	}
	copiedBytes = make([]byte, len(b))
	copy(copiedBytes, b)

	return
}

func Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}
