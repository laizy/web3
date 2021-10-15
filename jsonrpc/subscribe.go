package jsonrpc

import (
	"fmt"

	"github.com/laizy/web3"

	"github.com/laizy/web3/jsonrpc/transport"
)

// SubscriptionEnabled returns true if the subscription endpoints are enabled
func (c *Client) SubscriptionEnabled() bool {
	_, ok := c.transport.(transport.PubSubTransport)
	return ok
}

// Subscribe starts a new subscription
func (c *Client) Subscribe(method string, param interface{}, callback func(b []byte)) (func() error, error) {
	pub, ok := c.transport.(transport.PubSubTransport)
	if !ok {
		return nil, fmt.Errorf("Transport does not support the subscribe method")
	}
	close, err := pub.Subscribe(method, param, callback)
	return close, err
}

/*
Emits an event any time a new header is added to the chain, including during a chain reorganization.
When a chain reorganization occurs, this subscription will emit an event containing all new headers for the new chain. In particular, this means that you may see multiple headers emitted with the same height, and when this happens the later header should be taken as the correct one after a reorganization.
*/
func (c *Client) SubscribeNewHeads(callback func(b *web3.Block)) (func() error, error) {
	return c.Subscribe("newHeads", nil, func(b []byte) {
		var block web3.Block
		if err := block.UnmarshalJSON(b); err != nil {
			panic(fmt.Errorf("parse head msg error: %v, msg:%s", err, string(b)))
		}
		callback(&block)
	})
}

/*
Emits logs which are part of newly added blocks that match specified filter criteria.

When a chain reorganization occurs, logs which are part of blocks on the old chain will be emitted again with the property removed set to true. Further, logs which are part of the blocks on the new chain are emitted, meaning that it is possible to see logs for the same transaction multiple times in the case of a reorganization.

Parameters is an object with the following fields:
adddress (optional): either a string representing an address or an array of such strings.Only logs created from one of these addresses will be emitted.
topics: an array of topic specifiers. Each topic specifier is either null, a string representing a topic, or an array of strings.Each position in the array which is not null restricts the emitted logs to only those who have one of the given topics in that position.
Some examples of topic specifications:
[]: Any topics allowed.
[A]: A in first position (and anything after).
[null, B]: Anything in first position and B in second position (and anything after).
[A, B]: A in first position and B in second position (and anything after).
[[A, B], [A, B]]: (A or B) in first position and (A or B) in second position (and anything after).
*/
func (c *Client) SubscribeLogs(callback func(log *web3.Log), addresses []web3.Address, topics ...[][]web3.Hash) (func() error, error) {
	param := make(map[string]interface{})
	if len(addresses) > 0 {
		param["address"] = addresses
	}
	if len(topics) == 1 {
		param["topics"] = topics[0]
	}
	return c.Subscribe("logs", param, func(b []byte) {
		var log web3.Log
		if err := log.UnmarshalJSON(b); err != nil {
			panic(fmt.Errorf("parse head msg error: %v, msg:%s", err, string(b)))
		}
		callback(&log)
	})
}
