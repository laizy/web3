package tests

import (
	"fmt"
	"testing"

	"github.com/laizy/web3"
	"github.com/laizy/web3/jsonrpc"
	registry2 "github.com/laizy/web3/registry"
	web32 "github.com/laizy/web3/utils"
)

func TestEventRegistry_DumpLog(t *testing.T) {
	registry := &registry2.EventRegistry{}
	registry.RegisterPresetMainnet()

	web3.RegisterParser(registry)

	client, err := jsonrpc.NewClient("https://mainnet.infura.io/v3/99650ccb5bd14cf1884829c028826d16")
	web32.Ensure(err)
	receipt, err := client.Eth().GetTransactionReceipt(web3.HexToHash("0x5da4e1d62fab2d5182f0fb301c06d2bfd809b54e631244fbfd0a45fecf81ceb1"))
	web32.Ensure(err)

	fmt.Println(web32.JsonString(receipt))
}
