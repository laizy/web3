package transport

import (
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/laizy/web3"
	"github.com/laizy/web3/jsonrpc/codec"
	"github.com/valyala/fasthttp"
)

// HTTP is an http transport
type HTTP struct {
	addr   string
	client *fasthttp.Client
	nextId uint64
}

func newHTTP(addr string) *HTTP {
	return &HTTP{
		addr:   addr,
		client: &fasthttp.Client{},
	}
}

// Close implements the transport interface
func (h *HTTP) Close() error {
	return nil
}

func (c *HTTP) nextID() uint64 {
	id := atomic.AddUint64(&c.nextId, 1)
	return id
}

// Call implements the transport interface
func (h *HTTP) Call(method string, out interface{}, params ...interface{}) error {
	// Encode json-rpc request
	request := codec.Request{
		Method:  method,
		JsonRpc: "2.0",
		ID:      h.nextID(),
	}

	if len(params) > 0 {
		data, err := json.Marshal(params)
		if err != nil {
			return err
		}
		request.Params = data
	} else {
		request.Params = nil
	}
	raw, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	if web3.TraceRpc {
		fmt.Printf("http eth rpc request: %s\n", string(raw))
	}

	req.SetRequestURI(h.addr)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetBody(raw)

	if err := h.client.Do(req, res); err != nil {
		return err
	}

	// Decode json-rpc response
	var response codec.Response
	body := res.Body()
	if web3.TraceRpc {
		fmt.Printf("http eth rpc response: %s\n", string(body))
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}
	if response.Error != nil {
		return response.Error
	}

	if err := json.Unmarshal(response.Result, out); err != nil {
		return err
	}
	return nil
}
