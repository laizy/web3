package hardhat

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/laizy/web3/utils/common/hexutil"

	"github.com/laizy/web3/abi"
	"github.com/laizy/web3/registry"
	"github.com/laizy/web3/utils"
)


func GetArtifacts(artifactDirName ...string) (map[string]*Artifact, error) {
	name := ""
	if len(artifactDirName)!= 0 {
		name = artifactDirName[0]
	}
	
	pathes, err := getArtifactPathes(name)
	if err != nil {
		return nil, err
	}
	results := make(map[string]*Artifact)
	for name, path := range pathes {
		arti, err := getArtifactWithPath(path)
		if err != nil {
			return nil, err
		}
		results[name] = arti
	}

	return results, nil
}

func GetArtifact(name string, artifactDirName ...string) (*Artifact, error) {
	path, err := GetArtifactPath(name, artifactDirName...)
	if err != nil {
		return nil, err
	}

	return getArtifactWithPath(path)
}

func getArtifactWithPath(path string) (*Artifact, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	type artifact struct {
		ContractName     string        `json:"contractName"`
		SourceName       string        `json:"sourceName"`
		Abi              interface{}   `json:"abi"`
		Bytecode         hexutil.Bytes `json:"bytecode"`
		DeployedBytecode hexutil.Bytes `json:"deployedBytecode"`
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

func getArtifactPathes(artifactDirName string) (map[string]string, error) {
	if artifactDirName == "" {
		artifactDirName = "artifacts"
	}
	dir, err := GetProjectRoot()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	buildDir := filepath.Join(dir, artifactDirName)
	err = filepath.Walk(buildDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		base := filepath.Base(path)
		if strings.HasSuffix(base, ".dbg.json") {
			name := strings.TrimSuffix(base, ".dbg.json")
			contractFile := name + ".json"
			full := filepath.Join(filepath.Dir(path), contractFile)
			result[name] = full
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetArtifactPath(name string, artifactDirName ...string) (string, error) {
	artifactDir := ""
	if len(artifactDirName)!= 0 {
		artifactDir = artifactDirName[0]
	}
	pathes, err := getArtifactPathes(artifactDir)
	if err != nil {
		return "", err
	}

	path, ok := pathes[name]
	if !ok {
		return "", os.ErrNotExist
	}

	return path, nil
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
	ContractName     string        `json:"contractName"` // "DSProxy",
	SourceName       string        `json:"sourceName"`   // "contracts/proxy.sol"
	Abi              string        `json:"abi"`
	Bytecode         hexutil.Bytes `json:"bytecode"`         // 0x6080
	DeployedBytecode hexutil.Bytes `json:"deployedBytecode"` // 0x6080
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

func RegisterProjectEvents() error {
	artifacts, err := GetArtifacts()
	if err != nil {
		return err
	}
	for _, arti := range artifacts {
		artiAbi, err := abi.NewABI(arti.Abi)
		if err != nil {
			return err
		}
		registry.Instance().RegisterFromAbi(artiAbi)
	}

	return nil
}
