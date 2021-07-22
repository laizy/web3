package registry

import (
	"fmt"
	"testing"

	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/jsonrpc"
)

func TestEventRegistry_DumpLog(t *testing.T) {
	registry := &EventRegistry{}
	registry.RegisterPresetMainnet()

	client, err := jsonrpc.NewClient("https://mainnet.infura.io/v3/99650ccb5bd14cf1884829c028826d16")
	Ensure(err)
	receipt, err := client.Eth().GetTransactionReceipt(ethgo.HexToHash("0x5da4e1d62fab2d5182f0fb301c06d2bfd809b54e631244fbfd0a45fecf81ceb1"))
	Ensure(err)

	for _, log := range receipt.Logs {
		l := registry.DumpLog(log)
		fmt.Println(l)
	}
}
