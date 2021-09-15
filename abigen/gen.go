package abigen

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/laizy/web3/abi"
	"github.com/laizy/web3/compiler"
)

type Config struct {
	Package string
	Output  string
	Name    string
}

func cleanName(str string) string {
	return handleSnakeCase(strings.Trim(str, "_"))
}

// ToCamelCase converts an under-score string to a camel-case string
func toCamelCase(input string) string {
	return strings.Title(strings.Trim(input, "_"))
}

func outputArg(str string) string {
	if str == "" {

	}
	return str
}

func handleSnakeCase(str string) string {
	if !strings.Contains(str, "_") {
		return str
	}

	spl := strings.Split(str, "_")
	res := ""
	for indx, elem := range spl {
		if indx != 0 {
			elem = strings.Title(elem)
		}
		res += elem
	}
	return res
}

func funcName(str string) string {
	return strings.Title(handleSnakeCase(str))
}

func encodeSimpleArg(typ *abi.Type) string {
	switch typ.Kind() {
	case abi.KindAddress:
		return "web3.Address"

	case abi.KindString:
		return "string"

	case abi.KindBool:
		return "bool"

	case abi.KindInt:
		return typ.GoType().String()

	case abi.KindUInt:
		return typ.GoType().String()

	case abi.KindFixedBytes:
		return fmt.Sprintf("[%d]byte", typ.Size())

	case abi.KindBytes:
		return "[]byte"

	case abi.KindSlice:
		return "[]" + encodeSimpleArg(typ.Elem())

	case abi.KindTuple:
		return typ.RawName()
	default:
		return fmt.Sprintf("input not done for type: %s", typ.String())
	}
}

func encodeArg(str interface{}) string {
	arg, ok := str.(*abi.TupleElem)
	if !ok {
		panic("bad 1")
	}
	return encodeSimpleArg(arg.Elem)
}

func tupleLen(tuple interface{}) interface{} {
	if isNil(tuple) {
		return 0
	}
	arg, ok := tuple.(*abi.Type)
	if !ok {
		panic("bad tuple")
	}
	return len(arg.TupleElems())
}

func tupleElems(tuple interface{}) (res []interface{}) {
	if isNil(tuple) {
		return
	}

	arg, ok := tuple.(*abi.Type)
	if !ok {
		panic("bad tuple")
	}
	for _, i := range arg.TupleElems() {
		res = append(res, i)
	}
	return
}

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"title":       strings.Title,
		"clean":       cleanName,
		"arg":         encodeArg,
		"outputArg":   outputArg,
		"funcName":    funcName,
		"tupleElems":  tupleElems,
		"tupleLen":    tupleLen,
		"toCamelCase": toCamelCase,
	}
}

func isNil(c interface{}) bool {
	return c == nil || (reflect.ValueOf(c).Kind() == reflect.Ptr && reflect.ValueOf(c).IsNil())
}

func GenCode(artifacts map[string]*compiler.Artifact, config *Config) error {
	def, err := LoadStructDef(config.Output)
	if err != nil {
		return fmt.Errorf("read struct from json: %w", err)
	}

	generator := NewGenerator(config, artifacts)
	result, err := generator.Gen()
	if err != nil {
		return fmt.Errorf("generateGen: %v", err)
	}
	for _, file := range result.BinFiles {
		if err := ioutil.WriteFile(filepath.Join(generator.Config.Output, file.FileName+"_artifacts.go"), file.Code, 0644); err != nil {
			return err
		}
	}
	for _, file := range result.AbiFiles {
		if err := ioutil.WriteFile(filepath.Join(generator.Config.Output, file.FileName+".go"), file.Code, 0644); err != nil {
			return err
		}
	}

	for _, artifact := range artifacts {
		// parse abi
		abi, err := abi.NewABI(artifact.Abi)
		if err != nil {
			return err
		}
		def.ExtractFromAbi(abi)
	}

	return def.RenderGoCodeToFile(config.Package, config.Output)
}

