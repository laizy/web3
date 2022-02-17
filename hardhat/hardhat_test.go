package hardhat

import "testing"

var newJson = "{\n  \"abi\": [],\n  \"bytecode\": {\n    \"object\": \"0x60\",\n    \"sourceMap\": \"60:992:6:-:0;;;;;;;;;;;;;;;-1:-1:-1;;;60:992:6;;;;;;;;;;;;;;;;;\",\n    \"linkReferences\": {}\n  },\n  \"deployed_bytecode\": {\n    \"object\": \"0x7300\",\n    \"sourceMap\": \"60:992:6:-:0;;;;;;;;\",\n    \"linkReferences\": {}\n  }\n}"
var oldJson = "{\n  \"_format\": \"hh-sol-artifact-1\",\n  \"contractName\": \"xx\",\n  \"sourceName\": \"contracts/xx.sol\",\n  \"abi\": [],\n  \"bytecode\": \"0x6080\",\n  \"deployedBytecode\": \"0x6080\",\n  \"linkReferences\": {},\n  \"deployedLinkReferences\": {}\n}"

func TestDecodeOld(t *testing.T) {
	if _, err := decodeArtifact([]byte(oldJson)); err != nil {
		t.Fatalf("decode old err: %v", err)
	}
	if _, err := decodeArtifact([]byte(newJson)); err != nil {
		t.Fatalf("decode new err: %v", err)
	}
}
