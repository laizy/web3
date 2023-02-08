package contract

import (
	"testing"

	"github.com/laizy/web3"
	"github.com/laizy/web3/abi"
	"github.com/laizy/web3/jsonrpc"
	"github.com/laizy/web3/testutil"
	"github.com/stretchr/testify/assert"
)

func TestErr(t *testing.T) {
	t.Skip() //test locally
	abiStr := `[
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "_unlockTime",
          "type": "uint256"
        }
      ],
      "stateMutability": "payable",
      "type": "constructor"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "name": "CustomError",
      "type": "error"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "when",
          "type": "uint256"
        }
      ],
      "name": "Withdrawal",
      "type": "event"
    },
    {
      "inputs": [],
      "name": "expectErr",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "expectRevert",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "owner",
      "outputs": [
        {
          "internalType": "address payable",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "unlockTime",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "withdraw",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    }
  ]`

	s := testutil.NewTestServer(t, nil)
	defer s.Close()
	_abi := abi.MustNewABI(abiStr)
	c, err := jsonrpc.NewClient("http://localhost:8545")
	assert.NoError(t, err)
	cc := NewContract(web3.HexToAddress("0x5fbdb2315678afecb367f032d93f642f64180aa3"), _abi, c)
	_, err = cc.Call("expectErr", web3.Latest)
	assert.Equal(t, err.Error(), "{\"code\":-32603,\"message\":\"Error: VM Exception while processing transaction: reverted with custom error 'CustomError(\\\"0x0000000000000000000000000000000000000001\\\", \\\"0x0000000000000000000000000000000000000002\\\")'\",\"data\":{\"data\":\"0x0ebdea3800000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002\",\"message\":\"Error: VM Exception while processing transaction: reverted with custom error 'CustomError(\\\"0x0000000000000000000000000000000000000001\\\", \\\"0x0000000000000000000000000000000000000002\\\")'\"},\"decoded_message\":\"CustomError(0x0000000000000000000000000000000000000001,0x0000000000000000000000000000000000000002)\"}")
}
