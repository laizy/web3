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

//transfer indexed arg to hash
func encodeTopicArg(str interface{}) string {
	arg := encodeArg(str)
	return transferToTopic(arg)
}

func transferToTopic(s string) string {
	if s == "string" || s == "[]byte" {
		s = "web3.Hash"
	}
	return s
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

func buildSignature(name string, typ *abi.Type) string {
	types := make([]string, len(typ.TupleElems()))
	for i, input := range typ.TupleElems() {
		types[i] = input.Elem.String()
	}
	return fmt.Sprintf("%v(%v)", name, strings.Join(types, ","))
}

func nameToKey(name string, index int) string {
	return abi.NameToKey(name, index)
}
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"title":               strings.Title,
		"clean":               cleanName,
		"arg":                 encodeArg,
		"outputArg":           outputArg,
		"funcName":            funcName,
		"tupleElems":          tupleElems,
		"tupleLen":            tupleLen,
		"toCamelCase":         toCamelCase,
		"nameToKey":           nameToKey,
		"topic":               encodeTopicArg,
		"sig":                 buildSignature,
		"getTopicFilterParam": getTopicFilterParam,
		"getTopicFilterInput": getTopicFilterInput,
		"getFilterEventParam": getFilterEventParam,
	}
}

//optimizeEvent change inner empty name to arg%d.
func optimizeEvent(event *abi.Event) *abi.Event {
	ev := event.Copy()
	for j, e := range ev.Inputs.TupleElems() {
		if e.Name == "" {
			e.Name = fmt.Sprintf("arg%d", j)
		}
	}
	return event
}

//optimizeInput change inner empty name to arg%d.
func optimizeInput(m *abi.Method) *abi.Method {
	m = m.Copy()
	for j, e := range m.Inputs.TupleElems() {
		if e.Name == "" {
			e.Name = fmt.Sprintf("arg%d", j)
		}
	}
	return m
}

func getTopicFilterParam(event *abi.Event) string {
	var params []string
	for _, v := range event.Inputs.TupleElems() {
		if v.Indexed {
			params = append(params, fmt.Sprintf("%s []%s", cleanName(v.Name), encodeTopicArg(v)))
		}
	}
	return strings.Join(params, ",")
}

func getTopicFilterInput(event *abi.Event) string {
	var params []string
	for _, v := range event.Inputs.TupleElems() {
		if v.Indexed {
			params = append(params, cleanName(v.Name))
		}
	}
	return strings.Join(params, ",")
}

func getFilterEventParam(event *abi.Event) string {
	var params []string
	for _, v := range event.Inputs.TupleElems() {
		if v.Indexed {
			params = append(params, fmt.Sprintf("%s []%s", cleanName(v.Name), encodeTopicArg(v)))
		}
	}
	params = append(params, "startBlock uint64", "endBlock ...uint64")
	return strings.Join(params, ",")
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
	"github.com/mitchellh/mapstructure"
)

var (
	_ = json.Unmarshal
	_ = big.NewInt
	_ = fmt.Printf
	_ = utils.JsonStr
	_ = mapstructure.Decode
)

{{$cname := .Name}}
// {{.Name}} is a solidity contract
type {{.Name}} struct {
	c *contract.Contract
}
{{if .Contract.Bin}}
// Deploy{{.Name}} deploys a new {{.Name}} contract
func Deploy{{.Name}}(provider *jsonrpc.Client, from web3.Address {{if .Abi.Constructor}}{{range $index, $val := tupleElems .Abi.Constructor.Inputs}}, {{if .Name}}{{clean .Name}}{{else}}val{{$index}}{{end}} {{arg .}} {{end}}{{end}}) *contract.Txn {
	return contract.DeployContract(provider, from, abi{{.Name}}, bin{{.Name}}{{if .Abi.Constructor}} {{range $index, $val := tupleElems .Abi.Constructor.Inputs}}, {{if .Name}}{{clean .Name}}{{else}}val{{$index}}{{end}}{{end}}{{end}})
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

	out, err = {{$.Ptr}}.c.Call("{{$key}}", web3.EncodeBlock(block...){{range $index, $val := tupleElems .Inputs}}, {{if .Name}}{{clean .Name}}{{else}}val{{$index}}{{end}}{{end}})
	if err != nil {
		return
	}

	// decode outputs
	{{range $index,$val := tupleElems .Outputs}}
	if err = mapstructure.Decode(out["{{nameToKey .Name $index}}"],&retval{{$index}}); err != nil {
		err = fmt.Errorf("failed to encode output at index {{$index}}")
	}{{end}}

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

func({{$.Ptr}} *{{$.Name}}) {{title .Name}}TopicFilter({{getTopicFilterParam $value}})[][]web3.Hash{
	{{range $index, $input := tupleElems .Inputs}}
    {{if .Indexed}}var {{clean .Name}}Rule []interface{}
    for _, {{.Name}}Item := range {{clean .Name}} {
		{{clean .Name}}Rule = append({{clean .Name}}Rule, {{.Name}}Item)
	}
	{{end}}{{end}}

	var query [][]interface{}
	query = append(query,[]interface{}{ {{title .Name}}EventID} {{range $index, $input := tupleElems .Inputs}} {{if .Indexed}}, {{clean .Name}}Rule {{end}}{{end}})

	topics, err := contract.MakeTopics(query...)
	utils.Ensure(err)

	return topics
}

func ({{$.Ptr}} *{{$.Name}}) Filter{{title .Name}}Event({{getFilterEventParam $value}})([]*{{title .Name}}Event, error){
	topic :={{$.Ptr}}.{{title .Name}}TopicFilter({{getTopicFilterInput $value}})	

	logs, err := {{$.Ptr}}.c.FilterLogsWithTopic(topic, startBlock, endBlock...)
	if err != nil {
		return nil, err
	}
	res := make([]*{{title .Name}}Event, 0)
	evts := {{$.Ptr}}.c.Abi.Events["{{.Name}}"]
	for _, log := range logs {
		args, err := evts.ParseLog(log)
		if err != nil {
			return nil, err
		}
		var evtItem {{title .Name}}Event
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
