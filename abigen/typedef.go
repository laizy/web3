package abigen

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/laizy/web3/abi"
	"github.com/laizy/web3/utils"
	"github.com/laizy/web3/utils/common"
)

var ErrConflictDef = errors.New("conflict struct definition with same name")

type FieldDef struct {
	Name string
	Type string
}

type StructDef struct {
	Name   string
	Fields []*FieldDef
}

type StructDefExtractor struct {
	Defs map[string]*StructDef `json:"definitions"`
}

func NewStructDefExtractor() *StructDefExtractor {
	return &StructDefExtractor{Defs: make(map[string]*StructDef)}
}

func (self *StructDefExtractor) extractNormal(typ *abi.Type) string {
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
		return "[]" + self.ExtractFromType(typ.Elem())

	default:
		return fmt.Sprintf("input not done for type: %s", typ.String())
	}
}

func (self *StructDefExtractor) ExtractFromType(typ *abi.Type) string {
	switch typ.Kind() {
	case abi.KindTuple:
		name := typ.RawName()
		s := &StructDef{Name: name}
		for _, ty := range typ.TupleElems() {
			goType := self.ExtractFromType(ty.Elem)
			name := ty.Name
			s.Fields = append(s.Fields, &FieldDef{Name: name, Type: goType})
		}
		if name == "" {
			return ""
		}
		if old, exist := self.Defs[name]; exist { // check if two struct have same name but different struct, panic.
			if !reflect.DeepEqual(s, old) {
				panic(ErrConflictDef)
			}
		}
		self.Defs[name] = s
		return name
	default:
		return self.extractNormal(typ)
	}
}

func (self *StructDefExtractor) ExtractFromAbi(abi *abi.ABI) *StructDefExtractor {
	if abi.Constructor != nil {
		self.ExtractFromType(abi.Constructor.Inputs)
	}
	for _, method := range abi.Methods {
		self.ExtractFromType(method.Inputs)
	}
	for _, event := range abi.Events {
		self.ExtractFromType(optimizeEvent(event).Inputs)
	}

	return self
}

func LoadStructDef(outputDir string) (*StructDefExtractor, error) {
	def := &StructDefExtractor{
		Defs: make(map[string]*StructDef),
	}
	file := filepath.Join(outputDir, "structs.json")
	if common.FileExist(file) {
		err := utils.LoadJsonFile(filepath.Join(outputDir, "structs.json"), def)
		if err != nil {
			return nil, err
		}
	}

	return def, nil
}

func (self *StructDefExtractor) RenderGoCodeToFile(packageName string, outputDir string) error {
	code, err := self.RenderGoCode(packageName)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(outputDir, "structs.go"), []byte(code), 0644); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(outputDir, "structs.json"), []byte(utils.JsonString(self)), 0644); err != nil {
		return err
	}

	return nil
}

func (self *StructDefExtractor) RenderGoCode(packageName string) (string, error) {
	tempStruct, err := template.New("eth-structs").Funcs(map[string]interface{}{"title": strings.Title}).Parse(templateStructStr)
	utils.Ensure(err)

	input := map[string]interface{}{
		"Package": packageName,
		"Structs": self.Defs,
	}
	var b bytes.Buffer
	if err := tempStruct.Execute(&b, input); err != nil {
		return "", err
	}
	code, err := format.Source(b.Bytes())
	utils.Ensure(err)

	return string(code), nil
}

var templateStructStr = `
package {{.Package}}

import (
	"fmt"
	"math/big"

	"github.com/laizy/web3"
)

var (
	_ = big.NewInt
	_ = fmt.Printf
)

{{$structs := .Structs}}
{{range $structs}}
type {{.Name}} struct {
{{range .Fields}}
{{title .Name}}   {{.Type}} {{end}}
}
{{end}}
`
