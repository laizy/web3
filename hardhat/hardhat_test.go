package hardhat

import "testing"

var forgeJson = `
{
  "abi": [],
  "bytecode": {
    "object": "0x60",
    "sourceMap": "60:992:6:-:0;;;;;;;;;;;;;;;-1:-1:-1;;;60:992:6;;;;;;;;;;;;;;;;;",
    "linkReferences": {}
  },
  "deployed_bytecode": {
    "object": "0x7300",
    "sourceMap": "60:992:6:-:0;;;;;;;;",
    "linkReferences": {}
  }
}
`
var hardhatJson = `
{
  "_format": "hh-sol-artifact-1",
  "contractName": "xx",
  "sourceName": "contracts/xx.sol",
  "abi": [],
  "bytecode": "0x6080",
  "deployedBytecode": "0x6080",
  "linkReferences": {},
  "deployedLinkReferences": {}
}
`

func TestDecodeOld(t *testing.T) {
	if _, err := decodeArtifact([]byte(hardhatJson)); err != nil {
		t.Fatalf("decode old err: %v", err)
	}
	if _, err := decodeArtifact([]byte(forgeJson)); err != nil {
		t.Fatalf("decode new err: %v", err)
	}
}
