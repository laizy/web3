package abigen

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/laizy/web3/abi"
	"github.com/laizy/web3/compiler"
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
		fileName := strings.ToLower(name)
		input := map[string]interface{}{
			"Ptr":      "_a",
			"Config":   g.Config,
			"Contract": artifact,
			"Abi":      abi,
			"Name":     name,
		}
		abiCode, err := GenCodeToBytes("eth-abi", g.funcMap, templateAbiStr, input)
		if err != nil {
			return Result{}, fmt.Errorf("genCodeToBytes: %v", err)
		}
		res.AbiFiles = append(res.AbiFiles, CodeFile{fileName, abiCode})

		binCode, err := GenCodeToBytes("eth-bin", g.funcMap, templateAbiStr, input)
		if err != nil {
			return Result{}, fmt.Errorf("genCodeToBytes: %v", err)
		}
		res.BinFiles = append(res.BinFiles, CodeFile{fileName, binCode})
	}

	return
}
