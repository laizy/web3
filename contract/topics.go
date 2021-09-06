package contract

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/laizy/web3"
	"github.com/laizy/web3/crypto"
	"github.com/laizy/web3/utils/u256"
)

// MakeTopics converts a filter query argument list into a filter topic set.
func MakeTopics(query ...[]interface{}) ([][]web3.Hash, error) {
	topics := make([][]web3.Hash, len(query))
	for i, filter := range query {
		for _, rule := range filter {
			var topic web3.Hash

			// Try to generate the topic based on simple types
			switch rule := rule.(type) {
			case web3.Hash:
				topic = rule
			case web3.Address:
				topic = rule.ToHash()
			case bool:
				if rule {
					topic[web3.HashLength-1] = 1
				}
			case int8:
				topic = genIntType(int64(rule), 1)
			case int16:
				topic = genIntType(int64(rule), 2)
			case int32:
				topic = genIntType(int64(rule), 4)
			case int64:
				topic = genIntType(rule, 8)
			case uint8, uint16, uint32, uint64, *big.Int, u256.Int, *u256.Int:
				topic = u256.New(rule).Bytes32()
			case string:
				hash := crypto.Keccak256Hash([]byte(rule))
				copy(topic[:], hash[:])
			case []byte:
				hash := crypto.Keccak256Hash(rule)
				copy(topic[:], hash[:])

			default:
				// todo(rjl493456442) according solidity documentation, indexed event
				// parameters that are not value types i.e. arrays and structs are not
				// stored directly but instead a keccak256-hash of an encoding is stored.
				//
				// We only convert stringS and bytes to hash, still need to deal with
				// array(both fixed-size and dynamic-size) and struct.

				// Attempt to generate the topic from funky types
				val := reflect.ValueOf(rule)
				switch {
				// static byte array
				case val.Kind() == reflect.Array && reflect.TypeOf(rule).Elem().Kind() == reflect.Uint8:
					reflect.Copy(reflect.ValueOf(topic[:val.Len()]), val)
				default:
					return nil, fmt.Errorf("unsupported indexed type: %T", rule)
				}
			}
			topics[i] = append(topics[i], topic)
		}
	}
	return topics, nil
}

func genIntType(rule int64, size uint) [32]byte {
	var topic [web3.HashLength]byte
	if rule < 0 {
		// if a rule is negative, we need to put it into two's complement.
		// extended to common.HashLength bytes.
		topic = [web3.HashLength]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
	}
	for i := uint(0); i < size; i++ {
		topic[web3.HashLength-i-1] = byte(rule >> (i * 8))
	}
	return topic
}
