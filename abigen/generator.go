package abigen

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"
	"text/template"

	"github.com/laizy/web3/abi"
	"github.com/laizy/web3/compiler"
	"github.com/laizy/web3/utils"
)

type Generator struct {
	Artifacts map[string]*compiler.Artifact
	Config    *Config
	funcMap   template.FuncMap
}

type CodeFile struct {
	FileName string
	Code     []byte
}
type Result struct {
	AbiFiles []CodeFile
	BinFiles []CodeFile
}

func NewGenerator(cfg *Config, artifacts map[string]*compiler.Artifact) *Generator {
	return &Generator{
		artifacts,
		cfg,
		FuncMap(),
	}
}

func (g *Generator) Gen() (res Result, err error) {
	for name, artifact := range g.Artifacts {
		// parse abi
		abi, err := abi.NewABI(artifact.Abi)
		if err != nil {
			return Result{}, err
		}

		for n, e := range abi.Events { //replace old event to get event nil arg's name
			abi.Events[n] = optimizeEvent(e)
		}
		// replace old input to get input nil arg's name
		for n, m := range abi.Methods {
			abi.Methods[n] = optimizeInput(m)
		}
		fileName := strings.ToLower(name)
		input := map[string]interface{}{
			"Ptr":      "_a",
			"Config":   g.Config,
			"Contract": artifact,
			"Abi":      abi,
			"Name":     name,
		}
		abiCode, err := genCodeToBytes("eth-abi", g.funcMap, templateAbiStr, input)
		if err != nil {
			return Result{}, fmt.Errorf("genCodeToBytes: %v", err)
		}
		res.AbiFiles = append(res.AbiFiles, CodeFile{fileName, abiCode})

		binCode, err := genCodeToBytes("eth-bin", g.funcMap, templateBinStr, input)
		if err != nil {
			return Result{}, fmt.Errorf("genCodeToBytes: %v", err)
		}
		res.BinFiles = append(res.BinFiles, CodeFile{fileName, binCode})
	}

	return
}

func genCodeToBytes(name string, funcMap template.FuncMap, temp string, input map[string]interface{}) ([]byte, error) {
	tempExt, err := template.New(name).Funcs(funcMap).Parse(temp)
	utils.Ensure(err)
	buffer := bytes.NewBuffer(nil)
	if err := tempExt.Execute(buffer, input); err != nil {
		return nil, err
	}
	b, err := format.Source(buffer.Bytes())
	if err != nil {
		return nil, err
	}
	return b, nil
}
