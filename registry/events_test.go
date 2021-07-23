package registry

import (
	"fmt"
	"testing"

	web32 "github.com/umbracle/go-web3/utils"

	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/jsonrpc"
)

func TestEventRegistry_DumpLog(t *testing.T) {
	registry := &EventRegistry{}
	registry.RegisterPresetMainnet()

	web3.RegisterParser(registry)

	client, err := jsonrpc.NewClient("https://mainnet.infura.io/v3/99650ccb5bd14cf1884829c028826d16")
	web32.Ensure(err)
	receipt, err := client.Eth().GetTransactionReceipt(web3.HexToHash("0x5da4e1d62fab2d5182f0fb301c06d2bfd809b54e631244fbfd0a45fecf81ceb1"))
	web32.Ensure(err)

	fmt.Println(web32.JsonString(receipt))
}
