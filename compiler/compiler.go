package compiler

import "strings"

//
//import (
//	"fmt"
//	"strings"
//)
//
//type factory func(path string) Compiler
//
//var compilers = map[string]factory{
//	"solidity": NewSolidityCompiler,
//	"vyper":    NewVyperCompiler,
//}
//
//// Compiler is an Ethereum compiler
//type Compiler interface {
//	// Compile compiles a file
//	Compile(files ...string) (map[string]*Artifact, error)
//}
//
//// NewCompiler instantiates a new compiler
//func NewCompiler(name string, path string) (Compiler, error) {
//	factory, ok := compilers[name]
//	if !ok {
//		return nil, fmt.Errorf("unknown compiler '%s'", name)
//	}
//	return factory(path), nil
//}
//
//
func NewArtifact(abi, bin, binRuntime string) *Artifact {
	return &Artifact{
		Abi:        abi,
		Bin:        unifyHexString(bin),
		BinRuntime: unifyHexString(binRuntime),
	}
}

func unifyHexString(hex string) string {
	if !strings.HasPrefix(hex, "0x") {
		hex = "0x" + hex
	}
	return hex
}
