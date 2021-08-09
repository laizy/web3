package abi

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/umbracle/go-web3/utils"

	"github.com/umbracle/go-web3"

	"github.com/stretchr/testify/assert"
)

func TestType(t *testing.T) {
	cases := []struct {
		s   string
		a   *ArgumentStr
		t   *Type
		err bool
	}{
		{
			s: "bool",
			a: simpleType("bool"),
			t: &Type{kind: KindBool, t: boolT, raw: "bool"},
		},
		{
			s: "uint32",
			a: simpleType("uint32"),
			t: &Type{kind: KindUInt, size: 32, t: uint32T, raw: "uint32"},
		},
		{
			s: "int32",
			a: simpleType("int32"),
			t: &Type{kind: KindInt, size: 32, t: int32T, raw: "int32"},
		},
		{
			s: "int32[]",
			a: simpleType("int32[]"),
			t: &Type{kind: KindSlice, t: reflect.SliceOf(int32T), raw: "int32[]", elem: &Type{kind: KindInt, size: 32, t: int32T, raw: "int32"}},
		},
		{
			s: "bytes[2]",
			a: simpleType("bytes[2]"),
			t: &Type{
				kind: KindArray,
				t:    reflect.ArrayOf(2, dynamicBytesT),
				raw:  "bytes[2]",
				size: 2,
				elem: &Type{
					kind: KindBytes,
					t:    dynamicBytesT,
					raw:  "bytes",
				},
			},
		},
		{
			s: "string[]",
			a: simpleType("string[]"),
			t: &Type{
				kind: KindSlice,
				t:    reflect.SliceOf(stringT),
				raw:  "string[]",
				elem: &Type{
					kind: KindString,
					t:    stringT,
					raw:  "string",
				},
			},
		},
		{
			s: "string[2]",
			a: simpleType("string[2]"),
			t: &Type{
				kind: KindArray,
				size: 2,
				t:    reflect.ArrayOf(2, stringT),
				raw:  "string[2]",
				elem: &Type{
					kind: KindString,
					t:    stringT,
					raw:  "string",
				},
			},
		},

		{
			s: "string[2][]",
			a: simpleType("string[2][]"),
			t: &Type{
				kind: KindSlice,
				t:    reflect.SliceOf(reflect.ArrayOf(2, stringT)),
				raw:  "string[2][]",
				elem: &Type{
					kind: KindArray,
					size: 2,
					t:    reflect.ArrayOf(2, stringT),
					raw:  "string[2]",
					elem: &Type{
						kind: KindString,
						t:    stringT,
						raw:  "string",
					},
				},
			},
		},
		{
			s: "tuple(int64 indexed arg0)",
			a: &ArgumentStr{
				Type: "tuple",
				Components: []*ArgumentStr{
					{
						Name:    "arg0",
						Type:    "int64",
						Indexed: true,
					},
				},
			},
			t: &Type{
				kind: KindTuple,
				raw:  "(int64)",
				t:    tupleT,
				tuple: []*TupleElem{
					{
						Name: "arg0",
						Elem: &Type{
							kind: KindInt,
							size: 64,
							t:    int64T,
							raw:  "int64",
						},
						Indexed: true,
					},
				},
			},
		},
		{
			s: "tuple(int64 arg_0)[2]",
			a: &ArgumentStr{
				Type: "tuple[2]",
				Components: []*ArgumentStr{
					{
						Name: "arg_0",
						Type: "int64",
					},
				},
			},
			t: &Type{
				kind: KindArray,
				size: 2,
				raw:  "(int64)[2]",
				t:    reflect.ArrayOf(2, tupleT),
				elem: &Type{
					kind: KindTuple,
					raw:  "(int64)",
					t:    tupleT,
					tuple: []*TupleElem{
						{
							Name: "arg_0",
							Elem: &Type{
								kind: KindInt,
								size: 64,
								t:    int64T,
								raw:  "int64",
							},
						},
					},
				},
			},
		},
		{
			s: "tuple(int64 a)[]",
			a: &ArgumentStr{
				Type: "tuple[]",
				Components: []*ArgumentStr{
					{
						Name: "a",
						Type: "int64",
					},
				},
			},
			t: &Type{
				kind: KindSlice,
				raw:  "(int64)[]",
				t:    reflect.SliceOf(tupleT),
				elem: &Type{
					kind: KindTuple,
					raw:  "(int64)",
					t:    tupleT,
					tuple: []*TupleElem{
						{
							Name: "a",
							Elem: &Type{
								kind: KindInt,
								size: 64,
								t:    int64T,
								raw:  "int64",
							},
						},
					},
				},
			},
		},
		{
			s: "tuple(int32 indexed arg0, tuple(int32 c) b_2)",
			a: &ArgumentStr{
				Type: "tuple",
				Components: []*ArgumentStr{
					{
						Name:    "arg0",
						Type:    "int32",
						Indexed: true,
					},
					{
						Name: "b_2",
						Type: "tuple",
						Components: []*ArgumentStr{
							{
								Name: "c",
								Type: "int32",
							},
						},
					},
				},
			},
			t: &Type{
				kind: KindTuple,
				t:    tupleT,
				raw:  "(int32,(int32))",
				tuple: []*TupleElem{
					{
						Name: "arg0",
						Elem: &Type{
							kind: KindInt,
							size: 32,
							t:    int32T,
							raw:  "int32",
						},
						Indexed: true,
					},
					{
						Name: "b_2",
						Elem: &Type{
							kind: KindTuple,
							t:    tupleT,
							raw:  "(int32)",
							tuple: []*TupleElem{
								{
									Name: "c",
									Elem: &Type{
										kind: KindInt,
										size: 32,
										t:    int32T,
										raw:  "int32",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			s: "tuple()",
			a: &ArgumentStr{
				Type:       "tuple",
				Components: []*ArgumentStr{},
			},
			t: &Type{
				kind:  KindTuple,
				raw:   "()",
				t:     tupleT,
				tuple: []*TupleElem{},
			},
		},
		{
			s:   "int[[",
			err: true,
		},
		{
			s:   "int",
			err: true,
		},
		{
			s:   "tuple[](a int32)",
			err: true,
		},
		{
			s:   "int32[a]",
			err: true,
		},
		{
			s:   "tuple(a int32",
			err: true,
		},
		{
			s:   "tuple(a int32,",
			err: true,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			e0, err := NewType(c.s)
			if err != nil && !c.err {
				t.Fatal(err)
			}
			if err == nil && c.err {
				t.Fatal("it should have failed")
			}

			if !c.err {
				e1, err := NewTypeFromArgument(c.a)
				if err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(c.t, e0) {

					// fmt.Println(c.t.t)
					// fmt.Println(e0.t)

					t.Fatal("bad new type")
				}
				if !reflect.DeepEqual(c.t, e1) {
					t.Fatal("bad")
				}
			}
		})
	}
}

func TestSize(t *testing.T) {
	cases := []struct {
		Input string
		Size  int
	}{
		{
			"int32", 32,
		},
		{
			"int32[]", 32,
		},
		{
			"int32[2]", 32 * 2,
		},
		{
			"int32[2][2]", 32 * 2 * 2,
		},
		{
			"string", 32,
		},
		{
			"string[]", 32,
		},
		{
			"tuple(uint8 a, uint32 b)[1]",
			64,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			tt, err := NewType(c.Input)
			if err != nil {
				t.Fatal(err)
			}

			size := getTypeSize(tt)
			if size != c.Size {
				t.Fatalf("expected size %d but found %d", c.Size, size)
			}
		})
	}
}

func simpleType(s string) *ArgumentStr {
	return &ArgumentStr{
		Type: s,
	}
}

type LiquidityPoolView struct {
	Pid             *big.Int     `json:"pid"`
	LpToken         web3.Address `json:"lpToken"`
	AllocPoint      *big.Int     `json:"allocPoint"`
	LastRewardBlock *big.Int     `json:"lastRewardBlock"`
	RewardsPerBlock *big.Int     `json:"rewardsPerBlock"`
	AccKstPerShare  *big.Int     `json:"accKstPerShare"`
	AllocKstAmount  *big.Int     `json:"allocKstAmount"`
	AccKstAmount    *big.Int     `json:"accKstAmount"`
	TotalAmount     *big.Int     `json:"totalAmount"`
	Token0          web3.Address `json:"token0"`
	Symbol0         string       `json:"symbol0"`
	Name0           string       `json:"name0"`
	Decimals0       uint8        `json:"decimals0"`
	Token1          web3.Address `json:"token1"`
	Symbol1         string       `json:"symbol1"`
	Name1           string       `json:"name1"`
	Decimals1       uint8        `json:"decimals1"`
}

var TypeStr = "tuple(tuple(uint256 pid, address lptoken, uint256 allocpoint, uint256 lastrewardblock, uint256 rewardsperblock, uint256 acckstpershare, uint256 allockstamount, uint256 acckstamount, uint256 totalamount, address token0, string symbol0, string name0, uint8 decimals0, address token1, string symbol1, string name1, uint8 decimals1)[] views)"

func TestNewType(t *testing.T) {
	typ, err := NewType(TypeStr)
	assert.Nil(t, err)

	views := map[string][]*LiquidityPoolView{
		"views": {
			&LiquidityPoolView{
				Pid:             big.NewInt(1),
				LpToken:         web3.Address{},
				AllocPoint:      big.NewInt(2),
				LastRewardBlock: big.NewInt(3),
				RewardsPerBlock: big.NewInt(4),
				AccKstPerShare:  big.NewInt(5),
				AllocKstAmount:  big.NewInt(6),
				AccKstAmount:    big.NewInt(7),
				TotalAmount:     big.NewInt(8),
				Token0:          web3.Address{},
				Symbol0:         "symbol0",
				Name0:           "name0",
				Decimals0:       9,
				Token1:          web3.Address{},
				Symbol1:         "symbol1",
				Name1:           "name1",
				Decimals1:       10,
			},
		},
	}
	fmt.Println(typ.String())
	buf, err := typ.Encode(views)
	assert.Nil(t, err)
	val, err := typ.Decode(buf)
	assert.Nil(t, err)

	fmt.Println(utils.JsonString(val))

	output := "0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000700000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000042000000000000000000000000000000000000000000000000000000000000007600000000000000000000000000000000000000000000000000000000000000aa00000000000000000000000000000000000000000000000000000000000000de00000000000000000000000000000000000000000000000000000000000001120000000000000000000000000000000000000000000000000000000000000146000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002a20f39354702fadf7d2087edb8c0730bca87ca7000000000000000000000000000000000000000000000000000000000000001e000000000000000000000000000000000000000000000000000000000047a08000000000000000000000000000000000000000000000000029a2241af62c00000000000000000000000000000000000000000000000000000000004f1abd01b900000000000000000000000000000000000000000000664ed708a01f5f35845700000000000000000000000000000000000000000000a8364f448b9c53cf9e43000000000000000000000000000000000000000000611111c22bee65e01f8f220000000000000000000000000000000000000000008be2fd7dff27e4cde2a893000000000000000000000000382bb369d343125bfb2117af9c149795c6c65c5000000000000000000000000000000000000000000000000000000000000002400000000000000000000000000000000000000000000000000000000000000280000000000000000000000000000000000000000000000000000000000000001200000000000000000000000054e4622dc504176b3bb432dccaf504569699a7ff00000000000000000000000000000000000000000000000000000000000002c000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000012000000000000000000000000000000000000000000000000000000000000000455534454000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004555344540000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044254434b0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044254434b000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000e611e389c51b5772a5fd715b8d638d9c98b640ba000000000000000000000000000000000000000000000000000000000000001e000000000000000000000000000000000000000000000000000000000047a08000000000000000000000000000000000000000000000000029a2241af62c000000000000000000000000000000000000000000000000000000000020be9fc7350000000000000000000000000000000000000000000058b6fafbf0c4287e78fd00000000000000000000000000000000000000000000a8346c442d57d337399500000000000000000000000000000000000000000063c18d211581f3498c6d8c0000000000000000000000000000000000000000009b39591dd8f053f004e506000000000000000000000000382bb369d343125bfb2117af9c149795c6c65c50000000000000000000000000000000000000000000000000000000000000024000000000000000000000000000000000000000000000000000000000000002800000000000000000000000000000000000000000000000000000000000000012000000000000000000000000ef71ca2ee68f45b9ad6f72fbdb33d707b872315c00000000000000000000000000000000000000000000000000000000000002c000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000012000000000000000000000000000000000000000000000000000000000000000455534454000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004555344540000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044554484b0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044554484b0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000007d7b8526852e7ed0380351a56517b381dd3516a3000000000000000000000000000000000000000000000000000000000000001e000000000000000000000000000000000000000000000000000000000047a08000000000000000000000000000000000000000000000000029a2241af62c00000000000000000000000000000000000000000000000000000000005760b175cf00000000000000000000000000000000000000000000503325d2711ab1d6731b00000000000000000000000000000000000000000000a837ca52fa4c85facbfa00000000000000000000000000000000000000000050d683e6f90596862802890000000000000000000000000000000000000000008894d35c2121fc50a0d77e00000000000000000000000054e4622dc504176b3bb432dccaf504569699a7ff000000000000000000000000000000000000000000000000000000000000024000000000000000000000000000000000000000000000000000000000000002800000000000000000000000000000000000000000000000000000000000000012000000000000000000000000ef71ca2ee68f45b9ad6f72fbdb33d707b872315c00000000000000000000000000000000000000000000000000000000000002c00000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000001200000000000000000000000000000000000000000000000000000000000000044254434b0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044254434b0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044554484b0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044554484b000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003000000000000000000000000d346967e8874b9c4dcdd543a88ae47ee8c8bd21f000000000000000000000000000000000000000000000000000000000000001e000000000000000000000000000000000000000000000000000000000047a08400000000000000000000000000000000000000000000000029a2241af62c000000000000000000000000000000000000000000000000000000061c64b77f5d3e0000000000000000000000000000000000000000000057737f0ff14a17e799e400000000000000000000000000000000000000000000aa05d580d19ca4eedebd000000000000000000000000000000000000000000598790d0bd5fadcda82fa7000000000000000000000000000000000000000000a0535acaeee599f136e3de000000000000000000000000382bb369d343125bfb2117af9c149795c6c65c500000000000000000000000000000000000000000000000000000000000000240000000000000000000000000000000000000000000000000000000000000028000000000000000000000000000000000000000000000000000000000000000120000000000000000000000008f8526dbfd6e38e3d8307702ca8469bae6c56c1500000000000000000000000000000000000000000000000000000000000002c00000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000001200000000000000000000000000000000000000000000000000000000000000045553445400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000455534454000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004574f4b5400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000b57726170706564204f4b540000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000aea843a6715d091f7a5644ed4fcca479820bd18a000000000000000000000000000000000000000000000000000000000000001e000000000000000000000000000000000000000000000000000000000047a08000000000000000000000000000000000000000000000000029a2241af62c000000000000000000000000000000000000000000000000000117cd7d902a7d266400000000000000000000000000000000000000000000661ea6b4beb3d935251d00000000000000000000000000000000000000000000b39d845c66b4321b9b8300000000000000000000000000000000000000000053480bdcd84b67fe38068900000000000000000000000000000000000000000080bff8f0880daf4ab765150000000000000000000000008f8526dbfd6e38e3d8307702ca8469bae6c56c15000000000000000000000000000000000000000000000000000000000000024000000000000000000000000000000000000000000000000000000000000002800000000000000000000000000000000000000000000000000000000000000012000000000000000000000000df54b6c6195ea4d948d03bfd818d365cf175cfc200000000000000000000000000000000000000000000000000000000000002c0000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000000120000000000000000000000000000000000000000000000000000000000000004574f4b5400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000b57726170706564204f4b5400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000034f4b42000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000034f4b420000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000500000000000000000000000089824289ae1d431aef91bb39d666f6d0f635e1b9000000000000000000000000000000000000000000000000000000000000001e000000000000000000000000000000000000000000000000000000000047a08000000000000000000000000000000000000000000000000029a2241af62c00000000000000000000000000000000000000000000000000000000003f29d045c700000000000000000000000000000000000000000000574d6c2ce67d4cb3a48b00000000000000000000000000000000000000000000a8337356d85247e22e55000000000000000000000000000000000000000000515a971f2cbf54c16f9105000000000000000000000000000000000000000000866d3146ee20d95d846a45000000000000000000000000382bb369d343125bfb2117af9c149795c6c65c50000000000000000000000000000000000000000000000000000000000000024000000000000000000000000000000000000000000000000000000000000002800000000000000000000000000000000000000000000000000000000000000012000000000000000000000000df54b6c6195ea4d948d03bfd818d365cf175cfc200000000000000000000000000000000000000000000000000000000000002c000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000012000000000000000000000000000000000000000000000000000000000000000455534454000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004555344540000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000034f4b42000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000034f4b420000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000600000000000000000000000084ee6a98990010fe87d2c79822763fca584418e90000000000000000000000000000000000000000000000000000000000000005000000000000000000000000000000000000000000000000000000000047a08000000000000000000000000000000000000000000000000006f05b59d3b20000000000000000000000000000000000000000000000000000000000026bd2e55e000000000000000000000000000000000000000000000ddfcc4dcd07b2a2d4ae000000000000000000000000000000000000000000001912b7529b4f8644d7a300000000000000000000000000000000000000000020f7641b28ebae7805f29f000000000000000000000000000000000000000000456819af06902b00649568000000000000000000000000382bb369d343125bfb2117af9c149795c6c65c50000000000000000000000000000000000000000000000000000000000000024000000000000000000000000000000000000000000000000000000000000002800000000000000000000000000000000000000000000000000000000000000012000000000000000000000000ab0d1578216a545532882e420a8c61ea07b00b1200000000000000000000000000000000000000000000000000000000000002c000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000012000000000000000000000000000000000000000000000000000000000000000455534454000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004555344540000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000034b53540000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000b4b5377617020546f6b656e000000000000000000000000000000000000000000"

	raw, err := hex.DecodeString(output[2:])
	assert.Nil(t, err)
	val, err = typ.Decode(raw)
	assert.Nil(t, err)

	fmt.Println(utils.JsonString(val))
}
