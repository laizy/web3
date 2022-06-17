package jsonrpc

import (
	"github.com/laizy/web3/utils/l2"
)

// L2 is the l2 client namespace
type L2 struct {
	c *Client
}

// L2 returns the reference to the l2 namespace
func (c *Client) L2() *L2 {
	return c.endpoints.l
}

func (l *L2) GetPendingTxBatches() (*l2.RollupInputBatches, error) {
	var out l2.RollupInputBatches
	err := l.c.Call("l2_getPendingTxBatches", &out)
	return &out, err
}
