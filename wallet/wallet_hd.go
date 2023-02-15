package wallet

import (
	"fmt"
	"math/big"
	"strings"
)

type DerivationPath []uint32

// 0x800000
var decVal = big.NewInt(2147483648)

// DefaultDerivationPath is the default derivation path for Ethereum addresses
var DefaultDerivationPath = DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0, 0}

func parseDerivationPath(path string) (*DerivationPath, error) {
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return nil, fmt.Errorf("no derivation path")
	}

	// clean all the parts of any trim spaces
	for indx := range parts {
		parts[indx] = strings.TrimSpace(parts[indx])
	}

	// first part has to be an 'm'
	if parts[0] != "m" {
		return nil, fmt.Errorf("first has to be m")
	}

	result := DerivationPath{}
	for _, p := range parts[1:] {
		val := new(big.Int)
		if strings.HasSuffix(p, "'") {
			p = strings.TrimSuffix(p, "'")
			val.Add(val, decVal)
		}

		bigVal, ok := new(big.Int).SetString(p, 0)
		if !ok {
			return nil, fmt.Errorf("invalid path")
		}
		val.Add(val, bigVal)

		// TODO, limit to uint32
		if !val.IsUint64() {
			return nil, fmt.Errorf("bad")
		}
		result = append(result, uint32(val.Uint64()))
	}

	return &result, nil
}
