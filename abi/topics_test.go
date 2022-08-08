package abi

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/laizy/web3"
	"github.com/laizy/web3/testutil"
	"github.com/stretchr/testify/assert"
)

func TestTopicEncoding(t *testing.T) {
	cases := []struct {
		Type string
		Val  interface{}
	}{
		{
			Type: "bool",
			Val:  true,
		},
		{
			Type: "bool",
			Val:  false,
		},
		{
			Type: "uint64",
			Val:  uint64(20),
		},
		{
			Type: "uint256",
			Val:  big.NewInt(1000000),
		},
		{
			Type: "address",
			Val:  web3.Address{0x1},
		},
		{
			Type: "bytes32",
			Val:  [32]byte{31: 2},
		},
		{
			Type: "bytes24",
			Val:  [24]byte{0x1},
		},
		{
			Type: "bytes16",
			Val:  [16]byte{0x1},
		},
		{
			Type: "bytes8",
			Val:  [8]byte{0x1},
		},
	}

	for _, c := range cases {
		tt, err := NewType(c.Type)
		assert.NoError(t, err)

		res, err := EncodeTopic(tt, c.Val)
		assert.NoError(t, err)

		val, err := ParseTopic(tt, res)
		assert.NoError(t, err)

		assert.Equal(t, val, c.Val)
	}
}

func TestIntegrationTopics(t *testing.T) {
	s := testutil.NewTestServer(t, nil)
	defer s.Close()

	type field struct {
		typ    string
		indx   bool
		val    interface{}
		valStr string
	}

	cases := []struct {
		fields []field
	}{
		{
			fields: []field{
				{"uint32", false, uint32(1), "1"},
				{"uint8", true, uint8(10), "10"},
				{"bytes32", true, [32]byte{31: 2}, "bytes32(uint256(0x2))"},
				{"bytes8", true, [8]byte{7: 3}, "bytes8(uint64(0x3))"},
			},
		},
	}

	for _, c := range cases {
		cc := &testutil.Contract{}

		evnt := testutil.NewEvent("A")
		input := []string{}

		result := map[string]interface{}{}
		for indx, field := range c.fields {
			evnt.Add(field.typ, field.indx)
			input = append(input, field.valStr)
			result[fmt.Sprintf("val_%d", indx)] = field.val
		}

		cc.AddEvent(evnt)
		cc.EmitEvent("setA", "A", input...)

		// deploy the contract
		artifact, addr := s.DeployContract(cc)
		receipt := s.TxnTo(addr, "setA")

		// read the abi
		abi, err := NewABI(artifact.Abi)
		assert.NoError(t, err)

		// parse the logs
		found, err := ParseLog(abi.Events["A"].Inputs, receipt.Logs[0])
		assert.NoError(t, err)

		if !reflect.DeepEqual(found, result) {
			t.Log(found)
			t.Log(result)
			t.Fatal("not equal")
		}
	}
}
