package codec

import (
	"encoding/json"
	"fmt"

	"github.com/laizy/web3/registry"
	"github.com/laizy/web3/utils/common/hexutil"
)

// Request is a jsonrpc request
type Request struct {
	JsonRpc string          `json:"jsonrpc"`
	ID      uint64          `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// Response is a jsonrpc response
type Response struct {
	ID     uint64          `json:"id"`
	Result json.RawMessage `json:"result"`
	Error  *ErrorObject    `json:"error,omitempty"`
}

// ErrorObject is a jsonrpc error
type ErrorObject struct {
	Code           int    `json:"code"`
	Message        string `json:"message"`
	Data           string `json:"data,omitempty"`
	DecodedMessage string `json:"decoded_message,omitempty"`
}

// Subscription is a jsonrpc subscription
type Subscription struct {
	ID     string          `json:"subscription"`
	Result json.RawMessage `json:"result"`
}

// Error implements error interface
func (e *ErrorObject) Error() string {
	if len(e.Data) != 0 {
		info, err := registry.ErrInstance().ParseError(hexutil.MustDecode(e.Data))
		if err == nil {
			e.DecodedMessage = info
		}
	}
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("jsonrpc.internal marshal error: %v", err)
	}
	return string(data)
}
