package jsonrpc

// L2 is the l2 client namespace
type L2 struct {
	c *Client
}

// L2 returns the reference to the l2 namespace
func (c *Client) L2() *L2 {
	return c.endpoints.l
}

//tx batch data is already encoded as params of AppendBatch in RollupInputChain.sol, just add a func selector beyond it
//to invoke the AppendBatch is fine.
func (l *L2) GetPendingTxBatches() ([]byte, error) {
	var out []byte
	err := l.c.Call("l2_getPendingTxBatches", &out)
	return out, err
}