var templateAbiStr = `package {{.Config.Package}}

import (
    "encoding/json"
	"fmt"
	"math/big"

	"github.com/laizy/web3"
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/jsonrpc"
	"github.com/laizy/web3/utils"
)

var (
	_ = big.NewInt
	_ = fmt.Printf
	_ = utils.JsonStr
)

// {{.Name}} is a solidity contract
type {{.Name}} struct {
	c *contract.Contract
}
{{if .Contract.Bin}}
// Deploy{{.Name}} deploys a new {{.Name}} contract
func Deploy{{.Name}}(provider *jsonrpc.Client, from web3.Address, args ...interface{}) *contract.Txn {
	return contract.DeployContract(provider, from, abi{{.Name}}, bin{{.Name}}, args...)
}
{{end}}
// New{{.Name}} creates a new instance of the contract at a specific address
func New{{.Name}}(addr web3.Address, provider *jsonrpc.Client) *{{.Name}} {
	return &{{.Name}}{c: contract.NewContract(addr, abi{{.Name}}, provider)}
}

// Contract returns the contract object
func ({{.Ptr}} *{{.Name}}) Contract() *contract.Contract {
	return {{.Ptr}}.c
}

// calls
{{range $key, $value := .Abi.Methods}}{{if .Const}}
// {{funcName $key}} calls the {{$key}} method in the solidity contract
func ({{$.Ptr}} *{{$.Name}}) {{funcName $key}}({{range $index, $val := tupleElems .Inputs}}{{if .Name}}{{clean .Name}}{{else}}val{{$index}}{{end}} {{arg .}}, {{end}}block ...web3.BlockNumber) ({{range $index, $val := tupleElems .Outputs}}retval{{$index}} {{arg .}}, {{end}}err error) {
	var out map[string]interface{}
	_ = out // avoid not used compiler error

	{{ $length := tupleLen .Outputs }}{{ if ne $length 0 }}var ok bool{{ end }}

	out, err = {{$.Ptr}}.c.Call("{{$key}}", web3.EncodeBlock(block...){{range $index, $val := tupleElems .Inputs}}, {{if .Name}}{{clean .Name}}{{else}}val{{$index}}{{end}}{{end}})
	if err != nil {
		return
	}

	// decode outputs
	{{range $index, $val := tupleElems .Outputs}}retval{{$index}}, ok = out["{{if .Name}}{{.Name}}{{else}}{{$index}}{{end}}"].({{arg .}})
	if !ok {
		err = fmt.Errorf("failed to encode output at index {{$index}}")
		return
	}
{{end}}
	return
}
{{end}}{{end}}

// txns
{{range $key, $value := .Abi.Methods}}{{if not .Const}}
// {{funcName $key}} sends a {{$key}} transaction in the solidity contract
func ({{$.Ptr}} *{{$.Name}}) {{funcName $key}}({{range $index, $input := tupleElems .Inputs}}{{if $index}}, {{end}}{{clean .Name}} {{arg .}}{{end}}) *contract.Txn {
	return {{$.Ptr}}.c.Txn("{{$key}}"{{range $index, $elem := tupleElems .Inputs}}, {{clean $elem.Name}}{{end}})
}
{{end}}{{end}}

// events
{{range $key, $value := .Abi.Events}}{{if not .Anonymous}}
//{{.Name}}Event
type {{.Name}}Event struct { {{range $index, $input := tupleElems $value.Inputs}}
    {{toCamelCase .Name}} {{arg .}}{{end}}
	Raw *web3.Log
}

func ({{$.Ptr}} *{{$.Name}}) Filter{{.Name}}Event(opts *web3.FilterOpts{{range $index, $input := tupleElems .Inputs}}{{if .Indexed}}, {{clean .Name}} []{{arg .}}{{end}}{{end}})([]*{{.Name}}Event, error){
	{{range $index, $input := tupleElems .Inputs}}
    {{if .Indexed}}var {{.Name}}Rule []interface{}
    for _, {{.Name}}Item := range {{clean .Name}} {
		{{.Name}}Rule = append({{.Name}}Rule, {{.Name}}Item)
	}
	{{end}}{{end}}
	logs, err := {{$.Ptr}}.c.FilterLogs(opts, "{{.Name}}"{{range $index, $input := tupleElems .Inputs}}{{if .Indexed}}, {{.Name}}Rule{{end}}{{end}})
	if err != nil {
		return nil, err
	}
	res := make([]*{{.Name}}Event, 0)
	evts := {{$.Ptr}}.c.Abi.Events["{{.Name}}"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem {{.Name}}Event
		err = json.Unmarshal([]byte(utils.JsonStr(args)), &evtItem)
		if err != nil {
			return nil, err
		}
		evtItem.Raw = log
		res = append(res, &evtItem)
	}
	return res, nil
}
{{end}}{{end}}
`

var templateBinStr = `package {{.Config.Package}}

import (
	"encoding/hex"
	"fmt"

	"github.com/laizy/web3/abi"
)

var abi{{.Name}} *abi.ABI

// {{.Name}}Abi returns the abi of the {{.Name}} contract
func {{.Name}}Abi() *abi.ABI {
	return abi{{.Name}}
}

var bin{{.Name}} []byte
{{if .Contract.Bin}}
// {{.Name}}Bin returns the bin of the {{.Name}} contract
func {{.Name}}Bin() []byte {
	return bin{{.Name}}
}
{{end}}
var binRuntime{{.Name}} []byte
{{if .Contract.BinRuntime}}
// {{.Name}}BinRuntime returns the runtime bin of the {{.Name}} contract
func {{.Name}}BinRuntime() []byte {
	return binRuntime{{.Name}}
}
{{end}}
func init() {
	var err error
	abi{{.Name}}, err = abi.NewABI(abi{{.Name}}Str)
	if err != nil {
		panic(fmt.Errorf("cannot parse {{.Name}} abi: %v", err))
	}
	if len(bin{{.Name}}Str) != 0 {
		bin{{.Name}}, err = hex.DecodeString(bin{{.Name}}Str[2:])
		if err != nil {
			panic(fmt.Errorf("cannot parse {{.Name}} bin: %v", err))
		}
	}
	if len(binRuntime{{.Name}}Str) != 0 {
		binRuntime{{.Name}}, err = hex.DecodeString(binRuntime{{.Name}}Str[2:])
		if err != nil {
			panic(fmt.Errorf("cannot parse {{.Name}} bin runtime: %v", err))
		}
	}
}

var bin{{.Name}}Str = "{{.Contract.Bin}}"

var binRuntime{{.Name}}Str = "{{.Contract.BinRuntime}}"

var abi{{.Name}}Str = ` + "`" + `{{.Contract.Abi}}` + "`\n"
