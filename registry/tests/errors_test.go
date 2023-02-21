package tests

import (
	"testing"

	"github.com/laizy/web3"
	"github.com/laizy/web3/abi"
	registry2 "github.com/laizy/web3/registry"
	"github.com/stretchr/testify/assert"
)

func TestErrorRegistry(t *testing.T) {
	abiStr := `[{
     "inputs": [
       {
         "internalType": "address",
         "name": "have",
         "type": "address"
       },
       {
         "internalType": "address",
         "name": "want",
         "type": "address"
       }
     ],
     "name": "OnlyCoordinatorCanFulfill",
     "type": "error"
   }]`
	registry := registry2.NewErrorRegistry()
	_abi, err := abi.NewABI(abiStr)
	assert.Nil(t, err)
	registry.RegisterFromAbi(_abi)
	data, err := _abi.Errors["OnlyCoordinatorCanFulfill"].EncodeIDAndInput(web3.Address{19: 1}, web3.Address{18: 1})
	assert.Nil(t, err)
	info, err := registry.ParseError(data)
	assert.Nil(t, err)
	t.Log(info)
}
