package jsonrpc

import (
	"encoding/json"
	"fmt"

	"github.com/laizy/web3"
)

func (c *Client) SubscribePendingTx(watchAddr web3.Address, callback func(tx *web3.Transaction)) (func() error, error) {
	close, err := c.Subscribe("alchemy_filteredNewFullPendingTransactions", map[string]string{"address": watchAddr.String()}, func(b []byte) {
		var tx web3.Transaction
		err := json.Unmarshal(b, &tx)
		if err != nil {
			panic(fmt.Errorf("parse message error: %v, msg:%s", err, string(b)))
		}
		callback(&tx)
	})
	return close, err
}
