package hardhat

import (
	"encoding/json"
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/umbracle/go-web3/utils"
)

func GetArtifact(name string) (*Artifact, error) {
	path, err := GetArtifactPath(name)
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	type artifact struct {
		ContractName     string      `json:"contractName"`
		SourceName       string      `json:"sourceName"`
		Abi              interface{} `json:"abi"`
		Bytecode         string      `json:"bytecode"`
		DeployedBytecode string      `json:"deployedBytecode"`
	}
	var value artifact
	err = json.Unmarshal(buf, &value)
	if err != nil {
		return nil, err
	}

	return &Artifact{
		ContractName:     value.ContractName,
		SourceName:       value.SourceName,
		Abi:              utils.JsonStr(value.Abi),
		Bytecode:         value.Bytecode,
		DeployedBytecode: value.DeployedBytecode,
	}, nil
}

func GetArtifactPath(name string) (string, error) {
	dir, err := GetProjectRoot()
	if err != nil {
		return "", err
	}

	var breakError = errors.New("normal break")
	var result string
	buildDir := filepath.Join(dir, "artifacts")
	err = filepath.Walk(buildDir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Base(path) == name+".json" {
			result = path
			return breakError
		}

		return nil
	})
	if err == nil {
		return "", fs.ErrNotExist
	}
	if err == breakError {
		return result, nil
	}

	return "", err
}

func GetProjectRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if pathExists(filepath.Join(cwd, "hardhat.config.js")) ||
			pathExists(filepath.Join(cwd, "hardhat.config.ts")) ||
			pathExists(filepath.Join(cwd, "hardhat.config.json")) {
			return cwd, nil
		}
		parent := filepath.Dir(cwd)
		if parent == cwd {
			break
		}
		cwd = parent
	}

	return "", errors.New("project not found")
}

type Artifact struct {
	ContractName     string `json:"contractName"` // "DSProxy",
	SourceName       string `json:"sourceName"`   // "contracts/proxy.sol"
	Abi              string `json:"abi"`
	Bytecode         string `json:"bytecode"`         // 0x6080
	DeployedBytecode string `json:"deployedBytecode"` // 0x6080
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}
