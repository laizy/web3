package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
)

func Ensure(err error, msg ...interface{}) {
	if err != nil {
		if len(msg) == 0 {
			msg = []interface{}{""}
		}
		panic(fmt.Errorf("%v %v", err, msg[0]))
	}
}

func JsonString(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	Ensure(err)

	return string(b)
}

func JsonBytes(v interface{}) []byte {
	b, err := json.Marshal(v)
	Ensure(err)

	return b
}

func JsonStr(v interface{}) string {
	return string(JsonBytes(v))
}

func EnsureTrue(b bool, msg ...string) {
	if !b {
		panic("must be true:" + strings.Join(msg, ", "))
	}
}

func EnsureEqual(a, b interface{}, msg ...string) {
	if reflect.DeepEqual(a, b) == false {
		panic(fmt.Errorf("not equal: %v != %v, %s", a, b, strings.Join(msg, ", ")))
	}
}

func LoadJsonFile(file string, val interface{}) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(content, val); err != nil {
		if syntaxerr, ok := err.(*json.SyntaxError); ok {
			line := findLine(content, syntaxerr.Offset)
			return fmt.Errorf("JSON syntax error at %v:%v: %v", file, line, err)
		}
		return fmt.Errorf("JSON unmarshal error in %v: %v", file, err)
	}
	return nil
}

// findLine returns the line number for the given offset into data.
func findLine(data []byte, offset int64) (line int) {
	line = 1
	for i, r := range string(data) {
		if int64(i) >= offset {
			return
		}
		if r == '\n' {
			line++
		}
	}
	return
}
