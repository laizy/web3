package registry

import (
	"github.com/umbracle/go-web3"
)

func (self *EventRegistry) RegisterPresetMainnet() {
	wellKnowns := map[web3.Address]string{
		web3.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7"): "USDT",
		web3.HexToAddress("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"): "USDC",
		web3.HexToAddress("0xdb0f18081b505a7de20b18ac41856bcb4ba86a1a"): "pWING",
		web3.HexToAddress("0xcb46c550539ac3db72dc7af7c89b11c306c727c2"): "pONT",
		web3.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"): "WETH9",
	}

	for addr, name := range wellKnowns {
		self.RegisterContractAlias(addr, name)
	}
}
