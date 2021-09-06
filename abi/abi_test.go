package abi

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbi(t *testing.T) {
	cases := []struct {
		Input  string
		Output *ABI
	}{
		{
			Input: `[
				{
					"name": "abc",
					"type": "function"
				}
			]`,
			Output: &ABI{
				Methods: map[string]*Method{
					"abc": {
						Name:    "abc",
						Inputs:  &Type{kind: KindTuple, raw: "tuple", tuple: []*TupleElem{}},
						Outputs: &Type{kind: KindTuple, raw: "tuple", tuple: []*TupleElem{}},
					},
				},
				Events: map[string]*Event{},
			},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			abi, err := NewABI(c.Input)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, abi.Methods, c.Output.Methods)
		})
	}
}

func TestAbi_HumanReadable(t *testing.T) {
	cases := []string{
		"event Transfer(address from, address to, uint256 amount)",
		"function symbol() returns (string)",
	}
	vv, err := NewABIFromList(cases)
	assert.NoError(t, err)

	fmt.Println(vv.Methods["symbol"].Inputs.String())
}
