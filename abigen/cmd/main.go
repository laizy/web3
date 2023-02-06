package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/laizy/web3/utils"

	"github.com/laizy/web3/abigen"
	"github.com/laizy/web3/compiler"
	"github.com/laizy/web3/hardhat"
)

const (
	version = "0.1.0"
)

func main() {
	var source string
	var pckg string
	var output string
	var name string
	var onlyAbi bool

	flag.StringVar(&source, "source", "", "List of abi files")
	flag.StringVar(&pckg, "package", "main", "Name of the package")
	flag.StringVar(&output, "output", "", "Output directory")
	flag.StringVar(&name, "name", "", "name of the contract")
	flag.BoolVar(&onlyAbi, "abi", false, "only extract abi")

	flag.Parse()

	config := &abigen.Config{
		Package: pckg,
		Output:  output,
		Name:    name,
	}

	if _, err := os.Stat(output); errors.Is(err, os.ErrNotExist) { //create dir if not exist
		err := os.Mkdir(output, os.ModePerm)
		utils.Ensure(err)
	}

	if source == "" {
		fmt.Println(version)
		os.Exit(0)
	}

	sources := strings.Split(source, ",")
	for _, source := range sources {
		matches, err := filepath.Glob(source)
		if err != nil {
			fmt.Printf("Failed to read files: %v", err)
			os.Exit(1)
		}
		if len(matches) == 0 {
			fmt.Printf("ERROR: Have no file at: %v", sources)
			continue
		}

		for _, source := range matches {
			is, err := regexp.Match(".*\\.metadata\\.json$", []byte(source))
			utils.Ensure(err)
			if is {
				//ignore metadata file
				continue
			}
			artifacts, err := process(source, config)
			if err != nil {
				fmt.Printf("Failed to parse sources: %v", err)
				os.Exit(1)
			}
			if onlyAbi {
				for name, v := range artifacts {
					fmt.Println("name: ", name)
					if len(v.Abi) == 0 {
						fmt.Printf("No abi from %s\n", name)
						os.Exit(1)
					}
					filename := filepath.Join(output, name+"_abi.json")
					fmt.Println("write abi to: ", filename)
					if err := ioutil.WriteFile(filename, []byte(v.Abi), 0644); err != nil {
						panic(err)
					}
				}
			} else {
				if err := abigen.GenCode(artifacts, config); err != nil {
					fmt.Printf("Failed to generate sources: %v", err)
					os.Exit(1)
				}
			}
		}
	}
}

const (
	vyExt   = 0
	solExt  = 1
	abiExt  = 2
	jsonExt = 3
)

func process(sources string, config *abigen.Config) (map[string]*compiler.Artifact, error) {
	files := strings.Split(sources, ",")
	if len(files) == 0 {
		return nil, fmt.Errorf("input not found")
	}

	prev := -1
	for _, f := range files {
		var ext int
		switch extt := filepath.Ext(f); extt {
		case ".abi":
			ext = abiExt
		case ".sol":
			ext = solExt
		case ".vy", ".py":
			ext = vyExt
		case ".json":
			ext = jsonExt
		default:
			return nil, fmt.Errorf("file extension '%s' not found", extt)
		}

		if prev == -1 {
			prev = ext
		} else if ext != prev {
			return nil, fmt.Errorf("two file formats found")
		}
	}

	switch prev {
	case abiExt:
		return processAbi(files, config)
	case solExt:
		return processSolc(files)
	case vyExt:
		return processVyper(files)
	case jsonExt:
		return processJson(files)
	}

	return nil, nil
}

func processVyper(sources []string) (map[string]*compiler.Artifact, error) {
	c, err := compiler.NewCompiler("vyper", "vyper")
	if err != nil {
		return nil, err
	}
	raw, err := c.Compile(sources...)
	if err != nil {
		return nil, err
	}
	res := map[string]*compiler.Artifact{}
	for rawName, entry := range raw {
		_, name := filepath.Split(rawName)
		name = strings.TrimSuffix(name, ".vy")
		name = strings.TrimSuffix(name, ".v.py")
		res[strings.Title(name)] = entry
	}
	return res, nil
}

func processSolc(sources []string) (map[string]*compiler.Artifact, error) {
	c, err := compiler.NewCompiler("solidity", "solc")
	if err != nil {
		return nil, err
	}
	raw, err := c.Compile(sources...)
	if err != nil {
		return nil, err
	}
	res := map[string]*compiler.Artifact{}
	for rawName, entry := range raw {
		name := strings.Split(rawName, ":")[1]
		res[strings.Title(name)] = entry
	}
	return res, nil
}

func processAbi(sources []string, config *abigen.Config) (map[string]*compiler.Artifact, error) {
	artifacts := map[string]*compiler.Artifact{}

	for _, abiPath := range sources {
		content, err := ioutil.ReadFile(abiPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read abi file (%s): %v", abiPath, err)
		}

		// Use the name of the file to name the contract
		path, name := filepath.Split(abiPath)

		name = strings.TrimSuffix(name, filepath.Ext(name))
		binPath := filepath.Join(path, name+".bin")

		bin, err := ioutil.ReadFile(binPath)
		if err != nil {
			// bin not found
			bin = []byte{}
		}
		if len(sources) == 1 && config.Name != "" {
			name = config.Name
		}
		artifacts[strings.Title(name)] = &compiler.Artifact{
			Abi: string(content),
			Bin: string(bin),
		}
	}
	return artifacts, nil
}

func processJson(sources []string) (map[string]*compiler.Artifact, error) {
	artifacts := map[string]*compiler.Artifact{}
	for _, jsonPath := range sources {
		content, err := ioutil.ReadFile(jsonPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read abi file (%s): %v", jsonPath, err)
		}

		// Use the name of the file to name the contract
		_, name := filepath.Split(jsonPath)
		name = strings.TrimSuffix(name, ".json")

		art, err := hardhat.DecodeArtifact(content)
		if err != nil {
			return nil, err
		}

		artifacts[strings.Title(name)] = compiler.NewArtifact(string(art.Abi), art.Bytecode.String(), art.DeployedBytecode.String())
	}
	return artifacts, nil
}
